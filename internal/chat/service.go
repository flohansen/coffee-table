package chat

import (
	"context"
	"fmt"
	"sync"

	"github.com/flohansen/coffee-table/pkg/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//go:generate protoc --proto_path=../../proto --go_out=../../pkg/proto --go_opt=paths=source_relative --go-grpc_out=../../pkg/proto --go-grpc_opt=paths=source_relative chat.proto

type Service struct {
	proto.UnimplementedChatBrokerServer
	streams map[string]proto.ChatBroker_ConnectServer
	mu      sync.RWMutex
}

func NewService() *Service {
	return &Service{
		streams: make(map[string]proto.ChatBroker_ConnectServer),
	}
}

func (s *Service) Connect(req *proto.ConnectRequest, stream proto.ChatBroker_ConnectServer) error {
	s.mu.Lock()
	s.streams[req.Username] = stream
	s.mu.Unlock()

	s.Broadcast(context.Background(), &proto.Message{
		Sender: "System",
		Text:   fmt.Sprintf("%s joined the chat", req.Username),
	})

	<-stream.Context().Done()

	s.Broadcast(context.Background(), &proto.Message{
		Sender: "System",
		Text:   fmt.Sprintf("%s left the chat", req.Username),
	})

	s.mu.Lock()
	delete(s.streams, req.Username)
	s.mu.Unlock()

	return nil
}

func (s *Service) Broadcast(ctx context.Context, msg *proto.Message) (*proto.BroadcastResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for username, conn := range s.streams {
		if username == msg.Sender {
			continue
		}

		if err := conn.Send(msg); err != nil {
			return &proto.BroadcastResponse{}, status.Error(codes.Internal, "error broadcasting message")
		}
	}

	return &proto.BroadcastResponse{}, nil
}
