package grpc

import (
	"context"
	"errors"
	"testing"

	"github.com/Sadere/gophkeeper/internal/server/model"
	service "github.com/Sadere/gophkeeper/internal/server/service/mocks"
	"github.com/Sadere/gophkeeper/pkg/constants"
	pb "github.com/Sadere/gophkeeper/pkg/proto/keeper/v1"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
)

func TestRegister(t *testing.T) {
	tests := []struct {
		name    string
		req     *pb.RegisterRequestV1
		prepare func(s *service.MockIUserService)
		wantErr bool
	}{
		{
			name: "register success",
			prepare: func(s *service.MockIUserService) {
				s.EXPECT().RegisterUser(gomock.Any(), "login", "password").Return(&model.User{
					ID: 233,
				}, nil)
			},
			req: &pb.RegisterRequestV1{
				Login:    "login",
				Password: "password",
			},
			wantErr: false,
		},
		{
			name: "user already exists",
			prepare: func(s *service.MockIUserService) {
				s.EXPECT().RegisterUser(gomock.Any(), "login", "password").Return(nil, &model.ErrUserExists{Login: "login"})
			},
			req: &pb.RegisterRequestV1{
				Login:    "login",
				Password: "password",
			},
			wantErr: true,
		},
		{
			name: "register error",
			prepare: func(s *service.MockIUserService) {
				s.EXPECT().RegisterUser(gomock.Any(), "login", "password").Return(nil, errors.New("error"))
			},
			req: &pb.RegisterRequestV1{
				Login:    "login",
				Password: "password",
			},
			wantErr: true,
		},
		{
			name:    "invalid request",
			prepare: func(s *service.MockIUserService) {},
			req: &pb.RegisterRequestV1{
				Login:    "",
				Password: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, userMock, _, userCtrl, secretCtrl := NewTestServer(t)

			defer func() {
				userCtrl.Finish()
				secretCtrl.Finish()
			}()

			tt.prepare(userMock)

			_, err := server.RegisterV1(context.Background(), tt.req)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	tests := []struct {
		name    string
		req     *pb.LoginRequestV1
		prepare func(s *service.MockIUserService)
		wantErr bool
	}{
		{
			name: "login success",
			prepare: func(s *service.MockIUserService) {
				s.EXPECT().LoginUser(gomock.Any(), "login", "password").Return(&model.User{
					ID: 233,
				}, nil)
			},
			req: &pb.LoginRequestV1{
				Login:    "login",
				Password: "password",
			},
			wantErr: false,
		},
		{
			name: "bad credentials",
			prepare: func(s *service.MockIUserService) {
				s.EXPECT().LoginUser(gomock.Any(), "login", "password").Return(nil, model.ErrBadCredentials)
			},
			req: &pb.LoginRequestV1{
				Login:    "login",
				Password: "password",
			},
			wantErr: true,
		},
		{
			name: "login error",
			prepare: func(s *service.MockIUserService) {
				s.EXPECT().LoginUser(gomock.Any(), "login", "password").Return(nil, errors.New("error"))
			},
			req: &pb.LoginRequestV1{
				Login:    "login",
				Password: "password",
			},
			wantErr: true,
		},
		{
			name:    "invalid request",
			prepare: func(s *service.MockIUserService) {},
			req: &pb.LoginRequestV1{
				Login:    "",
				Password: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, userMock, _, userCtrl, secretCtrl := NewTestServer(t)

			defer func() {
				userCtrl.Finish()
				secretCtrl.Finish()
			}()

			tt.prepare(userMock)

			_, err := server.LoginV1(context.Background(), tt.req)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestExtractClientID(t *testing.T) {
	t.Run("failed to get metadata", func(t *testing.T) {
		_, err := extractClientID(context.Background())

		assert.EqualError(t, err, "failed to get metadata")
	})

	t.Run("missing client id", func(t *testing.T) {
		md := metadata.New(nil)

		mdCtx := metadata.NewIncomingContext(context.Background(), md)
		_, err := extractClientID(mdCtx)

		assert.EqualError(t, err, "missing client id metadata")
	})

	t.Run("non number id", func(t *testing.T) {
		md := metadata.New(map[string]string{
			constants.ClientIDHeader: "invalid",
		})

		mdCtx := metadata.NewIncomingContext(context.Background(), md)
		_, err := extractClientID(mdCtx)

		assert.Error(t, err)
	})
}
