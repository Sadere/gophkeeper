// Main package of server
package server

import (
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"

	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/Sadere/gophkeeper/internal/database"
	"github.com/Sadere/gophkeeper/internal/server/config"
	"github.com/Sadere/gophkeeper/internal/server/grpc"
	"github.com/Sadere/gophkeeper/internal/server/repository"
	"github.com/Sadere/gophkeeper/internal/server/service"
)

// Main server app which creates and registers grpc server according to provided config
type KeeperApp struct {
	config *config.Config
	log    *zap.SugaredLogger
}

// Returns instance of server app
func NewApp(cfg *config.Config, log *zap.SugaredLogger) *KeeperApp {
	return &KeeperApp{
		config: cfg,
		log:    log,
	}
}

// Start server
func (a *KeeperApp) Start() error {
	// Run migrations
	if err := database.MigrateUp(a.config.PostgresDSN); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	// Init DB
	db, err := database.NewConnection("pgx", a.config.PostgresDSN)
	if err != nil {
		return err
	}

	// Listen on address
	listen, err := net.Listen("tcp", a.config.Address)
	if err != nil {
		return err
	}

	// Set up services and repositories
	userRepo := repository.NewPgUserRepository(db)
	userService := service.NewUserService(userRepo)

	secretRepo := repository.NewPgSecretRepository(db)
	secretService := service.NewSecretService(secretRepo)

	// Create new gRPC server instance
	server := grpc.NewKeeperServer(a.config, a.log, userService, secretService)

	srv, err := server.Register()
	if err != nil {
		return err
	}

	// Run server in background
	go func() {
		if err := srv.Serve(listen); err != nil {
			a.log.Fatalf("gRPC serve error: %s\n", err)
		}
	}()

	// Catch quit signals and performing graceful shutdown
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit

	a.log.Infoln("gRPC server shutdown ...")

	srv.GracefulStop()

	return nil
}
