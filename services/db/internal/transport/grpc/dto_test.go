package grpc

import (
	"testing"
	"time"

	"github.com/dodocheck/go-pet-project-1/pkg/pb"
	"github.com/dodocheck/go-pet-project-1/services/db/internal/models"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestTaskImportDataFromPB(t *testing.T) {
	tests := []struct {
		name string
		in   *pb.TaskImportData
		want models.TaskImportData
	}{
		{
			name: "regular task",
			in: &pb.TaskImportData{
				Title: "my title",
				Text:  "my text",
			},
			want: models.TaskImportData{
				Title: "my title",
				Text:  "my text",
			},
		},
		{
			name: "nil task",
			in:   nil,
			want: models.TaskImportData{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := taskImportDataFromPB(tt.in)

			if diff := cmp.Diff(got, tt.want, protocmp.Transform()); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestTaskExportDataToPB(t *testing.T) {
	createdAtTS := time.Date(2025, 12, 10, 4, 6, 3, 2, time.UTC)
	finishedAtTS := time.Date(2025, 12, 10, 3, 5, 2, 1, time.UTC)
	tests := []struct {
		name string
		in   models.TaskExportData
		want *pb.TaskExportData
	}{
		{
			name: "unfinished task",
			in: models.TaskExportData{
				Id:         54,
				Title:      "some title",
				Text:       "some text",
				Finished:   false,
				CreatedAt:  createdAtTS,
				FinishedAt: nil,
			},
			want: &pb.TaskExportData{
				Id:         54,
				Title:      "some title",
				Text:       "some text",
				Finished:   false,
				CreatedAt:  timestamppb.New(createdAtTS),
				FinishedAt: nil,
			},
		},
		{
			name: "finished task",
			in: models.TaskExportData{
				Id:         678,
				Title:      "some title2",
				Text:       "some text2",
				Finished:   true,
				CreatedAt:  createdAtTS,
				FinishedAt: &finishedAtTS,
			},
			want: &pb.TaskExportData{
				Id:         678,
				Title:      "some title2",
				Text:       "some text2",
				Finished:   true,
				CreatedAt:  timestamppb.New(createdAtTS),
				FinishedAt: timestamppb.New(finishedAtTS),
			},
		},
		{
			name: "empty task",
			in:   models.TaskExportData{},
			want: &pb.TaskExportData{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := taskExportDataToPB(tt.in)
			if diff := cmp.Diff(got, tt.want, protocmp.Transform()); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestTaskSliceToPB(t *testing.T) {
	createdAtTS := time.Date(2025, 12, 10, 4, 6, 3, 2, time.UTC)
	finishedAtTS := time.Date(2025, 12, 10, 3, 5, 2, 1, time.UTC)
	tests := []struct {
		name string
		in   []models.TaskExportData
		want *pb.TaskList
	}{
		{
			name: "regular list",
			in: []models.TaskExportData{
				{
					Id:         54,
					Title:      "some title",
					Text:       "some text",
					Finished:   false,
					CreatedAt:  createdAtTS,
					FinishedAt: nil,
				},
				{
					Id:         678,
					Title:      "some title2",
					Text:       "some text2",
					Finished:   true,
					CreatedAt:  createdAtTS,
					FinishedAt: &finishedAtTS,
				},
			},
			want: &pb.TaskList{
				Tasks: []*pb.TaskExportData{
					{
						Id:         54,
						Title:      "some title",
						Text:       "some text",
						Finished:   false,
						CreatedAt:  timestamppb.New(createdAtTS),
						FinishedAt: nil,
					},
					{
						Id:         678,
						Title:      "some title2",
						Text:       "some text2",
						Finished:   true,
						CreatedAt:  timestamppb.New(createdAtTS),
						FinishedAt: timestamppb.New(finishedAtTS),
					},
				},
			},
		},
		{
			name: "empty list",
			in:   nil,
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := taskSliceToPB(tt.in)

			if diff := cmp.Diff(got, tt.want, protocmp.Transform()); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestTaskIdFromPB(t *testing.T) {
	tests := []struct {
		name string
		in   *pb.TaskId
		want int
	}{
		{
			name: "regular id",
			in:   &pb.TaskId{Id: 420},
			want: 420,
		},
		{
			name: "nil id",
			in:   nil,
			want: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := taskIdFromPB(tt.in)

			if diff := cmp.Diff(got, tt.want, protocmp.Transform()); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
