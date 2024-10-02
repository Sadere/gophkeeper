package screens

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/knz/catwalk"

	api "github.com/Sadere/gophkeeper/internal/client/api/mocks"
)

func TestChooseTypeModel_Back(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	c := api.NewMockIApiClient(ctrl)

	state := NewState(c)

	// Load list
	c.EXPECT().LoadPreviews(gomock.Any()).Return(nil, nil)

	m := NewChooseTypeModel(state)

	catwalk.RunModel(t, "testdata/choose_type/back_test", m)
}

func TestChooseTypeModel_Selects(t *testing.T) {
	tests := []struct {
		name     string
		testPath string
	}{
		{
			name:     "choose credentials",
			testPath: "testdata/choose_type/cred_test",
		},
		{
			name:     "choose text",
			testPath: "testdata/choose_type/text_test",
		},
		{
			name:     "choose card",
			testPath: "testdata/choose_type/card_test",
		},
		{
			name:     "choose upload file",
			testPath: "testdata/choose_type/file_test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewChooseTypeModel(&State{})

			catwalk.RunModel(t, tt.testPath, m)
		})
	}
}
