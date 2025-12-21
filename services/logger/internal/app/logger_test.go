package app

import (
	"context"
	"errors"
	"testing"

	"github.com/segmentio/kafka-go"
)

type fakeResult struct {
	msg kafka.Message
	err error
}

type fakeReader struct {
	closeCalled bool
	closeErr    error

	readMsgResult fakeResult
	readMsgCalled bool
}

func (fr *fakeReader) ReadMessage(ctx context.Context) (kafka.Message, error) {
	fr.readMsgCalled = true
	return fr.readMsgResult.msg, fr.readMsgResult.err
}

func (fr *fakeReader) Close() error {
	fr.closeCalled = true
	return fr.closeErr
}

func TestLogger_Close_DelegatesToMessageReader(t *testing.T) {
	wantErr := errors.New("my close err")
	fr := &fakeReader{closeErr: wantErr}
	logger := NewLogger(fr)

	err := logger.Close()

	if !fr.closeCalled {
		t.Fatalf("expected reader.Close to be called")
	}
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected %v, got %v", wantErr, err)
	}
}

func TestLogger_Run_StopsOnContextCanceled(t *testing.T) {
	ctx := context.Background()

	fr := &fakeReader{
		readMsgResult: fakeResult{
			msg: kafka.Message{},
			err: context.Canceled,
		},
	}
	logger := NewLogger(fr)

	err := logger.Run(ctx)

	if !fr.readMsgCalled {
		t.Fatalf("expected ReadMessage to be called")
	}
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestLogger_Run_ReturnsReadMessageError(t *testing.T) {
	ctx := context.Background()

	wantErr := errors.New("my error")
	fr := &fakeReader{
		readMsgResult: fakeResult{
			msg: kafka.Message{},
			err: wantErr,
		}}
	logger := NewLogger(fr)

	err := logger.Run(ctx)

	if !fr.readMsgCalled {
		t.Fatalf("expected ReadMessage to be called")
	}
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected %v, got %v", wantErr, err)
	}
}
