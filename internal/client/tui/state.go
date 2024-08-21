package tui

import (
	"time"

	"github.com/Sadere/gophkeeper/internal/client/api"
)

type State struct {
	client      api.IApiClient
	accessToken string
	startAt     time.Time
}

func NewState(cl api.IApiClient) *State {
	return &State{
		client:  cl,
		startAt: time.Now(),
	}
}
