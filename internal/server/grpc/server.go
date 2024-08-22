package grpc

import (
	"github.com/Sadere/gophkeeper/internal/server/config"
	"github.com/Sadere/gophkeeper/internal/server/grpc/interceptor"
	"github.com/Sadere/gophkeeper/internal/server/service"
	"go.uber.org/zap"

	"google.golang.org/grpc"

	pb "github.com/Sadere/gophkeeper/pkg/proto/keeper/v1"
)

type KeeperServer struct {
	pb.UnimplementedAuthServiceServer
	pb.UnimplementedSecretsServiceServer

	config        *config.Config
	userService   service.IUserService
	secretService service.ISecretService
	log           *zap.SugaredLogger
}

func NewKeeperServer(cfg *config.Config, log *zap.SugaredLogger, userService service.IUserService, secretService service.ISecretService) *KeeperServer {
	return &KeeperServer{
		config: cfg,
		log:    log,

		userService:   userService,
		secretService: secretService,
	}
}

func (s *KeeperServer) Register() (*grpc.Server, error) {
	srvInterceptors := make([]grpc.UnaryServerInterceptor, 0)

	// Log requests
	srvInterceptors = append(srvInterceptors, interceptor.Logger(s.log))

	// Authentication
	srvInterceptors = append(srvInterceptors, interceptor.Authentication([]byte(s.config.SecretKey)))

	srv := grpc.NewServer(grpc.ChainUnaryInterceptor(
		srvInterceptors...,
	))

	pb.RegisterAuthServiceServer(srv, s)
	pb.RegisterSecretsServiceServer(srv, s)

	return srv, nil
}
