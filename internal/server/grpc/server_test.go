package grpc

import (
	"testing"

	"github.com/Sadere/gophkeeper/internal/server/config"
	service "github.com/Sadere/gophkeeper/internal/server/service/mocks"
	"github.com/Sadere/gophkeeper/internal/server/utils"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func NewTestServer(t *testing.T) (*KeeperServer, *service.MockIUserService, *service.MockISecretService, *gomock.Controller, *gomock.Controller) {
	userCtrl := gomock.NewController(t)
	secretCtrl := gomock.NewController(t)

	userMock := service.NewMockIUserService(userCtrl)
	secretMock := service.NewMockISecretService(secretCtrl)

	cfg := &config.Config{
		SecretKey: "test_key",
	}

	log, err := utils.NewZapLogger("fatal")
	if err != nil {
		log.Fatal("failed to create logger ", err)
	}

	return NewKeeperServer(
			cfg,
			log,
			userMock,
			secretMock,
		),
		userMock,
		secretMock,
		userCtrl,
		secretCtrl
}

func TestServer(t *testing.T) {
	server, _, _, userCtrl, secretCtrl := NewTestServer(t)
	defer func() {
		userCtrl.Finish()
		secretCtrl.Finish()
	}()

	t.Run("register no tls", func(t *testing.T) {
		_, err := server.Register()

		assert.NoError(t, err)
	})

	t.Run("register tls", func(t *testing.T) {
		server.config.EnableTLS = true
		_, err := server.Register()

		assert.NoError(t, err)
	})
}

func TestLoadTLSConfig(t *testing.T) {
	t.Run("failed to read ca cert", func(t *testing.T) {
		_, err := loadTLSConfig("invalid", "", "")

		assert.EqualError(t, err, "failed to read CA cert: open invalid: file does not exist")
	})

	t.Run("failed to read server cert", func(t *testing.T) {
		_, err := loadTLSConfig("ca-cert.pem", "invalid", "")

		assert.EqualError(t, err, "failed to read server cert: open invalid: file does not exist")
	})

	t.Run("failed to read server key", func(t *testing.T) {
		_, err := loadTLSConfig("ca-cert.pem", "server-cert.pem", "invalid")

		assert.EqualError(t, err, "failed to read server key: open invalid: file does not exist")
	})
}
