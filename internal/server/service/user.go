package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Sadere/gophkeeper/internal/server/auth"
	"github.com/Sadere/gophkeeper/internal/server/model"
	"github.com/Sadere/gophkeeper/internal/server/repository"
)

//go:generate mockgen -source user.go -destination mocks/mock_user.go -package service

// User service interface
type IUserService interface {
	RegisterUser(ctx context.Context, login string, password string) (*model.User, error)
	LoginUser(ctx context.Context, login string, password string) (*model.User, error)
}

// User service implementation
type UserService struct {
	userRepo repository.UserRepository
}

// Returns new user service
func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// Registers new user
func (s *UserService) RegisterUser(ctx context.Context, login string, password string) (*model.User, error) {
	var newUser model.User
	user, err := s.userRepo.GetUserByLogin(ctx, login)

	// Checking if user exists
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	if user != nil {
		return nil, &model.ErrUserExists{Login: login}
	}

	// Hashing password
	passwordHash, err := auth.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to generate password hash: %w", err)
	}

	// Create new user
	newUser = model.User{
		Login:        login,
		PasswordHash: passwordHash,
		CreatedAt:    time.Now(),
	}

	var newUserID uint64
	newUserID, err = s.userRepo.Create(ctx, newUser)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	newUser.ID = newUserID

	return &newUser, nil
}

// Return user model if provided credentials are correct
func (s *UserService) LoginUser(ctx context.Context, login string, password string) (*model.User, error) {
	// Retrieve user
	user, err := s.userRepo.GetUserByLogin(ctx, login)

	if errors.Is(err, sql.ErrNoRows) {
		return user, model.ErrBadCredentials
	}

	if err != nil {
		return user, fmt.Errorf("failed to authenticate user: %w", err)
	}

	// Check password
	if !auth.CheckPassword(user.PasswordHash, password) {
		return user, model.ErrBadCredentials
	}

	return user, nil
}
