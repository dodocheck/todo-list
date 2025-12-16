package dbgrpc

import (
	"context"

	"github.com/dodocheck/go-pet-project-1/services/api/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

type DBClient struct {
	grpcClient pb.TasksServiceClient
}

func NewDBClient(grpcClient pb.TasksServiceClient) *DBClient {
	return &DBClient{grpcClient: grpcClient}
}

func (c *DBClient) AddTask(ctx context.Context, task pb.TaskImportData) (pb.TaskExportData, error) {
	createdTask, err := c.grpcClient.AddTask(ctx, &task)
	if err != nil {
		return pb.TaskExportData{}, err
	}
	return *createdTask, nil
}

func (c *DBClient) RemoveTask(ctx context.Context, id int) error {
	_, err := c.grpcClient.RemoveTask(ctx, &pb.TaskId{Id: int64(id)})
	if err != nil {
		return err
	}
	return nil
}

func (c *DBClient) ListAllTasks(ctx context.Context) ([]pb.TaskExportData, error) {
	taskList, err := c.grpcClient.ListAllTasks(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, err
	}

	taskListToReturn := []pb.TaskExportData{}
	for _, v := range taskList.GetTasks() {
		taskListToReturn = append(taskListToReturn, *v)
	}

	return taskListToReturn, nil
}

func (c *DBClient) MarkTaskFinished(ctx context.Context, id int) (pb.TaskExportData, error) {
	updatedTask, err := c.grpcClient.MarkTaskFinished(ctx, &pb.TaskId{Id: int64(id)})
	if err != nil {
		return pb.TaskExportData{}, err
	}
	return *updatedTask, nil
}
