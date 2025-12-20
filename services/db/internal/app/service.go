package app

import (
	"context"
	"log"

	"github.com/dodocheck/go-pet-project-1/services/db/internal/models"
)

type Service struct {
	dbController TaskRepository
}

func NewService(dbController TaskRepository) *Service {
	return &Service{
		dbController: dbController}
}

func (s *Service) AddTask(ctx context.Context, task models.TaskImportData) (models.TaskExportData, error) {
	log.Printf("IN: add task: %+v\n", task)

	createdTask, err := s.dbController.AddTask(ctx, task)

	if err != nil {
		log.Printf("OUT(ERR): add task: %v\n", err)
	} else {
		log.Printf("OUT(OK): add task: %+v\n", createdTask)
	}

	return createdTask, err
}

func (s *Service) DeleteTask(ctx context.Context, id int) error {
	log.Printf("IN: delete task with ID: %v\n", id)

	err := s.dbController.DeleteTask(ctx, id)

	if err != nil {
		log.Printf("OUT(ERR): delete task with ID %v: %v\n", id, err)
	} else {
		log.Printf("OUT(OK): delete task with ID %v\n", id)
	}

	return err
}

func (s *Service) ListAllTasks(ctx context.Context) ([]models.TaskExportData, error) {
	log.Println("IN: list tasks")

	tasks, err := s.dbController.ListAllTasks(ctx)

	if err != nil {
		log.Printf("OUT(ERR): list tasks: %v\n", err)
	} else {
		log.Printf("OUT(OK): list tasks: %+v\n", tasks)
	}

	return tasks, err
}

func (s *Service) MarkTaskFinished(ctx context.Context, id int) (models.TaskExportData, error) {
	log.Printf("IN: finish task with ID: %v\n", id)

	updatedTask, err := s.dbController.MarkTaskFinished(ctx, id)

	if err != nil {
		log.Printf("OUT(ERR): finish task with ID %v: %v\n", id, err)
	} else {
		log.Printf("OUT(OK): finish task with ID %v\n", id)
	}

	return updatedTask, err
}
