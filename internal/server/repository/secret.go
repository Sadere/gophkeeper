package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Sadere/gophkeeper/pkg/model"
	"github.com/jmoiron/sqlx"
)

type SecretRepository interface {
	GetUserSecrets(ctx context.Context, userID uint64) (model.Secrets, error)
	Create(ctx context.Context, secret *model.Secret) (uint64, error)
	Update(ctx context.Context, secret *model.Secret) error
	GetSecret(ctx context.Context, secretID uint64, userID uint64) (*model.Secret, error)
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

func (r *PgSecretRepository) Create(ctx context.Context, secret *model.Secret) (uint64, error) {
	var newSecretID uint64

	result := r.db.QueryRowContext(ctx,
		`INSERT INTO entries
			(user_id, metadata, ent_type, payload)
		VALUES ($1, $2, $3, $4)
		RETURNING id`,
		secret.UserID,
		secret.Metadata,
		secret.SType,
		secret.Payload,
	)

	err := result.Scan(&newSecretID)

	if err != nil {
		return 0, err
	}

	return newSecretID, nil
}

func (r *PgSecretRepository) Update(ctx context.Context, secret *model.Secret) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE entries
		SET
			updated_at = $1,
			metadata = $2,
			ent_type = $3,
			payload = $4
		WHERE id = $5`,

		secret.UpdatedAt,
		secret.Metadata,
		secret.SType,
		secret.Payload,
	)

	return err
}

func (r *PgSecretRepository) GetSecret(ctx context.Context, secretID uint64, userID uint64) (*model.Secret, error) {
	var secret model.Secret

	err := r.db.QueryRowxContext(ctx,
		`SELECT
			id, created_at, updated_at, metadata, ent_type, payload
		FROM users
		WHERE id = $1 AND user_id = $2`, secretID, userID).StructScan(&secret)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	return &secret, err
}