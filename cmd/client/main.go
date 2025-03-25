package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/newinh/grpc-hello-world/proto/gen/v1"
)

var (
	SECURE bool = false
)

func main() {
	opts := make([]grpc.DialOption, 0, 4)
	if SECURE {
		roots, err := x509.SystemCertPool()
		if err != nil {
			panic(err)
		}
		creds := credentials.NewTLS(&tls.Config{
			RootCAs: roots,
		})
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	conn, err := grpc.NewClient("localhost:50054", opts...)
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	helloServiceClient := pb.NewHelloServiceClient(conn)

	ctx := context.Background()
	res, err := helloServiceClient.SayHello(ctx, &pb.HelloRequest{Name: "World"})
	if err != nil {
		panic(err)
	}
	fmt.Println(res.Message)
}
