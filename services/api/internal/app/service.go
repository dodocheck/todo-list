package app

import (
	"context"

	"github.com/dodocheck/go-pet-project-1/services/api/internal/models"
)

type Service struct {
	store TaskStore
}

func NewService(store TaskStore) *Service {
	return &Service{
		store: store}
}

func (s *Service) AddTask(ctx context.Context, task models.TaskImportData) (models.TaskExportData, error) {
	return s.store.AddTask(ctx, task)
}
func (s *Service) RemoveTask(ctx context.Context, id int) error {
	return s.store.RemoveTask(ctx, id)
}

func (s *Service) ListAllTasks(ctx context.Context) ([]models.TaskExportData, error) {
	return s.store.ListAllTasks(ctx)
}

func (s *Service) MarkTaskFinished(ctx context.Context, id int) (models.TaskExportData, error) {
	return s.store.MarkTaskFinished(ctx, id)
}
