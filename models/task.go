package models

import (
	"time"
)

type Task struct {
	Title      string
	Text       string
	finished   bool
	createdAt  time.Time
	finishedAt *time.Time
}

func NewTask(title string, text string) *Task {
	return &Task{
		Title:      title,
		Text:       text,
		finished:   false,
		createdAt:  time.Now(),
		finishedAt: nil}
}

func (t *Task) MarkFinished() {
	t.finished = true
	finishedAt := time.Now()
	t.finishedAt = &finishedAt
}
