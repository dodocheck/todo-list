package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/dodocheck/go-pet-project-1/services/api/internal/app"
	"github.com/dodocheck/go-pet-project-1/services/api/internal/models"
)

type fakeDBClient struct {
	addFn    func(ctx context.Context, task models.TaskImportData) (models.TaskExportData, error)
	removeFn func(ctx context.Context, id int) error
	listFn   func(ctx context.Context) ([]models.TaskExportData, error)
	doneFn   func(ctx context.Context, id int) (models.TaskExportData, error)

	addCalls    int
	removeCalls int
	listCalls   int
	doneCalls   int

	gotAddTask models.TaskImportData
	gotAddCtx  context.Context

	gotRemoveID int
	gotDoneID   int
}

func (f *fakeDBClient) AddTask(ctx context.Context, task models.TaskImportData) (models.TaskExportData, error) {
	f.addCalls++
	f.gotAddCtx = ctx
	f.gotAddTask = task
	if f.addFn == nil {
		panic("AddTask called but addFn not set")
	}
	return f.addFn(ctx, task)
}

func (f *fakeDBClient) RemoveTask(ctx context.Context, id int) error {
	f.removeCalls++
	f.gotRemoveID = id
	if f.removeFn == nil {
		panic("RemoveTask called but removeFn not set")
	}
	return f.removeFn(ctx, id)
}

func (f *fakeDBClient) ListAllTasks(ctx context.Context) ([]models.TaskExportData, error) {
	f.listCalls++
	if f.listFn == nil {
		panic("ListAllTasks called but listFn not set")
	}
	return f.listFn(ctx)
}

func (f *fakeDBClient) MarkTaskFinished(ctx context.Context, id int) (models.TaskExportData, error) {
	f.doneCalls++
	f.gotDoneID = id
	if f.doneFn == nil {
		panic("MarkTaskFinished called but doneFn not set")
	}
	return f.doneFn(ctx, id)
}

func TestHandleAddTask_BadJSON_Returns400_AndDoesNotCallDB(t *testing.T) {
	db := &fakeDBClient{}
	svc := app.NewService(db)
	h := NewHttpHandlers(svc)

	req := httptest.NewRequest(http.MethodPost, "/create", strings.NewReader(`{bad-json}`))
	rr := httptest.NewRecorder()

	h.handleAddTask(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected code %d, got %d, body=%s", http.StatusBadRequest, rr.Code, rr.Body.String())
	}
	if db.addCalls != 0 {
		t.Fatalf("expected AddTask not called, got calls=%d", db.addCalls)
	}
}

func TestHandleAddTask_Returns400_OnServiceError(t *testing.T) {
	wantErr := errors.New("my error")
	db := &fakeDBClient{
		addFn: func(ctx context.Context, task models.TaskImportData) (models.TaskExportData, error) {
			return models.TaskExportData{}, wantErr
		},
	}
	svc := app.NewService(db)
	h := NewHttpHandlers(svc)

	req := httptest.NewRequest(http.MethodPost, "/create", strings.NewReader(`{"title":"t","text":"x"}`))
	rr := httptest.NewRecorder()

	h.handleAddTask(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected code %d, got %d, body=%s", http.StatusBadRequest, rr.Code, rr.Body.String())
	}
	if db.addCalls != 1 {
		t.Fatalf("expected AddTask calls=1, got %d", db.addCalls)
	}
}

func TestHandleAddTask_Success_Returns201AndTaskJSON(t *testing.T) {
	fixedTime := time.Date(2025, 12, 21, 12, 0, 0, 0, time.UTC)
	db := &fakeDBClient{
		addFn: func(ctx context.Context, task models.TaskImportData) (models.TaskExportData, error) {
			return models.TaskExportData{
				Id:        1,
				Title:     task.Title,
				Text:      task.Text,
				Finished:  false,
				CreatedAt: fixedTime,
			}, nil
		},
	}
	svc := app.NewService(db)
	h := NewHttpHandlers(svc)

	req := httptest.NewRequest(http.MethodPost, "/create", strings.NewReader(`{"title":"Buy milk","text":"2 bottles"}`))
	rr := httptest.NewRecorder()

	h.handleAddTask(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected code %d, got %d, body=%s", http.StatusCreated, rr.Code, rr.Body.String())
	}
	if db.addCalls != 1 {
		t.Fatalf("expected AddTask calls=1, got calls=%d", db.addCalls)
	}
	var got models.TaskExportData
	if err := json.Unmarshal(rr.Body.Bytes(), &got); err != nil {
		t.Fatalf("bad json response: %v, body=%s", err, rr.Body.String())
	}
	if got.Id != 1 || got.Title != "Buy milk" || got.Text != "2 bottles" || got.Finished || !got.CreatedAt.Equal(fixedTime) {
		t.Fatalf("unexpected created task response %+v", got)
	}
}

