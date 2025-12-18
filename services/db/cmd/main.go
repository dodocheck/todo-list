package main

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/dodocheck/go-pet-project-1/services/db/internal/app"
	"github.com/dodocheck/go-pet-project-1/services/db/internal/dbcontroller/cachedrepository"
	"github.com/dodocheck/go-pet-project-1/services/db/internal/dbcontroller/postgres"
	"github.com/dodocheck/go-pet-project-1/services/db/internal/transport/grpc"
)

func main() {
	postgresController := postgres.NewPostgresController()
	ttlSeconds, _ := strconv.Atoi(os.Getenv("REDIS_TTL_SECONDS"))
	cachedController, err := cachedrepository.NewCachedRepository(context.Background(), postgresController, "redis:6379", ttlSeconds)
	if err != nil {
		log.Fatalf("failed to create db cached controller: %v", err)
	}

	service := app.NewService(cachedController)

	server := grpc.NewServer(service)

	dbGrpcServerAddress := ":" + os.Getenv("DB_SERVICE_INTERNAL_PORT")
	if err := server.StartServer(dbGrpcServerAddress); err != nil {
		log.Fatal("Failed to start server:", err)
	}

}
