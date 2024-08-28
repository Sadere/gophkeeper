package grpc

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand/v2"
	"os"
	"path/filepath"
	"sync"

	"github.com/Sadere/gophkeeper/cert"
	"github.com/Sadere/gophkeeper/internal/client/api"
	"github.com/Sadere/gophkeeper/internal/client/api/grpc/interceptor"
	"github.com/Sadere/gophkeeper/internal/client/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/Sadere/gophkeeper/pkg/constants"
	"github.com/Sadere/gophkeeper/pkg/convert"
	"github.com/Sadere/gophkeeper/pkg/model"

	pb "github.com/Sadere/gophkeeper/pkg/proto/keeper/v1"
)

type GRPCClient struct {
	config         *config.Config
	authClient     pb.AuthServiceClient
	secretsClient  pb.SecretsServiceClient
	notifyClient   pb.NotificationServiceClient
	accessToken    string
	masterPassword string
	chunkSize      int
	clientID       int32 // Unique ID to distinguish between multiple running clients for same user
	previews       sync.Map
}

var _ api.IApiClient = &GRPCClient{}

func NewGRPCClient(cfg *config.Config) (*GRPCClient, error) {
	var opts []grpc.DialOption

	newClient := GRPCClient{
		config:    cfg,
		chunkSize: constants.ChunkSize,
		clientID:  int32(rand.IntN(math.MaxInt32)),
	}

	// Unary interceptors
	opts = append(
		opts,
		grpc.WithChainUnaryInterceptor(
			interceptor.Timeout(constants.DefaultClientTimeout),
			interceptor.AddAuth(&newClient.accessToken, newClient.clientID),
		),
	)

	// Stream interceptor
	opts = append(
		opts,
		grpc.WithStreamInterceptor(interceptor.AddAuthStream(&newClient.accessToken, newClient.clientID)),
	)

	// TLS
	if cfg.EnableTLS {
		tlsCredential, err := loadTLSConfig("ca-cert.pem", "client-cert.pem", "client-key.pem")
		if err != nil {
			return nil, fmt.Errorf("failed to load TLS config: %w", err)
		}

		opts = append(
			opts,
			grpc.WithTransportCredentials(
				tlsCredential,
			),
		)
	} else {
		opts = append(
			opts,
			grpc.WithTransportCredentials(
				insecure.NewCredentials(),
			),
		)
	}

	// create gRPC client
	c, err := grpc.NewClient(
		cfg.ServerAddress,
		opts...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC client: %w", err)
	}

	// register services
	newClient.authClient = pb.NewAuthServiceClient(c)
	newClient.secretsClient = pb.NewSecretsServiceClient(c)
	newClient.notifyClient = pb.NewNotificationServiceClient(c)

	return &newClient, nil
}

func loadTLSConfig(caCertFile, clientCertFile, clientKeyFile string) (credentials.TransportCredentials, error) {
	// Read CA cert
	caPem, err := cert.Cert.ReadFile(caCertFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read CA cert: %w", err)
	}

	// Read client cert
	clientCertPEM, err := cert.Cert.ReadFile(clientCertFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read client cert: %w", err)
	}

	// Read client key
	clientKeyPEM, err := cert.Cert.ReadFile(clientKeyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read client key: %w", err)
	}

	// Create key pair
	clientCert, err := tls.X509KeyPair(clientCertPEM, clientKeyPEM)
	if err != nil {
		return nil, fmt.Errorf("failed to load x509 key pair: %w", err)
	}

	// Create cert pool and append CA's cert
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(caPem) {
		return nil, fmt.Errorf("failed to append CA cert to cert pool: %w", err)
	}

	// Create config
	config := &tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      certPool,
	}

	return credentials.NewTLS(config), nil
}

func (c *GRPCClient) Login(ctx context.Context, login string, password string) (string, error) {
	req := &pb.LoginV1Request{
		Login:    login,
		Password: password,
	}

	response, err := c.authClient.LoginV1(ctx, req)
	if err != nil {
		return "", err
	}

	c.accessToken = response.AccessToken
	c.masterPassword = password

	return response.AccessToken, nil
}

func (c *GRPCClient) Register(ctx context.Context, login string, password string) (string, error) {
	req := &pb.RegisterV1Request{
		Login:    login,
		Password: password,
	}

	response, err := c.authClient.RegisterV1(ctx, req)
	if err != nil {
		return "", err
	}

	c.accessToken = response.AccessToken

	return response.AccessToken, nil
}

func (c *GRPCClient) LoadPreviews(ctx context.Context) (model.SecretPreviews, error) {
	var previews model.SecretPreviews

	response, err := c.secretsClient.SecretPreviewsV1(ctx, &emptypb.Empty{})
	if err != nil {
		if e, ok := status.FromError(err); ok {
			// Return empty list if no secrets added yet
			if e.Code() == codes.NotFound {
				return nil, nil
			}
		}
		return nil, fmt.Errorf("failed to retrieve secrets: %w", err)
	}

	for _, pbPreview := range response.Previews {
		preview := convert.ProtoToPreview(pbPreview)

		// Set status
		v, ok := c.previews.Load(preview.ID)
		if ok {
			status, ok := v.(model.SecretPreviewStatus)
			if ok {
				preview.Status = status
			}
		}

		previews = append(previews, preview)
	}

	return previews, nil
}

