package repository

import (
	"context"

	"github.com/Sadere/gophkeeper/pkg/model"
	"github.com/jmoiron/sqlx"
)

type SecretRepository interface {
	GetUserSecrets(ctx context.Context, userID uint64) (model.Secrets, error)
}

type PgSecretRepository struct {
	db *sqlx.DB
}

func NewPgSecretRepository(db *sqlx.DB) *PgSecretRepository {
	return &PgSecretRepository{
		db: db,
	}
}

func (r *PgSecretRepository) GetUserSecrets(ctx context.Context, userID uint64) (model.Secrets, error) {
	var secrets model.Secrets

	sql := "SELECT * FROM entries WHERE user_id = $1 ORDER BY updated_at DESC"
	err := r.db.SelectContext(ctx, &secrets, sql, userID)
	if err != nil {
		return nil, err
	}

	return secrets, nil
}
