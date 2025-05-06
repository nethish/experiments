
# gRPC Streams
* Send won't wait for anything. As long as there is buffer, send just sends all of it.
* Limit number of streams with `grpc.MaxConcurrentStreams(X)` option
* Interceptors - Only one interceptors, chaining can be done at the caller
* Client deadline exceeded -> Server context cancelled

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

## TODO
* Currently the server and client communiate via tcp
* Establish the same with unix socket?
