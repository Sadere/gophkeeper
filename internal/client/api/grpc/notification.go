package grpc

import (
	"context"
	"log"
	"time"

	"github.com/Sadere/gophkeeper/internal/client/tui/msg"
	"github.com/Sadere/gophkeeper/pkg/model"
	pb "github.com/Sadere/gophkeeper/pkg/proto/keeper/v1"
	tea "github.com/charmbracelet/bubbletea"
)

func (c *GRPCClient) Notifications(p *tea.Program) {
	var (
		stream pb.NotificationService_SubscribeV1Client
		err    error
		status model.SecretPreviewStatus
	)

	for {
		// Subscribe to notifications
		if stream == nil {
			if stream, err = c.subscribe(); err != nil {
				log.Printf("failed to subscribe: %v", err)
				c.sleep()
				continue
			}
		}

		response, err := stream.Recv()
		if err != nil {
			log.Printf("failed to recv msg: %v", err)
			stream = nil
			c.sleep()

			// Retry
			continue
		}

		// Process notification
		if response.Updated {
			status = model.SecretPreviewUpdated
		} else {
			status = model.SecretPreviewNew
		}

		c.previews.Store(response.Id, status)

		// Trigger secret list reload
		if p != nil {
			p.Send(msg.ReloadSecretList{})
		}
	}
}

func (c *GRPCClient) sleep() {
	time.Sleep(time.Second * 2)
}

func (c *GRPCClient) subscribe() (pb.NotificationService_SubscribeV1Client, error) {
	return c.notifyClient.SubscribeV1(context.Background(), &pb.SubscribeV1Request{
		Id: c.clientID,
	})
}
