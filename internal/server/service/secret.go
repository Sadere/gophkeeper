// Service layer of server
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
	"github.com/Sadere/gophkeeper/internal/server/utils"

	pkgModel "github.com/Sadere/gophkeeper/pkg/model"
)

//go:generate mockgen -source secret.go -destination mocks/mock_secret.go -package service

// Interface for secrets management
type ISecretService interface {
	GetUserSecrets(ctx context.Context, userID uint64) (pkgModel.Secrets, error)
	SaveSecret(ctx context.Context, password string, secret *pkgModel.Secret) (uint64, error)
	GetSecret(ctx context.Context, password string, ID uint64, userID uint64) (*pkgModel.Secret, error)
}

// Secret service implementation
type SecretService struct {
	secretRepo repository.SecretRepository
}

// Returns new secret service
func NewSecretService(secretRepo repository.SecretRepository) *SecretService {
	return &SecretService{
		secretRepo: secretRepo,
	}
}

// Returns list of secrets belonging to given user ID
func (s *SecretService) GetUserSecrets(ctx context.Context, userID uint64) (pkgModel.Secrets, error) {
	secrets, err := s.secretRepo.GetUserSecrets(ctx, userID)

	if err != nil {
		return nil, err
	}

	if len(secrets) == 0 {
		return nil, model.ErrNoSecrets
	}

	return secrets, nil
}

// Creates new or updates existsing secret and encrypts underlying secret data
func (s *SecretService) SaveSecret(ctx context.Context, password string, secret *pkgModel.Secret) (uint64, error) {
	var (
		err                error
		payload, encrypted []byte
		secretID           uint64
	)

	if ok := validateType(secret.SType); !ok {
		return 0, model.ErrWrongSecretType
	}

	// marshal corresponding data
	switch secret.SType {
	case string(pkgModel.CredSecret):
		payload, err = json.Marshal(secret.Creds)
	case string(pkgModel.TextSecret):
		payload, err = json.Marshal(secret.Text)
	case string(pkgModel.CardSecret):
		payload, err = json.Marshal(secret.Card)
	case string(pkgModel.BlobSecret):
		payload, err = json.Marshal(secret.Blob)
	}

	if err != nil {
		return 0, fmt.Errorf("failed to save secret data: %w", err)
	}

	// validate card number using Luhn's algo
	if secret.SType == string(pkgModel.CardSecret) {
		if ok := utils.CheckLuhn(secret.Card.Number); !ok {
			return 0, model.ErrNumberInvaliod
		}
	}

	// Encrypt secret data
	encrypted, err = crypto.Encrypt(password, payload)
	if err != nil {
		return 0, fmt.Errorf("failed to encrypt data: %w", err)
	}

	secret.Payload = encrypted

	// Store or update secret
	if secret.ID > 0 {
		// Check authority over existing secret
		_, err = s.secretRepo.GetSecret(ctx, secret.ID, secret.UserID)
		if err != nil {
			// Secret not found for this user
			return 0, model.ErrSecretNotFound
		}

		err = s.secretRepo.Update(ctx, secret)
		secretID = secret.ID
	} else {
		secretID, err = s.secretRepo.Create(ctx, secret)
	}

	if err != nil {
		return 0, fmt.Errorf("failed to store secret: %w", err)
	}

	return secretID, nil
}

// Returns decrypted secret
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
	case string(pkgModel.TextSecret):
		secret.Text = &pkgModel.Text{}
		if err := json.Unmarshal(decrypted, secret.Text); err != nil {
			return nil, fmt.Errorf("failed to extract text: %w", err)
		}
	case string(pkgModel.CardSecret):
		secret.Card = &pkgModel.Card{}
		if err := json.Unmarshal(decrypted, secret.Card); err != nil {
			return nil, fmt.Errorf("failed to extract card: %w", err)
		}
	case string(pkgModel.BlobSecret):
		secret.Blob = &pkgModel.Blob{}
		if err := json.Unmarshal(decrypted, secret.Blob); err != nil {
			return nil, fmt.Errorf("failed to extract blob: %w", err)
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
