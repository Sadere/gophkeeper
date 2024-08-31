package repository

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Sadere/gophkeeper/pkg/model"
	"github.com/stretchr/testify/assert"
)

func newTestSecret(userID uint64) *model.Secret {
	return &model.Secret{
		ID:        111,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    uint64(userID),
		Metadata:  "Test metadata",
		SType:     string(model.CardSecret),
		Payload:   []byte{},
	}
}

func TestGetUserSecrets(t *testing.T) {
	db, mock := NewMock(t)
	repo := NewPgSecretRepository(db)

	userID := 333
	expectedSecret := newTestSecret(uint64(userID))

	rows := sqlmock.NewRows([]string{
		"id",
		"created_at",
		"updated_at",
		"user_id",
		"metadata",
		"ent_type",
		"payload",
	}).
		AddRow(
			expectedSecret.ID,
			expectedSecret.CreatedAt,
			expectedSecret.UpdatedAt,
			expectedSecret.UserID,
			expectedSecret.Metadata,
			expectedSecret.SType,
			expectedSecret.Payload,
		)

	mock.ExpectQuery("SELECT (.+) FROM entries WHERE user_id").
		WithArgs(userID).
		WillReturnRows(rows)

	// Test function
	actualSecrets, err := repo.GetUserSecrets(context.Background(), uint64(userID))

	assert.NoError(t, err)

	assert.Len(t, actualSecrets, 1)

	if !reflect.DeepEqual(expectedSecret, actualSecrets[0]) {
		t.Errorf("unexpected secret want = %v got = %v", expectedSecret, actualSecrets[0])
	}
}

func TestSecretCreateSuccess(t *testing.T) {
	db, mock := NewMock(t)
	repo := NewPgSecretRepository(db)

	inSecret := newTestSecret(333)

	expectedSecretID := uint64(444)

	rows := sqlmock.NewRows([]string{"id"}).AddRow(expectedSecretID)
	mock.ExpectQuery("INSERT INTO entries").WithArgs(
		inSecret.UserID,
		inSecret.Metadata,
		inSecret.SType,
		inSecret.Payload,
	).WillReturnRows(rows)

	// Test function
	actualSecretID, err := repo.Create(context.Background(), inSecret)

	assert.NoError(t, err)

	assert.Equal(t, expectedSecretID, actualSecretID)
}

func TestSecretUpdate(t *testing.T) {
	db, mock := NewMock(t)
	repo := NewPgSecretRepository(db)

	inSecret := newTestSecret(333)

	result := sqlmock.NewResult(0, 1)
	mock.ExpectExec("UPDATE entries").WithArgs(
		inSecret.UpdatedAt,
		inSecret.Metadata,
		inSecret.SType,
		inSecret.Payload,
		inSecret.ID,
	).WillReturnResult(result)

	// Test function
	err := repo.Update(context.Background(), inSecret)

	assert.NoError(t, err)
}

func TestGetSecret(t *testing.T) {
	db, mock := NewMock(t)
	repo := NewPgSecretRepository(db)

	userID := 333
	expectedSecret := newTestSecret(uint64(userID))

	rows := sqlmock.NewRows([]string{
		"id",
		"created_at",
		"updated_at",
		"user_id",
		"metadata",
		"ent_type",
		"payload",
	}).
		AddRow(
			expectedSecret.ID,
			expectedSecret.CreatedAt,
			expectedSecret.UpdatedAt,
			expectedSecret.UserID,
			expectedSecret.Metadata,
			expectedSecret.SType,
			expectedSecret.Payload,
		)

	mock.ExpectQuery("SELECT (.+) FROM entries WHERE id").
		WithArgs(expectedSecret.ID, userID).
		WillReturnRows(rows)

	// Test function
	actualSecret, err := repo.GetSecret(context.Background(), expectedSecret.ID, uint64(userID))

	assert.NoError(t, err)

	if !reflect.DeepEqual(expectedSecret, actualSecret) {
		t.Errorf("unexpected secret want = %v got = %v", expectedSecret, actualSecret)
	}
}
