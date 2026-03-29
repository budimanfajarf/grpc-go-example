package main

import (
	"context"
	"flag"
	"io"
	"log"
	"time"

	"github.com/budimanfajarf/grpc-go-example/catalog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
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
		if status.Code(err) == codes.NotFound {
			log.Printf("Store not found")
		} else {
			log.Fatalf("unexpected error: %v", err)
		}
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
	products := productsResponse.GetData()
	for _, product := range products {
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

	productNames := map[string]string{}
	reserveStocksRequest := &catalog.ReserveStocksRequest{}
	for i, product := range products {
		reserveStock := &catalog.ReserveStock{
			ProductUuid: product.GetUuid(),
			Quantity:    int32(i + 1),
		}
		productNames[product.GetUuid()] = product.GetName()
		reserveStocksRequest.Data = append(reserveStocksRequest.Data, reserveStock)
	}

	reserveCtx, reserveCancel := context.WithTimeout(context.Background(), time.Second)
	defer reserveCancel()
	reserveStocksResponse, err := c.ReserveStocks(reserveCtx, reserveStocksRequest)
	if err != nil {
		log.Fatalf("could not reserve stocks: %v", err)
	}
	reserved := reserveStocksResponse.GetReserved()
	if !reserved {
		for _, item := range reserveStocksResponse.GetData() {
			status := item.GetStatus()
			if status == catalog.ReserveStatus_RESERVE_STATUS_SUCCESS {
				continue
			}

			switch status {
			case catalog.ReserveStatus_RESERVE_STATUS_INSUFFICIENT_STOCK:
				log.Fatalf("Insufficient stock for product: %s", productNames[item.GetProductUuid()])
			case catalog.ReserveStatus_RESERVE_STATUS_NOT_FOUND:
				log.Fatalf("Product not found: %s", productNames[item.GetProductUuid()])
			default:
				log.Fatalf("Unknown reserve status for product: %s, status: %v", productNames[item.GetProductUuid()], status)
			}
		}
	}
	log.Println("Stocks reserved successfully")
}
