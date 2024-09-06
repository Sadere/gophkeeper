package screens

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/knz/catwalk"

	api "github.com/Sadere/gophkeeper/internal/client/api/mocks"
	"github.com/Sadere/gophkeeper/pkg/model"
)

func TestCardModel_Submit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	c := api.NewMockIApiClient(ctrl)

	state := NewState(c)

	c.EXPECT().SaveCard(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	c.EXPECT().LoadPreviews(gomock.Any()).Return(nil, nil)

	m := NewCardModel(state, 0)

	catwalk.RunModel(t, "testdata/card/submit_tests", m)
}

func TestCardModel_Edit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	c := api.NewMockIApiClient(ctrl)

	state := NewState(c)
	secretID := uint64(111)
	secret := &model.Secret{
		ID:        secretID,
		Metadata:  "card meta",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		SType:     string(model.CardSecret),
		Card: &model.Card{
			Number:   "2444",
			ExpYear:  26,
			ExpMonth: 11,
			Cvv:      555,
		},
	}

	c.EXPECT().LoadSecret(gomock.Any(), secretID).Return(secret, nil)

	m := NewCardModel(state, secretID)

	catwalk.RunModel(t, "testdata/card/edit_tests", m)
}
