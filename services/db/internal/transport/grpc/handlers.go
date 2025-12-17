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

	taskFromPB := taskImportDataFromPB(task)

	createdTask, err := s.dbController.AddTask(ctx, taskFromPB)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "add task error: %v", err)
	}

	taskToPB := taskExportDataToPB(createdTask)

	return taskToPB, nil
}

func (s *Server) RemoveTask(ctx context.Context, id *pb.TaskId) (*emptypb.Empty, error) {
	if id == nil {
		return nil, status.Error(codes.InvalidArgument, "received empty id")
	}

	if err := s.dbController.DeleteTask(ctx, taskIdFromPB(id)); err != nil {
		return nil, status.Errorf(codes.Internal, "remove task error: %v", err)
	}

	return nil, nil
}

func (s *Server) ListAllTasks(ctx context.Context, _ *emptypb.Empty) (*pb.TaskList, error) {
	allTasks, err := s.dbController.ListAllTasks(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list tasks error: %v", err)
	}

	return taskSliceToPB(allTasks), nil
}

func (s *Server) MarkTaskFinished(ctx context.Context, id *pb.TaskId) (*pb.TaskExportData, error) {
	updatedTask, err := s.dbController.MarkTaskFinished(ctx, taskIdFromPB(id))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "finish task error: %v", err)
	}
	return taskExportDataToPB(updatedTask), nil
}
