package main

import (
    "google.golang.org/grpc"
    "os"

    pb "grpc-tutorial/proto"
)

func main() {
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    srv := grpc.NewServer()
    pb.RegisterCalculatorServiceServer(srv, NewServer())

    listen, err :=
}
