package app

import (
	"context"

	"github.com/dodocheck/go-pet-project-1/services/db/pb"
)

type DBController interface {
	AddTask(ctx context.Context, task pb.TaskImportData) (pb.TaskExportData, error)
	DeleteTask(ctx context.Context, id int) error
	ListAllTasks(ctx context.Context) ([]pb.TaskExportData, error)
	MarkTaskFinished(ctx context.Context, id int) (pb.TaskExportData, error)
}
