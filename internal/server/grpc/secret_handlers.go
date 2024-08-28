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

// Saves new secret or updates existing one
func (s *KeeperServer) SaveUserSecretV1(ctx context.Context, in *pb.SaveUserSecretV1Request) (*emptypb.Empty, error) {
	userID, err := extractUserID(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Validate request
	if err := validateRequest(in); err != nil {
		return nil, err
	}

	secret := convert.ProtoToSecret(in.Secret)
	secret.UserID = userID

	// Save secret
	secretID, errSave := s.secretService.SaveSecret(ctx, in.MasterPassword, secret)

	if errors.Is(errSave, model.ErrWrongSecretType) {
		return nil, status.Error(codes.InvalidArgument, errSave.Error())
	}

	if errSave != nil {
		return nil, status.Error(codes.Internal, errSave.Error())
	}

	// Send notifications
	clientID, err := extractClientID(ctx)
	if err == nil {
		if secret.ID > 0 {
			// Update notification
			err = s.notifyClients(userID, clientID, secret.ID, true)
		} else {
			// New secret notification
			err = s.notifyClients(userID, clientID, secretID, false)
		}

		if err != nil {
			s.log.Error("failed to notify clients: ", err)
		}
	}

	return &emptypb.Empty{}, nil
}

func (s *KeeperServer) GetUserSecretV1(ctx context.Context, in *pb.GetUserSecretV1Request) (*pb.GetUserSecretV1Response, error) {
	var response pb.GetUserSecretV1Response

	userID, err := extractUserID(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Validate request
	if err := validateRequest(in); err != nil {
		return nil, err
	}

	// Acquire secret
	secret, err := s.secretService.GetSecret(ctx, in.MasterPassword, in.Id, userID)
	if errors.Is(err, model.ErrSecretNotFound) {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	// Other errors
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	response.Secret = convert.SecretToProto(secret)

	return &response, nil
}
