package db

import "pet1/pkg/contracts"

type Controller interface {
	AddTask(task contracts.TaskImportData) (contracts.TaskExportData, error)
	DeleteTask(id int) error
	ListAllTasks() ([]contracts.TaskExportData, error)
	MarkTaskFinished(id int) (contracts.TaskExportData, error)
	Close()
}
