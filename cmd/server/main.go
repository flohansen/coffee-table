package main

import (
	"log"
	"net"

	"github.com/flohansen/coffee-table/internal/chat"
	"github.com/flohansen/coffee-table/pkg/proto"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":5050")
	if err != nil {
		log.Fatalf("could not create tcp listener: %s", err)
	}

	server := grpc.NewServer()
	proto.RegisterChatBrokerServer(server, chat.NewService())
	if err := server.Serve(lis); err != nil {
		log.Fatalf("server error: %s", err)
	}
}
