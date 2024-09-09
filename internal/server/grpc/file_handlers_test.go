package grpc

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/Sadere/gophkeeper/pkg/constants"
	"github.com/Sadere/gophkeeper/pkg/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	pb "github.com/Sadere/gophkeeper/pkg/proto/keeper/v1"
	mock "github.com/Sadere/gophkeeper/pkg/proto/keeper/v1/mocks"
)

func TestUploadFile(t *testing.T) {
	// Assign variables
	userID := uint64(111)
	masterPw := "password"
	ctxUser := context.WithValue(context.Background(), constants.CtxUserIDKey, userID)
	testFileName := "test.txt"

	// Create temp dir
	tempDir, err := os.MkdirTemp("", "uploaddir")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	streamCtrl := gomock.NewController(t)
	defer streamCtrl.Finish()

	stream := mock.NewMockSecretsService_UploadFileV1Server(streamCtrl)

	// Setup stream mock
	stream.EXPECT().Context().Return(ctxUser)
	stream.EXPECT().Recv().Return(&pb.UploadFileRequestV1{
		Metadata:       "metadata",
		FileName:       testFileName,
		MasterPassword: masterPw,
		Chunk:          []byte("test file\n"),
	}, nil)
	stream.EXPECT().Recv().Return(nil, io.EOF)
	stream.EXPECT().SendAndClose(gomock.Any()).Return(nil)

	server, _, secretMock, userCtrl, secretCtrl := NewTestServer(t)
	defer func() {
		userCtrl.Finish()
		secretCtrl.Finish()
	}()

	// Upload to temp dir
	server.config.UploadDir = tempDir

	secretMock.EXPECT().SaveSecret(gomock.Any(), masterPw, gomock.Any()).MinTimes(1).Return(uint64(0), nil)

	err = server.UploadFileV1(stream)

	assert.NoError(t, err)
}

func TestDownloadFile(t *testing.T) {
	// Assign variables
	userID := uint64(111)
	request := &pb.DownloadFileRequestV1{
		Id:             uint64(123),
		MasterPassword: "password",
	}
	ctxUser := context.WithValue(context.Background(), constants.CtxUserIDKey, userID)

	// Create temp dir
	tempDir, err := os.MkdirTemp("", "downloaddir")
	require.NoError(t, err)

	// Create user dir in temp dir
	userDir := fmt.Sprintf("%s/%d", tempDir, userID)
	err = os.MkdirAll(userDir, 0700)
	require.NoError(t, err)

	// Remove temp dir
	defer os.RemoveAll(tempDir)

	// Create temp file
	tempFile, err := os.CreateTemp(userDir, "test.txt")
	require.NoError(t, err)

	// Write temp file
	_, err = tempFile.WriteString("test file\n")
	require.NoError(t, err)

	// Get tmp file name
	tempFileName := filepath.Base(tempFile.Name())

	err = tempFile.Close()
	require.NoError(t, err)

	streamCtrl := gomock.NewController(t)
	defer streamCtrl.Finish()

	stream := mock.NewMockSecretsService_DownloadFileV1Server(streamCtrl)

	// Setup stream mock
	stream.EXPECT().Context().Return(ctxUser)
	stream.EXPECT().Send(&pb.DownloadFileResponseV1{
		Chunk: []byte("test file\n"),
	}).Return(nil)

	server, _, secretMock, userCtrl, secretCtrl := NewTestServer(t)
	defer func() {
		userCtrl.Finish()
		secretCtrl.Finish()
	}()

	// Upload to temp dir
	server.config.UploadDir = tempDir

	// Setup secret mock
	secret := &model.Secret{
		Blob: &model.Blob{
			FileName: tempFileName,
		},
	}
	secretMock.EXPECT().GetSecret(gomock.Any(), request.MasterPassword, request.Id, userID).Return(secret, nil)

	err = server.DownloadFileV1(request, stream)

	assert.NoError(t, err)
}
