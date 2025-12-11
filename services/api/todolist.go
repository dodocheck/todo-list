package api

import (
	"pet1/models"
	"pet1/services/db"
)

type ToDoList struct {
	dbController db.Controller
}

func NewToDoList(controller db.Controller) *ToDoList {
	return &ToDoList{dbController: controller}
}

func (l *ToDoList) AddTask(task models.TaskImportData) (models.TaskExportData, error) {
	return l.dbController.AddTask(task)
}

func (l *ToDoList) RemoveTask(id int) error {
	return l.dbController.DeleteTask(id)
}

func (l *ToDoList) ListAllTasks() ([]models.TaskExportData, error) {
	return l.dbController.ListAllTasks()
}

func (l *ToDoList) MarkTaskFinished(id int) (models.TaskExportData, error) {
	return l.dbController.MarkTaskFinished(id)
}
