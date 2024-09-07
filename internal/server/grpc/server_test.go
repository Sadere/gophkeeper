package grpc

import (
	"testing"

	"github.com/Sadere/gophkeeper/internal/server/config"
	service "github.com/Sadere/gophkeeper/internal/server/service/mocks"
	"github.com/Sadere/gophkeeper/internal/server/utils"
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
