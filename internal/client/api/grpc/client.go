package grpc

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/Sadere/gophkeeper/internal/client/api"
	"github.com/Sadere/gophkeeper/internal/client/api/grpc/interceptor"
	"github.com/Sadere/gophkeeper/internal/client/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	accessToken    string
	masterPassword string
	chunkSize      int
}

var _ api.IApiClient = &GRPCClient{}

func NewGRPCClient(cfg *config.Config) (*GRPCClient, error) {
	newClient := GRPCClient{
		config:    cfg,
		chunkSize: constants.ChunkSize,
	}

	// create gRPC client
	c, err := grpc.NewClient(
		cfg.ServerAddress,
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),

		// Unary interceptors
		grpc.WithChainUnaryInterceptor(
			interceptor.Timeout(constants.DefaultClientTimeout),
			interceptor.AddToken(&newClient.accessToken),
		),

		// Stream interceptor
		grpc.WithStreamInterceptor(interceptor.AddTokenStream(&newClient.accessToken)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC client: %w", err)
	}

	// register services
	authClient := pb.NewAuthServiceClient(c)
	secretsClient := pb.NewSecretsServiceClient(c)

	newClient.authClient = authClient
	newClient.secretsClient = secretsClient

	return &newClient, nil
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
		return nil, fmt.Errorf("failed to retrieve secrets: %w", err)
	}

	for _, preview := range response.Previews {
		previews = append(previews, convert.ProtoToPreview(preview))
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
