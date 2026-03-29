package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

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
func (s *server) GetStore(_ context.Context, in *catalog.GetStoreRequest) (*catalog.GetStoreResponse, error) {
	uuid := in.GetUuid()
	log.Printf("Received: %v", uuid)

	if uuid == "" {
		return nil, fmt.Errorf("uuid required")
	}

	// return nil, nil // simulate not found

	return &catalog.GetStoreResponse{
		Data: &catalog.Store{
			Uuid: uuid,
			Name: "Example Store",
		},
	}, nil
}

func (s *server) ListProducts(_ context.Context, in *catalog.ListProductsRequest) (*catalog.ListProductsResponse, error) {
	uuids := in.GetUuids()
	if len(uuids) == 0 {
		return nil, fmt.Errorf("uuids required")
	}
	log.Printf("Received: %v", uuids)

	data := []*catalog.Product{}
	for i, uuid := range uuids {
		product := &catalog.Product{Uuid: uuid, Name: fmt.Sprintf("Example Product %d", i+1), Price: int64(10000 + (i * 2000))}
		data = append(data, product)
	}
	return &catalog.ListProductsResponse{Data: data}, nil
}

func (s *server) StreamProducts(in *catalog.StreamProductsRequest, stream catalog.Catalog_StreamProductsServer) error {
	uuids := in.GetUuids()
	if len(uuids) == 0 {
		return fmt.Errorf("uuids required")
	}
	log.Printf("Received: %v", uuids)

	for i, uuid := range uuids {
		product := &catalog.Product{Uuid: uuid, Name: fmt.Sprintf("Example Product %d", i+1), Price: int64(10000 + (i * 2000))}
		log.Printf("Sending product: %v", product.GetUuid())
		// if i == 1 {
		// 	log.Printf("Simulating error for product: %v", product.GetUuid())
		// 	return fmt.Errorf("simulated error for product: %v", product.GetUuid())
		// }
		time.Sleep(500 * time.Millisecond) // add a small delay to simulate processing time
		if err := stream.Send(product); err != nil {
			return err
		}
	}
	return nil
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
