package screens

import (
	"time"

	"github.com/Sadere/gophkeeper/internal/client/api"
)

type State struct {
	client      api.IApiClient
	accessToken string
	startAt     time.Time

	windowHeight int
	windowWidth  int
}

func NewState(cl api.IApiClient) *State {
	return &State{
		client:  cl,
		startAt: time.Now(),
	}
}

func (s *State) SetSize(w, h int) {
	s.windowHeight = h
	s.windowWidth = w
}

func (s *State) Height() int {
	return s.windowHeight
}

func (s *State) Width() int {
	return s.windowWidth
}
