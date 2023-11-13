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
	"google.golang.org/grpc/credentials/insecure"
)

type Peer struct {
	gRPC.UnimplementedModelServer

	name  string
	port  string
	mutex sync.Mutex
}

var port = ""
var peerIndex = 0
var neighbor gRPC.ModelClient
var neighborConn *grpc.ClientConn

func (p *Peer) GiveKey(ctx context.Context, key *gRPC.Key) (*gRPC.Ack, error) {
	log.Printf("Recieved key")
	ack := &gRPC.Ack{Status: "sucess!"}
	return ack, nil
}

func (p *Peer) ChangeNeighbor(ctx context.Context, in *gRPC.NeighborDetails) (*gRPC.Ack, error) {		
	print("Changing neighbor to: ", in.Port)
	set_neighbor(in.Port)
	ack := &gRPC.Ack{Status: "sucess in changing neighbor!"}
	_, err := neighbor.GiveKey(context.Background(), &gRPC.Key{})
	if err != nil {
		log.Fatal(err)
	}
	return ack, nil
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
	var list net.Listener
	port_i := 0

	//attemp to connect to port untill sucess
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
			peerIndex = port_i
		}

	}
	grpcServer := grpc.NewServer()
	peer := &Peer{
		name: "PeerX",
		port: port,
	}
	gRPC.RegisterModelServer(grpcServer, peer)
	log.Println("Peer upstart sucessfull on port ", port)

	if peerIndex > 0 {	
		set_neighbor(get_port(0))
		previous_port := get_port(peerIndex - 1)
		p_conn, err := connectToPeer(previous_port)
		if err != nil {
			log.Fatal(err)
		}
		own_details := &gRPC.NeighborDetails{
			Port: port,
		}
		println("attempting to change ", previous_port, "'s port to ", port)
		ack, _ := p_conn.ChangeNeighbor(context.Background(), own_details)
		print(ack.Status)
	}
	if err := grpcServer.Serve(list); err != nil {
		log.Fatalf("failed to serve %v", err)
	}
}

func set_neighbor(_port string) {
	_neighbor, err := connectToPeer(_port)
	if err != nil {
		print(err)
	}
	println("set neighbor to " + _port)
	neighbor = _neighbor
}

func connectToPeer(_port string) (gRPC.ModelClient, error) {
	if neighborConn != nil {
		println("Trying to close neighborConn")
		neighborConn.Close()
	}
	println("connect to peer entering")
	opts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	println("after opts")
	conn, err := grpc.Dial(":"+_port, opts...)
	neighborConn = conn
	println("after dial")
	if err != nil {
		log.Println("Error connecting to peer:", err)
		return nil, err
	} else {
		log.Println("Connected to neighbor at port: ", _port)
	}
	println("connect to peer exiting")

	return gRPC.NewModelClient(conn), nil
}

func main() {
	go launchPeer()
	for {
	}
}
