package main

import (
	"context"
	"errors"
	"fmt"
	pb "github.com/KFN002/distributed-arithmetic-expression-evaluator.git/proto"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct {
	pb.AgentServiceServer
}

func (s *server) Calculate(ctx context.Context, req *pb.CalculationRequest) (*pb.CalculationResponse, error) {
	log.Println("gRPC server live!")

	number1 := float64(req.GetFirstNumber())
	number2 := float64(req.GetSecondNumber())
	operation := req.GetOperation()

	var result float64

	switch operation {
	case "+":
		result = number1 + number2
	case "-":
		result = number1 - number2
	case "*":
		result = number1 * number2
	case "/":
		if number2 == 0 {
			return nil, errors.New("division by zero")
		}
		result = number1 / number2
	default:
		return nil, errors.New("invalid operation")
	}

	log.Println(result)

	return &pb.CalculationResponse{Result: float32(result)}, nil
}

func main() {
	host := "localhost"
	port := "8050"

	addr := fmt.Sprintf("%s:%s", host, port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("error starting TCP listener: %v", err)
	}

	log.Printf("TCP listener started at port: %s", port)

	grpcServer := grpc.NewServer()
	pb.RegisterAgentServiceServer(grpcServer, &server{})

	log.Printf("gRPC server listening at %s", addr)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("error serving gRPC: %v", err)
	}
}
