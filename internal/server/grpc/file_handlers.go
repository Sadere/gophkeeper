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
	"github.com/Sadere/gophkeeper/pkg/constants"
	pkgModel "github.com/Sadere/gophkeeper/pkg/model"
	pb "github.com/Sadere/gophkeeper/pkg/proto/keeper/v1"
)

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

func (s *KeeperServer) DownloadFileV1(in *pb.DownloadFileV1Request, srv pb.SecretsService_DownloadFileV1Server) error {
	ctx := srv.Context()

	userID, err := extractUserID(ctx)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	// Load secret
	secret, err := s.secretService.GetSecret(ctx, in.MasterPassword, in.Id, userID)

	if errors.Is(err, model.ErrSecretNotFound) {
		return status.Error(codes.NotFound, err.Error())
	}

	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	// Open file
	path := fmt.Sprintf("%s/%d/%s", s.config.UploadDir, userID, secret.Blob.FileName)

	f, err := os.Open(path)
	if err != nil {
		return status.Error(codes.NotFound, err.Error())
	}

	buf := make([]byte, constants.ChunkSize)

	for {
		n, err := f.Read(buf)

		// File is complete
		if err == io.EOF {
			break
		}

		if err != nil {
			return status.Error(codes.Internal, err.Error())
		}

		chunk := buf[:n]

		resp := &pb.DownloadFileV1Response{
			Chunk: chunk,
		}

		if err := srv.Send(resp); err != nil {
			return status.Error(codes.Internal, err.Error())
		}
	}

	return nil
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
