package app

import "github.com/dodocheck/go-pet-project-1/services/db/internal/models"

type Controller interface {
	AddTask(task models.TaskImportData) (models.TaskExportData, error)
	DeleteTask(id int) error
	ListAllTasks() ([]models.TaskExportData, error)
	MarkTaskFinished(id int) (models.TaskExportData, error)
	Close()
}
