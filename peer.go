package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"

	gRPC "github.com/DHLarsen/ChittyChat/proto"
	"google.golang.org/genproto/googleapis/rpc/context"
	"google.golang.org/grpc"
)

type Peer struct {
	gRPC.UnimplementedModelServer

	name  string
	port  string
	mutex sync.Mutex
}

var port = ""

func (p *Peer) giveKey(ctx context.Context, key *gRPC.Key) (*gRPC.Ack, error){

ack := &gRPC.Ack{
	status: "1"
	}
return ack, nil
}




func get_port(_i int) string {
	body, err := os.ReadFile("ports.txt")
	if err != nil {
		log.Fatalf("unable to read file: %v", err)
	}
	ports := strings.Split(string(body), "\n")
	return ports[_i]
}

func launchPeer() {
	var list net.Listener
	port_i := 0
	for port == "" {
		_port := get_port(port_i)
		_list, err := net.Listen("tcp", "localhost:"+strings.Trim(_port, " "))
		if err != nil {
			fmt.Println("failed to listen ", err)
			port_i++
		} else {
			port = _port
			list = _list
		}
		
	}	
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	peer := &Peer{
		name: "PeerX",
		port: port,
	}
	gRPC.RegisterModelServer(grpcServer, peer)
	log.Println("Peer upstart sucessfull on port ", port)
	if err := grpcServer.Serve(list); err != nil {
		log.Fatalf("failed to serve %v", err)
	}
}

func main() {
	launchPeer()
}

