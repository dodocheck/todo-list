package app

import (
	"context"
	"errors"
	"log"

	"github.com/segmentio/kafka-go"
)

type MessageReader interface {
	ReadMessage(ctx context.Context) (kafka.Message, error)
	Close() error
}

type Logger struct {
	reader MessageReader
}

func NewLogger(reader MessageReader) *Logger {
	return &Logger{
		reader: reader}
}

func (l *Logger) Close() error {
	return l.reader.Close()
}

func (l *Logger) Run(ctx context.Context) error {
	for {
		msg, err := l.reader.ReadMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return nil
			}
			return err
		}
		log.Println(string(msg.Value))
	}
}
