package screens

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/knz/catwalk"

	api "github.com/Sadere/gophkeeper/internal/client/api/mocks"
	"github.com/Sadere/gophkeeper/pkg/model"
)

func testPreviews() model.SecretPreviews {
	return model.SecretPreviews{
		&model.SecretPreview{
			ID:       1,
			Metadata: "creds",
			SType:    string(model.CredSecret),
			Status:   model.SecretPreviewNew,
		},
		&model.SecretPreview{
			ID:       2,
			Metadata: "card",
			SType:    string(model.CardSecret),
			Status:   model.SecretPreviewUpdated,
		},
		&model.SecretPreview{
			ID:       3,
			Metadata: "text",
			SType:    string(model.TextSecret),
		},
		&model.SecretPreview{
			ID:       4,
			Metadata: "blob",
			SType:    string(model.BlobSecret),
		},
	}
}

func TestSecretListModel_ListItems(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	c := api.NewMockIApiClient(ctrl)

	state := NewState(c)

	previews := testPreviews()

	// Load list
	c.EXPECT().LoadPreviews(gomock.Any()).Return(previews, nil)

	m := NewSecretListModel(state)

	catwalk.RunModel(t, "testdata/secret_list/list_tests", m)
}

func TestSecretListModel_Edits(t *testing.T) {
	tests := []struct {
		name     string
		testPath string
		prepare  func(c *api.MockIApiClient)
	}{
		{
			name: "edit creds",
			prepare: func(c *api.MockIApiClient) {
				secret := &model.Secret{
					ID:       uint64(1),
					Metadata: "cred meta",
					SType:    string(model.CredSecret),
					Creds: &model.Credentials{
						Login:    "login",
						Password: "password",
					},
				}

				c.EXPECT().LoadSecret(gomock.Any(), uint64(1)).Return(secret, nil)
			},
			testPath: "testdata/secret_list/cred_test",
		},
		{
			name: "edit text",
			prepare: func(c *api.MockIApiClient) {
				secret := &model.Secret{
					ID:       uint64(3),
					Metadata: "text meta",
					SType:    string(model.TextSecret),
					Text: &model.Text{
						Content: "text content",
					},
				}

				c.EXPECT().LoadSecret(gomock.Any(), uint64(3)).Return(secret, nil)
			},
			testPath: "testdata/secret_list/text_test",
		},
		{
			name: "edit card",
			prepare: func(c *api.MockIApiClient) {
				secret := &model.Secret{
					ID:       uint64(2),
					Metadata: "card meta",
					SType:    string(model.CardSecret),
					Card: &model.Card{
						Number:   "2444",
						ExpYear:  26,
						ExpMonth: 11,
						Cvv:      555,
					},
				}

				c.EXPECT().LoadSecret(gomock.Any(), uint64(2)).Return(secret, nil)
			},
			testPath: "testdata/secret_list/card_test",
		},
		{
			name:     "download file",
			prepare:  func(c *api.MockIApiClient) {},
			testPath: "testdata/secret_list/file_test",
		},
		{
			name: "reload list",
			prepare: func(c *api.MockIApiClient) {
				c.EXPECT().LoadPreviews(gomock.Any()).Return(nil, nil)
			},
			testPath: "testdata/secret_list/reload_test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			c := api.NewMockIApiClient(ctrl)

			state := NewState(c)

			previews := testPreviews()

			// Load list
			c.EXPECT().LoadPreviews(gomock.Any()).Return(previews, nil)

			// Prepare client mock
			tt.prepare(c)

			m := NewSecretListModel(state)

			// Tests
			catwalk.RunModel(t, tt.testPath, m)
		})
	}
}

func TestSecretListModel_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	c := api.NewMockIApiClient(ctrl)

	state := NewState(c)

	// Load list
	c.EXPECT().LoadPreviews(gomock.Any()).Return(nil, errors.New("error list"))

	m := NewSecretListModel(state)

	catwalk.RunModel(t, "testdata/secret_list/error_test", m)
}
