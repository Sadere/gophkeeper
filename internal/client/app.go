package client

import (
	"github.com/Sadere/gophkeeper/internal/client/api/grpc"
	"github.com/Sadere/gophkeeper/internal/client/config"
	"github.com/Sadere/gophkeeper/internal/client/tui"
)

type KeeperClient struct {
	Root *tui.RootModel
}

func NewKeeperClient(cfg *config.Config) (*KeeperClient, error) {
	gClient, err := grpc.NewGRPCClient(cfg)
	if err != nil {
		return nil, err
	}

	state := tui.NewState(gClient)

	root := tui.NewRootModel(state)

	return &KeeperClient{
		Root: root,
	}, nil
}
