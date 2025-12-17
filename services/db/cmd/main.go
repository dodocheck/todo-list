package main

import (
	"log"
	"os"

	"github.com/dodocheck/go-pet-project-1/services/db/internal/dbcontroller/postgres"
	"github.com/dodocheck/go-pet-project-1/services/db/internal/transport/grpc"
)

func main() {
	dbController := postgres.NewPostgresController()
	server := grpc.NewServer(dbController)

	dbGrpcServerAddress := os.Getenv("GRPC_SERVER_ADDR")
	if err := server.StartServer(dbGrpcServerAddress); err != nil {
		log.Fatal("Failed to start server:", err)
	}

}
