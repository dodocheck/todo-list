package models

import (
	"time"
)

type TaskExportData struct {
	Id         int
	Title      string
	Text       string
	Finished   bool
	CreatedAt  time.Time
	FinishedAt *time.Time
}

type TaskImportData struct {
	Title string
	Text  string
}
