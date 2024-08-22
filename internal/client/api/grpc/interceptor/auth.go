package interceptor

import (
	"context"

	"github.com/Sadere/gophkeeper/pkg/constants"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func AddToken(token *string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		// pass request if token is empty
		if len(*token) == 0 {
			return invoker(ctx, method, req, reply, cc, opts...)
		}

		// add access token to metadata
		md := metadata.New(map[string]string{
			constants.AccessTokenHeader: *token,
		})

		mdCtx := metadata.NewOutgoingContext(ctx, md)
		return invoker(mdCtx, method, req, reply, cc, opts...)
	}
}
