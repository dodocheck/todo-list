package main

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/dodocheck/go-pet-project-1/services/db/internal/dbcontroller/cachedcontroller"
	"github.com/dodocheck/go-pet-project-1/services/db/internal/transport/grpc"
)

func main() {
	ttlSeconds, _ := strconv.Atoi(os.Getenv("REDIS_TTL_SECONDS"))
	cachedController, err := cachedcontroller.NewCachedController(context.Background(), "redis:6379", ttlSeconds)
	if err != nil {
		log.Fatalf("failed to create db cached controller: %v", err)
	}

	server := grpc.NewServer(cachedController)

	dbGrpcServerAddress := ":" + os.Getenv("DB_SERVICE_INTERNAL_PORT")
	if err := server.StartServer(dbGrpcServerAddress); err != nil {
		log.Fatal("Failed to start server:", err)
	}

}
