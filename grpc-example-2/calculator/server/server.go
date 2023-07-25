package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	pb "grpc-example-2/calculator/calculatorpb"
	"io"
	"log"
	"net"
)

// server implements CalculatorServer ...
type server struct {
	pb.UnimplementedCalculatorServer
}

func (s *server) Sum(ctx context.Context, req *pb.SumRequest) (*pb.SumResponse, error) {
	fmt.Println("call: Sum()")
	resp := &pb.SumResponse{
		Result: req.Num1 + req.Num2,
	}

	return resp, nil
}

func (s *server) Average(stream pb.Calculator_AverageServer) error {
	fmt.Println("call: Average()")
	var avg float32
	var count int

	for {
		req, err := stream.Recv()

		avg += req.Num
		count++

		if err == io.EOF {
			return stream.SendAndClose(&pb.AverageResponse{
				Result: avg / float32(count),
			})
		}
		if err != nil {
			log.Fatalf("err while recv average %v", err)
			return err
		}
	}
}

func (s *server) PrimeNumberDecomposition(req *pb.PNDRequest, stream pb.Calculator_PrimeNumberDecompositionServer) error {
	fmt.Println("call: PrimeNumberDecomposition()")

	k := int32(2)
	N := req.GetNumber()
	for N > 1 {
		if N%k == 0 {
			N = N / k
			stream.Send(&pb.PNDResponse{
				Result: k,
			})
		} else {
			k++
			log.Printf("k increase to %v\n", k)
		}
	}

	return nil
}

func (s *server) FindMax(stream pb.Calculator_FindMaxServer) error {
	fmt.Println("call: FindMax()")
	max := int32(0)

	for {
		req, err := stream.Recv()

		if err == io.EOF {
			log.Println("EOF ...")
			return nil
		}
		if err != nil {
			log.Fatalf("err while recv FindMax %v", err)
			return err
		}

		num := req.Num
		if num > max {
			max = num
		}
		err = stream.Send(&pb.FindMaxResponse{
			Result: max,
		})
		if err != nil {
			log.Fatalf("err while send %v", err)
			return err
		}
	}
}

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:8089")
	if err != nil {
		log.Fatalf("create listen err %s", err)
	}

	s := grpc.NewServer()
	pb.RegisterCalculatorServer(s, &server{})
	fmt.Println("Server is running ...")

	// Serve accepts incoming connections on the listener lis
	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("err while serve %s", err)
	}
}
