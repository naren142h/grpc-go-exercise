package server

import (
	"context"
	pb "grpc-tutorial/proto"
)

var _ pb.CalculatorServiceServer = (*Server)(nil)

type Server struct {
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Add(ctx context.Context, req *pb.Request) (*pb.Response, error) {

}

func (s *Server) Multiply(ctx context.Context, req *pb.Request) (*pb.Response, error) {

}
