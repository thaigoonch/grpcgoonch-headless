package main

import (
	"fmt"
	"net"
	"time"

	grpcgoonch "github.com/thaigoonch/grpcgoonch/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/keepalive"
)

func main() {
	fmt.Println("grpcgoonch waiting for client requests...")
	port := 9000
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		grpclog.Fatalf("Failed to listen on port %d: %v", port, err)
	}

	s := grpcgoonch.Server{}

	grpcServer := grpc.NewServer(grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionAge: time.Minute * 6,
	}))
	grpcgoonch.RegisterServiceServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		grpclog.Fatalf("Failed to serve gRPC server over port %d: %v", port, err)
	}
}
