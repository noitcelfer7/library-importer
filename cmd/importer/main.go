package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	library_proto "github.com/noitcelfer7/library-proto/gen/go/proto/library"

	"library_importer/internal/importer/config"
	"library_importer/internal/importer/server"
)

func main() {
	data, err := os.ReadFile("config.json")

	if err != nil {
		panic(fmt.Sprintf("os.ReadFile Error: %v", err))
	}

	var config = config.Config{}

	err = json.Unmarshal(data, &config)

	if err != nil {
		panic(fmt.Sprintf("json.Unmarshal Error: %v", err))
	}

	target := net.JoinHostPort(config.Grpc.Client.Host, config.Grpc.Client.Port)

	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		panic(fmt.Sprintf("grpc.NewClient Error: %v", err))
	}

	defer conn.Close()

	cc := library_proto.NewDataExchangeServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	server.Serve(&config, cc, ctx)
}
