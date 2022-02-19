package grpcgoonch

import (
	"log"

	"golang.org/x/net/context"
)

type Server struct {
	ServiceServer
}

func (s *Server) CryptoRequest(ctx context.Context, input *Request) (*DecryptedText, error) {
	log.Printf("Received text from client: %s", input.Text)
	log.Printf("Received key from client: %s", input.Key)
	return &DecryptedText{Result: "Hello From the Goonch Server!"}, nil
}
