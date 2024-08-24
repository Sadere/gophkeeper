package service

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"

	"github.com/Sadere/gophkeeper/internal/server/crypto"
	"github.com/Sadere/gophkeeper/internal/server/model"
	"github.com/Sadere/gophkeeper/internal/server/repository"

	pkgModel "github.com/Sadere/gophkeeper/pkg/model"
)

type ISecretService interface {
	GetUserSecrets(ctx context.Context, userID uint64) (pkgModel.Secrets, error)
	AddSecret(ctx context.Context, password string, secret *pkgModel.Secret) error
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

func (s *SecretService) AddSecret(ctx context.Context, password string, secret *pkgModel.Secret) error {
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

	// encrypt secret data
	encrypted, err = crypto.Encrypt(password, payload)
	if err != nil {
		return fmt.Errorf("failed to encrypt data: %w", err)
	}

	secret.Payload = encrypted

	// store secret
	_, err = s.secretRepo.Create(ctx, secret)
	if err != nil {
		return fmt.Errorf("failed to store secret: %w", err)
	}

	return nil
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
