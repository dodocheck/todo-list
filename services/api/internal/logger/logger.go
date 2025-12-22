package logger

import (
	"context"
	"encoding/json"
	"log"

	"github.com/dodocheck/go-pet-project-1/services/api/internal/models"
	"github.com/segmentio/kafka-go"
)

type MessageWriter interface {
	Close() error
	WriteMessages(ctx context.Context, msgs ...kafka.Message) error
}

type Logger struct {
	writer     MessageWriter
	logChannel <-chan models.ActionLog
}

func NewLogger(writer MessageWriter, logChannel <-chan models.ActionLog) *Logger {
	return &Logger{
		writer:     writer,
		logChannel: logChannel,
	}
}

func (l *Logger) Close() error {
	return l.writer.Close()
}

func (l *Logger) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case newLog := <-l.logChannel:
			logBytes, err := json.Marshal(newLog)
			if err != nil {
				log.Printf("failed to marshal action log: %v\n", err)
			}
			if err := l.writer.WriteMessages(ctx, kafka.Message{Value: logBytes}); err == nil {
				log.Printf("failed to write msg to kafka: %v\n", err)
			}
		}
	}
}
