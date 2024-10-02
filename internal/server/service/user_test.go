package service

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/Sadere/gophkeeper/internal/server/auth"
	"github.com/Sadere/gophkeeper/internal/server/model"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	repository "github.com/Sadere/gophkeeper/internal/server/repository/mocks"
)

func TestRegister(t *testing.T) {
	login := "test_login"
	password := "test_pw"

	tests := []struct {
		name    string
		prepare func(m *repository.MockUserRepository)
		wantErr bool
	}{
		{
			name: "register success",
			prepare: func(m *repository.MockUserRepository) {
				m.EXPECT().GetUserByLogin(gomock.Any(), login).Return(nil, sql.ErrNoRows)

				m.EXPECT().Create(gomock.Any(), gomock.Any()).Return(uint64(333), nil)
			},
			wantErr: false,
		},
		{
			name: "register ger user error",
			prepare: func(m *repository.MockUserRepository) {
				m.EXPECT().GetUserByLogin(gomock.Any(), login).Return(nil, errors.New("internal"))
			},
			wantErr: true,
		},
		{
			name: "register user exists",
			prepare: func(m *repository.MockUserRepository) {
				m.EXPECT().GetUserByLogin(gomock.Any(), login).Return(&model.User{}, nil)
			},
			wantErr: true,
		},
		{
			name: "register error create",
			prepare: func(m *repository.MockUserRepository) {
				m.EXPECT().GetUserByLogin(gomock.Any(), login).Return(nil, sql.ErrNoRows)

				m.EXPECT().Create(gomock.Any(), gomock.Any()).Return(uint64(0), errors.New("internal"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := repository.NewMockUserRepository(ctrl)
			s := NewUserService(m)

			tt.prepare(m)

			_, err := s.RegisterUser(context.Background(), login, password)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	login := "test_login"
	password := "test_pw"

	hash, err := auth.HashPassword(password)

	assert.NoError(t, err)

	user := &model.User{
		ID:           99,
		CreatedAt:    time.Now(),
		Login:        login,
		PasswordHash: hash,
	}

	tests := []struct {
		name     string
		login    string
		password string
		prepare  func(m *repository.MockUserRepository)
		wantErr  bool
	}{
		{
			name:     "success",
			login:    login,
			password: password,
			prepare: func(m *repository.MockUserRepository) {
				m.EXPECT().GetUserByLogin(gomock.Any(), login).Return(user, nil)

			},
			wantErr: false,
		},
		{
			name:     "user not found",
			login:    login,
			password: password,
			prepare: func(m *repository.MockUserRepository) {
				m.EXPECT().GetUserByLogin(gomock.Any(), login).Return(nil, sql.ErrNoRows)

			},
			wantErr: true,
		},
		{
			name:     "get user error",
			login:    login,
			password: password,
			prepare: func(m *repository.MockUserRepository) {
				m.EXPECT().GetUserByLogin(gomock.Any(), login).Return(nil, errors.New("internal"))

			},
			wantErr: true,
		},
		{
			name:     "wrong password",
			login:    login,
			password: "invalid",
			prepare: func(m *repository.MockUserRepository) {
				m.EXPECT().GetUserByLogin(gomock.Any(), login).Return(user, nil)

			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := repository.NewMockUserRepository(ctrl)
			s := NewUserService(m)

			tt.prepare(m)

			_, err := s.LoginUser(context.Background(), tt.login, tt.password)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
