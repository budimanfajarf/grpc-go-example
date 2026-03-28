package main

import (
	"context"
	"flag"
	"log"
	"time"

	catalogpb "github.com/budimanfajarf/grpc-go-example/catalog"
	hellopb "github.com/budimanfajarf/grpc-go-example/helloworld"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	HelloAddr   = flag.String("hello_addr", "localhost:50051", "the address to connect to hello server")
	CatalogAddr = flag.String("catalog_addr", "localhost:50052", "the address to connect to catalog server")
)

func main() {
	flag.Parse()

	helloConn, err := grpc.NewClient(*HelloAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect to hello server: %v", err)
	}
	defer helloConn.Close()
	greeterClient := hellopb.NewGreeterClient(helloConn)

	catalogConn, err := grpc.NewClient(*CatalogAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect to catalog server: %v", err)
	}
	defer catalogConn.Close()
	catalogClient := catalogpb.NewCatalogClient(catalogConn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	helloReply, err := greeterClient.SayHello(ctx, &hellopb.HelloRequest{Name: "Budi"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", helloReply.GetMessage())

	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	store, err := catalogClient.GetStore(ctx, &catalogpb.GetStoreRequest{Uuid: "019d32d0-41cb-71a8-b71a-3d1089974b45"})
	if err != nil {
		log.Fatalf("could not get store: %v", err)
	}
	log.Printf("Store: %s", store.GetName())
}
