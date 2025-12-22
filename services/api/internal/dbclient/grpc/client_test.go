package dbgrpc

import (
	"context"
	"errors"
	"testing"

	"github.com/dodocheck/go-pet-project-1/services/api/internal/models"
	"github.com/dodocheck/go-pet-project-1/services/api/pb"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type fakeGrpcClient struct {
	addFn    func(ctx context.Context, in *pb.TaskImportData, opts ...grpc.CallOption) (*pb.TaskExportData, error)
	removeFn func(ctx context.Context, in *pb.TaskId, opts ...grpc.CallOption) (*emptypb.Empty, error)
	listFn   func(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*pb.TaskList, error)
	doneFn   func(ctx context.Context, in *pb.TaskId, opts ...grpc.CallOption) (*pb.TaskExportData, error)

	addCalls    int
	removeCalls int
	listCalls   int
	doneCalls   int

	gotAddCtx  context.Context
	gotAddTask *pb.TaskImportData

	gotRemoveCtx context.Context
	gotRemoveId  *pb.TaskId

	gotListCtx context.Context

	gotDoneCtx context.Context
	gotDoneId  *pb.TaskId
}

func (f *fakeGrpcClient) AddTask(ctx context.Context, in *pb.TaskImportData, opts ...grpc.CallOption) (*pb.TaskExportData, error) {
	f.addCalls++
	f.gotAddCtx = ctx
	f.gotAddTask = in

	if f.addFn == nil {
		panic("AddTask called but addFn not set")
	}

	return f.addFn(ctx, in)
}

func (f *fakeGrpcClient) RemoveTask(ctx context.Context, in *pb.TaskId, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	f.removeCalls++
	f.gotRemoveCtx = ctx
	f.gotRemoveId = in

	if f.removeFn == nil {
		panic("RemoveTask called but removeFn not set")
	}

	return f.removeFn(ctx, in)
}

func (f *fakeGrpcClient) ListAllTasks(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*pb.TaskList, error) {
	f.listCalls++
	f.gotListCtx = ctx

	if f.listFn == nil {
		panic("ListAllTasks called but listFn not set")
	}

	return f.listFn(ctx, in)
}

func (f *fakeGrpcClient) MarkTaskFinished(ctx context.Context, in *pb.TaskId, opts ...grpc.CallOption) (*pb.TaskExportData, error) {
	f.doneCalls++
	f.gotDoneCtx = ctx
	f.gotDoneId = in

	if f.doneFn == nil {
		panic("MarkTaskFinished called but doneFn not set")
	}

	return f.doneFn(ctx, in)
}

func TestAddTask_DelegatesToGrpcClient(t *testing.T) {
	wantTask := &pb.TaskExportData{
		Id:    1,
		Title: "my title",
		Text:  "my text",
	}
	wantErr := errors.New("boom")
	fakeClient := &fakeGrpcClient{
		addFn: func(ctx context.Context, in *pb.TaskImportData, opts ...grpc.CallOption) (*pb.TaskExportData, error) {
			return wantTask, wantErr
		},
	}
	dbClient := NewDBClient(fakeClient)

	gotTask, gotErr := dbClient.AddTask(context.Background(), models.TaskImportData{})

	if !errors.Is(gotErr, wantErr) {
		t.Fatalf("expected err %v, got %v", wantErr, gotErr)
	}
	wantTaskConv := taskExportDataFromPB(wantTask)
	if gotTask.Id != wantTaskConv.Id {
		t.Fatalf("expected Id %d, got %d", wantTaskConv.Id, gotTask.Id)
	}
	if gotTask.Title != wantTaskConv.Title {
		t.Fatalf("expected Title %q, got %q", wantTaskConv.Title, gotTask.Title)
	}
	if gotTask.Text != wantTaskConv.Text {
		t.Fatalf("expected Text %q, got %q", wantTaskConv.Text, gotTask.Text)
	}
}

func TestRemoveTask_DelegatesToGrpcClient(t *testing.T) {
	wantErr := errors.New("boom")
	fakeClient := &fakeGrpcClient{
		removeFn: func(ctx context.Context, in *pb.TaskId, opts ...grpc.CallOption) (*emptypb.Empty, error) {
			return nil, wantErr
		},
	}
	dbClient := NewDBClient(fakeClient)

	gotErr := dbClient.RemoveTask(context.Background(), 1)

	if !errors.Is(gotErr, wantErr) {
		t.Fatalf("expected err %v, got %v", wantErr, gotErr)
	}
}

func TestListAllTasks_DelegatesToGrpcClient(t *testing.T) {
	wantList := &pb.TaskList{
		Tasks: []*pb.TaskExportData{
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
		},
	}
	wantErr := errors.New("boom")
	fakeClient := &fakeGrpcClient{
		listFn: func(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*pb.TaskList, error) {
			return wantList, wantErr
		},
	}
	dbClient := NewDBClient(fakeClient)

	gotTask, gotErr := dbClient.ListAllTasks(context.Background())

	if !errors.Is(gotErr, wantErr) {
		t.Fatalf("expected err %v, got %v", wantErr, gotErr)
	}

	wantListConv := taskSliceFromPB(wantList)
	if len(wantListConv) != 2 {
		t.Fatalf("expected len task list =2, got %d", len(wantListConv))
	}
	if gotTask[0].Id != wantListConv[0].Id ||
		gotTask[1].Id != wantListConv[1].Id {
		t.Fatalf("field Id mismatch")
	}
	if gotTask[0].Title != wantListConv[0].Title ||
		gotTask[1].Title != wantListConv[1].Title {
		t.Fatalf("field Title mismatch")
	}
	if gotTask[0].Text != wantListConv[0].Text ||
		gotTask[1].Text != wantListConv[1].Text {
		t.Fatalf("field Text mismatch")
	}
}

func TestMarkTaskFinished_DelegatesToGrpcClient(t *testing.T) {
	wantTask := &pb.TaskExportData{
		Id:    1,
		Title: "my title",
		Text:  "my text",
	}
	wantErr := errors.New("boom")
	fakeClient := &fakeGrpcClient{
		doneFn: func(ctx context.Context, in *pb.TaskId, opts ...grpc.CallOption) (*pb.TaskExportData, error) {
			return wantTask, wantErr
		},
	}
	dbClient := NewDBClient(fakeClient)

	gotTask, gotErr := dbClient.MarkTaskFinished(context.Background(), 1)

	if !errors.Is(gotErr, wantErr) {
		t.Fatalf("expected err %v, got %v", wantErr, gotErr)
	}
	wantTaskConv := taskExportDataFromPB(wantTask)
	if gotTask.Id != wantTaskConv.Id {
		t.Fatalf("expected Id %d, got %d", wantTaskConv.Id, gotTask.Id)
	}
	if gotTask.Title != wantTaskConv.Title {
		t.Fatalf("expected Title %q, got %q", wantTaskConv.Title, gotTask.Title)
	}
	if gotTask.Text != wantTaskConv.Text {
		t.Fatalf("expected Text %q, got %q", wantTaskConv.Text, gotTask.Text)
	}
}
