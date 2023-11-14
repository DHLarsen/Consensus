package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	gRPC "github.com/DHLarsen/ChittyChat/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Peer struct {
	gRPC.UnimplementedModelServer

	name string
	port string
}

var port = ""
var peerIndex = 0
var neighbor gRPC.ModelClient

var hasKey = false
var wantsKey = false
var writeString = ""
var writeStringMutex sync.Mutex

var acessChan = make(chan int)
var acessAckChan = make(chan bool)

func (p *Peer) GiveKey(ctx context.Context, key *gRPC.Key) (*gRPC.Ack, error) {
	//uncomment this line to print every time the peer recieves the key
	//log.Printf("Recieved key")
	ack := &gRPC.Ack{Status: "Give key sucess!"}
	hasKey = true
	return ack, nil
}

func sendKey() {
	for {
		if hasKey && neighbor != nil {
			if wantsKey {
				//send wait time to acessChan, which starts acess in critical section
				acessChan <- 1000
				//waits for acessAck
				<-acessAckChan
			}
			hasKey = false
			key := &gRPC.Key{
				Status: "key",
			}
			neighbor.GiveKey(context.Background(), key)
		}
		time.Sleep(1 * time.Second)
	}
}

func criticalAcess() {
	for {
		a := <-acessChan
		time.Sleep(time.Duration(a) * time.Millisecond)
		file, err := os.OpenFile("CriticalSection.txt", os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			log.Fatal(err)
		}

		_, werr := file.WriteString(writeString)
		if werr != nil {
			log.Fatalf("failed writing to file: %s", err)
		}

		file.Close()
		writeString = ""
		wantsKey = false
		acessAckChan <- true
	}
}

func (p *Peer) ChangeNeighbor(ctx context.Context, in *gRPC.NeighborDetails) (*gRPC.Ack, error) {
	ack := &gRPC.Ack{Status: in.Port + " is changing its neighbor to you!"}
	go set_neighbor(in.Port)
	return ack, nil
}

/*
// Is called with go ChangeNeighborSeperateThread.
func ChangeNeighborSeperateThread(ctx context.Context, in *gRPC.NeighborDetails) {
	set_neighbor(in.Port)
	key, err := neighbor.GiveKey(context.Background(), &gRPC.Key{})
	println(key.Status)
	if err != nil {
		log.Fatal(err)
	}
}*/

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

	//attemp to connect to port until sucess
	for port == "" {
		_port := get_port(port_i)
		_port = strings.Trim(_port, " ")
		_list, err := net.Listen("tcp", "localhost:"+strings.Trim(_port, " "))
		if err != nil {
			//fmt.Println("failed to listen ", err)
			port_i++
		} else {
			port = _port
			list = _list
			peerIndex = port_i
		}

	}
	grpcServer := grpc.NewServer()
	peer := &Peer{
		name: "Peer" + string(peerIndex),
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
		log.Println("attempting to change ", previous_port, "'s neighbor to port ", port)
		ack, _ := p_conn.ChangeNeighbor(context.Background(), own_details)
		log.Println(ack.Status)
	} else {
		hasKey = true
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
	neighbor = _neighbor
}

func connectToPeer(_port string) (gRPC.ModelClient, error) {
	opts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
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
	go launchPeer()
	go sendKey()
	go criticalAcess()
	c := 0
	reader := bufio.NewReader(os.Stdin)
	for {
		//Read input into var input and any errors into err
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		input = strings.TrimSpace(input) //Trim input

		writeStringMutex.Lock()
		writeString += fmt.Sprint(port) + ": " + input + "\n"
		writeStringMutex.Unlock()
		wantsKey = true
		c++
	}
}
