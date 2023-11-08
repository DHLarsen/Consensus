package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"

	// This has to be the same as the go.mod module,
	// followed by the path to the folder the proto file is in.
	gRPC "github.com/DHLarsen/ChittyChat/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var port = "8889"

var vTime []int64
var vTimeIndex int

var client gRPC.ModelClient
var ServerConn *grpc.ClientConn

var clientName string

var timeMutex sync.Mutex

func ConnectToServer() {
	opts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	conn, err := grpc.Dial(":"+port, opts...)
	if err != nil {
		print(err)
	}

	client = gRPC.NewModelClient(conn)
	ServerConn = conn

	log.Println("the connection is: ", conn.GetState().String())

	go updateListen()
}

func updateListen() {
	updateRequest := gRPC.UpdateRequest{ClientName: clientName}

	stream, err := client.GetUpdate(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	stream.Send(&updateRequest)

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			//panic(err)
		}
		if resp != nil {
			updateVTime(resp.Time)
			if vTimeIndex == 0 {
				vTimeIndex = len(vTime) - 1
				fmt.Println("Set time index to: ", vTimeIndex)
			}
			if vTimeIndex != 0 {
				timeMutex.Lock()
				vTime[vTimeIndex]++
				timeMutex.Unlock()
			}
			printOutput(resp)
		}
	}
}

func updateVTime(newVTime []int64) {
	for len(vTime) < len(newVTime) {
		vTime = append(vTime, 0)
	}
	for i, time := range newVTime {
		if time > vTime[i] {
			vTime[i] = time
		}
	}
}

func main() {
	startup()
	ConnectToServer()

	defer ServerConn.Close()
	reader := bufio.NewReader(os.Stdin)

	stream, err := client.SendMessage(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	stream.Send(&gRPC.Message{ClientName: clientName, Message: "Client " + clientName + " has joined the room!", Time: vTime})

	//Infinite loop to listen for clients input.
	for {
		prepareInput()

		//Read input into var input and any errors into err
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		input = strings.TrimSpace(input) //Trim input

		timeMutex.Lock()
		vTime[vTimeIndex]++
		timeMutex.Unlock()
		if len(input) > 128 {
			log.Println("Message too long, max 128 characters")
		} else {
			stream.Send(&gRPC.Message{ClientName: clientName, Message: input, Time: vTime})
		}
		continue
	}
}

func prepareInput() {
	fmt.Print("-> ")
}

func printOutput(msg *gRPC.Message) {
	name := msg.ClientName
	if name == clientName {
		name = "You"
	}
	fmt.Println(name, ": ", msg.Message, " | time: ", vTime)
	prepareInput()
}

func startup() {
	fmt.Print("Please enter your name: ")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
	clientName = strings.TrimSpace(input.Text())
}
