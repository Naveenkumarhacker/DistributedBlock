package main

import (
	"DistributedBlock/constants"
	"DistributedBlock/database"
	"DistributedBlock/pb"
	"DistributedBlock/servers"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {

	db := database.InitMySqlDB()

	// Initialize gRPC server
	grpcServer := grpc.NewServer()

	// Register server
	s, err := servers.NewBlockServer(db)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}
	pb.RegisterBlockServiceServer(grpcServer, s)

	// Listen to gRPC requests
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", constants.BlockServicePort))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Println("gRPC server is listening on port 50051")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
