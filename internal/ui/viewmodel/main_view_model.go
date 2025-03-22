package viewmodel

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"time"

	"github.com/flohansen/coffee-table/pkg/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type View int

const (
	ViewLogin View = iota
	ViewChat
)

type MainViewModel struct {
	Error          *Observer[string]
	CurrentMessage *Observer[string]
	Users          *Observer[[]*proto.User]
	Message        *Observer[string]
	CurrentView    *Observer[View]
	Secure         *Observer[bool]
	username       string
	serverURL      string
	client         proto.ChatBrokerClient
	messageStream  proto.ChatBroker_ConnectClient
}

func NewMainViewModel() *MainViewModel {
	return &MainViewModel{
		Error:          NewObserver(""),
		Message:        NewObserver(""),
		CurrentView:    NewObserver(ViewLogin),
		Users:          NewObserver([]*proto.User{}),
		CurrentMessage: NewObserver(""),
		Secure:         NewObserver(false),
	}
}

func (v *MainViewModel) UpdateUsername(username string) {
	v.username = username
}

func (v *MainViewModel) UpdateServerURL(url string) {
	v.serverURL = url
}

func (v *MainViewModel) Connect() {
	var clientOptions []grpc.DialOption

	if v.Secure.Get() {
		caPem, err := os.ReadFile("ca-cert.pem")
		if err != nil {
			v.Error.Set(err.Error())
			return
		}

		certPool := x509.NewCertPool()
		if !certPool.AppendCertsFromPEM(caPem) {
			v.Error.Set("could not add CA's certificate")
			return
		}

		clientOptions = append(clientOptions, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
			RootCAs: certPool,
		})))
	} else {
		clientOptions = append(clientOptions, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	client, err := grpc.NewClient(v.serverURL, clientOptions...)
	if err != nil {
		v.Error.Set(err.Error())
		return
	}

	v.client = proto.NewChatBrokerClient(client)
	v.messageStream, err = v.client.Connect(context.Background(), &proto.ConnectRequest{
		Username: v.username,
	})
	if err != nil {
		v.Error.Set(err.Error())
		return
	}

	v.CurrentView.Set(ViewChat)
	go v.ReceiveMessages()
}

func (v *MainViewModel) SendMessage(message string) {
	v.client.Broadcast(context.Background(), &proto.Message{
		Sender: v.username,
		Text:   message,
	})

	v.Message.Set("")
}

func (v *MainViewModel) ReceiveMessages() {
	for {
		msg, err := v.messageStream.Recv()
		if err != nil {
			panic(err)
		}

		userStream, err := v.client.GetUsers(context.Background(), &proto.GetUsersRequest{})
		if err != nil {
			panic(err)
		}

		var users []*proto.User
		for {
			user, err := userStream.Recv()
			if err != nil {
				break
			}

			users = append(users, user)
		}
		v.Users.Set(users)

		textColor := "white"
		color := "green"

		switch msg.Sender {
		case v.username:
			textColor = "white"
			color = "blue"
		case "System":
			textColor = "gray"
			color = "gray"
		}

		currentMessage := fmt.Sprintf("[%s]%s %s:[%s] %s\n", color, msg.TimeSent.AsTime().Format(time.DateTime), msg.Sender, textColor, msg.Text)
		v.CurrentMessage.Set(currentMessage)
	}
}
