package logger

import (
	"context"
	"errors"
	"testing"

	"github.com/dodocheck/go-pet-project-1/services/api/internal/models"
	"github.com/segmentio/kafka-go"
)

type fakeWriter struct {
	closeCalled bool
	closeErr    error

	writeMsgCalled bool
	gotCtx         context.Context
	gotCtxCancel   context.CancelFunc
}

func (fw *fakeWriter) Close() error {
	fw.closeCalled = true
	return fw.closeErr
}

func (fw *fakeWriter) WriteMessages(ctx context.Context, msgs ...kafka.Message) error {
	fw.writeMsgCalled = true
	fw.gotCtxCancel()
	return nil
}

func TestClose_DelegatesToMessageWriter(t *testing.T) {
	wantErr := errors.New("my error")
	fw := &fakeWriter{
		closeErr: wantErr,
	}
	logCh := make(chan models.ActionLog)
	logger := NewLogger(fw, logCh)

	err := logger.Close()
	if !fw.closeCalled {
		t.Fatalf("expected Close to be called")
	}
	if err != wantErr {
		t.Fatalf("expected %v, got %v", wantErr, err)
	}
}

func TestRun_StopsOnContextCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	fw := &fakeWriter{}
	logCh := make(chan models.ActionLog)
	logger := NewLogger(fw, logCh)

	err := logger.Run(ctx)
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestRun_RetriesToWriteMessage(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	fw := &fakeWriter{
		gotCtx:       ctx,
		gotCtxCancel: cancel,
	}
	logCh := make(chan models.ActionLog, 10)
	logCh <- CreateListTasksLog()
	logCh <- CreateTaskAddedLog()
	logCh <- CreateTaskDeletedLog()
	logCh <- CreateTaskDoneLog()
	logger := NewLogger(fw, logCh)

	err := logger.Run(ctx)
	if !fw.writeMsgCalled {
		t.Fatalf("expected WriteMessage to be called")
	}
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

}
