package dbgrpc

import (
	"github.com/dodocheck/go-pet-project-1/services/api/internal/models"
	"github.com/dodocheck/go-pet-project-1/services/api/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func taskImportDataToPB(task models.TaskImportData) *pb.TaskImportData {
	return &pb.TaskImportData{
		Title: task.Title,
		Text:  task.Text,
	}
}

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
		Id:        int64(task.Id),
		Title:     task.Title,
		Text:      task.Text,
		Finished:  task.Finished,
		CreatedAt: timestamppb.New(task.CreatedAt),
	}

	if task.FinishedAt != nil {
		out.FinishedAt = timestamppb.New(*task.FinishedAt)
	}

	return out
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

func taskSliceToPB(tasks []models.TaskExportData) *pb.TaskList {
	taskList := &pb.TaskList{}
	for _, v := range tasks {
		taskList.Tasks = append(taskList.Tasks, taskExportDataToPB(v))
	}
	return taskList
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

func taskIdFromPB(id *pb.TaskId) int {
	if id == nil {
		return -1
	}

	return int(id.GetId())
}
