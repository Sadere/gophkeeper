package client

import (
	"github.com/Sadere/gophkeeper/internal/client/api/grpc"
	"github.com/Sadere/gophkeeper/internal/client/config"
	"github.com/Sadere/gophkeeper/internal/client/tui/screens"
)

// Main client app struct
type KeeperClient struct {
	Root   *screens.RootModel
	Client *grpc.GRPCClient
}

// Returns instance of client app
func NewKeeperClient(cfg *config.Config) (*KeeperClient, error) {
	gClient, err := grpc.NewGRPCClient(cfg)
	if err != nil {
		return nil, err
	}

	state := screens.NewState(gClient)

	root := screens.NewRootModel(state)

	return &KeeperClient{
		Root:   root,
		Client: gClient,
	}, nil
}
