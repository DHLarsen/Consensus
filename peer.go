package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
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

var port = ""
var neighbor gRPC.ModelClient

func (p *Peer) GiveKey(ctx context.Context, key *gRPC.Key) (*gRPC.Ack, error) {
	log.Printf("Recieved key")
	return nil, nil
}

func get_port(_i int) string {
	body, err := os.ReadFile("ports.txt")
	if err != nil {
		log.Fatalf("unable to read file: %v", err)
	}
	ports := strings.Split(string(body), "\n")

	// Use strings.TrimSpace to remove leading and trailing whitespaces
	return strings.TrimSpace(ports[_i])
}

func launchPeer() {
	log.Println("2")

	var list net.Listener
	port_i := 0
	log.Println("3")

	for port == "" {
		_port := get_port(port_i)
		log.Println(_port + "before")
		_port = strings.Trim(_port, " ")
		log.Println(_port + "after")
		_list, err := net.Listen("tcp", "localhost:"+strings.Trim(_port, " "))
		if err != nil {
			fmt.Println("failed to listen ", err)
			port_i++
		} else {
			port = _port
			list = _list
		}

	}
	grpcServer := grpc.NewServer()
	peer := &Peer{
		name: "PeerX",
		port: port,
	}
	gRPC.RegisterModelServer(grpcServer, peer)
	log.Println("Peer upstart sucessfull on port ", port)

	if port == "8909" {
		_neighbor, err := connectToPeer("8910")
		if err != nil {
			print(err)
		}
		neighbor = _neighbor
		ack, err := neighbor.GiveKey(context.Background(), &gRPC.Key{})
		println("ack: ", ack)
	}
	if err := grpcServer.Serve(list); err != nil {
		log.Fatalf("failed to serve %v", err)
	}
}

func connectToPeer(_port string) (gRPC.ModelClient, error) {
	opts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithInsecure(),
	}
	conn, err := grpc.Dial(":"+_port, opts...)
	if err != nil {
		log.Println("Error connecting to peer:", err)
		return nil, err
	} else {
		log.Println("Connected to neighbor at port: ", _port)
	}

	return gRPC.NewModelClient(conn), nil
}

func main() {
	log.Println("1")
	go launchPeer()
	for {
	}
}
