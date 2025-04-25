
# gRPC Streams
* Send won't wait for anything. As long as there is buffer, send just sends all of it.

Generate server code
```
protoc \
  --proto_path=protos \
  --go_out=gen \
  --go-grpc_out=gen \
  protos/api.proto
```

Bring up the server
```
go run ./server/main.go
```

Run the client
```
go run ./client/client.go
```
