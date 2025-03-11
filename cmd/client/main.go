package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/flohansen/coffee-table/pkg/proto"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	username := os.Args[1]

	client, err := grpc.NewClient(":5050", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not create client: %s", err)
	}

	c := proto.NewChatBrokerClient(client)
	conn, err := c.Connect(context.Background(), &proto.ConnectRequest{Username: username})
	if err != nil {
		log.Fatalf("could not connect to chat server: %s", err)
	}

	go func() {
		for {
			msg, err := conn.Recv()
			if err != nil {
				log.Fatalf("error receiving message: %s", err)
			}

			fmt.Printf("[%s]: %s\n", msg.Sender, msg.Text)
		}
	}()

	for {
		time.Sleep(time.Second)
		c.Broadcast(context.Background(), &proto.Message{
			Sender: username,
			Text:   uuid.New().String(),
		})
	}
}
