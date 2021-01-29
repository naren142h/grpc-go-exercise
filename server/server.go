package main

import (
	"context"
	pb "grpc-tutorial/proto"
)

var _ pb.CalculatorServiceServer = (*Server)(nil)

// Server is a struct implements the pb.CalculatorServiceServer
type Server struct {
}

// NewServer returns a new Server
func NewServer() *Server {
	return &Server{}
}

// Add adds two numbers
func (s *Server) Add(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	a := req.GetA()
	b := req.GetB()
	res := a + b
	return &pb.Response{Result: res}, nil
}

//Multiply multiplies two numbers
func (s *Server) Multiply(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	a := req.GetA()
	b := req.GetB()
	res := a * b
	return &pb.Response{Result: res}, nil
}
