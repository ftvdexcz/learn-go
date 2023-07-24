package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	pb "grpc-example-2/calculator/calculatorpb"
	"io"
	"log"
)

func main() {
	conn, err := grpc.Dial("localhost:8089", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("err while dial %v", err)
	}

	defer conn.Close()

	client := pb.NewCalculatorClient(conn)

	//callSum(client)
	//callPrimeNumberDecomposition(client)
	callAverage(client)
}

func callSum(c pb.CalculatorClient) {
	resp, err := c.Sum(context.Background(), &pb.SumRequest{
		Num1: 5,
		Num2: 3,
	})

	if err != nil {
		log.Fatalf("call Sum err: %s", err)
	}

	fmt.Printf("Result: %d", resp.Result)
}

func callPrimeNumberDecomposition(c pb.CalculatorClient) {
	stream, err := c.PrimeNumberDecomposition(context.Background(), &pb.PNDRequest{
		Number: 120,
	})

	if err != nil {
		log.Fatalf("err when decomposition: %s", err)
	}

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			log.Println("server finish streaming")
			return
		}

		fmt.Printf("Result: %d\n", resp.Result)
	}
}

func callAverage(c pb.CalculatorClient) {
	stream, err := c.Average(context.Background())
	if err != nil {
		log.Fatalf("err when call average: %s", err)
	}

	err = stream.Send(&pb.AverageRequest{
		Num: 5,
	})
	if err != nil {
		log.Fatalf("err when send request: %s", err)
	}

	err = stream.Send(&pb.AverageRequest{
		Num: 2,
	})
	if err != nil {
		log.Fatalf("err when send request: %s", err)
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("err when receive: %s", err)
	}

	log.Printf("average: %v", resp)
}
