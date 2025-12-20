package main

import (
	"context"
	"log"

	"github.com/dodocheck/go-pet-project-1/services/logger/internal/app"
	"github.com/dodocheck/go-pet-project-1/services/logger/internal/kafka"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	reader := kafka.NewKafkaReader("kafka:9092", "actions-log", "loggerGroupId")
	service := app.NewService(reader)

	if err := service.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
