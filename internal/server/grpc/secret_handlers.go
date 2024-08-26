package grpc

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

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
	_, errSave := s.secretService.SaveSecret(ctx, in.MasterPassword, secret)

	if errors.Is(errSave, model.ErrWrongSecretType) {
		return nil, status.Error(codes.InvalidArgument, errSave.Error())
	}

	if errSave != nil {
		return nil, status.Error(codes.Internal, errSave.Error())
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

func (s *KeeperServer) UploadFileV1(stream pb.SecretsService_UploadFileV1Server) error {
	var (
		f        *os.File
		masterPw string
		secret   *pkgModel.Secret
		secretID uint64
		err      error
	)

	// Close file in the end
	defer func() {
		if f != nil {
			err := f.Close()
			if err != nil {
				s.log.Error(err)
			}
		}
	}()

	ctx := stream.Context()

	userID, err := extractUserID(ctx)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	for {
		// Read stream
		req, err := stream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			return status.Error(codes.Internal, err.Error())
		}

		// Save master pw
		if len(masterPw) == 0 {
			masterPw = req.MasterPassword
		}

		// Load secret
		if secret == nil {
			secret = &pkgModel.Secret{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				UserID:    userID,
				Metadata:  req.Metadata,
				SType:     string(pkgModel.BlobSecret),
				Blob: &pkgModel.Blob{
					FileName: req.FileName,
				},
			}
			secretID, err = s.secretService.SaveSecret(context.Background(), req.MasterPassword, secret)
			if err != nil {
				return status.Error(codes.Internal, err.Error())
			}

			secret.ID = secretID

			// Open file
			path := fmt.Sprintf("%s/%d", s.config.UploadDir, userID)

			f, err = s.openFile(path, secret.Blob.FileName)
			if err != nil {
				return status.Error(codes.Internal, err.Error())
			}
		}

		// Write chunk to file
		chunk := req.Chunk
		_, err = f.Write(chunk)
		if err != nil {
			return status.Error(codes.Internal, fmt.Sprintf("failed to write chunk: %s", err.Error()))
		}
	}

	// Update secret
	if secret != nil {
		secret.Blob.IsDone = true

		_, err = s.secretService.SaveSecret(context.Background(), masterPw, secret)
		if err != nil {
			return status.Error(codes.Internal, fmt.Sprintf("failed to update secret: %s", err.Error()))
		}
	}

	return stream.SendAndClose(&emptypb.Empty{})
}

func (s *KeeperServer) openFile(path, fileName string) (*os.File, error) {
	var f *os.File

	// Create upload dir if not exists
	_, err := os.Stat(path)

	if os.IsNotExist(err) {
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return nil, fmt.Errorf("failed to create upload dir: %w", err)
		}
	}

	if err != nil {
		return nil, err
	}

	// Open file
	filePath := fmt.Sprintf("%s/%s", path, fileName)

	f, err = os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	return f, nil
}
