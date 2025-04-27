package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

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

	caCert, err := os.ReadFile("ca-cert.pem")
	if err != nil {
		log.Fatal(err)
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(caCert) {
		log.Fatal("Failed to add CA certificate")
	}

		// Конфигурация TLS
		tlsConfig := &tls.Config{
			RootCAs: certPool,
		}

	creds := credentials.NewTLS(tlsConfig)


	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(creds))

	if err != nil {
		panic(fmt.Sprintf("grpc.NewClient Error: %v", err))
	}

	defer conn.Close()

	cc := library_proto.NewDataExchangeServiceClient(conn)

	server.Serve(&config, cc)
}
