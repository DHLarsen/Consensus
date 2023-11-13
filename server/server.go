package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"

	gRPC "github.com/DHLarsen/ChittyChat/proto"
	"google.golang.org/grpc"
)

type Peer struct {
	gRPC.UnimplementedModelServer

	name  string
	port  string
	mutex sync.Mutex
}

var port = "8889"


func (s *Server) SendMessage(msgStream gRPC.Model_SendMessageServer) error {
}
func (s *Server) GetUpdate(updateStream gRPC.Model_GetUpdateServer) error {
}

func launchPeer() {
	list, err := net.Listen("tcp", "localhost:"+port)
	if err != nil {
		log.Fatalf("failed to listen %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	peer := &Peer{
		name: "PeerX",
		port: port,
	}
	gRPC.RegisterModelServer(grpcServer, peer)
	log.Println("Peer upstart sucessfull")
	if err := grpcServer.Serve(list); err != nil {
		log.Fatalf("failed to serve %v", err)
	}
}

func main() {
	launchPeer()
}
