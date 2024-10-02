package screens

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/knz/catwalk"

	api "github.com/Sadere/gophkeeper/internal/client/api/mocks"
	"github.com/Sadere/gophkeeper/pkg/model"
)

func TestTextModel_Submit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	c := api.NewMockIApiClient(ctrl)

	state := NewState(c)

	// First submit erroneus
	c.EXPECT().SaveText(gomock.Any(), gomock.Any(), "metadata", "content").Return(errors.New("fail"))

	// Successful submit
	c.EXPECT().SaveText(gomock.Any(), gomock.Any(), "metadata", "content").Return(nil)
	c.EXPECT().LoadPreviews(gomock.Any()).Return(nil, nil)

	m := NewTextModel(state, 0)

	catwalk.RunModel(t, "testdata/text/submit_tests", m)
}

func TestTextModel_Edit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	c := api.NewMockIApiClient(ctrl)

	state := NewState(c)

	secretID := uint64(111)
	secret := &model.Secret{
		ID:       secretID,
		Metadata: "text meta",
		SType:    string(model.TextSecret),
		Text: &model.Text{
			Content: "text content",
		},
	}

	c.EXPECT().LoadSecret(gomock.Any(), secretID).Return(secret, nil)
	c.EXPECT().LoadPreviews(gomock.Any()).Return(nil, nil)

	m := NewTextModel(state, secretID)

	catwalk.RunModel(t, "testdata/text/edit_tests", m)
}
