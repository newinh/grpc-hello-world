package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	pb "github.com/newinh/grpc-hello-world/proto/gen/v1"
)

func main() {
	grpcServer := grpc.NewServer()
	healthChecker := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthChecker)

	reflection.Register(grpcServer)
	pb.RegisterHelloServiceServer(grpcServer, &helloServer{})

	port := 50053
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}

	go func() {
		fmt.Printf("gRPC server is running on port %d\n", port)
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
