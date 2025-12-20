package kafka

import (
	"context"
	"errors"
	"log"

	"github.com/segmentio/kafka-go"
)

type KafkaReader struct {
	reader *kafka.Reader
}

func NewKafkaReader(broker, topic, groupId string) *KafkaReader {
	return &KafkaReader{
		reader: kafka.NewReader(
			kafka.ReaderConfig{
				Brokers: []string{broker},
				Topic:   topic,
				GroupID: groupId})}
}

func (kr *KafkaReader) Close() error {
	return kr.reader.Close()
}

func (kr *KafkaReader) Run(ctx context.Context) error {
	for {
		msg, err := kr.reader.ReadMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return nil
			}
			return err
		}
		log.Println(string(msg.Value))
	}
}
