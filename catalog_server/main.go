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

func (s *server) StreamProducts(in *catalog.StreamProductsRequest, stream catalog.Catalog_StreamProductsServer) error {
	log.Printf("Received: %v", in.GetUuids())
	for i, uuid := range in.GetUuids() {
		product := &catalog.Product{Uuid: uuid, Name: fmt.Sprintf("Example Product %d", i+1)}
		log.Printf("Sending product: %v", product.GetUuid())
		if err := stream.Send(product); err != nil {
			return err
		}
	}
	return nil
}

func (s *server) ListProducts(_ context.Context, in *catalog.ListProductsRequest) (*catalog.ListProductsResponse, error) {
	log.Printf("Received: %v", in.GetUuids())
	data := []*catalog.Product{}
	for i, uuid := range in.GetUuids() {
		product := &catalog.Product{Uuid: uuid, Name: fmt.Sprintf("Example Product %d", i+1)}
		data = append(data, product)
	}
	return &catalog.ListProductsResponse{Data: data}, nil
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