func TestHandleDeleteTask_BadId_Returns400_AndDoesNotCallDB(t *testing.T) {
	db := &fakeDBClient{}
	svc := app.NewService(db)
	h := NewHttpHandlers(svc)

	req := httptest.NewRequest(http.MethodDelete, "/delete", strings.NewReader(`{bad-json}`))
	rr := httptest.NewRecorder()

	h.handleDeleteTask(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected code %d, got %d, body=%s", http.StatusBadRequest, rr.Code, rr.Body.String())
	}
	if db.removeCalls != 0 {
		t.Fatalf("expected AddTask not called, got calls=%d", db.removeCalls)
	}
}

func TestHandleDeleteTask_Returns500(t *testing.T) {
	wantErr := errors.New("my error")
	db := &fakeDBClient{
		removeFn: func(ctx context.Context, id int) error {
			return wantErr
		},
	}
	svc := app.NewService(db)
	h := NewHttpHandlers(svc)

	req := httptest.NewRequest(http.MethodDelete, "/delete", strings.NewReader(`{"id":1}`))
	rr := httptest.NewRecorder()

	h.handleDeleteTask(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected code %d, got %d, body=%s", http.StatusInternalServerError, rr.Code, rr.Body.String())
	}
	if db.removeCalls != 1 {
		t.Fatalf("expected DeleteTask calls=1, got %d", db.removeCalls)
	}
}

func TestHandleDeleteTask_Success_Returns204(t *testing.T) {
	db := &fakeDBClient{
		removeFn: func(ctx context.Context, id int) error {
			return nil
		},
	}
	svc := app.NewService(db)
	h := NewHttpHandlers(svc)

	req := httptest.NewRequest(http.MethodDelete, "/delete", strings.NewReader(`{"id":1}`))
	rr := httptest.NewRecorder()

	h.handleDeleteTask(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Fatalf("expected code %d, got %d, body=%s", http.StatusNoContent, rr.Code, rr.Body.String())
	}
	if db.removeCalls != 1 {
		t.Fatalf("expected DeleteTask calls=1, got calls=%d", db.removeCalls)
	}
	if db.gotRemoveID != 1 {
		t.Fatalf("expected RemoveTask id=1, got %d", db.gotRemoveID)
	}
}

func TestHandleListAllTasks_Returns500(t *testing.T) {
	wantErr := errors.New("my error")
	db := &fakeDBClient{
		listFn: func(ctx context.Context) ([]models.TaskExportData, error) {
			return nil, wantErr
		},
	}
	svc := app.NewService(db)
	h := NewHttpHandlers(svc)

	req := httptest.NewRequest(http.MethodGet, "/list", strings.NewReader(``))
	rr := httptest.NewRecorder()

	h.handleListAllTasks(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected code %d, got %d, body=%s", http.StatusInternalServerError, rr.Code, rr.Body.String())
	}
	if db.listCalls != 1 {
		t.Fatalf("expected ListAllTasks calls=1, got %d", db.listCalls)
	}
}

