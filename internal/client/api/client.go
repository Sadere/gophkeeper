package api

import (
	"context"
)

type IApiClient interface {
	Register(ctx context.Context, login string, password string) (string, error)
	Login(ctx context.Context, login string, password string) (string, error)
}
