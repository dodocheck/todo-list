package grpc

import (
	"context"

	"github.com/dodocheck/go-pet-project-1/services/db/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Server) AddTask(ctx context.Context, task *pb.TaskImportData) (*pb.TaskExportData, error) {
	if task == nil {
		return nil, status.Error(codes.InvalidArgument, "received empty task")
	}

	createdTask, err := s.dbController.AddTask(ctx, *task)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "add task error: %v", err)
	}

	return &createdTask, nil
}

func (s *Server) RemoveTask(ctx context.Context, id *pb.TaskId) (*emptypb.Empty, error) {
	if id == nil {
		return nil, status.Error(codes.InvalidArgument, "received empty id")
	}

	if err := s.dbController.DeleteTask(ctx, int(id.GetId())); err != nil {
		return nil, status.Errorf(codes.Internal, "remove task error: %v", err)
	}

	return nil, nil
}

func (s *Server) ListAllTasks(ctx context.Context, _ *emptypb.Empty) (*pb.TaskList, error) {

}

func (s *Server) MarkTaskFinished(ctx context.Context, id *pb.TaskId) (*pb.TaskExportData, error) {

}
