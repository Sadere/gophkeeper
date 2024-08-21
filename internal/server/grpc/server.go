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

	config      *config.Config
	userService service.IUserService
	log         *zap.SugaredLogger
}

func NewKeeperServer(cfg *config.Config, userService service.IUserService, log *zap.SugaredLogger) *KeeperServer {
	return &KeeperServer{
		config:      cfg,
		userService: userService,
		log:         log,
	}
}

func (s *KeeperServer) Register() (*grpc.Server, error) {
	srvInterceptors := make([]grpc.UnaryServerInterceptor, 0)

	// Log requests
	srvInterceptors = append(srvInterceptors, interceptor.Logger(s.log))

	srv := grpc.NewServer(grpc.ChainUnaryInterceptor(
		srvInterceptors...,
	))

	pb.RegisterAuthServiceServer(srv, s)

	return srv, nil
}
