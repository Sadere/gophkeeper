package screens

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/knz/catwalk"

	api "github.com/Sadere/gophkeeper/internal/client/api/mocks"
)

func TestRegisterModel(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	c := api.NewMockIApiClient(ctrl)

	state := NewState(c)

	// First submit with error
	c.EXPECT().Register(gomock.Any(), "login", "password").Return("", errors.New("error"))

	// Successful submit
	c.EXPECT().Register(gomock.Any(), "login", "password").Return("", nil)
	c.EXPECT().LoadPreviews(gomock.Any()).Return(nil, nil)

	m := NewRegisterModel(state)

	catwalk.RunModel(t, "testdata/register/tests", m)
}
