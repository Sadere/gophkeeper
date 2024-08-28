package grpc

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"

	"github.com/Sadere/gophkeeper/cert"
	"github.com/Sadere/gophkeeper/internal/server/config"
	"github.com/Sadere/gophkeeper/internal/server/grpc/interceptor"
	"github.com/Sadere/gophkeeper/internal/server/service"
	"go.uber.org/zap"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

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
	var opts []grpc.ServerOption

	srvInterceptors := make([]grpc.UnaryServerInterceptor, 0)

	// Log requests
	srvInterceptors = append(srvInterceptors, interceptor.Logger(s.log))

	// Authentication
	srvInterceptors = append(srvInterceptors, interceptor.Authentication([]byte(s.config.SecretKey)))

	// Chain of unary interceptors
	opts = append(
		opts,
		grpc.ChainUnaryInterceptor(
			srvInterceptors...,
		),
	)

	// Stream interceptor
	opts = append(
		opts,
		grpc.StreamInterceptor(
			interceptor.StreamAuthentication([]byte(s.config.SecretKey)),
		),
	)

	// TLS config
	if s.config.EnableTLS {
		tlsCreds, err := loadTLSConfig("ca-cert.pem", "server-cert.pem", "server-key.pem")
		if err != nil {
			return nil, fmt.Errorf("failed to load TLS config: %w", err)
		}

		// Append TLS credentials to server options
		opts = append(opts, grpc.Creds(tlsCreds))
	}

	srv := grpc.NewServer(opts...)

	pb.RegisterAuthServiceServer(srv, s)
	pb.RegisterSecretsServiceServer(srv, s)

	return srv, nil
}

func loadTLSConfig(caCertFile, serverCertFile, serverKeyFile string) (credentials.TransportCredentials, error) {
	// Read CA cert
	caPem, err := cert.Cert.ReadFile(caCertFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read CA cert: %w", err)
	}

	// Read server cert
	serverCertPEM, err := cert.Cert.ReadFile(serverCertFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read server cert: %w", err)
	}

	// Read server key
	serverKeyPEM, err := cert.Cert.ReadFile(serverKeyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read server key: %w", err)
	}

	// Create key pair
	serverCert, err := tls.X509KeyPair(serverCertPEM, serverKeyPEM)
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
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	}

	return credentials.NewTLS(config), nil
}
