package grpc

import (
	"context"
	"fmt"

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
}

var _ api.IApiClient = &GRPCClient{}

func NewGRPCClient(cfg *config.Config) (*GRPCClient, error) {
	newClient := GRPCClient{
		config: cfg,
	}

	// create gRPC client
	c, err := grpc.NewClient(
		cfg.ServerAddress,
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),

		// interceptors...
		grpc.WithChainUnaryInterceptor(
			interceptor.Timeout(constants.DefaultClientTimeout),
			interceptor.AddToken(&newClient.accessToken),
		),
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

func (c *GRPCClient) SaveCredential(ctx context.Context, metadata, login, password string) error {
	// form gRPC request
	request := &pb.SaveUserSecretV1Request{
		MasterPassword: password,
		Secret: &pb.Secret{
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
