package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	grpcgoonch "github.com/thaigoonch/grpcgoonch/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/keepalive"
)

var (
	reg                 = prometheus.NewRegistry()
	grpcMetrics         = grpc_prometheus.NewServerMetrics()
	customMetricCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "grpcgoonch_server_handle_count",
		Help: "Total number of RPCs handled on the goonch server.",
	}, []string{"name"})
)

func init() {
	reg.MustRegister(grpcMetrics, customMetricCounter)
}

func main() {
	fmt.Println("grpcgoonch waiting for client requests...")
	port := 9000
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		grpclog.Fatalf("Failed to listen on port %d: %v", port, err)
	}

	// Create an http server for prometheus
	httpServer := &http.Server{
		Handler: promhttp.HandlerFor(reg, promhttp.HandlerOpts{}),
		Addr:    fmt.Sprintf("0.0.0.0:%d", port)}

	// Create a gRPC server
	s := grpcgoonch.Server{}
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpcMetrics.UnaryServerInterceptor()),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionAge: time.Minute * 6,
		}))
	grpcgoonch.RegisterServiceServer(grpcServer, &s)
	grpcMetrics.InitializeMetrics(grpcServer)

	// Start http server for prometheus
	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			log.Fatal("Unable to start an http server.")
		}
	}()

	if err := grpcServer.Serve(lis); err != nil {
		grpclog.Fatalf("Failed to serve gRPC server over port %d: %v", port, err)
	}
}
