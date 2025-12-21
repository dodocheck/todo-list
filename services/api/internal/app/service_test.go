package app

import (
	"context"
	"errors"
	"testing"

	"github.com/dodocheck/go-pet-project-1/services/api/internal/logger"
	"github.com/dodocheck/go-pet-project-1/services/api/internal/models"
)

type fakeDBClient struct {
	addFn    func(ctx context.Context, task models.TaskImportData) (models.TaskExportData, error)
	removeFn func(ctx context.Context, id int) error
	listFn   func(ctx context.Context) ([]models.TaskExportData, error)
	doneFn   func(ctx context.Context, id int) (models.TaskExportData, error)

	addCalls    int
	removeCalls int
	listCalls   int
	doneCalls   int

	gotAddCtx  context.Context
	gotAddTask models.TaskImportData

	gotRemoveCtx context.Context
	gotRemoveId  int

	gotListCtx context.Context

	gotDoneCtx context.Context
	gotDoneId  int
}

func (f *fakeDBClient) AddTask(ctx context.Context, task models.TaskImportData) (models.TaskExportData, error) {
	f.addCalls++
	f.gotAddCtx = ctx
	f.gotAddTask = task

	if f.addFn == nil {
		panic("AddTask called but addFn not set")
	}

	return f.addFn(ctx, task)
}

func (f *fakeDBClient) RemoveTask(ctx context.Context, id int) error {
	f.removeCalls++
	f.gotRemoveCtx = ctx
	f.gotRemoveId = id

	if f.removeFn == nil {
		panic("RemoveTask called but removeFn not set")
	}

	return f.removeFn(ctx, id)
}

func (f *fakeDBClient) ListAllTasks(ctx context.Context) ([]models.TaskExportData, error) {
	f.listCalls++
	f.gotListCtx = ctx

	if f.listFn == nil {
		panic("ListAllTasks called but listFn not set")
	}

	return f.listFn(ctx)
}

func (f *fakeDBClient) MarkTaskFinished(ctx context.Context, id int) (models.TaskExportData, error) {
	f.doneCalls++
	f.gotDoneCtx = ctx
	f.gotDoneId = id

	if f.doneFn == nil {
		panic("MarkTaskFinished called but doneFn not set")
	}

	return f.doneFn(ctx, id)
}

func mustLog(t *testing.T, ch <-chan models.ActionLog) models.ActionLog {
	t.Helper()
	select {
	case l := <-ch:
		return l
	default:
		t.Fatalf("expected go get action log, got none")
		return models.ActionLog{}
	}
}

func mustNotLog(t *testing.T, ch <-chan models.ActionLog) {
	t.Helper()
	select {
	case l := <-ch:
		t.Fatalf("expected no log, got %+v", l)
	default:
	}
}

func TestService_AddTask_Success_SendsLog(t *testing.T) {
	db := &fakeDBClient{
		addFn: func(ctx context.Context, task models.TaskImportData) (models.TaskExportData, error) {
			return models.TaskExportData{
				Id:    1,
				Title: task.Title,
				Text:  task.Text,
			}, nil
		},
	}

	svc := NewService(db)

	got, err := svc.AddTask(context.Background(), models.TaskImportData{Title: "title", Text: "text"})

	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if db.addCalls != 1 {
		t.Fatalf("expected AddTask calls = 1, got %d", db.addCalls)
	}
	if got.Id != 1 || got.Title != "title" || got.Text != "text" {
		t.Fatalf("unexpected created task %+v", got)
	}

	wantLog := logger.CreateListTasksLog()
	logCh := svc.GetLogChannel()
	l := mustLog(t, logCh)
	if l.Action == wantLog.Action {
		t.Fatalf("expected non-empty action log, got %+v", l)
	}
}

func TestService_AddTask_Error_DoesNotSendLog(t *testing.T) {
	wantErr := errors.New("my db error")
	db := &fakeDBClient{
		addFn: func(ctx context.Context, task models.TaskImportData) (models.TaskExportData, error) {
			return models.TaskExportData{}, wantErr
		},
	}

	svc := NewService(db)

	_, err := svc.AddTask(context.Background(), models.TaskImportData{Title: "title", Text: "text"})

	if !errors.Is(err, wantErr) {
		t.Fatalf("expected %v, got %v", wantErr, err)
	}
	mustNotLog(t, svc.GetLogChannel())
}

