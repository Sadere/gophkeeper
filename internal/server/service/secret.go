package service

import (
	"context"

	"github.com/Sadere/gophkeeper/internal/server/model"
	"github.com/Sadere/gophkeeper/internal/server/repository"

	pkgModel "github.com/Sadere/gophkeeper/pkg/model"
)

type ISecretService interface {
	GetUserSecrets(ctx context.Context, userID uint64) (pkgModel.Secrets, error)
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
