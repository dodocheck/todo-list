package grpc

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/dodocheck/go-pet-project-1/services/db/internal/app"
	"github.com/dodocheck/go-pet-project-1/services/db/internal/models"
	"github.com/dodocheck/go-pet-project-1/services/db/pb"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/testing/protocmp"
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

func TestAddTask_NilTask_ReturnsInvalidArgument(t *testing.T) {
	srv := NewServer(app.NewService(&fakeRepo{}))

	got, err := srv.AddTask(context.Background(), nil)

	if got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
	if status.Code(err) != codes.InvalidArgument {
		t.Fatalf("code=%v want=%v got=%v", status.Code(err), codes.InvalidArgument, err)
	}
}

func TestAddTask_ServiceError_ReportsInternalError(t *testing.T) {
	fr := &fakeRepo{addTaskErr: errors.New("boom")}
	srv := NewServer(app.NewService(fr))

	got, err := srv.AddTask(context.Background(), &pb.TaskImportData{})

	if got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
	if status.Code(err) != codes.Internal {
		t.Fatalf("code=%v want=%v got=%v", status.Code(err), codes.Internal, err)
	}
}

func TestAddTask_OK_DelegatesToService(t *testing.T) {
	ctx := context.Background()
	wantTaskIn := &pb.TaskImportData{
		Title: "my in title",
		Text:  "my in text",
	}
	createdAtTS := time.Date(2025, 12, 10, 4, 6, 3, 2, time.UTC)
	finishedAtTS := time.Date(2025, 12, 10, 3, 5, 2, 1, time.UTC)
	wantTaskOut := models.TaskExportData{
		Id:         4,
		Title:      "my title",
		Text:       "my text",
		Finished:   true,
		CreatedAt:  createdAtTS,
		FinishedAt: &finishedAtTS,
	}
	fr := &fakeRepo{
		addTaskRet: wantTaskOut,
	}
	srv := NewServer(app.NewService(fr))

	got, _ := srv.AddTask(ctx, wantTaskIn)

	if fr.addTaskCalls != 1 {
		t.Fatalf("expected AddTask calls=1, got=%d", fr.addTaskCalls)
	}
	if fr.addTaskCtx != ctx {
		t.Fatal("context mismatch")
	}
	if diff := cmp.Diff(fr.addTaskIn, taskImportDataFromPB(wantTaskIn), protocmp.Transform()); diff != "" {
		t.Fatal(diff)
	}
	if diff := cmp.Diff(got, taskExportDataToPB(wantTaskOut), protocmp.Transform()); diff != "" {
		t.Fatal(diff)
	}

}
