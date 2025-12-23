package app

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/dodocheck/go-pet-project-1/services/db/internal/models"
)

type fakeRepo struct {
	addTaskCalls int
	addTaskCtx   context.Context
	addTaskIn    models.TaskImportData
	addTaskRet   models.TaskExportData
	addTaskErr   error

	deleteTaskCalls int
	deleteTaskCtx   context.Context
	deleteTaskIn    int
	deleteTaskErr   error

	listAllTasksCalls int
	listAllTasksCtx   context.Context
	listAllTasksRet   []models.TaskExportData
	listAllTasksErr   error

	markTaskFinishedCalls int
	markTaskFinishedCtx   context.Context
	markTaskFinishedIn    int
	markTaskFinishedRet   models.TaskExportData
	markTaskFinishedErr   error

	closeCalled int
	closeErr    error
}

func (f *fakeRepo) AddTask(ctx context.Context, task models.TaskImportData) (models.TaskExportData, error) {
	f.addTaskCalls++
	f.addTaskCtx = ctx
	f.addTaskIn = task
	return f.addTaskRet, f.addTaskErr
}

func (f *fakeRepo) DeleteTask(ctx context.Context, id int) error {
	f.deleteTaskCalls++
	f.deleteTaskCtx = ctx
	f.deleteTaskIn = id
	return f.deleteTaskErr
}

func (f *fakeRepo) ListAllTasks(ctx context.Context) ([]models.TaskExportData, error) {
	f.listAllTasksCalls++
	f.listAllTasksCtx = ctx
	return f.listAllTasksRet, f.listAllTasksErr
}

func (f *fakeRepo) MarkTaskFinished(ctx context.Context, id int) (models.TaskExportData, error) {
	f.markTaskFinishedCalls++
	f.markTaskFinishedCtx = ctx
	f.markTaskFinishedIn = id
	return f.markTaskFinishedRet, f.markTaskFinishedErr
}

func (f *fakeRepo) Close() error {
	f.closeCalled++
	return f.closeErr
}

func TestServiceAddTask_DelegatesToTaskRepo(t *testing.T) {
	ctx := context.Background()
	createdAtTS := time.Date(2025, 12, 10, 4, 6, 3, 2, time.UTC)
	finishedAtTS := time.Date(2025, 12, 10, 3, 5, 2, 1, time.UTC)
	wantTaskIn := models.TaskImportData{
		Title: "some title",
		Text:  "some text",
	}
	wantTaskOut := models.TaskExportData{
		Id:         1,
		Title:      "my title",
		Text:       "my text",
		Finished:   true,
		CreatedAt:  createdAtTS,
		FinishedAt: &finishedAtTS,
	}
	wantErr := errors.New("boom")
	fakeRepo := &fakeRepo{
		addTaskRet: wantTaskOut,
		addTaskErr: wantErr,
	}
	svc := NewService(fakeRepo)

	gotTask, gotErr := svc.AddTask(ctx, wantTaskIn)

	if fakeRepo.addTaskCalls != 1 {
		t.Fatalf("expected AddTask called=1, got %d", fakeRepo.addTaskCalls)
	}
	if !errors.Is(gotErr, wantErr) {
		t.Fatalf("expected err %v, got %v", wantErr, gotErr)
	}
	if fakeRepo.addTaskCtx != ctx {
		t.Fatalf("context mismatch")
	}
	if !reflect.DeepEqual(fakeRepo.addTaskIn, wantTaskIn) {
		t.Fatalf("mismatch task in: got:%+v want: %+v", fakeRepo.addTaskIn, wantTaskIn)
	}
	if !reflect.DeepEqual(gotTask, wantTaskOut) {
		t.Fatalf("mismatch task out: got:%+v want: %+v", gotTask, wantTaskOut)
	}
}

func TestServiceDeleteTask_DelegatesToTaskRepo(t *testing.T) {
	ctx := context.Background()
	wantId := 23
	wantErr := errors.New("boom")
	fakeRepo := &fakeRepo{
		deleteTaskErr: wantErr,
	}
	svc := NewService(fakeRepo)

	gotErr := svc.DeleteTask(ctx, wantId)

	if fakeRepo.deleteTaskCalls != 1 {
		t.Fatalf("expected DeleteTask called=1, got %d", fakeRepo.deleteTaskCalls)
	}
	if !errors.Is(gotErr, wantErr) {
		t.Fatalf("expected err %v, got %v", wantErr, gotErr)
	}
	if fakeRepo.deleteTaskCtx != ctx {
		t.Fatalf("context mismatch")
	}
	if fakeRepo.deleteTaskIn != wantId {
		t.Fatalf("expected id=%d, got %d", wantId, fakeRepo.deleteTaskIn)
	}
}

func TestServiceListAllTasks_DelegatesToTaskRepo(t *testing.T) {
	ctx := context.Background()
	wantList := []models.TaskExportData{
		{
			Id:    1,
			Title: "my title1",
			Text:  "my text1",
		},
		{
			Id:    2,
			Title: "my title2",
			Text:  "my text2",
		},
	}

	wantErr := errors.New("boom")
	fakeRepo := &fakeRepo{
		listAllTasksRet: wantList,
		listAllTasksErr: wantErr,
	}
	svc := NewService(fakeRepo)

	got, gotErr := svc.ListAllTasks(context.Background())

	if fakeRepo.listAllTasksCalls != 1 {
		t.Fatalf("expected ListAllTasks called=1, got %d", fakeRepo.listAllTasksCalls)
	}
	if !errors.Is(gotErr, wantErr) {
		t.Fatalf("expected err %v, got %v", wantErr, gotErr)
	}
	if fakeRepo.listAllTasksCtx != ctx {
		t.Fatalf("context mismatch")
	}
	if !reflect.DeepEqual(got, wantList) {
		t.Fatalf("mismatch task list: want=%+v got=%+v", wantList, got)
	}
}

func TestServiceMarkTaskFinished_DelegatesToTaskRepo(t *testing.T) {
	ctx := context.Background()
	wantId := 234
	createdAtTS := time.Date(2025, 12, 10, 4, 6, 3, 2, time.UTC)
	finishedAtTS := time.Date(2025, 12, 10, 3, 5, 2, 1, time.UTC)
	wantTask := models.TaskExportData{
		Id:         wantId,
		Title:      "my title",
		Text:       "my text",
		Finished:   true,
		CreatedAt:  createdAtTS,
		FinishedAt: &finishedAtTS,
	}
	wantErr := errors.New("boom")
	fakeRepo := &fakeRepo{
		markTaskFinishedRet: wantTask,
		markTaskFinishedErr: wantErr,
	}
	svc := NewService(fakeRepo)

	gotTask, gotErr := svc.MarkTaskFinished(context.Background(), wantId)

	if fakeRepo.markTaskFinishedCalls != 1 {
		t.Fatalf("expected MarkTaskFinished called=1, got %d", fakeRepo.markTaskFinishedCalls)
	}
	if !errors.Is(gotErr, wantErr) {
		t.Fatalf("expected err %v, got %v", wantErr, gotErr)
	}
	if fakeRepo.markTaskFinishedCtx != ctx {
		t.Fatalf("context mismatch")
	}
	if fakeRepo.markTaskFinishedIn != wantId {
		t.Fatalf("expected id=%d, got=%d", wantId, fakeRepo.markTaskFinishedIn)
	}
	if !reflect.DeepEqual(gotTask, wantTask) {
		t.Fatalf("mismatch task: want %+v got %+v", wantTask, gotTask)
	}
}
