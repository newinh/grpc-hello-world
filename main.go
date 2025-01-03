package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	pb "github.com/newinh/grpc-hello-world/proto/gen/v1"
)

func main() {
	grpcServer := grpc.NewServer()

	grpcLogger := grpclog.NewLoggerV2(os.Stdout, os.Stdout, os.Stderr)
	grpclog.SetLoggerV2(grpcLogger)

	reflection.Register(grpcServer)
	grpc_health_v1.RegisterHealthServer(grpcServer, health.NewServer())
	pb.RegisterHelloServiceServer(grpcServer, &helloServer{})

	port := os.Getenv("PORT")
	if port == "" {
		port = "50054"
	}
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		panic(err)
	}

	go func() {
		fmt.Printf("gRPC server is running on port %s\n", port)
		if err := grpcServer.Serve(lis); err != nil {
			panic(err)
		}
	}()
	defer grpcServer.GracefulStop()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
}

var (
	_ pb.HelloServiceServer = (*helloServer)(nil)
)

type helloServer struct {
	pb.UnimplementedHelloServiceServer
}

func (s *helloServer) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{Message: "Hello " + req.Name}, nil
}
