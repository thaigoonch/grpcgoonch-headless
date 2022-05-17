package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	grpcgoonch "github.com/thaigoonch/grpcgoonch-headless/service"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/keepalive"
)

var (
	grpcPort     = 9000
	promPort     = 9092
	reg          = prometheus.NewRegistry()
	grpcMetrics  = grpc_prometheus.NewServerMetrics()
	grpcReqCount = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "grpcgoonchheadless_server_handle_count",
		Help: "Total number of RPCs handled on the goonch server.",
	})
)

type Server struct {
	grpcgoonch.ServiceServer
}

func init() {
	reg.MustRegister(grpcMetrics, grpcReqCount)
	_, err := reg.Gather()
	if err != nil {
		log.Fatalf("Prometheus metric registration error: %v", err)
	}
}

func (s *Server) CryptoRequest(ctx context.Context, input *grpcgoonch.Request) (*grpcgoonch.DecryptedText, error) {
	log.Printf("Received text from client: %s", input.Text)

	encrypted, err := encrypt(input.Key, input.Text)
	if err != nil {
		return &grpcgoonch.DecryptedText{Result: ""},
			fmt.Errorf("error during encryption: %v", err)
	}
	result, err := decrypt(input.Key, encrypted)
	if err != nil {
		return &grpcgoonch.DecryptedText{Result: ""},
			fmt.Errorf("error during decryption: %v", err)
	}

	grpcReqCount.Inc() // increment prometheus metric
	time.Sleep(4 * time.Second)
	return &grpcgoonch.DecryptedText{Result: result}, nil
}

func encrypt(key []byte, text string) (string, error) {
	plaintext := []byte(text)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func decrypt(key []byte, cryptoText string) (string, error) {
	ciphertext, _ := base64.URLEncoding.DecodeString(cryptoText)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return fmt.Sprintf("%v", ciphertext), nil
}

func main() {
	fmt.Println("grpcgoonch-headless waiting for client requests...")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		grpclog.Fatalf("Failed to listen on port %d: %v", grpcPort, err)
	}

	// Create an http server for prometheus
	httpServer := &http.Server{
		Handler: promhttp.HandlerFor(reg, promhttp.HandlerOpts{}),
		Addr:    fmt.Sprintf(":%d", promPort)}

	// Create a gRPC server
	s := Server{}
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpcMetrics.UnaryServerInterceptor()),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionAge: time.Second * 31,
		}),
	)
	grpcgoonch.RegisterServiceServer(grpcServer, &s)
	grpcMetrics.InitializeMetrics(grpcServer)

	// Start http server for prometheus
	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			log.Fatalf("Unable to start an http server on port %d: %v", promPort, err)
		}
	}()

	if err := grpcServer.Serve(lis); err != nil {
		grpclog.Fatalf("Failed to serve gRPC server over port %d: %v", grpcPort, err)
	}
}
