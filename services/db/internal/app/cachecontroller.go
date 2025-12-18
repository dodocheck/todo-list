package app

import (
	"context"

	"github.com/dodocheck/go-pet-project-1/services/db/internal/models"
)

type CacheController interface {
	CacheTaskList(ctx context.Context, tasks []models.TaskExportData) error
	DeleteTaskList(ctx context.Context) error
	GetTaskList(ctx context.Context) ([]models.TaskExportData, error)
	Close() error
}
