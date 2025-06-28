
## setup
```bash
apt update
apt install -y protobuf-compiler

go install google.golang.org/protobuf/cmd/protoc-gen-go
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc

export PATH="$PATH:$(go env GOPATH)/bin"
```

## generate grpc code

```bash
protoc --go_out=pkg/go-grpc/order --go_opt=paths=source_relative \
       --go-grpc_out=pkg/go-grpc/order --go-grpc_opt=paths=source_relative \
       protos/order.proto
```
