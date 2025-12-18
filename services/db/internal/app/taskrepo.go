package app

import (
	"context"

	"github.com/dodocheck/go-pet-project-1/services/db/internal/models"
)

type TaskRepository interface {
	AddTask(ctx context.Context, task models.TaskImportData) (models.TaskExportData, error)
	DeleteTask(ctx context.Context, id int) error
	ListAllTasks(ctx context.Context) ([]models.TaskExportData, error)
	MarkTaskFinished(ctx context.Context, id int) (models.TaskExportData, error)
	Close() error
}
