package todolist

import (
	"time"
)

type Task struct {
	title      string
	text       string
	finished   bool
	createdAt  time.Time
	finishedAt *time.Time
}

func NewTask(title string, text string) *Task {
	return &Task{
		title:      title,
		text:       text,
		finished:   false,
		createdAt:  time.Now(),
		finishedAt: nil}
}

func (t *Task) MarkFinished() {
	t.finished = true
	finishedAt := time.Now()
	t.finishedAt = &finishedAt
}
