package interceptor

import (
	"context"
	"errors"
	"testing"

	"github.com/Sadere/gophkeeper/pkg/constants"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var (
	errUnauth = errors.New("no token")
	testToken = "AccessToken"
)

func checkCtx(ctx context.Context) error {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		return errUnauth
	}

	values := md.Get(constants.AccessTokenHeader)
	if len(values) == 0 {
		return errUnauth
	}

	if values[0] != testToken {
		return errUnauth
	}

	return nil
}

func TestAddAuth(t *testing.T) {
	invoker := func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
		return checkCtx(ctx)
	}

	t.Run("empty token", func(t *testing.T) {
		token := ""
		interceptor := AddAuth(&token, 11)

		err := interceptor(context.Background(), "SomeMethod", nil, nil, nil, invoker)

		assert.EqualError(t, err, errUnauth.Error())
	})

	t.Run("pass token", func(t *testing.T) {
		token := testToken
		interceptor := AddAuth(&token, 11)

		err := interceptor(context.Background(), "SomeMethod", nil, nil, nil, invoker)

		assert.NoError(t, err)
	})
}

func TestAddAuthStream(t *testing.T) {
	streamer := func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		return nil, checkCtx(ctx)
	}

	t.Run("empty token", func(t *testing.T) {
		token := ""
		interceptor := AddAuthStream(&token, 11)

		_, err := interceptor(context.Background(), nil, nil, "SomeMethod", streamer)

		assert.EqualError(t, err, errUnauth.Error())
	})

	t.Run("pass token", func(t *testing.T) {
		token := testToken
		interceptor := AddAuthStream(&token, 11)

		_, err := interceptor(context.Background(), nil, nil, "SomeMethod", streamer)

		assert.NoError(t, err)
	})
}