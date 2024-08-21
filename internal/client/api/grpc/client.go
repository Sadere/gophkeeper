package grpc

import (
	"context"
	"fmt"

	"github.com/Sadere/gophkeeper/internal/client/api"
	"github.com/Sadere/gophkeeper/internal/client/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/Sadere/gophkeeper/pkg/proto/keeper/v1"
)

type GRPCClient struct {
	config      *config.Config
	authClient  pb.AuthServiceV1Client
	accessToken string
}

var _ api.IApiClient = &GRPCClient{}

func NewGRPCClient(cfg *config.Config) (*GRPCClient, error) {
	newClient := GRPCClient{
		config: cfg,
	}

	c, err := grpc.NewClient(cfg.ServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC client: %w", err)
	}

	authClient := pb.NewAuthServiceV1Client(c)

	newClient.authClient = authClient

	return &newClient, nil
}

func (c *GRPCClient) Login(ctx context.Context, login string, password string) (string, error) {
	req := &pb.AuthRequestV1{
		Login:    login,
		Password: password,
	}

	response, err := c.authClient.LoginV1(ctx, req)
	if err != nil {
		return "", err
	}

	c.accessToken = response.AccessToken

	return response.AccessToken, nil
}

func (c *GRPCClient) Register(ctx context.Context, login string, password string) (string, error) {
	req := &pb.AuthRequestV1{
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
