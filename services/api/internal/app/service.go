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
	log.Printf("IN: add task: %+v\n", task)

	actionLog := logger.CreateTaskAddedLog()

	createdTask, err := s.dbClient.AddTask(ctx, task)

	if err == nil {
		s.logAction(actionLog)
		log.Printf("OUT(OK): add task: %+v\n", createdTask)
	} else {
		log.Printf("OUT(ERR): add task: %v\n", err)
	}

	return createdTask, err

}

func (s *Service) RemoveTask(ctx context.Context, id int) error {
	log.Printf("IN: remove task with ID: %v\n", id)

	actionLog := logger.CreateTaskDeletedLog()

	err := s.dbClient.RemoveTask(ctx, id)

	if err == nil {
		s.logAction(actionLog)
		log.Printf("OUT(OK): remove task with ID %v\n", id)
	} else {
		log.Printf("OUT(ERR): remove task with ID %v: %v\n", id, err)
	}

	return err
}

func (s *Service) ListAllTasks(ctx context.Context) ([]models.TaskExportData, error) {
	log.Println("IN: list tasks")

	actionLog := logger.CreateListTasksLog()

	tasks, err := s.dbClient.ListAllTasks(ctx)

	if err == nil {
		s.logAction(actionLog)
		log.Printf("OUT(OK): list tasks: %+v\n", tasks)
	} else {
		log.Printf("OUT(ERR): list tasks: %v\n", err)
	}

	return tasks, err
}

func (s *Service) MarkTaskFinished(ctx context.Context, id int) (models.TaskExportData, error) {
	log.Printf("IN: finish task with ID: %v\n", id)

	actionLog := logger.CreateTaskDoneLog()

	updatedTask, err := s.dbClient.MarkTaskFinished(ctx, id)

	if err == nil {
		s.logAction(actionLog)
		log.Printf("OUT(OK): finish task with ID %v\n", id)
	} else {
		log.Printf("OUT(ERR): finish task with ID %v: %v\n", id, err)
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
