package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	pb "github.com/naren142h/grpc-go-exercise/proto/calculator"
)

//var _ pb.CalculatorServiceServer = (*server)(nil)

// Server is a struct implements the pb.CalculatorServiceServer
type server struct {
	pb.UnimplementedCalculatorServiceServer
}

// NewServer returns a new Server
func newServer() *server {
	return &server{}
}

// Add adds two numbers
func (s *server) Add(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	a := req.GetA()
	b := req.GetB()
	res := a + b
	return &pb.Response{Result: res}, nil
}

//Multiply multiplies two numbers
func (s *server) Multiply(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	a := req.GetA()
	b := req.GetB()
	res := a * b
	return &pb.Response{Result: res}, nil
}

func main() {
	grpcEndpoint := fmt.Sprintf(":8090")

	grpcServer := grpc.NewServer()
	pb.RegisterCalculatorServiceServer(grpcServer, &server{})

	listen, err := net.Listen("tcp", grpcEndpoint)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Starting gRPC listener [%s]\n", grpcEndpoint)
	go func() {
		log.Fatal(grpcServer.Serve(listen))
	}()
	// Create a client connection to the gRPC server we just started
	// This is where the gRPC-Gateway proxies the requests
	conn, err := grpc.DialContext(
		context.Background(),
		grpcEndpoint,
		grpc.WithBlock(),
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	gwmux := runtime.NewServeMux()
	// Register Greeter
	err = pb.RegisterCalculatorServiceHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	httpEndpoint := fmt.Sprintf(":%s", port)

	gwServer := &http.Server{
		Addr:    httpEndpoint,
		Handler: gwmux,
	}

	log.Println("Serving gRPC-Gateway on http://0.0.0.0" + httpEndpoint)
	log.Fatalln(gwServer.ListenAndServe())

}
