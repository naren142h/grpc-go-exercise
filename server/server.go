package server

import (
	"context"
	pb "grpc-tutorial/proto"
	"log"
	"fmt"
)

var _ pb.CalculatorServiceServer = (*Server)(nil)

type Server struct {
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Add(ctx context.Context, req *pb.Request) (*pb.Response, error) {
    a := req.GetA()
    b := req.GetB()
    res := a + b
    return &pb.Response{result: res}, nil
}

func (s *Server) Multiply(ctx context.Context, req *pb.Request) (*pb.Response, error) {
    a := req.GetA()
    b := req.GetB()
    res := a * b
    return &pb.Response{result: res}, nil
}
