package grpc

import (
	"github.com/dodocheck/go-pet-project-1/pkg/pb"
	"github.com/dodocheck/go-pet-project-1/services/db/internal/models"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func taskImportDataFromPB(task *pb.TaskImportData) models.TaskImportData {
	if task == nil {
		return models.TaskImportData{}
	}

	return models.TaskImportData{
		Title: task.GetTitle(),
		Text:  task.GetText(),
	}
}

func taskExportDataToPB(task models.TaskExportData) *pb.TaskExportData {
	out := &pb.TaskExportData{
		Id:       int64(task.Id),
		Title:    task.Title,
		Text:     task.Text,
		Finished: task.Finished,
	}

	if !task.CreatedAt.IsZero() {
		out.CreatedAt = timestamppb.New(task.CreatedAt)
	}
	if task.FinishedAt != nil && !task.FinishedAt.IsZero() {
		out.FinishedAt = timestamppb.New(*task.FinishedAt)
	}

	return out
}

func taskSliceToPB(tasks []models.TaskExportData) *pb.TaskList {
	if tasks == nil {
		return nil
	}

	taskList := &pb.TaskList{}
	for _, v := range tasks {
		taskList.Tasks = append(taskList.Tasks, taskExportDataToPB(v))
	}
	return taskList
}

func taskIdFromPB(id *pb.TaskId) int {
	if id == nil {
		return -1
	}

	return int(id.GetId())
}
