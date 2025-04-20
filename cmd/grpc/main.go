package main

import (
	"github.com/51mans0n/avito-pvz-task/internal/db"
	grpcserver "github.com/51mans0n/avito-pvz-task/internal/grpc"
	pvz_v1 "github.com/51mans0n/avito-pvz-task/pkg/proto/pvz/v1"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	conn, err := db.InitDB()
	if err != nil {
		log.Fatalf("db: %v", err)
	}
	repo := db.NewRepo(conn)

	s := grpc.NewServer()
	pvz_v1.RegisterPVZServiceServer(s, grpcserver.New(repo))

	lis, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatalf("listen: %v", err)
	}
	log.Println("gRPC server :3000")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("serve: %v", err)
	}
}
