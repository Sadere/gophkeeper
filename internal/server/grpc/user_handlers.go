package grpc

import (
	"context"
	"errors"
	"time"

	"github.com/Sadere/gophkeeper/internal/server/auth"
	"github.com/Sadere/gophkeeper/internal/server/model"
	"github.com/bufbuild/protovalidate-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/Sadere/gophkeeper/pkg/constants"
	pb "github.com/Sadere/gophkeeper/pkg/proto/keeper/v1"
)

func (s *KeeperServer) RegisterV1(ctx context.Context, in *pb.RegisterV1Request) (*pb.RegisterV1Response, error) {
	var response pb.RegisterV1Response

	// Validate request
	if err := validateRequest(in); err != nil {
		return nil, err
	}

	// Register user
	user, err := s.userService.RegisterUser(ctx, in.Login, in.Password)

	// Check if user exists
	if errors.Is(err, &model.ErrUserExists{Login: in.Login}) {
		return nil, status.Error(codes.AlreadyExists, err.Error())
	}

	// Other errors
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	// Generate access token
	token, err := s.authUser(user.ID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to auth: %v", err)
	}

	response.AccessToken = token

	return &response, nil
}

func (s *KeeperServer) LoginV1(ctx context.Context, in *pb.LoginV1Request) (*pb.LoginV1Response, error) {
	var response pb.LoginV1Response

	// Validate request
	if err := validateRequest(in); err != nil {
		return nil, err
	}

	// Login user
	user, err := s.userService.LoginUser(ctx, in.Login, in.Password)

	// Check credentials
	if errors.Is(err, model.ErrBadCredentials) {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	// Other errors
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	// Generate access token
	token, err := s.authUser(user.ID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to auth: %v", err)
	}

	response.AccessToken = token

	return &response, nil
}

func validateRequest(in protoreflect.ProtoMessage) error {
	v, err := protovalidate.New()
	if err != nil {
		return status.Errorf(codes.Internal, "failed to init validator: %s", err)
	}

	if err = v.Validate(in); err != nil {
		return status.Errorf(codes.InvalidArgument, "failed to validate request: %s", err)
	}

	return nil
}

func (s *KeeperServer) authUser(userID uint64) (string, error) {
	token, err := auth.CreateToken(userID, time.Now().Add(time.Hour*24), []byte(s.config.SecretKey))
	if err != nil {
		return "", err
	}

	return token, nil
}

func extractUserID(ctx context.Context) (uint64, error) {
	uid := ctx.Value(constants.CtxUserIDKey)

	userID, ok := uid.(uint64)
	if !ok {
		return 0, errors.New("failed to extract user id from context")
	}

	return userID, nil
}
