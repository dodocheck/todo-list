package app

import (
	"context"
	"log"

	"github.com/dodocheck/go-pet-project-1/services/api/internal/logger"
	"github.com/dodocheck/go-pet-project-1/services/api/internal/models"
)

type Service struct {
	dbClient   DBClient
	logChannel chan models.ActionLog
}

func NewService(dbClient DBClient) *Service {
	return &Service{
		dbClient:   dbClient,
		logChannel: make(chan models.ActionLog, 100)}
}

func (s *Service) GetLogChannel() <-chan models.ActionLog {
	return s.logChannel
}

func (s *Service) AddTask(ctx context.Context, task models.TaskImportData) (models.TaskExportData, error) {
	actionLog := logger.CreateTaskAddedLog()

	createdTask, err := s.dbClient.AddTask(ctx, task)

	if err == nil {
		s.logAction(actionLog)
	}

	return createdTask, err

}

func (s *Service) RemoveTask(ctx context.Context, id int) error {
	actionLog := logger.CreateTaskDeletedLog()

	err := s.dbClient.RemoveTask(ctx, id)

	if err == nil {
		s.logAction(actionLog)
	}

	return err
}

func (s *Service) ListAllTasks(ctx context.Context) ([]models.TaskExportData, error) {
	actionLog := logger.CreateListTasksLog()

	tasks, err := s.dbClient.ListAllTasks(ctx)

	if err == nil {
		s.logAction(actionLog)
	}

	return tasks, err
}

func (s *Service) MarkTaskFinished(ctx context.Context, id int) (models.TaskExportData, error) {
	actionLog := logger.CreateTaskDoneLog()

	updatedTask, err := s.dbClient.MarkTaskFinished(ctx, id)

	if err == nil {
		s.logAction(actionLog)
	}

	return updatedTask, err
}

func (s *Service) logAction(actionLog models.ActionLog) {
	select {
	case s.logChannel <- actionLog:
	default:
		log.Println("dropped user action log - channel queue is full")
	}
}
