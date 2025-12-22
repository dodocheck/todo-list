package main

import (
	"context"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/dodocheck/go-pet-project-1/services/api/internal/app"
	dbgrpc "github.com/dodocheck/go-pet-project-1/services/api/internal/dbclient/grpc"
	"github.com/dodocheck/go-pet-project-1/services/api/internal/logger"
	"github.com/dodocheck/go-pet-project-1/services/api/internal/transport/http"
	"github.com/dodocheck/go-pet-project-1/services/api/pb"
	"github.com/segmentio/kafka-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	logPath := os.Getenv("LOG_FILE_PATH")

	_ = os.MkdirAll(filepath.Dir(logPath), 0o755)

	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		log.Fatal("open log file error:", err)
	}
	defer f.Close()

	log.SetOutput(io.MultiWriter(os.Stdout, f))
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

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
	kafkaWriter := kafka.NewWriter(
		kafka.WriterConfig{
			Brokers: []string{"kafka:9092"},
			Topic:   topic,
		})
	kafkaWriter.AllowAutoTopicCreation = true
	userActionLogger := logger.NewLogger(kafkaWriter, service.GetLogChannel())
	go userActionLogger.Run(context.Background())

	httpServer := http.NewHttpServer(service)

	if err := httpServer.StartServer(); err != nil {
		log.Fatal("Failed to start http web server:", err)
	}
}
