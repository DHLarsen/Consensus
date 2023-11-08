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

type Server struct {
	gRPC.UnimplementedModelServer

	name  string
	port  string
	mutex sync.Mutex
}

var port = "8889"

var vTime = []int64{0}
var vTimeIndex = 0

var updateChansMutex sync.Mutex
var updateChans = []chan *gRPC.Message{}

func (s *Server) SendMessage(msgStream gRPC.Model_SendMessageServer) error {
	clientName := ""
	for {
		msg, err := msgStream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			updateChansMutex.Lock()
			incrementVTime()
			updateChansMutex.Unlock()
			log.Println(clientName, " disconnected at vector time: ", vTime)
			broadcastMsg(&gRPC.Message{
				ClientName: "",
				Message:    fmt.Sprint(clientName, " disconnected at server time: ", vTime),
				Time:       vTime,
			})
			return err
		} else if clientName == "" {
			clientName = msg.ClientName
		}
		updateVTime(msg.Time)
		updateChansMutex.Lock()
		incrementVTime()
		updateChansMutex.Unlock()

		log.Printf("Received message from %s: %s vTime: %v", msg.ClientName, msg.Message, vTime)
		broadcastMsg(msg)
	}

	return nil
}
func (s *Server) GetUpdate(updateStream gRPC.Model_GetUpdateServer) error {
	updateChan := make(chan *gRPC.Message)
	updateChans = append(updateChans, updateChan)
	for {
		var msg = <-updateChan
		sendMessage(msg, updateStream)
	}
}

func broadcastMsg(msg *gRPC.Message) {
	for _, updateChan := range updateChans {
		updateChan <- msg
	}
}

func sendMessage(msg *gRPC.Message, updateStream gRPC.Model_GetUpdateServer) {
	updateChansMutex.Lock()
	incrementVTime()
	msg.Time = vTime
	log.Println("Sending: ", msg.Message, "vTime: ", vTime)
	updateStream.Send(msg)
	updateChansMutex.Unlock()
}

func updateVTime(newVTime []int64) {
	updateChansMutex.Lock()
	if newVTime != nil {
		for len(vTime) < len(newVTime) {
			vTime = append(vTime, 0)
		}
		for i, time := range newVTime {
			if time > vTime[i] {
				vTime[i] = time
			}
		}
	} else {
		vTime = append(vTime, 0)
	}
	log.Println("Set time to: ", vTime)
	updateChansMutex.Unlock()
}

func incrementVTime() {
	vTime[vTimeIndex]++
	log.Println("Set time to: ", vTime)
}

func launchServer() {
	list, err := net.Listen("tcp", "localhost:"+port)
	if err != nil {
		log.Fatalf("failed to listen %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	server := &Server{
		name: "Server",
		port: port,
	}
	gRPC.RegisterModelServer(grpcServer, server)
	log.Println("Server upstart sucessfull")
	if err := grpcServer.Serve(list); err != nil {
		log.Fatalf("failed to serve %v", err)
	}
}

func main() {
	launchServer()
}
