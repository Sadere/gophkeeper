package service

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/Sadere/gophkeeper/internal/server/crypto"
	repository "github.com/Sadere/gophkeeper/internal/server/repository/mocks"
	"github.com/Sadere/gophkeeper/pkg/model"
)

func TestGetUserSecrets(t *testing.T) {
	userID := uint64(111)

	tests := []struct {
		name    string
		prepare func(m *repository.MockSecretRepository)
		wantErr bool
	}{
		{
			name: "success",
			prepare: func(m *repository.MockSecretRepository) {
				m.EXPECT().GetUserSecrets(gomock.Any(), userID).Return(model.Secrets{
					&model.Secret{},
				}, nil)
			},
			wantErr: false,
		},
		{
			name: "empty list",
			prepare: func(m *repository.MockSecretRepository) {
				m.EXPECT().GetUserSecrets(gomock.Any(), userID).Return(model.Secrets{}, nil)
			},
			wantErr: true,
		},
		{
			name: "error",
			prepare: func(m *repository.MockSecretRepository) {
				m.EXPECT().GetUserSecrets(gomock.Any(), userID).Return(nil, errors.New("error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := repository.NewMockSecretRepository(ctrl)
			s := NewSecretService(m)

			tt.prepare(m)

			_, err := s.GetUserSecrets(context.Background(), userID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSaveSecret(t *testing.T) {
	tests := []struct {
		name    string
		secret  *model.Secret
		prepare func(m *repository.MockSecretRepository)
		wantErr bool
	}{
		{
			name: "create creds",
			secret: &model.Secret{
				SType: string(model.CredSecret),
				Creds: &model.Credentials{},
			},
			prepare: func(m *repository.MockSecretRepository) {
				m.EXPECT().Create(gomock.Any(), gomock.Any()).Return(uint64(222), nil)
			},
			wantErr: false,
		},
		{
			name: "create creds error",
			secret: &model.Secret{
				SType: string(model.CredSecret),
				Creds: &model.Credentials{},
			},
			prepare: func(m *repository.MockSecretRepository) {
				m.EXPECT().Create(gomock.Any(), gomock.Any()).Return(uint64(0), errors.New("error"))
			},
			wantErr: true,
		},
		{
			name: "update creds",
			secret: &model.Secret{
				ID:     444,
				UserID: 333,
				SType:  string(model.CredSecret),
				Creds:  &model.Credentials{},
			},
			prepare: func(m *repository.MockSecretRepository) {
				m.EXPECT().GetSecret(gomock.Any(), uint64(444), uint64(333)).Return(nil, nil)
				m.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "update creds failed",
			secret: &model.Secret{
				ID:     444,
				UserID: 333,
				SType:  string(model.CredSecret),
				Creds:  &model.Credentials{},
			},
			prepare: func(m *repository.MockSecretRepository) {
				m.EXPECT().GetSecret(gomock.Any(), uint64(444), uint64(333)).Return(nil, errors.New("error"))
			},
			wantErr: true,
		},
		{
			name: "create text",
			secret: &model.Secret{
				SType: string(model.TextSecret),
				Text:  &model.Text{},
			},
			prepare: func(m *repository.MockSecretRepository) {
				m.EXPECT().Create(gomock.Any(), gomock.Any()).Return(uint64(222), nil)
			},
			wantErr: false,
		},
		{
			name: "create card",
			secret: &model.Secret{
				SType: string(model.CardSecret),
				Card: &model.Card{
					Number: "374245455400126",
				},
			},
			prepare: func(m *repository.MockSecretRepository) {
				m.EXPECT().Create(gomock.Any(), gomock.Any()).Return(uint64(222), nil)
			},
			wantErr: false,
		},
		{
			name: "invalid card number",
			secret: &model.Secret{
				SType: string(model.CardSecret),
				Card: &model.Card{
					Number: "10019",
				},
			},
			prepare: func(m *repository.MockSecretRepository) {},
			wantErr: true,
		},
		{
			name: "create blob",
			secret: &model.Secret{
				SType: string(model.BlobSecret),
				Blob:  &model.Blob{},
			},
			prepare: func(m *repository.MockSecretRepository) {
				m.EXPECT().Create(gomock.Any(), gomock.Any()).Return(uint64(222), nil)
			},
			wantErr: false,
		},
		{
			name: "wrong secret type",
			secret: &model.Secret{
				SType: "invalid",
			},
			prepare: func(m *repository.MockSecretRepository) {},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := repository.NewMockSecretRepository(ctrl)
			s := NewSecretService(m)

			tt.prepare(m)

			_, err := s.SaveSecret(context.Background(), "password", tt.secret)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetSecret(t *testing.T) {
	ID := uint64(111)
	userID := uint64(222)
	password := "password"

	tests := []struct {
		name      string
		stype     string
		payload   string
		errReturn error
		wantErr   bool
	}{
		{
			name:      "get creds",
			stype:     string(model.CredSecret),
			payload:   `{"login":"login","password":"pw"}`,
			errReturn: nil,
			wantErr:   false,
		},
		{
			name:      "get creds error",
			stype:     string(model.CredSecret),
			payload:   `invalid`,
			errReturn: nil,
			wantErr:   true,
		},
		{
			name:      "get text",
			stype:     string(model.TextSecret),
			payload:   `{"content":"text content"}`,
			errReturn: nil,
			wantErr:   false,
		},
		{
			name:      "get text error",
			stype:     string(model.TextSecret),
			payload:   `invalid`,
			errReturn: nil,
			wantErr:   true,
		},
		{
			name:      "get card",
			stype:     string(model.CardSecret),
			payload:   `{"number":"111","exp_year":10,"exp_month":12,"cvv":444}`,
			errReturn: nil,
			wantErr:   false,
		},
		{
			name:      "get card error",
			stype:     string(model.CardSecret),
			payload:   `invalid`,
			errReturn: nil,
			wantErr:   true,
		},
		{
			name:      "get blob",
			stype:     string(model.BlobSecret),
			payload:   `{"file_name":"file.txt"}`,
			errReturn: nil,
			wantErr:   false,
		},
		{
			name:      "get blob error",
			stype:     string(model.BlobSecret),
			payload:   `invalid`,
			errReturn: nil,
			wantErr:   true,
		},
		{
			name:      "not found",
			errReturn: sql.ErrNoRows,
			wantErr:   true,
		},
		{
			name:      "random error",
			errReturn: errors.New("error"),
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := repository.NewMockSecretRepository(ctrl)
			s := NewSecretService(m)

			// Setup secret
			secret := &model.Secret{
				SType: tt.stype,
			}
			payload, err := crypto.Encrypt(password, []byte(tt.payload))

			assert.NoError(t, err)

			secret.Payload = payload

			// Setup mock
			ctx := context.Background()
			m.EXPECT().GetSecret(ctx, ID, userID).Return(secret, tt.errReturn)

			// Test function
			_, errGet := s.GetSecret(ctx, password, ID, userID)

			// Assert result
			if tt.wantErr {
				assert.Error(t, errGet)
			} else {
				assert.NoError(t, errGet)
			}
		})
	}
}
