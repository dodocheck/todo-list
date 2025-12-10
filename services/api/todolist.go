package api

import (
	"database/sql"
	"pet1/models"
	"pet1/services/db/postgres"
)

func AddTask(db *sql.DB, newTask models.Task) int64 {
	taskInsertData := postgres.TaskInsertData{
		Title: newTask.Title,
		Text:  newTask.Text}

	return postgres.InsertTask(db, taskInsertData)
}

func (l *ToDoList) RemoveTask(taskId int) error {
	l.mtx.Lock()
	defer l.mtx.Unlock()

	if _, ok := l.tasks[taskId]; !ok {
		return ErrorTaskNotFound
	}

	delete(l.tasks, taskId)
	return nil
}

func (l *ToDoList) GetTaskById(taskId int) (Task, error) {
	l.mtx.RLock()
	defer l.mtx.RUnlock()

	foundedTask, ok := l.tasks[taskId]
	if !ok {
		return Task{}, ErrorTaskNotFound
	}

	return foundedTask, nil
}

func (l *ToDoList) GetAllTasks() map[int]Task {
	l.mtx.RLock()
	defer l.mtx.RUnlock()

	tempMap := make(map[int]Task)
	for k, v := range l.tasks {
		tempMap[k] = v
	}

	return tempMap
}

func (l *ToDoList) MarkTaskFinished(taskId int) (Task, error) {
	l.mtx.Lock()
	defer l.mtx.Unlock()

	taskToFinish, ok := l.tasks[taskId]
	if !ok {
		return Task{}, ErrorTaskNotFound
	}

	taskToFinish.MarkFinished()
	l.tasks[taskId] = taskToFinish

	return taskToFinish, nil
}
