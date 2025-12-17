package app

import (
	"context"

	"github.com/dodocheck/go-pet-project-1/services/db/internal/models"
)

type Service struct {
	dbController DBController
}

func NewService(dbController DBController) *Service {
	return &Service{
		dbController: dbController}
}

func (s *Service) AddTask(ctx context.Context, task models.TaskImportData) (models.TaskExportData, error) {
	return s.dbController.AddTask(ctx, task)
}

func (s *Service) DeleteTask(ctx context.Context, id int) error {
	return s.dbController.DeleteTask(ctx, id)
}

func (s *Service) ListAllTasks(ctx context.Context) ([]models.TaskExportData, error) {
	return s.dbController.ListAllTasks(ctx)
}

func (s *Service) MarkTaskFinished(ctx context.Context, id int) (models.TaskExportData, error) {
	return s.dbController.MarkTaskFinished(ctx, id)
}
