//go:generate bash -c "docker run -v $(pwd):/app -w /app golang:1.17"
//go:generate bash -c "protoc --go_out=service --proto_path=$(pwd) --go_opt=paths=source_relative --go-grpc_out=service --go-grpc_opt=paths=source_relative service.proto"
//go:generate bash -c "./generate.sh"
package grpcgoonch
