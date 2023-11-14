Made by tpep, dhla and habr

How to run:
Starting the server:
From a terminal navigate to the server folder and run the command: "go run server.go"
Or from project root directory run "go run server/server.go"

repeat these steps for every node(peer):

Starting the client:
open a new terminal. Run the command: "go run peer.go"
client.go

Send messages:
Type the message and press enter.

If your ports are occupied, please change the ports in the ports.txt file.

protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/Model.proto
