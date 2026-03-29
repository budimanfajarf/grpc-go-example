package main

import (
	"context"
	"flag"
	"io"
	"log"
	"time"

	"github.com/budimanfajarf/grpc-go-example/catalog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr = flag.String("addr", "localhost:50052", "the address to connect to")
)

func main() {
	flag.Parse()
	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := catalog.NewCatalogClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	storeResponse, err := c.GetStore(ctx, &catalog.GetStoreRequest{Uuid: "019d32d0-41cb-71a8-b71a-3d1089974b45"})
	if err != nil {
		log.Fatalf("could not get store: %v", err)
	}
	store := storeResponse.GetData()
	log.Printf("Store: %s", store.GetName())

	productUuids := []string{
		"019d32d0-41cb-71a8-b71a-3d1089974b45",
		"019d32d0-41cb-71a8-b71a-3d1089974b46",
		"019d32d0-41cb-71a8-b71a-3d1089974b47",
	}

	productsResponse, err := c.ListProducts(ctx, &catalog.ListProductsRequest{Uuids: productUuids})
	if err != nil {
		log.Fatalf("could not list products: %v", err)
	}
	for _, product := range productsResponse.GetData() {
		log.Printf("Product: %s", product.GetName())
	}

	streamCtx, streamCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer streamCancel()
	streamProducts, err := c.StreamProducts(streamCtx, &catalog.StreamProductsRequest{Uuids: productUuids})
	if err != nil {
		log.Fatalf("StreamProducts failed: %v", err)
	}
	log.Printf("StreamProducts started")
	for {
		product, err := streamProducts.Recv()
		if err == io.EOF {
			log.Printf("StreamProducts ended")
			break
		}
		if err != nil {
			log.Fatalf("StreamProducts failed: %v", err)
		}
		log.Printf("Stream Product: %s", product.GetName())
	}
}
