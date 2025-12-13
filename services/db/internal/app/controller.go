package db

import "github.com/dodocheck/go-pet-project-1/shared/contracts"

type Controller interface {
	AddTask(task contracts.TaskImportData) (contracts.TaskExportData, error)
	DeleteTask(id int) error
	ListAllTasks() ([]contracts.TaskExportData, error)
	MarkTaskFinished(id int) (contracts.TaskExportData, error)
	Close()
}
