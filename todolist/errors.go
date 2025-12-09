package todolist

import "errors"

var (
	ErrorTaskAlreadyExists = errors.New("task already exists")
	ErrorTaskNotFound      = errors.New("task not found")
)
