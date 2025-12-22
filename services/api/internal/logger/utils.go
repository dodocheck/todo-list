package logger

import (
	"time"

	"github.com/dodocheck/go-pet-project-1/services/api/internal/models"
)

func CreateTaskAddedLog() models.ActionLog {
	return models.ActionLog{
		Action: "task created",
		Time:   time.Now(),
	}
}

func CreateTaskDeletedLog() models.ActionLog {
	return models.ActionLog{
		Action: "task deleted",
		Time:   time.Now(),
	}
}

func CreateTaskDoneLog() models.ActionLog {
	return models.ActionLog{
		Action: "task done",
		Time:   time.Now(),
	}
}

func CreateListTasksLog() models.ActionLog {
	return models.ActionLog{
		Action: "list tasks",
		Time:   time.Now(),
	}
}
