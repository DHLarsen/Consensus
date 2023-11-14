# Mutual exclusion example

Made by tpep, dhla and habr

This project was developed with our ChittyChat project as a baseline

# Starting a peer:
open a new terminal. Run the command: "go run peer.go"

this can be repeated for up to 10 peers

to change what ports are being used, or increase the amount of peers able to run, please edit the "ports.txt" file

# Usage
To write a message to the critical section simply input it in the console while the peer is running in it

Do this command to update the protoc generated files:

protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/Model.proto
