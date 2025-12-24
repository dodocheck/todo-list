package dbgrpc

import (
	"testing"
	"time"

	"github.com/dodocheck/go-pet-project-1/pkg/pb"
	"github.com/dodocheck/go-pet-project-1/services/api/internal/models"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestTaskImportDataToPB(t *testing.T) {
	tests := []struct {
		name string
		in   models.TaskImportData
		want *pb.TaskImportData
	}{
		{
			name: "regular task",
			in: models.TaskImportData{
				Title: "Buy milk",
				Text:  "2 bottles",
			},
			want: &pb.TaskImportData{
				Title: "Buy milk",
				Text:  "2 bottles",
			},
		},
		{
			name: "empty task",
			in:   models.TaskImportData{},
			want: &pb.TaskImportData{
				Title: "",
				Text:  "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := taskImportDataToPB(tt.in)
			if got == nil {
				t.Fatalf("got nil, want non-nil")
			}
			if got.Title != tt.want.Title {
				t.Fatalf("expected title %q, got %q", tt.want.Title, got.Title)
			}
			if got.Text != tt.want.Text {
				t.Fatalf("expected text %q, got %q", tt.want.Text, got.Text)
			}
		})
	}
}

func TestTaskExportDataFromPB(t *testing.T) {
	createdAtTS := time.Date(2025, 12, 11, 1, 2, 3, 4, time.UTC)
	finishedAtTS := time.Date(2025, 11, 22, 2, 3, 4, 5, time.UTC)

	tests := []struct {
		name string
		in   *pb.TaskExportData
		want models.TaskExportData
	}{
		{
			name: "unfinished task",
			in: &pb.TaskExportData{
				Id:         1,
				Title:      "Read book",
				Text:       "And learn it",
				Finished:   false,
				CreatedAt:  timestamppb.New(createdAtTS),
				FinishedAt: nil,
			},
			want: models.TaskExportData{
				Id:         1,
				Title:      "Read book",
				Text:       "And learn it",
				Finished:   false,
				CreatedAt:  createdAtTS,
				FinishedAt: nil,
			},
		},
		{
			name: "finished task",
			in: &pb.TaskExportData{
				Id:         2,
				Title:      "Learn Golang",
				Text:       "Do pet project",
				Finished:   true,
				CreatedAt:  timestamppb.New(createdAtTS),
				FinishedAt: timestamppb.New(finishedAtTS),
			},
			want: models.TaskExportData{
				Id:         2,
				Title:      "Learn Golang",
				Text:       "Do pet project",
				Finished:   true,
				CreatedAt:  createdAtTS,
				FinishedAt: &finishedAtTS,
			},
		},
		{
			name: "nil task",
			in:   nil,
			want: models.TaskExportData{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := taskExportDataFromPB(tt.in)

			if tt.in != nil {
				if got.Id != tt.want.Id {
					t.Fatalf("expected Id %d, got %d", tt.want.Id, got.Id)
				}
				if got.Title != tt.want.Title {
					t.Fatalf("expected Title %q, got %q", tt.want.Title, got.Title)
				}
				if got.Text != tt.want.Text {
					t.Fatalf("expected Text %q, got %q", tt.want.Text, got.Text)
				}
				if got.Finished != tt.want.Finished {
					t.Fatalf("expected Finished %t, got %t", tt.want.Finished, got.Finished)
				}
				if !got.CreatedAt.Equal(tt.want.CreatedAt) {
					t.Fatalf("expected CreatedAt %v, got %v", tt.want.CreatedAt, got.CreatedAt)
				}
				if got.FinishedAt != nil && tt.want.FinishedAt != nil &&
					!got.FinishedAt.Equal(*tt.want.FinishedAt) {
					t.Fatalf("expected FinishedAt %v, got %v", tt.want.FinishedAt, got.FinishedAt)
				}
			}
		})
	}
}

func TestTaskSliceFromPB(t *testing.T) {
	createdAtTS := time.Date(2025, 12, 11, 1, 2, 3, 4, time.UTC)
	finishedAtTS := time.Date(2025, 11, 22, 2, 3, 4, 5, time.UTC)

	tests := []struct {
		name string
		in   *pb.TaskList
		want []models.TaskExportData
	}{
		{
			name: "regular list",
			in: &pb.TaskList{
				Tasks: []*pb.TaskExportData{
					{
						Id:         1,
						Title:      "title1",
						Text:       "text1",
						Finished:   false,
						CreatedAt:  timestamppb.New(createdAtTS),
						FinishedAt: nil,
					},
					{
						Id:         2,
						Title:      "title2",
						Text:       "text2",
						Finished:   true,
						CreatedAt:  timestamppb.New(createdAtTS),
						FinishedAt: timestamppb.New(finishedAtTS),
					},
				},
			},
			want: []models.TaskExportData{
				{
					Id:         1,
					Title:      "title1",
					Text:       "text1",
					Finished:   false,
					CreatedAt:  createdAtTS,
					FinishedAt: &finishedAtTS,
				},
				{
					Id:         2,
					Title:      "title2",
					Text:       "text2",
					Finished:   true,
					CreatedAt:  createdAtTS,
					FinishedAt: &finishedAtTS,
				},
			},
		},
		{
			name: "nil list",
			in:   nil,
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := taskSliceFromPB(tt.in)
			if got == nil && tt.in != nil {
				t.Fatalf("expected non-nil, got nil")
			}
			if got != nil && tt.in == nil {
				t.Fatalf("expected nil, got %+v", got)
			}
			if got != nil && tt.in != nil {
				if len(got) != 2 {
					t.Fatalf("expected len=2, got %d", len(got))
				}
				if got[0].Id != tt.want[0].Id || got[1].Id != tt.want[1].Id {
					t.Fatalf("field Id mismatch")
				}
				if got[0].Title != tt.want[0].Title || got[1].Title != tt.want[1].Title {
					t.Fatalf("field Title mismatch")
				}
				if got[0].Text != tt.want[0].Text || got[1].Text != tt.want[1].Text {
					t.Fatalf("field Text mismatch")
				}
				if got[0].Finished != tt.want[0].Finished || got[1].Finished != tt.want[1].Finished {
					t.Fatalf("field Finished mismatch")
				}
				if got[0].CreatedAt != tt.want[0].CreatedAt || got[1].CreatedAt != tt.want[1].CreatedAt {
					t.Fatalf("field CreatedAt mismatch")
				}
				if got[0].FinishedAt != nil && tt.want[0].FinishedAt != nil &&
					!got[0].FinishedAt.Equal(*tt.want[0].FinishedAt) {
					t.Fatalf("field FinishedAt mismatch")
				}
				if got[1].FinishedAt != nil && tt.want[1].FinishedAt != nil &&
					!got[1].FinishedAt.Equal(*tt.want[1].FinishedAt) {
					t.Fatalf("field FinishedAt mismatch")
				}
			}
		})
	}
}

func TestTaskIdToPB(t *testing.T) {
	tests := []struct {
		name string
		in   int
		want *pb.TaskId
	}{
		{
			name: "regular id",
			in:   123,
			want: &pb.TaskId{Id: 123},
		},
		{
			name: "bad id",
			in:   -342,
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := taskIdToPB(tt.in)
			if tt.want == nil && got != nil {
				t.Fatalf("expected nil, got %+v", got)
			}
			if got != nil && tt.want != nil && got.Id != tt.want.Id {
				t.Fatalf("expected Id %d, got %d", tt.want.Id, got.Id)
			}

		})
	}
}
