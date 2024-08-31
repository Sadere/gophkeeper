package repository

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Sadere/gophkeeper/internal/server/model"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func NewMock(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("failed to create mock db: %s", err)
	}

	dbx := sqlx.NewDb(db, "sqlmock")

	return dbx, mock
}

func TestUserCreateSuccess(t *testing.T) {
	db, mock := NewMock(t)
	repo := NewPgUserRepository(db)

	inUser := model.User{
		Login:        "test_login",
		PasswordHash: "test_pw",
		CreatedAt:    time.Now(),
	}

	expectedUserID := uint64(444)

	rows := sqlmock.NewRows([]string{"id"}).AddRow(expectedUserID)
	mock.ExpectQuery("INSERT INTO users").WithArgs(inUser.Login, inUser.PasswordHash, inUser.CreatedAt).WillReturnRows(rows)

	// Test function
	actualUserID, err := repo.Create(context.Background(), inUser)

	assert.NoError(t, err)

	assert.Equal(t, expectedUserID, actualUserID)
}

func TestUserCreateError(t *testing.T) {
	db, mock := NewMock(t)
	repo := NewPgUserRepository(db)

	inUser := model.User{
		Login:        "test_login",
		PasswordHash: "test_pw",
		CreatedAt:    time.Now(),
	}

	expectedUserID := uint64(0)

	mock.ExpectQuery("INSERT INTO users").WithArgs(inUser.Login, inUser.PasswordHash, inUser.CreatedAt).WillReturnError(errors.New("err"))

	// Test function
	actualUserID, err := repo.Create(context.Background(), inUser)

	assert.Error(t, err)

	assert.Equal(t, expectedUserID, actualUserID)
}

func TestGetUserByID(t *testing.T) {
	db, mock := NewMock(t)
	repo := NewPgUserRepository(db)

	expectedUser := &model.User{
		ID:           uint64(555),
		Login:        "test_login",
		PasswordHash: "test_pw",
		CreatedAt:    time.Now(),
	}

	rows := sqlmock.NewRows([]string{"id", "login", "created_at", "password"}).
		AddRow(expectedUser.ID, expectedUser.Login, expectedUser.CreatedAt, expectedUser.PasswordHash)
	mock.ExpectQuery("SELECT (.+) FROM users WHERE id").WithArgs(expectedUser.ID).WillReturnRows(rows)

	// Test function
	actualUser, err := repo.GetUserByID(context.Background(), expectedUser.ID)

	assert.NoError(t, err)

	if !reflect.DeepEqual(expectedUser, actualUser) {
		t.Errorf("unexpected user want = %v got = %v", expectedUser, actualUser)
	}
}

func TestGetUserByLogin(t *testing.T) {
	db, mock := NewMock(t)
	repo := NewPgUserRepository(db)

	expectedUser := &model.User{
		ID:           uint64(555),
		Login:        "test_login",
		PasswordHash: "test_pw",
		CreatedAt:    time.Now(),
	}

	rows := sqlmock.NewRows([]string{"id", "login", "created_at", "password"}).
		AddRow(expectedUser.ID, expectedUser.Login, expectedUser.CreatedAt, expectedUser.PasswordHash)
	mock.ExpectQuery("SELECT (.+) FROM users WHERE login").WithArgs(expectedUser.Login).WillReturnRows(rows)

	// Test function
	actualUser, err := repo.GetUserByLogin(context.Background(), expectedUser.Login)

	assert.NoError(t, err)

	if !reflect.DeepEqual(expectedUser, actualUser) {
		t.Errorf("unexpected user want = %v got = %v", expectedUser, actualUser)
	}
}
