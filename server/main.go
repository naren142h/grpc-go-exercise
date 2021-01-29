package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"

	pb "grpc-tutorial/proto"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	grpcEndpoint := fmt.Sprintf(":%s", port)

	grpcServer := grpc.NewServer()
	pb.RegisterCalculatorServiceServer(grpcServer, NewServer())

	listen, err := net.Listen("tcp", grpcEndpoint)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Starting gRPC listener [%s]\n", grpcEndpoint)
	log.Fatal(grpcServer.Serve(listen))

}
