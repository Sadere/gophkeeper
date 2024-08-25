package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"slices"

	"github.com/Sadere/gophkeeper/internal/server/crypto"
	"github.com/Sadere/gophkeeper/internal/server/model"
	"github.com/Sadere/gophkeeper/internal/server/repository"

	pkgModel "github.com/Sadere/gophkeeper/pkg/model"
)

type ISecretService interface {
	GetUserSecrets(ctx context.Context, userID uint64) (pkgModel.Secrets, error)
	SaveSecret(ctx context.Context, password string, secret *pkgModel.Secret) error
	GetSecret(ctx context.Context, password string, ID uint64, userID uint64) (*pkgModel.Secret, error)
}

type SecretService struct {
	secretRepo repository.SecretRepository
}

func NewSecretService(secretRepo repository.SecretRepository) *SecretService {
	return &SecretService{
		secretRepo: secretRepo,
	}
}

func (s *SecretService) GetUserSecrets(ctx context.Context, userID uint64) (pkgModel.Secrets, error) {
	secrets, err := s.secretRepo.GetUserSecrets(ctx, userID)

	if len(secrets) == 0 {
		return nil, model.ErrNoSecrets
	}

	if err != nil {
		return nil, err
	}

	return secrets, nil
}

func (s *SecretService) SaveSecret(ctx context.Context, password string, secret *pkgModel.Secret) error {
	var (
		err                error
		payload, encrypted []byte
	)

	if ok := validateType(secret.SType); !ok {
		return model.ErrWrongSecretType
	}

	// marshal corresponding data
	switch secret.SType {
	case string(pkgModel.CredSecret):
		payload, err = json.Marshal(secret.Creds)
	}

	if err != nil {
		return fmt.Errorf("failed to save secret data: %w", err)
	}

	// Encrypt secret data
	encrypted, err = crypto.Encrypt(password, payload)
	if err != nil {
		return fmt.Errorf("failed to encrypt data: %w", err)
	}

	secret.Payload = encrypted

	// Store or update secret
	if secret.ID > 0 {
		// Check authority over existing secret
		_, err = s.secretRepo.GetSecret(ctx, secret.ID, secret.UserID)
		if err != nil {
			// Secret not found for this user
			return model.ErrSecretNotFound
		}

		err = s.secretRepo.Update(ctx, secret)
	} else {
		_, err = s.secretRepo.Create(ctx, secret)
	}

	if err != nil {
		return fmt.Errorf("failed to store secret: %w", err)
	}

	return nil
}

func (s *SecretService) GetSecret(ctx context.Context, password string, ID uint64, userID uint64) (*pkgModel.Secret, error) {
	secret, err := s.secretRepo.GetSecret(ctx, ID, userID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, model.ErrSecretNotFound
	}

	if err != nil {
		return nil, err
	}

	// Decrypt secret
	decrypted, err := crypto.Decrypt(password, secret.Payload)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt secret: %w", err)
	}

	// Unmarshal corresponding struct
	switch secret.SType {
	case string(pkgModel.CredSecret):
		secret.Creds = &pkgModel.Credentials{}
		if err := json.Unmarshal(decrypted, secret.Creds); err != nil {
			return nil, fmt.Errorf("failed to extract credentials: %w", err)
		}
	}

	return secret, nil
}

func validateType(sType string) bool {
	allowedTypes := []string{
		string(pkgModel.CredSecret),
		string(pkgModel.TextSecret),
		string(pkgModel.BlobSecret),
		string(pkgModel.CardSecret),
	}

	return slices.Contains(allowedTypes, sType)
}
