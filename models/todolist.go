package todolist

import (
	"math/rand/v2"
	"sync"
)

type ToDoList struct {
	tasks map[int]Task
	mtx   sync.RWMutex
}

func NewToDoList() ToDoList {
	return ToDoList{
		tasks: make(map[int]Task),
		mtx:   sync.RWMutex{},
	}
}

func (l *ToDoList) AddTask(newTask Task) (int, error) {
	l.mtx.Lock()
	defer l.mtx.Unlock()

	for _, v := range l.tasks {
		if v.title == newTask.title {
			return 0, ErrorTaskAlreadyExists
		}
	}

	newTaskId := rand.Int()
	l.tasks[newTaskId] = newTask
	return newTaskId, nil
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
