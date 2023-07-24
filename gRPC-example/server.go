package main

import (
	"context"
	"gRPC-example/pb"
	"google.golang.org/grpc"
	"log"
	"net"
)

type myInvoicerServer struct {
	pb.UnimplementedInvoicerServer
}

func (receiver myInvoicerServer) Create(context.Context, *pb.CreateRequest) (*pb.CreateResponse, error) {
	return &pb.CreateResponse{
		Pdf:  nil,
		Docx: nil,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":8088")
	if err != nil {
		log.Fatalf("cannot create listener: %s", err)
	}

	serverRegistrar := grpc.NewServer()
	service := myInvoicerServer{}

	pb.RegisterInvoicerServer(serverRegistrar, service)
}
