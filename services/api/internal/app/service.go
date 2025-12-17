package app

import (
	"context"

	"github.com/dodocheck/go-pet-project-1/services/api/internal/models"
)

type Service struct {
	dbClient DBClient
}

func NewService(dbClient DBClient) *Service {
	return &Service{
		dbClient: dbClient}
}

func (s *Service) AddTask(ctx context.Context, task models.TaskImportData) (models.TaskExportData, error) {
	return s.dbClient.AddTask(ctx, task)
}

func (s *Service) RemoveTask(ctx context.Context, id int) error {
	return s.dbClient.RemoveTask(ctx, id)
}

func (s *Service) ListAllTasks(ctx context.Context) ([]models.TaskExportData, error) {
	return s.dbClient.ListAllTasks(ctx)
}

func (s *Service) MarkTaskFinished(ctx context.Context, id int) (models.TaskExportData, error) {
	return s.dbClient.MarkTaskFinished(ctx, id)
}
