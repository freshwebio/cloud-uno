# Contributing help and pointers

## Generating code from proto files

We use [protocol buffer files (*.proto)](https://developers.google.com/protocol-buffers) to define the gRPC communication between the hosts agent and the server.

Enter in to the `pkg/hosts` directory containing the .proto file we want to generate Go code for in your terminal.

Then run the following command:
```bash
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative hosts.proto
```