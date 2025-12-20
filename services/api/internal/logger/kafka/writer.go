package kafka

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/dodocheck/go-pet-project-1/services/api/internal/models"
	"github.com/segmentio/kafka-go"
)

type KafkaWriter struct {
	writer     *kafka.Writer
	logChannel <-chan models.ActionLog
}

func NewKafkaWriter(broker, topic string, logChannel <-chan models.ActionLog) *KafkaWriter {
	w := kafka.NewWriter(
		kafka.WriterConfig{
			Brokers: []string{broker},
			Topic:   topic,
		})
	w.AllowAutoTopicCreation = true

	return &KafkaWriter{
		writer:     w,
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
					log.Printf("failed to marshal action log: %v\n", err)
					continue
				}

				msg := kafka.Message{Value: logBytes}
				retries := 3
				for i := range retries {
					if err := kw.writer.WriteMessages(ctx, msg); err == nil {
						break
					} else {
						log.Printf("failed to write msg to kafka (attempt %d/%d): %v\n", i+1, retries, err)
						time.Sleep(200 * time.Millisecond)
					}
				}

			}
		}
	}()
}
