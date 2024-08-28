package grpc

import (
	"errors"
	"slices"

	"github.com/Sadere/gophkeeper/internal/server/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/Sadere/gophkeeper/pkg/proto/keeper/v1"
)

type sub struct {
	stream   pb.NotificationService_SubscribeV1Server
	id       int32
	finished chan<- bool
}

func (s *KeeperServer) SubscribeV1(in *pb.SubscribeV1Request, stream pb.NotificationService_SubscribeV1Server) error {
	var subs []sub

	ctx := stream.Context()

	userID, err := extractUserID(ctx)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	s.log.Infoln("received subscribe from client #", in.Id, "user ID", userID)

	fin := make(chan bool)

	v, ok := s.subscribers.Load(userID)
	if ok {
		// Try to cast value into sub slice
		subs, ok = v.([]sub)
		if !ok {
			return status.Error(codes.Internal, "failed to cast subscribers")
		}
	}

	// Append to subscribers slice
	subs = append(
		subs,
		sub{
			stream:   stream,
			id:       in.Id,
			finished: fin,
		},
	)

	// Store subs in map
	s.subscribers.Store(userID, subs)

	for {
		select {
		case <-fin:
			s.log.Infof("closing stream for client #%d", in.Id)
			return nil
		case <-ctx.Done():
			s.log.Infof("client #%d has disconnected", in.Id)
			return nil
		}
	}
}

func (s *KeeperServer) notifyClients(userID uint64, clientID int32, ID uint64, updated bool) error {
	v, ok := s.subscribers.Load(userID)
	if !ok {
		return model.ErrNoSubscribers
	}

	subs, ok := v.([]sub)
	if !ok {
		return errors.New("failed to cast to subs")
	}

	var unsubs []int

	for i, sub := range subs {
		if sub.id == clientID {
			// Skip originating client
			continue
		}

		resp := &pb.SubscribeV1Response{
			Id:      ID,
			Updated: updated,
		}

		if err := sub.stream.Send(resp); err != nil {
			s.log.Error("failed to send notification to client: %v", err)
			select {
			case sub.finished <- true:
				s.log.Info("client unsubscribed: %v", clientID)
			default:
			}

			// Mark as unsub
			unsubs = append(unsubs, i)
		}
	}

	// Delete unsubs from slice
	for _, unsub := range unsubs {
		subs = slices.Delete(subs, unsub, unsub)
	}

	if len(subs) > 0 {
		s.subscribers.Store(userID, subs)
	} else {
		// All clients unsubscribed
		s.subscribers.Delete(userID)
	}

	return nil
}