func TestHandleListAllTasks_Success_Returns200_AndSliceOfTasks(t *testing.T) {
	db := &fakeDBClient{
		listFn: func(ctx context.Context) ([]models.TaskExportData, error) {
			return []models.TaskExportData{
				{
					Id:    1,
					Title: "title1",
					Text:  "text1",
				},
				{
					Id:    2,
					Title: "title2",
					Text:  "text2",
				},
			}, nil
		},
	}
	svc := app.NewService(db)
	h := NewHttpHandlers(svc)

	req := httptest.NewRequest(http.MethodGet, "/list", strings.NewReader(``))
	rr := httptest.NewRecorder()

	h.handleListAllTasks(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected code %d, got %d, body=%s", http.StatusOK, rr.Code, rr.Body.String())
	}
	if db.listCalls != 1 {
		t.Fatalf("expected ListAllTasks calls=1, got calls=%d", db.listCalls)
	}
	var got []models.TaskExportData
	if err := json.Unmarshal(rr.Body.Bytes(), &got); err != nil {
		t.Fatalf("bad json response: %v, body=%s", err, rr.Body.String())
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 tasks, got %d: %+v", len(got), got)
	}
	if got[0].Id != 1 || got[0].Title != "title1" || got[0].Text != "text1" ||
		got[1].Id != 2 || got[1].Title != "title2" || got[1].Text != "text2" {
		t.Fatalf("unexpected task list response %+v", got)
	}
}

func TestHandleMarkTaskFinished_BadId_Returns400_AndDoesNotCallDB(t *testing.T) {
	db := &fakeDBClient{}
	svc := app.NewService(db)
	h := NewHttpHandlers(svc)

	req := httptest.NewRequest(http.MethodPut, "/done", strings.NewReader(`{bad-json}`))
	rr := httptest.NewRecorder()

	h.handleFinishTask(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected code %d, got %d, body=%s", http.StatusBadRequest, rr.Code, rr.Body.String())
	}
	if db.doneCalls != 0 {
		t.Fatalf("expected MarkTaskFinished not called, got calls=%d", db.doneCalls)
	}
}

func TestHandleFinishTask_Returns500(t *testing.T) {
	wantErr := errors.New("my error")
	db := &fakeDBClient{
		doneFn: func(ctx context.Context, id int) (models.TaskExportData, error) {
			return models.TaskExportData{}, wantErr
		},
	}
	svc := app.NewService(db)
	h := NewHttpHandlers(svc)

	req := httptest.NewRequest(http.MethodPut, "/done", strings.NewReader(`{"id":1}`))
	rr := httptest.NewRecorder()

	h.handleFinishTask(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected code %d, got %d, body=%s", http.StatusInternalServerError, rr.Code, rr.Body.String())
	}
	if db.doneCalls != 1 {
		t.Fatalf("expected MarkTaskFinished calls=1, got %d", db.doneCalls)
	}
}

func TestHandleMarkTaskFinished_Success_Returns200AndTaskJSON(t *testing.T) {
	fixedCreatedTime := time.Date(2025, 12, 21, 12, 0, 0, 0, time.UTC)
	fixedFinishedTime := time.Date(2025, 12, 22, 12, 0, 0, 0, time.UTC)
	db := &fakeDBClient{
		doneFn: func(ctx context.Context, id int) (models.TaskExportData, error) {
			return models.TaskExportData{
				Id:         id,
				Title:      "Buy bread",
				Text:       "and carrots",
				Finished:   true,
				CreatedAt:  fixedCreatedTime,
				FinishedAt: &fixedFinishedTime,
			}, nil
		},
	}
	svc := app.NewService(db)
	h := NewHttpHandlers(svc)

	req := httptest.NewRequest(http.MethodPut, "/done", strings.NewReader(`{"id":1}`))
	rr := httptest.NewRecorder()

	h.handleFinishTask(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected code %d, got %d, body=%s", http.StatusOK, rr.Code, rr.Body.String())
	}
	if db.doneCalls != 1 {
		t.Fatalf("expected MarkTaskFinished calls=1, got calls=%d", db.doneCalls)
	}
	if db.gotDoneID != 1 {
		t.Fatalf("expected MarkTaskFinished id=1, got %d", db.gotDoneID)
	}
	var got models.TaskExportData
	if err := json.Unmarshal(rr.Body.Bytes(), &got); err != nil {
		t.Fatalf("bad json response: %v, body=%s", err, rr.Body.String())
	}
	if got.Id != 1 || got.Title != "Buy bread" || got.Text != "and carrots" ||
		!got.Finished || !got.CreatedAt.Equal(fixedCreatedTime) {
		t.Fatalf("unexpected created task response %+v", got)
	}
	if got.FinishedAt == nil {
		t.Fatalf("expected FinishedAt not nil")
	}
	if !got.FinishedAt.Equal(fixedFinishedTime) {
		t.Fatalf("expected FinishedAt=%v, got %v", fixedFinishedTime, *got.FinishedAt)
	}
}
