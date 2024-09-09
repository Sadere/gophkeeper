package grpc

import (
	"context"
	"errors"
	"testing"

	"github.com/Sadere/gophkeeper/internal/client/config"
	"github.com/Sadere/gophkeeper/pkg/model"
	mock "github.com/Sadere/gophkeeper/pkg/proto/keeper/v1/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/Sadere/gophkeeper/pkg/proto/keeper/v1"
)

type TestClient struct {
	Client *GRPCClient

	AuthMock    *mock.MockAuthServiceClient
	SecretsMock *mock.MockSecretsServiceClient
	NotifMock   *mock.MockNotificationServiceClient

	AuthCtrl    *gomock.Controller
	SecretsCtrl *gomock.Controller
	NotifCtrl   *gomock.Controller
}

func (tc *TestClient) Finish() {
	tc.AuthCtrl.Finish()
	tc.SecretsCtrl.Finish()
	tc.NotifCtrl.Finish()
}

func NewTestClient(t *testing.T) *TestClient {
	authCtrl := gomock.NewController(t)
	authMock := mock.NewMockAuthServiceClient(authCtrl)

	secretsCtrl := gomock.NewController(t)
	secretsMock := mock.NewMockSecretsServiceClient(secretsCtrl)

	notifCtrl := gomock.NewController(t)
	notifMock := mock.NewMockNotificationServiceClient(notifCtrl)

	client := &GRPCClient{
		config:        &config.Config{},
		authClient:    authMock,
		secretsClient: secretsMock,
		notifyClient:  notifMock,
	}

	return &TestClient{
		Client: client,

		AuthMock:    authMock,
		SecretsMock: secretsMock,
		NotifMock:   notifMock,

		AuthCtrl:    authCtrl,
		SecretsCtrl: secretsCtrl,
		NotifCtrl:   notifCtrl,
	}
}

func TestLoadTLSConfig(t *testing.T) {
	t.Run("failed to read ca cert", func(t *testing.T) {
		_, err := loadTLSConfig("invalid", "", "")

		assert.EqualError(t, err, "failed to read CA cert: open invalid: file does not exist")
	})

	t.Run("failed to read client cert", func(t *testing.T) {
		_, err := loadTLSConfig("ca-cert.pem", "invalid", "")

		assert.EqualError(t, err, "failed to read client cert: open invalid: file does not exist")
	})

	t.Run("failed to read client key", func(t *testing.T) {
		_, err := loadTLSConfig("ca-cert.pem", "client-cert.pem", "invalid")

		assert.EqualError(t, err, "failed to read client key: open invalid: file does not exist")
	})
}


func TestNewClient(t *testing.T) {
	t.Run("client no tls", func(t *testing.T) {
		cfg := &config.Config{}
		client, err := NewGRPCClient(cfg)

		require.NoError(t, err)
		assert.GreaterOrEqual(t, client.clientID, int32(0))
	})

	t.Run("client with tls", func(t *testing.T) {
		cfg := &config.Config{
			EnableTLS: true,
		}
		_, err := NewGRPCClient(cfg)

		require.NoError(t, err)
	})
}

