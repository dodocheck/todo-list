package app

import (
	"context"

	"github.com/dodocheck/go-pet-project-1/pb"
)

type TaskStore interface {
	AddTask(ctx context.Context, task pb.TaskImportData) (pb.TaskExportData, error)
	RemoveTask(ctx context.Context, id int) error
	ListAllTasks(ctx context.Context) ([]pb.TaskExportData, error)
	MarkTaskFinished(ctx context.Context, id int) (pb.TaskExportData, error)
}
