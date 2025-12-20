package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/dodocheck/go-pet-project-1/services/api/internal/models"
	"github.com/segmentio/kafka-go"
)

type KafkaWriter struct {
	writer     *kafka.Writer
	logChannel <-chan models.ActionLog
}

func NewKafkaWriter(broker, topic string, logChannel <-chan models.ActionLog) *KafkaWriter {
	return &KafkaWriter{
		writer: kafka.NewWriter(
			kafka.WriterConfig{
				Brokers: []string{broker},
				Topic:   topic,
			}),
		logChannel: logChannel,
	}
}

func (kw *KafkaWriter) Close() error {
	return kw.writer.Close()
}

func (kw *KafkaWriter) Run(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case newLog := <-kw.logChannel:
				logBytes, err := json.Marshal(newLog)
				if err != nil {
					log.Println("failed to marshal action log")
					return
				}

				msg := kafka.Message{Value: logBytes}
				if err := kw.writer.WriteMessages(ctx, msg); err != nil {
					log.Println("failed to write msg to kafka")
					return
				}
			}
		}
	}()
}
