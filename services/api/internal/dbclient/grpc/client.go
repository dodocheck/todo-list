package dbgrpc

import (
	"context"

	"github.com/dodocheck/go-pet-project-1/services/api/internal/models"
	"github.com/dodocheck/go-pet-project-1/services/api/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

type DBClient struct {
	grpcClient pb.TasksServiceClient
}

func NewDBClient(grpcClient pb.TasksServiceClient) *DBClient {
	return &DBClient{grpcClient: grpcClient}
}

func (c *DBClient) AddTask(ctx context.Context, task models.TaskImportData) (models.TaskExportData, error) {
	createdTask, err := c.grpcClient.AddTask(ctx, taskImportDataToPB(task))
	if err != nil {
		return models.TaskExportData{}, err
	}
	return taskExportDataFromPB(createdTask), nil
}

func (c *DBClient) RemoveTask(ctx context.Context, id int) error {
	_, err := c.grpcClient.RemoveTask(ctx, taskIdToPB(id))
	if err != nil {
		return err
	}
	return nil
}

func (c *DBClient) ListAllTasks(ctx context.Context) ([]models.TaskExportData, error) {
	taskList, err := c.grpcClient.ListAllTasks(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, err
	}

	return taskSliceFromPB(taskList), nil
}

func (c *DBClient) MarkTaskFinished(ctx context.Context, id int) (models.TaskExportData, error) {
	updatedTask, err := c.grpcClient.MarkTaskFinished(ctx, taskIdToPB(id))
	if err != nil {
		return models.TaskExportData{}, err
	}
	return taskExportDataFromPB(updatedTask), nil
}
