package main

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Starting gRPC service on :3000...")

	lis, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	// TODO: Зарегать услугу тут
	log.Fatal(s.Serve(lis))
}
