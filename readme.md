Made by tpep, dhla and habr


repeat these steps for every node(peer):

Starting a peer:
open a new terminal. Run the command: "go run peer.go"

To write a message to the critical section simply input it in the console while the peer is running in it

If your ports are occupied, please change the ports in the ports.txt file.

protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/Model.proto
