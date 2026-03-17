package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	pb "github.com/L1irik259/TestForOzon/internal/transport/proto/github.com/L1irik259/TestForOzon/transport/genetation/go/v1"
	transport "github.com/L1irik259/TestForOzon/internal/transport/service"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()
	server := &transport.Server{}
	pb.RegisterOzonServiceServer(grpcServer, server)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
