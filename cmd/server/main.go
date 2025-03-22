package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net"

	"github.com/flohansen/coffee-table/internal/chat"
	"github.com/flohansen/coffee-table/pkg/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Flags struct {
	UseTLS   bool
	CertFile string
	KeyFile  string
}

func main() {
	var flags Flags
	flag.BoolVar(&flags.UseTLS, "tls", false, "If TLS should be enabled to secure chats")
	flag.StringVar(&flags.CertFile, "cert", "cert.pem", "The file path of the certificate")
	flag.StringVar(&flags.KeyFile, "key", "key.pem", "The file path of the key")
	flag.Parse()

	var serverOpts []grpc.ServerOption
	if flags.UseTLS {
		cert, err := tls.LoadX509KeyPair(flags.CertFile, flags.KeyFile)
		if err != nil {
			log.Fatalf("could not load certificate: %s", err)
		}

		serverOpts = append(serverOpts, grpc.Creds(credentials.NewTLS(&tls.Config{
			Certificates: []tls.Certificate{cert},
			ClientAuth:   tls.NoClientCert,
		})))
	}

	lis, err := net.Listen("tcp", ":5050")
	if err != nil {
		log.Fatalf("could not create tcp listener: %s", err)
	}

	server := grpc.NewServer(serverOpts...)
	proto.RegisterChatBrokerServer(server, chat.NewService())
	if err := server.Serve(lis); err != nil {
		log.Fatalf("server error: %s", err)
	}
}
