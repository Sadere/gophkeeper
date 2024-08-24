package grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/Sadere/gophkeeper/internal/server/model"
	"github.com/Sadere/gophkeeper/pkg/convert"

	pkgModel "github.com/Sadere/gophkeeper/pkg/model"
	pb "github.com/Sadere/gophkeeper/pkg/proto/keeper/v1"
)

func (s *KeeperServer) SecretPreviewsV1(ctx context.Context, in *emptypb.Empty) (*pb.SecretPreviewsV1Response, error) {
	var response pb.SecretPreviewsV1Response

	userID, err := extractUserID(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	secrets, err := s.secretService.GetUserSecrets(ctx, userID)
	if errors.Is(err, model.ErrNoSecrets) {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	for _, secret := range secrets {
		preview := &pkgModel.SecretPreview{
			ID:        secret.ID,
			CreatedAt: secret.CreatedAt,
			UpdatedAt: secret.UpdatedAt,
			Metadata:  secret.Metadata,
			SType:     secret.SType,
			Status:    pkgModel.SecretPreviewRead,
		}
		response.Previews = append(response.Previews, convert.PreviewToProto(preview))
	}

	return &response, nil
}

func (s *KeeperServer) SaveUserSecretV1(ctx context.Context, in *pb.SaveUserSecretV1Request) (*emptypb.Empty, error) {
	userID, err := extractUserID(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Validate request
	if err := validateRequest(in.Secret); err != nil {
		return nil, err
	}

	secret := convert.ProtoToSecret(in.Secret)
	secret.UserID = userID

	errAdd := s.secretService.AddSecret(ctx, in.MasterPassword, secret)

	if errors.Is(errAdd, model.ErrWrongSecretType) {
		return nil, status.Error(codes.InvalidArgument, errAdd.Error())
	}

	if errAdd != nil {
		return nil, status.Error(codes.Internal, errAdd.Error())
	}

	return &emptypb.Empty{}, nil
}