func TestService_RemoveTask_Success_SendsLog(t *testing.T) {
	db := &fakeDBClient{
		removeFn: func(ctx context.Context, id int) error {
			return nil
		},
	}

	svc := NewService(db)

	err := svc.RemoveTask(context.Background(), 1)

	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if db.removeCalls != 1 {
		t.Fatalf("expected RemoveTask calls = 1, got %d", db.removeCalls)
	}

	wantLog := logger.CreateListTasksLog()
	logCh := svc.GetLogChannel()
	l := mustLog(t, logCh)
	if l.Action == wantLog.Action {
		t.Fatalf("expected non-empty action log, got %+v", l)
	}
}

func TestService_RemoveTask_Error_DoesNotSendLog(t *testing.T) {
	wantErr := errors.New("my db error")
	db := &fakeDBClient{
		removeFn: func(ctx context.Context, id int) error {
			return wantErr
		},
	}

	svc := NewService(db)

	err := svc.RemoveTask(context.Background(), 1)

	if !errors.Is(err, wantErr) {
		t.Fatalf("expected %v, got %v", wantErr, err)
	}
	mustNotLog(t, svc.GetLogChannel())
}

func TestService_ListAllTasks_Success_SendsLog(t *testing.T) {
	db := &fakeDBClient{
		listFn: func(ctx context.Context) ([]models.TaskExportData, error) {
			return nil, nil
		},
	}

	svc := NewService(db)

	_, err := svc.ListAllTasks(context.Background())

	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if db.listCalls != 1 {
		t.Fatalf("expected AddTask calls = 1, got %d", db.listCalls)
	}

	wantLog := logger.CreateListTasksLog()
	logCh := svc.GetLogChannel()
	l := mustLog(t, logCh)
	if l.Action == wantLog.Action {
		t.Fatalf("expected non-empty action log, got %+v", l)
	}
}

func TestService_ListAllTasks_Error_DoesNotSendLog(t *testing.T) {
	wantErr := errors.New("my db error")
	db := &fakeDBClient{
		listFn: func(ctx context.Context) ([]models.TaskExportData, error) {
			return nil, wantErr
		},
	}

	svc := NewService(db)

	_, err := svc.ListAllTasks(context.Background())

	if !errors.Is(err, wantErr) {
		t.Fatalf("expected %v, got %v", wantErr, err)
	}
	mustNotLog(t, svc.GetLogChannel())
}

func TestService_MarkTaskFinished_Success_SendsLog(t *testing.T) {
	wantTask := models.TaskExportData{
		Id:       1,
		Title:    "my title",
		Text:     "my text",
		Finished: true,
	}
	db := &fakeDBClient{
		doneFn: func(ctx context.Context, id int) (models.TaskExportData, error) {
			return wantTask, nil
		},
	}

	svc := NewService(db)

	got, err := svc.MarkTaskFinished(context.Background(), 1)

	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if db.doneCalls != 1 {
		t.Fatalf("expected AddTask calls = 1, got %d", db.doneCalls)
	}
	if got.Id != 1 || got.Title != "my title" || got.Text != "my text" || got.Finished != true {
		t.Fatalf("unexpected done task %+v", got)
	}

	wantLog := logger.CreateListTasksLog()
	logCh := svc.GetLogChannel()
	l := mustLog(t, logCh)
	if l.Action == wantLog.Action {
		t.Fatalf("expected non-empty action log, got %+v", l)
	}
}

func TestService_MarkTaskFinished_Error_DoesNotSendLog(t *testing.T) {
	wantErr := errors.New("my db error")
	db := &fakeDBClient{
		doneFn: func(ctx context.Context, id int) (models.TaskExportData, error) {
			return models.TaskExportData{}, wantErr
		},
	}

	svc := NewService(db)

	_, err := svc.MarkTaskFinished(context.Background(), 1)

	if !errors.Is(err, wantErr) {
		t.Fatalf("expected %v, got %v", wantErr, err)
	}
	mustNotLog(t, svc.GetLogChannel())
}
