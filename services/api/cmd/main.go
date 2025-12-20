package main

import (
	"context"
	"log"
	"os"

	"github.com/dodocheck/go-pet-project-1/services/api/internal/app"
	dbgrpc "github.com/dodocheck/go-pet-project-1/services/api/internal/dbclient/grpc"
	"github.com/dodocheck/go-pet-project-1/services/api/internal/logger/kafka"
	"github.com/dodocheck/go-pet-project-1/services/api/internal/transport/http"
	"github.com/dodocheck/go-pet-project-1/services/api/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	dbAddr := "db-service:" + os.Getenv("DB_SERVICE_INTERNAL_PORT")
	conn, err := grpc.NewClient(dbAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("Failed to dial grpc db:", err)
	}
	defer conn.Close()

	grpcClient := pb.NewTasksServiceClient(conn)

	dbClient := dbgrpc.NewDBClient(grpcClient)

	service := app.NewService(dbClient)

	topic := os.Getenv("KAFKA_TOPIC_NAME")
	userActionLogger := kafka.NewKafkaWriter("kafka:9092", topic, service.GetLogChannel())
	userActionLogger.Run(context.Background())

	httpServer := http.NewHttpServer(service)

	if err := httpServer.StartServer(); err != nil {
		log.Fatal("Failed to start http web server:", err)
	}
}
