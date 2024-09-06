package screens

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/knz/catwalk"

	api "github.com/Sadere/gophkeeper/internal/client/api/mocks"
	"github.com/Sadere/gophkeeper/pkg/model"
)

func TestCredModel_Submit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	c := api.NewMockIApiClient(ctrl)

	state := NewState(c)

	// First submit erroneus
	c.EXPECT().SaveCredential(gomock.Any(), gomock.Any(), "metadata", "login", "password").Return(errors.New("fail"))

	// Successful submit
	c.EXPECT().SaveCredential(gomock.Any(), gomock.Any(), "metadata", "login", "password").Return(nil)
	c.EXPECT().LoadPreviews(gomock.Any()).Return(nil, nil)

	m := NewCredentialModel(state, 0)

	catwalk.RunModel(t, "testdata/credential_submit_tests", m)
}

func TestCredModel_Edit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	c := api.NewMockIApiClient(ctrl)

	state := NewState(c)

	secretID := uint64(111)
	secret := &model.Secret{
		ID:       secretID,
		Metadata: "cred meta",
		SType:    string(model.CredSecret),
		Creds: &model.Credentials{
			Login:    "login",
			Password: "password",
		},
	}

	c.EXPECT().LoadSecret(gomock.Any(), secretID).Return(secret, nil)
	c.EXPECT().LoadPreviews(gomock.Any()).Return(nil, nil)

	m := NewCredentialModel(state, secretID)

	catwalk.RunModel(t, "testdata/credential_edit_tests", m)
}
