package dbgrpc

import (
	"github.com/dodocheck/go-pet-project-1/pkg/pb"
	"github.com/dodocheck/go-pet-project-1/services/api/internal/models"
)

func taskImportDataToPB(task models.TaskImportData) *pb.TaskImportData {
	return &pb.TaskImportData{
		Title: task.Title,
		Text:  task.Text,
	}
}

func taskExportDataFromPB(task *pb.TaskExportData) models.TaskExportData {
	if task == nil {
		return models.TaskExportData{}
	}

	out := models.TaskExportData{
		Id:       int(task.GetId()),
		Title:    task.GetTitle(),
		Text:     task.GetText(),
		Finished: task.GetFinished(),
	}

	if task.GetCreatedAt() != nil {
		out.CreatedAt = task.CreatedAt.AsTime()
	}

	if task.GetFinishedAt() != nil {
		ts := task.GetFinishedAt().AsTime()
		out.FinishedAt = &ts
	}

	return out
}

func taskSliceFromPB(tasks *pb.TaskList) []models.TaskExportData {
	if tasks == nil {
		return nil
	}

	taskSlice := make([]models.TaskExportData, 0, len(tasks.Tasks))
	for _, v := range tasks.GetTasks() {
		taskSlice = append(taskSlice, taskExportDataFromPB(v))
	}
	return taskSlice
}

func taskIdToPB(id int) *pb.TaskId {
	if id < 0 {
		return nil
	}

	return &pb.TaskId{
		Id: int64(id),
	}
}