func TestLogin(t *testing.T) {
	login := "login"
	password := "password"

	tests := []struct {
		name    string
		prepare func(*mock.MockAuthServiceClient)
		wantErr bool
	}{
		{
			name: "login success",
			prepare: func(m *mock.MockAuthServiceClient) {
				m.EXPECT().LoginV1(
					gomock.Any(),
					&pb.LoginRequestV1{
						Login:    login,
						Password: password,
					},
				).Return(&pb.LoginResponseV1{
					AccessToken: "access_token",
				}, nil)
			},
			wantErr: false,
		},
		{
			name: "login failed",
			prepare: func(m *mock.MockAuthServiceClient) {
				m.EXPECT().LoginV1(
					gomock.Any(),
					&pb.LoginRequestV1{
						Login:    login,
						Password: password,
					},
				).Return(nil, errors.New("error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := NewTestClient(t)
			defer tc.Finish()

			tt.prepare(tc.AuthMock)

			_, err := tc.Client.Login(context.Background(), login, password)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRegister(t *testing.T) {
	login := "login"
	password := "password"

	tests := []struct {
		name    string
		prepare func(*mock.MockAuthServiceClient)
		wantErr bool
	}{
		{
			name: "register success",
			prepare: func(m *mock.MockAuthServiceClient) {
				m.EXPECT().RegisterV1(
					gomock.Any(),
					&pb.RegisterRequestV1{
						Login:    login,
						Password: password,
					},
				).Return(&pb.RegisterResponseV1{
					AccessToken: "access_token",
				}, nil)
			},
			wantErr: false,
		},
		{
			name: "register failed",
			prepare: func(m *mock.MockAuthServiceClient) {
				m.EXPECT().RegisterV1(
					gomock.Any(),
					&pb.RegisterRequestV1{
						Login:    login,
						Password: password,
					},
				).Return(nil, errors.New("error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := NewTestClient(t)
			defer tc.Finish()

			tt.prepare(tc.AuthMock)

			_, err := tc.Client.Register(context.Background(), login, password)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestLoadPreviews(t *testing.T) {
	type want struct {
		len int
		err string
	}
	tests := []struct {
		name    string
		prepare func(tc *TestClient)
		want    want
	}{
		{
			name: "load single preview",
			prepare: func(tc *TestClient) {
				tc.SecretsMock.EXPECT().SecretPreviewsV1(gomock.Any(), &emptypb.Empty{}).
					Return(
						&pb.SecretPreviewsResponseV1{
							Previews: []*pb.SecretPreview{
								{
									Id:       111,
									Metadata: "creds",
									Type:     pb.SecretType_SECRET_TYPE_CREDENTIAL,
								},
							},
						}, nil)
			},
			want: want{
				len: 1,
			},
		},
		{
			name: "load two previews",
			prepare: func(tc *TestClient) {
				tc.SecretsMock.EXPECT().SecretPreviewsV1(gomock.Any(), &emptypb.Empty{}).
					Return(
						&pb.SecretPreviewsResponseV1{
							Previews: []*pb.SecretPreview{
								{
									Id:       111,
									Metadata: "creds",
									Type:     pb.SecretType_SECRET_TYPE_CREDENTIAL,
								},
								{
									Id:       222,
									Metadata: "creds",
									Type:     pb.SecretType_SECRET_TYPE_CREDENTIAL,
								},
							},
						}, nil)
			},
			want: want{
				len: 2,
			},
		},
		{
			name: "load empty",
			prepare: func(tc *TestClient) {
				tc.SecretsMock.EXPECT().SecretPreviewsV1(gomock.Any(), &emptypb.Empty{}).
					Return(nil, status.Error(codes.NotFound, "not found"))
			},
			want: want{
				len: 0,
			},
		},
		{
			name: "load error",
			prepare: func(tc *TestClient) {
				tc.SecretsMock.EXPECT().SecretPreviewsV1(gomock.Any(), &emptypb.Empty{}).
					Return(nil, status.Error(codes.Internal, "error"))
			},
			want: want{
				len: 0,
				err: "failed to retrieve secrets: rpc error: code = Internal desc = error",
			},
		},
		{
			name: "set status",
			prepare: func(tc *TestClient) {
				secretID := uint64(111)

				// Set status in map
				tc.Client.previews.Store(secretID, model.SecretPreviewUpdated)

				tc.SecretsMock.EXPECT().SecretPreviewsV1(gomock.Any(), &emptypb.Empty{}).
					Return(
						&pb.SecretPreviewsResponseV1{
							Previews: []*pb.SecretPreview{
								{
									Id:       secretID,
									Metadata: "creds",
									Type:     pb.SecretType_SECRET_TYPE_CREDENTIAL,
								},
							},
						}, nil)
			},
			want: want{
				len: 1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := NewTestClient(t)
			defer tc.Finish()

			tt.prepare(tc)

			response, err := tc.Client.LoadPreviews(context.Background())

			assert.Len(t, response, tt.want.len)

			if len(tt.want.err) > 0 {
				assert.EqualError(t, err, tt.want.err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestLoadSecret(t *testing.T) {
	secretID := uint64(111)
	password := "password"

	tests := []struct {
		name    string
		prepare func(*mock.MockSecretsServiceClient)
		wantErr bool
	}{
		{
			name: "load success",
			prepare: func(m *mock.MockSecretsServiceClient) {
				m.EXPECT().GetUserSecretV1(
					gomock.Any(),
					&pb.GetUserSecretRequestV1{
						MasterPassword: password,
						Id:             secretID,
					},
				).Return(&pb.GetUserSecretResponseV1{
					Secret: &pb.Secret{
						Id:   secretID,
						Type: pb.SecretType_SECRET_TYPE_TEXT,
						Content: &pb.Secret_Text{
							Text: &pb.Text{
								Text: "text content",
							},
						},
					},
				}, nil)
			},
			wantErr: false,
		},
		{
			name: "load error",
			prepare: func(m *mock.MockSecretsServiceClient) {
				m.EXPECT().GetUserSecretV1(
					gomock.Any(),
					&pb.GetUserSecretRequestV1{
						MasterPassword: password,
						Id:             secretID,
					},
				).Return(nil, errors.New("error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := NewTestClient(t)
			defer tc.Finish()

			// Set master pw
			tc.Client.masterPassword = password

			tt.prepare(tc.SecretsMock)

			_, err := tc.Client.LoadSecret(context.Background(), secretID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSaveCredential(t *testing.T) {
	tc := NewTestClient(t)
	defer tc.Finish()

	ctx := context.Background()

	// Setup mock
	tc.SecretsMock.EXPECT().SaveUserSecretV1(ctx, gomock.Any()).Return(nil, nil)

	// Test function
	err := tc.Client.SaveCredential(ctx, uint64(111), "metadata", "login", "password")

	// Assert
	assert.NoError(t, err)
}

func TestSaveText(t *testing.T) {
	tc := NewTestClient(t)
	defer tc.Finish()

	ctx := context.Background()

	// Setup mock
	tc.SecretsMock.EXPECT().SaveUserSecretV1(ctx, gomock.Any()).Return(nil, nil)

	// Test function
	err := tc.Client.SaveText(ctx, uint64(111), "metadata", "text content")

	// Assert
	assert.NoError(t, err)
}

func TestSaveCard(t *testing.T) {
	tc := NewTestClient(t)
	defer tc.Finish()

	ctx := context.Background()

	// Setup mock
	tc.SecretsMock.EXPECT().SaveUserSecretV1(ctx, gomock.Any()).Return(nil, nil)

	// Test function
	err := tc.Client.SaveCard(ctx, uint64(111), "metadata", "1111111", 11, 24, 555)

	// Assert
	assert.NoError(t, err)
}
