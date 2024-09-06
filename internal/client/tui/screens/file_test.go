package screens

import (
	"errors"
	"io"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/golang/mock/gomock"
	"github.com/knz/catwalk"

	api "github.com/Sadere/gophkeeper/internal/client/api/mocks"
	"github.com/Sadere/gophkeeper/pkg/model"
)

func TestFileModel_Upload(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	c := api.NewMockIApiClient(ctrl)

	state := NewState(c)

	m := NewFileModel(state, 0)

	observer := func(w io.Writer, m tea.Model) error {
		fileModel, ok := m.(FileModel)
		if !ok {
			return errors.New("not file model")
		}

		if len(fileModel.filepicker.CurrentDirectory) == 0 {
			return errors.New("file picker failed to init")
		}

		return nil
	}

	catwalk.RunModel(t, "testdata/file/upload_tests", m, catwalk.WithObserver("file", observer))
}

func TestFileModel_Download(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	c := api.NewMockIApiClient(ctrl)

	state := NewState(c)
	secretID := uint64(222)
	secret := &model.Secret{
		ID:        secretID,
		Metadata:  "file meta",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		SType:     string(model.BlobSecret),
		Blob: &model.Blob{
			FileName: "file_download.txt",
		},
	}

	c.EXPECT().LoadSecret(gomock.Any(), secretID).Return(secret, nil)
	c.EXPECT().DownloadFile(gomock.Any(), secretID, secret.Blob.FileName).Return(nil)

	m := NewFileModel(state, secretID)

	catwalk.RunModel(t, "testdata/file/download_tests", m)
}
