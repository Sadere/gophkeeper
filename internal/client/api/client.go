package api

import (
	"context"

	"github.com/Sadere/gophkeeper/pkg/model"
)

type IApiClient interface {
	Register(ctx context.Context, login string, password string) (string, error)
	Login(ctx context.Context, login string, password string) (string, error)

	LoadPreviews(ctx context.Context) (model.SecretPreviews, error)
	SaveCredential(ctx context.Context, metadata, login, password string) error
}
