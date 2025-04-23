package main

import (
	"context"
	"fmt"
	"log"
	"time"

	library_proto "github.com/noitcelfer7/library-proto/gen/go/proto/library"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("127.0.0.1:12345", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		fmt.Printf("NewClient Error, %v", err)

		return
	}

	defer conn.Close()

	client := library_proto.NewDataExchangeServiceClient(conn)

	request := &library_proto.ExchangeRequest{
		AuthorFirstName: "Kee",
	}

	ctx, _ := context.WithTimeout(context.Background(), time.Second)

	response, _ := client.Exchange(ctx, request)

	log.Printf("%v", response)
}
