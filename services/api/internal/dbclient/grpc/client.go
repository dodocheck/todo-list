package dbgrpc

import (
	"context"

	"github.com/dodocheck/go-pet-project-1/pkg/pb"
	"github.com/dodocheck/go-pet-project-1/services/api/internal/models"
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
	return taskExportDataFromPB(createdTask), err
}

func (c *DBClient) RemoveTask(ctx context.Context, id int) error {
	_, err := c.grpcClient.RemoveTask(ctx, taskIdToPB(id))
	return err
}

func (c *DBClient) ListAllTasks(ctx context.Context) ([]models.TaskExportData, error) {
	taskList, err := c.grpcClient.ListAllTasks(ctx, &emptypb.Empty{})
	return taskSliceFromPB(taskList), err
}

func (c *DBClient) MarkTaskFinished(ctx context.Context, id int) (models.TaskExportData, error) {
	updatedTask, err := c.grpcClient.MarkTaskFinished(ctx, taskIdToPB(id))
	return taskExportDataFromPB(updatedTask), err
}
