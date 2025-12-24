package main

import (
	"context"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/dodocheck/go-pet-project-1/services/logger/internal/app"
	"github.com/segmentio/kafka-go"
)

func main() {
	logPath := os.Getenv("LOG_FILE_PATH")

	_ = os.MkdirAll(filepath.Dir(logPath), 0o755)

	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		log.Fatal("open log file error:", err)
	}
	defer func() { _ = f.Close() }()

	log.SetOutput(io.MultiWriter(os.Stdout, f))
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	topic := os.Getenv("KAFKA_TOPIC_NAME")
	kafkaReader := kafka.NewReader(
		kafka.ReaderConfig{
			Brokers: []string{"kafka:9092"},
			Topic:   topic,
			GroupID: "loggerGroupId"})

	logger := app.NewLogger(kafkaReader)

	if err := logger.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
