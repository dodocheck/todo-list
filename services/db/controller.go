package db

import "pet1/models"

type Controller interface {
	AddTask(task models.TaskImportData) (models.TaskExportData, error)
	DeleteTask(id int) error
	ListAllTasks() ([]models.TaskExportData, error)
	MarkTaskFinished(id int) (models.TaskExportData, error)
	Close()
}
