package grpc

import (
	"context"
	"errors"
	"strconv"
	"testing"
	"time"

	"github.com/Sadere/gophkeeper/internal/server/model"
	service "github.com/Sadere/gophkeeper/internal/server/service/mocks"
	"github.com/Sadere/gophkeeper/pkg/constants"
	"github.com/Sadere/gophkeeper/pkg/convert"
	pkgModel "github.com/Sadere/gophkeeper/pkg/model"
	pb "github.com/Sadere/gophkeeper/pkg/proto/keeper/v1"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestSecretPreviews(t *testing.T) {
	userID := uint64(111)

	tests := []struct {
		name    string
		ctx     context.Context
		prepare func(s *service.MockISecretService)
		wantErr bool
	}{
		{
			name: "success",
			ctx:  context.WithValue(context.Background(), constants.CtxUserIDKey, userID),
			prepare: func(s *service.MockISecretService) {
				s.EXPECT().GetUserSecrets(gomock.Any(), userID).Return(pkgModel.Secrets{&pkgModel.Secret{}}, nil)
			},
			wantErr: false,
		},
		{
			name: "empty list",
			ctx:  context.WithValue(context.Background(), constants.CtxUserIDKey, userID),
			prepare: func(s *service.MockISecretService) {
				s.EXPECT().GetUserSecrets(gomock.Any(), userID).Return(nil, model.ErrNoSecrets)
			},
			wantErr: true,
		},
		{
			name: "error",
			ctx:  context.WithValue(context.Background(), constants.CtxUserIDKey, userID),
			prepare: func(s *service.MockISecretService) {
				s.EXPECT().GetUserSecrets(gomock.Any(), userID).Return(nil, errors.New("error"))
			},
			wantErr: true,
		},
		{
			name:    "no user id",
			ctx:     context.Background(),
			prepare: func(s *service.MockISecretService) {},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, _, secretMock, userCtrl, secretCtrl := NewTestServer(t)

			defer func() {
				userCtrl.Finish()
				secretCtrl.Finish()
			}()

			tt.prepare(secretMock)

			_, err := server.SecretPreviewsV1(tt.ctx, &emptypb.Empty{})

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSaveUserSecret(t *testing.T) {
	userID := uint64(111)
	clientID := uint64(222)
	password := "password"

	ctxUser := context.WithValue(context.Background(), constants.CtxUserIDKey, userID)

	md := metadata.New(map[string]string{
		constants.ClientIDHeader: strconv.Itoa(int(clientID)),
	})

	mdCtx := metadata.NewIncomingContext(ctxUser, md)

	timeNow := time.Now()
	pbSecret := &pb.Secret{
		Id:        222,
		CreatedAt: timestamppb.New(timeNow),
		UpdatedAt: timestamppb.New(timeNow),
		Type:      pb.SecretType_SECRET_TYPE_CREDENTIAL,
		Content: &pb.Secret_Credential{
			Credential: &pb.Credential{
				Login:    "login",
				Password: "password",
			},
		},
	}
	secret := convert.ProtoToSecret(pbSecret)
	secret.UserID = userID

	tests := []struct {
		name    string
		ctx     context.Context
		req     *pb.SaveUserSecretRequestV1
		prepare func(s *service.MockISecretService)
		wantErr bool
	}{
		{
			name: "success",
			ctx:  mdCtx,
			req: &pb.SaveUserSecretRequestV1{
				MasterPassword: password,
				Secret:         pbSecret,
			},
			prepare: func(s *service.MockISecretService) {
				s.EXPECT().SaveSecret(mdCtx, password, secret).Return(uint64(333), nil)
			},
			wantErr: false,
		},
		{
			name: "wrong secret type",
			ctx:  mdCtx,
			req: &pb.SaveUserSecretRequestV1{
				MasterPassword: password,
				Secret:         pbSecret,
			},
			prepare: func(s *service.MockISecretService) {
				s.EXPECT().SaveSecret(mdCtx, password, secret).Return(uint64(0), model.ErrWrongSecretType)
			},
			wantErr: true,
		},
		{
			name: "error",
			ctx:  mdCtx,
			req: &pb.SaveUserSecretRequestV1{
				MasterPassword: password,
				Secret:         pbSecret,
			},
			prepare: func(s *service.MockISecretService) {
				s.EXPECT().SaveSecret(mdCtx, password, secret).Return(uint64(0), errors.New("error"))
			},
			wantErr: true,
		},
		{
			name: "invalid request",
			ctx:  mdCtx,
			req: &pb.SaveUserSecretRequestV1{
				MasterPassword: password,
				Secret: &pb.Secret{
					Type: pb.SecretType_SECRET_TYPE_CREDENTIAL,
					Content: &pb.Secret_Credential{
						Credential: &pb.Credential{
							Login:    "",
							Password: "",
						},
					},
				},
			},
			prepare: func(s *service.MockISecretService) {},
			wantErr: true,
		},
		{
			name:    "no user id",
			ctx:     context.Background(),
			prepare: func(s *service.MockISecretService) {},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, _, secretMock, userCtrl, secretCtrl := NewTestServer(t)

			defer func() {
				userCtrl.Finish()
				secretCtrl.Finish()
			}()

			tt.prepare(secretMock)

			_, err := server.SaveUserSecretV1(tt.ctx, tt.req)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetUserSecret(t *testing.T) {
	userID := uint64(111)
	secretID := uint64(222)
	password := "password"

	ctxUser := context.WithValue(context.Background(), constants.CtxUserIDKey, userID)

	tests := []struct {
		name    string
		ctx     context.Context
		req     *pb.GetUserSecretRequestV1
		prepare func(s *service.MockISecretService)
		wantErr bool
	}{
		{
			name: "success",
			ctx:  ctxUser,
			req: &pb.GetUserSecretRequestV1{
				MasterPassword: password,
				Id:             secretID,
			},
			prepare: func(s *service.MockISecretService) {
				s.EXPECT().GetSecret(ctxUser, password, secretID, userID).Return(&pkgModel.Secret{}, nil)
			},
			wantErr: false,
		},
		{
			name: "secret not found",
			ctx:  ctxUser,
			req: &pb.GetUserSecretRequestV1{
				MasterPassword: password,
				Id:             secretID,
			},
			prepare: func(s *service.MockISecretService) {
				s.EXPECT().GetSecret(ctxUser, password, secretID, userID).Return(nil, model.ErrSecretNotFound)
			},
			wantErr: true,
		},
		{
			name: "error",
			ctx:  ctxUser,
			req: &pb.GetUserSecretRequestV1{
				MasterPassword: password,
				Id:             secretID,
			},
			prepare: func(s *service.MockISecretService) {
				s.EXPECT().GetSecret(ctxUser, password, secretID, userID).Return(nil, errors.New("error"))
			},
			wantErr: true,
		},
		{
			name: "invalid request",
			ctx:  ctxUser,
			req: &pb.GetUserSecretRequestV1{
				MasterPassword: "",
			},
			prepare: func(s *service.MockISecretService) {},
			wantErr: true,
		},
		{
			name:    "no user id",
			ctx:     context.Background(),
			prepare: func(s *service.MockISecretService) {},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, _, secretMock, userCtrl, secretCtrl := NewTestServer(t)

			defer func() {
				userCtrl.Finish()
				secretCtrl.Finish()
			}()

			tt.prepare(secretMock)

			_, err := server.GetUserSecretV1(tt.ctx, tt.req)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
