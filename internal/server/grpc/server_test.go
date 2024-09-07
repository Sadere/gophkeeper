package grpc

import (
	"testing"

	"github.com/Sadere/gophkeeper/internal/server/config"
	service "github.com/Sadere/gophkeeper/internal/server/service/mocks"
	"github.com/golang/mock/gomock"
)

func NewTestServer(t *testing.T) (*KeeperServer, *service.MockIUserService, *service.MockISecretService, *gomock.Controller, *gomock.Controller) {
	userCtrl := gomock.NewController(t)
	secretCtrl := gomock.NewController(t)

	userMock := service.NewMockIUserService(userCtrl)
	secretMock := service.NewMockISecretService(secretCtrl)

	cfg := &config.Config{
		SecretKey: "test_key",
	}

	return NewKeeperServer(
			cfg,
			nil,
			userMock,
			secretMock,
		),
		userMock,
		secretMock,
		userCtrl,
		secretCtrl
}
