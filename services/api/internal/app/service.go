package app

import (
	"context"

	"github.com/dodocheck/go-pet-project-1/pb"
)

type Service struct {
	store TaskStore
}

func NewService(store TaskStore) *Service {
	return &Service{
		store: store}
}

func (s *Service) AddTask(ctx context.Context, task pb.TaskImportData) (pb.TaskExportData, error) {
	return s.store.AddTask(ctx, task)
}
func (s *Service) RemoveTask(ctx context.Context, id int) error {
	return s.store.RemoveTask(ctx, id)
}

func (s *Service) ListAllTasks(ctx context.Context) ([]pb.TaskExportData, error) {
	return s.store.ListAllTasks(ctx)
}

func (s *Service) MarkTaskFinished(ctx context.Context, id int) (pb.TaskExportData, error) {
	return s.store.MarkTaskFinished(ctx, id)
}
