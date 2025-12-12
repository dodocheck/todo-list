package app

import (
	"context"
	"pet1/pkg/contracts"
)

type TaskStore interface {
	AddTask(ctx context.Context, task contracts.TaskImportData) (contracts.TaskExportData, error)
	RemoveTask(ctx context.Context, id int) error
	ListAllTasks(ctx context.Context) ([]contracts.TaskExportData, error)
	MarkTaskFinished(ctx context.Context, id int) (contracts.TaskExportData, error)
}
