package app

import (
	"context"
	"errors"
	"testing"

	"github.com/dodocheck/go-pet-project-1/services/db/internal/models"
)

type fakeRepo struct {
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

	closeCalled bool
	closeErr    error
}

func (f *fakeRepo) AddTask(ctx context.Context, task models.TaskImportData) (models.TaskExportData, error) {
	f.addCalls++
	f.gotAddCtx = ctx
	f.gotAddTask = task

	if f.addFn == nil {
		panic("AddTask called but addFn not set")
	}

	return f.addFn(ctx, task)
}

func (f *fakeRepo) DeleteTask(ctx context.Context, id int) error {
	f.removeCalls++
	f.gotRemoveCtx = ctx
	f.gotRemoveId = id

	if f.removeFn == nil {
		panic("DeleteTask called but removeFn not set")
	}

	return f.removeFn(ctx, id)
}

func (f *fakeRepo) ListAllTasks(ctx context.Context) ([]models.TaskExportData, error) {
	f.listCalls++
	f.gotListCtx = ctx

	if f.listFn == nil {
		panic("ListAllTasks called but listFn not set")
	}

	return f.listFn(ctx)
}

func (f *fakeRepo) MarkTaskFinished(ctx context.Context, id int) (models.TaskExportData, error) {
	f.doneCalls++
	f.gotDoneCtx = ctx
	f.gotDoneId = id

	if f.doneFn == nil {
		panic("MarkTaskFinished called but doneFn not set")
	}

	return f.doneFn(ctx, id)
}

func (f *fakeRepo) Close() error {
	f.closeCalled = true
	return f.closeErr
}

func TestAddTask_DelegatesToTaskRepo(t *testing.T) {
	wantTask := models.TaskExportData{
		Id:    1,
		Title: "my title",
		Text:  "my text",
	}
	wantErr := errors.New("boom")
	fakeRepo := &fakeRepo{
		addFn: func(ctx context.Context, task models.TaskImportData) (models.TaskExportData, error) {
			return wantTask, wantErr
		},
	}
	svc := NewService(fakeRepo)

	gotTask, gotErr := svc.AddTask(context.Background(), models.TaskImportData{})

	if !errors.Is(gotErr, wantErr) {
		t.Fatalf("expected err %v, got %v", wantErr, gotErr)
	}
	if gotTask.Id != wantTask.Id {
		t.Fatalf("expected Id %d, got %d", wantTask.Id, gotTask.Id)
	}
	if gotTask.Title != wantTask.Title {
		t.Fatalf("expected Title %q, got %q", wantTask.Title, gotTask.Title)
	}
	if gotTask.Text != wantTask.Text {
		t.Fatalf("expected Text %q, got %q", wantTask.Text, gotTask.Text)
	}
}

func TestDeleteTask_DelegatesToTaskRepo(t *testing.T) {
	wantErr := errors.New("boom")
	fakeRepo := &fakeRepo{
		removeFn: func(ctx context.Context, id int) error {
			return wantErr
		},
	}
	svc := NewService(fakeRepo)

	gotErr := svc.DeleteTask(context.Background(), 1)

	if !errors.Is(gotErr, wantErr) {
		t.Fatalf("expected err %v, got %v", wantErr, gotErr)
	}
}

func TestListAllTasks_DelegatesToTaskRepo(t *testing.T) {
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
		listFn: func(ctx context.Context) ([]models.TaskExportData, error) {
			return wantList, wantErr
		},
	}
	svc := NewService(fakeRepo)

	gotTask, gotErr := svc.ListAllTasks(context.Background())

	if !errors.Is(gotErr, wantErr) {
		t.Fatalf("expected err %v, got %v", wantErr, gotErr)
	}

	if len(wantList) != 2 {
		t.Fatalf("expected len task list =2, got %d", len(wantList))
	}
	if gotTask[0].Id != wantList[0].Id ||
		gotTask[1].Id != wantList[1].Id {
		t.Fatalf("field Id mismatch")
	}
	if gotTask[0].Title != wantList[0].Title ||
		gotTask[1].Title != wantList[1].Title {
		t.Fatalf("field Title mismatch")
	}
	if gotTask[0].Text != wantList[0].Text ||
		gotTask[1].Text != wantList[1].Text {
		t.Fatalf("field Text mismatch")
	}
}

func TestMarkTaskFinished_DelegatesToTaskRepo(t *testing.T) {
	wantTask := models.TaskExportData{
		Id:    1,
		Title: "my title",
		Text:  "my text",
	}
	wantErr := errors.New("boom")
	fakeRepo := &fakeRepo{
		doneFn: func(ctx context.Context, id int) (models.TaskExportData, error) {
			return wantTask, wantErr
		},
	}
	svc := NewService(fakeRepo)

	gotTask, gotErr := svc.MarkTaskFinished(context.Background(), 1)

	if !errors.Is(gotErr, wantErr) {
		t.Fatalf("expected err %v, got %v", wantErr, gotErr)
	}
	if gotTask.Id != wantTask.Id {
		t.Fatalf("expected Id %d, got %d", wantTask.Id, gotTask.Id)
	}
	if gotTask.Title != wantTask.Title {
		t.Fatalf("expected Title %q, got %q", wantTask.Title, gotTask.Title)
	}
	if gotTask.Text != wantTask.Text {
		t.Fatalf("expected Text %q, got %q", wantTask.Text, gotTask.Text)
	}
}
