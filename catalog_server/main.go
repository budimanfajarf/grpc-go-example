package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/budimanfajarf/grpc-go-example/catalog"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50052, "The server port")
)

type server struct {
	catalog.UnimplementedCatalogServer
}

// GetStore implements catalog.CatalogServer
func (s *server) GetStore(_ context.Context, in *catalog.GetStoreRequest) (*catalog.Store, error) {
	log.Printf("Received: %v", in.GetUuid())
	return &catalog.Store{Uuid: in.GetUuid(), Name: "Example Store"}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	catalog.RegisterCatalogServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
