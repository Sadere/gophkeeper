package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Sadere/gophkeeper/internal/server/model"
	"github.com/jmoiron/sqlx"
)

//go:generate mockgen -source user.go -destination mocks/mock_user.go -package repository
type UserRepository interface {
	Create(ctx context.Context, user model.User) (uint64, error)
	GetUserByID(ctx context.Context, ID uint64) (*model.User, error)
	GetUserByLogin(ctx context.Context, login string) (*model.User, error)
}

type PgUserRepository struct {
	db *sqlx.DB
}

func NewPgUserRepository(db *sqlx.DB) *PgUserRepository {
	return &PgUserRepository{
		db: db,
	}
}

// Creates new user and returns new user id
func (r *PgUserRepository) Create(ctx context.Context, user model.User) (uint64, error) {
	var newUserID uint64

	result := r.db.QueryRowContext(ctx, "INSERT INTO users (login, password, created_at) VALUES ($1, $2, $3) RETURNING id",
		user.Login,
		user.PasswordHash,
		user.CreatedAt,
	)

	err := result.Scan(&newUserID)

	if err != nil {
		return 0, err
	}

	return newUserID, nil
}

// Gets user by user ID
func (r *PgUserRepository) GetUserByID(ctx context.Context, ID uint64) (*model.User, error) {
	var user model.User

	err := r.db.QueryRowxContext(ctx, "SELECT id, login, created_at, password FROM users WHERE id = $1", ID).StructScan(&user)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	return &user, err
}

// Gets user by login
func (r *PgUserRepository) GetUserByLogin(ctx context.Context, login string) (*model.User, error) {
	var user model.User

	err := r.db.QueryRowxContext(ctx, "SELECT id, login, created_at, password FROM users WHERE login = $1", login).StructScan(&user)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	return &user, err
}
