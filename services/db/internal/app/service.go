package app

import (
	"context"

	"github.com/dodocheck/go-pet-project-1/services/db/pb"
)

type Service struct {
	dbController DBController
}

func NewService(dbController DBController) *Service {
	return &Service{
		dbController: dbController}
}

func (s *Service) AddTask(ctx context.Context, task pb.TaskImportData) (pb.TaskExportData, error) {
	return s.dbController.AddTask(ctx, task)
}

func (s *Service) DeleteTask(ctx context.Context, id int) error {
	return s.dbController.DeleteTask(ctx, id)
}

func (s *Service) ListAllTasks(ctx context.Context) ([]pb.TaskExportData, error) {
	return s.dbController.ListAllTasks(ctx)
}

func (s *Service) MarkTaskFinished(ctx context.Context, id int) (pb.TaskExportData, error) {
	return s.dbController.MarkTaskFinished(ctx, id)
}
