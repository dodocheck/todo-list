package app

import "github.com/dodocheck/go-pet-project-1/pb"

type Controller interface {
	AddTask(task pb.TaskImportData) (pb.TaskExportData, error)
	DeleteTask(id int) error
	ListAllTasks() ([]pb.TaskExportData, error)
	MarkTaskFinished(id int) (pb.TaskExportData, error)
	Close()
}