func (c *GRPCClient) LoadSecret(ctx context.Context, ID uint64) (*model.Secret, error) {
	// form gRPC request
	request := &pb.GetUserSecretV1Request{
		MasterPassword: c.masterPassword,
		Id:             ID,
	}

	// performing gRPC call
	response, err := c.secretsClient.GetUserSecretV1(context.Background(), request)
	if err != nil {
		return nil, err
	}

	secret := convert.ProtoToSecret(response.Secret)

	return secret, nil
}

func (c *GRPCClient) SaveCredential(ctx context.Context, ID uint64, metadata, login, password string) error {
	// form gRPC request
	request := &pb.SaveUserSecretV1Request{
		MasterPassword: c.masterPassword,
		Secret: &pb.Secret{
			Id:        ID,
			CreatedAt: timestamppb.Now(),
			UpdatedAt: timestamppb.Now(),
			Metadata:  metadata,
			Type:      pb.SecretType_SECRET_TYPE_CREDENTIAL,
			Content: &pb.Secret_Credential{
				Credential: &pb.Credential{
					Login:    login,
					Password: password,
				},
			},
		},
	}

	// performing gRPC call
	_, err := c.secretsClient.SaveUserSecretV1(context.Background(), request)

	return err
}

func (c *GRPCClient) SaveText(ctx context.Context, ID uint64, metadata, text string) error {
	// form gRPC request
	request := &pb.SaveUserSecretV1Request{
		MasterPassword: c.masterPassword,
		Secret: &pb.Secret{
			Id:        ID,
			CreatedAt: timestamppb.Now(),
			UpdatedAt: timestamppb.Now(),
			Metadata:  metadata,
			Type:      pb.SecretType_SECRET_TYPE_TEXT,
			Content: &pb.Secret_Text{
				Text: &pb.Text{
					Text: text,
				},
			},
		},
	}

	// performing gRPC call
	_, err := c.secretsClient.SaveUserSecretV1(context.Background(), request)

	return err
}

func (c *GRPCClient) SaveCard(ctx context.Context, ID uint64, metadata, number string, expMonth, expYear, cvv uint32) error {
	// form gRPC request
	request := &pb.SaveUserSecretV1Request{
		MasterPassword: c.masterPassword,
		Secret: &pb.Secret{
			Id:        ID,
			CreatedAt: timestamppb.Now(),
			UpdatedAt: timestamppb.Now(),
			Metadata:  metadata,
			Type:      pb.SecretType_SECRET_TYPE_CARD,
			Content: &pb.Secret_Card{
				Card: &pb.Card{
					Number:   number,
					ExpMonth: expMonth,
					ExpYear:  expYear,
					Cvv:      cvv,
				},
			},
		},
	}

	// performing gRPC call
	_, err := c.secretsClient.SaveUserSecretV1(context.Background(), request)

	return err
}

func (c *GRPCClient) UploadFile(ctx context.Context, metadata, filePath string) error {
	// Open file
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}

	defer func() {
		err := f.Close()
		if err != nil {
			log.Println(fmt.Errorf("failed to close file: %w", err))
		}
	}()

	fileName := filepath.Base(filePath)

	stream, err := c.secretsClient.UploadFileV1(ctx)
	if err != nil {
		return err
	}

	buf := make([]byte, c.chunkSize)

	for {
		n, err := f.Read(buf)

		// File is done uploading
		if err == io.EOF {
			break
		}

		// I/O error
		if err != nil {
			return err
		}

		chunk := buf[:n]

		// Send chunk
		err = stream.Send(&pb.UploadFileV1Request{
			Metadata:       metadata,
			FileName:       fileName,
			MasterPassword: c.masterPassword,
			Chunk:          chunk,
		})
		if err != nil {
			return err
		}
	}

	// Close stream
	_, err = stream.CloseAndRecv()

	return err
}

func (c *GRPCClient) DownloadFile(ctx context.Context, ID uint64, fileName string) error {
	// open file
	f, err := c.openFile(c.config.DownloadPath, fileName)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}

	defer func() {
		err := f.Close()
		if err != nil {
			log.Println("failed to close file: ", err)
		}
	}()

	req := &pb.DownloadFileV1Request{
		Id:             ID,
		MasterPassword: c.masterPassword,
	}

	srv, err := c.secretsClient.DownloadFileV1(ctx, req)
	if err != nil {
		return fmt.Errorf("unable to establish connection: %w", err)
	}

	// Start download
	for {
		res, err := srv.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			return fmt.Errorf("transfer session was interrupted: %w", err)
		}

		// Write chunk to file
		_, err = f.Write(res.Chunk)
		if err != nil {
			return fmt.Errorf("error writing chunk: %w", err)
		}
	}

	return nil
}

func (c *GRPCClient) openFile(path, fileName string) (*os.File, error) {
	var f *os.File

	// Create download dir if not exists
	_, err := os.Stat(path)

	if os.IsNotExist(err) {
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return nil, fmt.Errorf("failed to create download dir: %w", err)
		}
	}

	if err != nil {
		return nil, err
	}

	// Open file
	filePath := fmt.Sprintf("%s/%s", path, fileName)

	f, err = os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	return f, nil
}
