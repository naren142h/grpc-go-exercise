package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/open-policy-agent/opa/rego"
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

	module := getRego()

	query, err := rego.New(
		rego.Query("y = data.example.authz.allow"),
		rego.Module("example.rego", module),
	).PrepareForEval(ctx)
	if err != nil {
		// Handle error.
	}

	input := map[string]interface{}{
		"method":   "POST",
		"function": "Add",
		"user":     req.GetUser(),
	}

	results, err := query.Eval(ctx, rego.EvalInput(input))

	if err != nil {
		// Handle evaluation error.
	} else if len(results) == 0 {
		// Handle undefined result.
	} else if result, ok := results[0].Bindings["y"].(bool); !ok {
		// Handle unexpected result type.
		fmt.Printf("Unexpected result: " + strconv.FormatBool(result))
	}

	// Handle result/decision.
	if result, _ := results[0].Bindings["y"].(bool); result {
		a := req.GetA()
		b := req.GetB()
		res := a + b
		return &pb.Response{Result: res}, nil
	}

	return &pb.Response{}, errors.New("Access denied! You are not allowed to Add!")

}

//Multiply multiplies two numbers
func (s *server) Multiply(ctx context.Context, req *pb.Request) (*pb.Response, error) {

	module := getRego()

	query, err := rego.New(
		rego.Query("y = data.example.authz.allow"),
		rego.Module("example.rego", module),
	).PrepareForEval(ctx)
	if err != nil {
		// Handle error.
	}

	input := map[string]interface{}{
		"method":   "POST",
		"function": "Multiply",
		"user":     req.GetUser(),
	}

	results, err := query.Eval(ctx, rego.EvalInput(input))

	if err != nil {
		// Handle evaluation error.
	} else if len(results) == 0 {
		// Handle undefined result.
	} else if result, ok := results[0].Bindings["y"].(bool); !ok {
		// Handle unexpected result type.
		fmt.Printf("Unexpected result: " + strconv.FormatBool(result))
	}

	// Handle result/decision.
	if result, _ := results[0].Bindings["y"].(bool); result {
		a := req.GetA()
		b := req.GetB()
		res := a * b
		return &pb.Response{Result: res}, nil
	}

	return &pb.Response{}, errors.New("Access denied! You are not allowed to Multiply!")
}

func getRego() string {
	return `
	package example.authz

	default allow = false

	allow {
		input.method = "POST"
		input.function = "Add"
		input.user = "Naren"
	}

	allow {
		input.method = "POST"
		input.function = "Multiply"
		input.user = "Sebastian"
	}
	`
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
