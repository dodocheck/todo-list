package app

import (
	"context"
	"pet1/pkg/contracts"
)

type Service struct {
	store TaskStore
}

func NewService(store TaskStore) *Service {
	return &Service{
		store: store}
}

func (s *Service) AddTask(ctx context.Context, task contracts.TaskImportData) (contracts.TaskExportData, error) {
	return s.store.AddTask(ctx, task)
}
func (s *Service) RemoveTask(ctx context.Context, id int) error {
	return s.store.RemoveTask(ctx, id)
}

func (s *Service) ListAllTasks(ctx context.Context) ([]contracts.TaskExportData, error) {
	return s.store.ListAllTasks(ctx)
}

func (s *Service) MarkTaskFinished(ctx context.Context, id int) (contracts.TaskExportData, error) {
	return s.store.MarkTaskFinished(ctx, id)
}
