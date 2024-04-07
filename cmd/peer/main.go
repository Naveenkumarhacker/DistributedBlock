package main

import (
	"DistributedBlock/constants"
	"DistributedBlock/database"
	"DistributedBlock/pb"
	"DistributedBlock/pkg/crypto"
	"DistributedBlock/servers"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	crypto.Init()

	nodeServer := servers.InitNode()
	defer nodeServer.Close()

	// Initialize a Gorp DbMap
	constants.DbMap = database.InitDb()
	defer constants.DbMap.Db.Close()

	conn, err := grpc.Dial(constants.BlockGRPCServiceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// Create gRPC client
	client := pb.NewBlockServiceClient(conn)

	servers.InitHttpServer(client, nodeServer)

	// TODO: Waiting until stop signal come
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGTERM)
	<-interruptChan
	fmt.Println("Server stopped gracefully.")

}
