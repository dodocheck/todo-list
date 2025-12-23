package app

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/dodocheck/go-pet-project-1/services/db/internal/models"
	"github.com/google/go-cmp/cmp"
)

type fakeCacheController struct {
	cacheTaskListCalls int
	cacheTaskListCtx   context.Context
	cacheTaskListTasks []models.TaskExportData
	cacheTaskListErr   error

	deleteTaskListCalls int
	deleteTaskListCtx   context.Context
	deleteTaskListErr   error

	getTaskListCalls int
	getTaskListCtx   context.Context
	getTaskListRet   []models.TaskExportData
	getTaskListErr   error

	cacheTaskCalls int
	cacheTaskCtx   context.Context
	cacheTaskIn    []models.TaskExportData
	cacheTaskErr   error

	deleteTaskByIdCalls int
	deleteTaskByIdCtx   context.Context
	deleteTaskByIdId    int
	deleteTaskByIdErr   error

	getTaskByIdCalls int
	getTaskByIdCtx   context.Context
	getTaskByIdId    int
	getTaskByIdRet   models.TaskExportData
	getTaskByIdErr   error

	flushAllDataCalls int
	flushAllDataCtx   context.Context
	flushAllDataErr   error

	closeCalls int
	closeErr   error
}

func (fcc *fakeCacheController) CacheTaskList(ctx context.Context, tasks []models.TaskExportData) error {
	fcc.cacheTaskListCalls++
	fcc.cacheTaskListCtx = ctx
	fcc.cacheTaskListTasks = tasks
	return fcc.cacheTaskListErr
}

func (fcc *fakeCacheController) DeleteTaskList(ctx context.Context) error {
	fcc.deleteTaskListCalls++
	fcc.deleteTaskListCtx = ctx
	return fcc.deleteTaskListErr
}

func (fcc *fakeCacheController) GetTaskList(ctx context.Context) ([]models.TaskExportData, error) {
	fcc.getTaskListCalls++
	fcc.getTaskListCtx = ctx
	return fcc.getTaskListRet, fcc.getTaskListErr
}

func (fcc *fakeCacheController) CacheTask(ctx context.Context, task models.TaskExportData) error {
	fcc.cacheTaskCalls++
	fcc.cacheTaskCtx = ctx
	fcc.cacheTaskIn = append(fcc.cacheTaskIn, task)
	return fcc.cacheTaskErr
}

func (fcc *fakeCacheController) DeleteTaskById(ctx context.Context, id int) error {
	fcc.deleteTaskByIdCalls++
	fcc.deleteTaskByIdCtx = ctx
	fcc.deleteTaskByIdId = id
	return fcc.deleteTaskByIdErr
}

func (fcc *fakeCacheController) GetTaskById(ctx context.Context, id int) (models.TaskExportData, error) {
	fcc.getTaskByIdCalls++
	fcc.getTaskByIdCtx = ctx
	fcc.getTaskByIdId = id
	return fcc.getTaskByIdRet, fcc.getTaskByIdErr
}

func (fcc *fakeCacheController) FlushAllData(ctx context.Context) error {
	fcc.flushAllDataCalls++
	fcc.flushAllDataCtx = ctx
	return fcc.flushAllDataErr
}

func (fcc *fakeCacheController) Close() error {
	fcc.closeCalls++
	return fcc.closeErr
}

func TestCacheRepoClose_DelegatesToTaskRepo(t *testing.T) {
	wantErr := errors.New("my error")
	cr := NewCachedRepository(
		&fakeRepo{closeErr: wantErr},
		&fakeCacheController{})
	err := cr.Close()
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected %v, got %v", wantErr, err)
	}
}

func TestCacheRepoAddTask_DelegatesToTaskRepo(t *testing.T) {
	ctx := context.Background()
	wantTaskIn := models.TaskImportData{
		Title: "my in title",
		Text:  "my in text",
	}
	createdAtTS := time.Date(2025, 12, 10, 4, 6, 3, 2, time.UTC)
	finishedAtTS := time.Date(2025, 12, 10, 3, 5, 2, 1, time.UTC)
	wantTaskOut := models.TaskExportData{
		Id:         4,
		Title:      "my title",
		Text:       "my text",
		Finished:   true,
		CreatedAt:  createdAtTS,
		FinishedAt: &finishedAtTS,
	}
	wantErr := errors.New("my error")
	fr := &fakeRepo{
		addTaskRet: wantTaskOut,
		addTaskErr: wantErr,
	}
	cr := NewCachedRepository(
		fr,
		&fakeCacheController{})

	got, err := cr.AddTask(ctx, wantTaskIn)

	if fr.addTaskCalls != 1 {
		t.Fatalf("expected AddTask called=1, got=%d", fr.addTaskCalls)
	}
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected %v, got %v", wantErr, err)
	}
	if fr.addTaskCtx != ctx {
		t.Fatalf("context mismatch")
	}
	if !reflect.DeepEqual(fr.addTaskIn, wantTaskIn) {
		t.Fatalf("task in mismatch: want %+v got %+v", wantTaskIn, fr.addTaskIn)
	}
	if !reflect.DeepEqual(got, wantTaskOut) {
		t.Fatalf("task out mismatch: want %+v got %+v", wantTaskOut, got)
	}
}

func TestCacheRepoAddTask_Success_CallsCacheController(t *testing.T) {
	ctx := context.Background()
	createdAtTS := time.Date(2025, 12, 10, 4, 6, 3, 2, time.UTC)
	finishedAtTS := time.Date(2025, 12, 10, 3, 5, 2, 1, time.UTC)
	wantTaskOut := models.TaskExportData{
		Id:         4,
		Title:      "my title",
		Text:       "my text",
		Finished:   true,
		CreatedAt:  createdAtTS,
		FinishedAt: &finishedAtTS,
	}
	fcr := &fakeCacheController{}
	cr := NewCachedRepository(
		&fakeRepo{addTaskRet: wantTaskOut},
		fcr)

	cr.AddTask(ctx, models.TaskImportData{})

	if fcr.cacheTaskCalls != 1 {
		t.Fatalf("expected CacheTask called once, got %d calls", fcr.cacheTaskCalls)
	}
	if fcr.cacheTaskCtx != ctx {
		t.Fatalf("context mismatch")
	}
	if diff := cmp.Diff(wantTaskOut, fcr.cacheTaskIn[0]); diff != "" {
		t.Fatal(diff)
	}
	if fcr.deleteTaskListCalls != 1 {
		t.Fatalf("expected DeleteTaskList called once, got %d calls", fcr.deleteTaskListCalls)
	}
	if fcr.deleteTaskListCtx != ctx {
		t.Fatalf("context mismatch")
	}
}

func TestCacheRepoAddTask_Error_DoesNotCallCacheController(t *testing.T) {
	wantErr := errors.New("boom")
	fcr := &fakeCacheController{}
	cr := NewCachedRepository(
		&fakeRepo{addTaskErr: wantErr},
		fcr)

	cr.AddTask(context.Background(), models.TaskImportData{})
	if fcr.cacheTaskCalls != 0 {
		t.Fatalf("expected CacheTask not called, got %d calls", fcr.cacheTaskCalls)
	}
	if fcr.deleteTaskListCalls != 0 {
		t.Fatalf("expected DeleteTaskList not called, got %d calls", fcr.deleteTaskListCalls)
	}
}

func TestCacheRepoDeleteTask_DelegatesToTaskRepo(t *testing.T) {
	ctx := context.Background()
	wantId := 1
	wantErr := errors.New("my error")
	fr := &fakeRepo{
		deleteTaskErr: wantErr,
	}
	cr := NewCachedRepository(
		fr,
		&fakeCacheController{})

	err := cr.DeleteTask(ctx, wantId)

	if fr.deleteTaskCalls != 1 {
		t.Fatalf("expected DeleteTask called=1, got=%d", fr.deleteTaskCalls)
	}
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected %v, got %v", wantErr, err)
	}
	if fr.deleteTaskCtx != ctx {
		t.Fatalf("context mismatch")
	}
	if fr.deleteTaskIn != wantId {
		t.Fatalf("expected id %d, got %d", wantId, fr.deleteTaskIn)
	}
}

func TestCacheRepoDeleteTask_Success_CallsCacheController(t *testing.T) {
	ctx := context.Background()
	wantId := 53
	fcr := &fakeCacheController{}
	cr := NewCachedRepository(
		&fakeRepo{},
		fcr)

	cr.DeleteTask(ctx, wantId)

	if fcr.deleteTaskByIdCalls != 1 {
		t.Fatalf("expected DeleteTaskById called once, got %d calls", fcr.deleteTaskByIdCalls)
	}
	if fcr.deleteTaskByIdCtx != ctx {
		t.Fatal("context mismatch")
	}
	if fcr.deleteTaskByIdId != wantId {
		t.Fatalf("expected id=%d, got=%d", wantId, fcr.deleteTaskByIdId)
	}
	if fcr.deleteTaskListCalls != 1 {
		t.Fatalf("expected DeleteTaskList called once, got %d calls", fcr.deleteTaskListCalls)
	}
	if fcr.deleteTaskListCtx != ctx {
		t.Fatal("context mismatch")
	}
}

func TestCacheRepoDeleteTask_Error_DoesNotCallCacheController(t *testing.T) {
	wantErr := errors.New("boom")
	fcr := &fakeCacheController{}
	cr := NewCachedRepository(
		&fakeRepo{deleteTaskErr: wantErr},
		fcr)

	cr.DeleteTask(context.Background(), 1)
	if fcr.deleteTaskByIdCalls != 0 {
		t.Fatalf("expected DeleteTaskById not called, got %d calls", fcr.deleteTaskByIdCalls)
	}
	if fcr.deleteTaskListCalls != 0 {
		t.Fatalf("expected DeleteTaskList not called, got %d calls", fcr.deleteTaskListCalls)
	}
}

func TestCacheRepoListAllTasks_CacheMiss_DelegatesToTaskRepo(t *testing.T) {
	ctx := context.Background()
	createdAtTS := time.Date(2025, 12, 10, 4, 6, 3, 2, time.UTC)
	finishedAtTS := time.Date(2025, 12, 10, 3, 5, 2, 1, time.UTC)
	wantTasksOut := []models.TaskExportData{
		{
			Id:         4,
			Title:      "my title",
			Text:       "my text",
			Finished:   false,
			CreatedAt:  createdAtTS,
			FinishedAt: nil,
		},
		{
			Id:         54,
			Title:      "my title2",
			Text:       "my text2",
			Finished:   true,
			CreatedAt:  createdAtTS,
			FinishedAt: &finishedAtTS,
		},
	}
	fr := &fakeRepo{
		listAllTasksRet: wantTasksOut,
	}
	cr := NewCachedRepository(
		fr,
		&fakeCacheController{
			getTaskListErr: errors.New("cache miss"),
		})

	got, _ := cr.ListAllTasks(ctx)
	if fr.listAllTasksCalls != 1 {
		t.Fatalf("expected ListAllTasks called=1, got %d", fr.listAllTasksCalls)
	}
	if fr.listAllTasksCtx != ctx {
		t.Fatalf("context mismatch")
	}
	if !reflect.DeepEqual(got, wantTasksOut) {
		t.Fatalf("task list mismatch: want %+v got %+v", wantTasksOut, got)
	}

	wantErr := errors.New("my error")
	fr.listAllTasksErr = wantErr
	_, err := cr.ListAllTasks(context.Background())
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected %v, got %v", wantErr, err)
	}
}

func TestCacheRepoListAllTasks_CacheHit_DelegatesToCacheController(t *testing.T) {
	ctx := context.Background()
	createdAtTS := time.Date(2025, 12, 10, 4, 6, 3, 2, time.UTC)
	finishedAtTS := time.Date(2025, 12, 10, 3, 5, 2, 1, time.UTC)
	wantTasksOut := []models.TaskExportData{
		{
			Id:         4,
			Title:      "my title",
			Text:       "my text",
			Finished:   false,
			CreatedAt:  createdAtTS,
			FinishedAt: nil,
		},
		{
			Id:         54,
			Title:      "my title2",
			Text:       "my text2",
			Finished:   true,
			CreatedAt:  createdAtTS,
			FinishedAt: &finishedAtTS,
		},
	}
	fcr := &fakeCacheController{
		getTaskListRet: wantTasksOut,
	}
	fr := &fakeRepo{}
	cr := NewCachedRepository(
		fr,
		fcr,
	)

	got, _ := cr.ListAllTasks(ctx)
	if fcr.getTaskListCalls != 1 {
		t.Fatalf("expected GetTaskList called=1, got %d", fcr.getTaskListCalls)
	}
	if fcr.getTaskListCtx != ctx {
		t.Fatalf("context mismatch")
	}
	if !reflect.DeepEqual(got, wantTasksOut) {
		t.Fatalf("task list mismatch: want %+v, got %+v", wantTasksOut, got)
	}
	if fr.listAllTasksCalls != 0 {
		t.Fatalf("mainDB expected not called, got called=%d", fr.listAllTasksCalls)
	}
}

func TestCacheRepoListAllTasks_CacheMissTaskRepoSuccess_CallsCacheController(t *testing.T) {
	ctx := context.Background()
	createdAtTS := time.Date(2025, 12, 10, 4, 6, 3, 2, time.UTC)
	finishedAtTS := time.Date(2025, 12, 10, 3, 5, 2, 1, time.UTC)
	wantTasksOut := []models.TaskExportData{
		{
			Id:         4,
			Title:      "my title",
			Text:       "my text",
			Finished:   false,
			CreatedAt:  createdAtTS,
			FinishedAt: nil,
		},
		{
			Id:         54,
			Title:      "my title2",
			Text:       "my text2",
			Finished:   true,
			CreatedAt:  createdAtTS,
			FinishedAt: &finishedAtTS,
		},
	}
	fcr := &fakeCacheController{getTaskListErr: errors.New("cache miss")}
	cr := NewCachedRepository(
		&fakeRepo{listAllTasksRet: wantTasksOut},
		fcr)

	cr.ListAllTasks(ctx)
	if fcr.cacheTaskListCalls != 1 {
		t.Fatalf("expected CacheTaskList called once, got %d calls", fcr.cacheTaskListCalls)
	}
	if fcr.cacheTaskListCtx != ctx {
		t.Fatal("context mismatch")
	}
	if diff := cmp.Diff(fcr.cacheTaskListTasks, wantTasksOut); diff != "" {
		t.Fatal(diff)
	}
	if fcr.cacheTaskCalls != len(wantTasksOut) {
		t.Fatalf("expected CacheTask calls=%d, got %d calls", len(wantTasksOut), fcr.cacheTaskCalls)
	}
	if diff := cmp.Diff(fcr.cacheTaskIn, wantTasksOut); diff != "" {
		t.Fatal(diff)
	}
}

func TestCacheRepoListAllTasks_CacheMissTaskRepoError_DoesNotCallCacheController(t *testing.T) {
	fcr := &fakeCacheController{getTaskListErr: errors.New("cache miss")}
	cr := NewCachedRepository(
		&fakeRepo{listAllTasksErr: errors.New("my error")},
		fcr)

	cr.ListAllTasks(context.Background())
	if fcr.cacheTaskListCalls != 0 {
		t.Fatalf("expected CacheTaskList not called, got %d calls", fcr.cacheTaskListCalls)
	}
	if fcr.cacheTaskCalls != 0 {
		t.Fatalf("expected CacheTask not called, got %d calls", fcr.cacheTaskCalls)
	}
}

func TestCacheRepoMarkTaskFinished_DelegatesToTaskRepo(t *testing.T) {
	ctx := context.Background()
	wantId := 5
	createdAtTS := time.Date(2025, 12, 10, 4, 6, 3, 2, time.UTC)
	finishedAtTS := time.Date(2025, 12, 10, 3, 5, 2, 1, time.UTC)
	wantTaskOut := models.TaskExportData{
		Id:         46,
		Title:      "my title",
		Text:       "my text",
		Finished:   true,
		CreatedAt:  createdAtTS,
		FinishedAt: &finishedAtTS,
	}
	wantErr := errors.New("my error")
	fr := &fakeRepo{
		markTaskFinishedRet: wantTaskOut,
		markTaskFinishedErr: wantErr,
	}
	cr := NewCachedRepository(
		fr,
		&fakeCacheController{})

	got, err := cr.MarkTaskFinished(context.Background(), wantId)

	if fr.markTaskFinishedCalls != 1 {
		t.Fatalf("expected MarkTaskFinished called=1, got=%d", fr.markTaskFinishedCalls)
	}
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected %v, got %v", wantErr, err)
	}
	if fr.markTaskFinishedCtx != ctx {
		t.Fatalf("context mismatch")
	}
	if fr.markTaskFinishedIn != wantId {
		t.Fatalf("expected id %d, got %d", wantId, fr.markTaskFinishedIn)
	}
	if !reflect.DeepEqual(got, wantTaskOut) {
		t.Fatalf("task out mismatch: want %+v got %+v", wantTaskOut, got)
	}
}

func TestCacheRepoMarkTaskFinished_Success_CallsCacheController(t *testing.T) {
	ctx := context.Background()
	createdAtTS := time.Date(2025, 12, 10, 4, 6, 3, 2, time.UTC)
	finishedAtTS := time.Date(2025, 12, 10, 3, 5, 2, 1, time.UTC)
	wantTaskOut := models.TaskExportData{
		Id:         46,
		Title:      "my title",
		Text:       "my text",
		Finished:   true,
		CreatedAt:  createdAtTS,
		FinishedAt: &finishedAtTS,
	}
	fcr := &fakeCacheController{}
	cr := NewCachedRepository(
		&fakeRepo{markTaskFinishedRet: wantTaskOut},
		fcr)

	cr.MarkTaskFinished(ctx, 1)
	if fcr.cacheTaskCalls != 1 {
		t.Fatalf("expected CacheTask called once, got %d calls", fcr.cacheTaskCalls)
	}
	if fcr.cacheTaskCtx != ctx {
		t.Fatal("context mismatch")
	}
	if diff := cmp.Diff(fcr.cacheTaskIn[0], wantTaskOut); diff != "" {
		t.Fatal(diff)
	}
	if fcr.deleteTaskListCalls != 1 {
		t.Fatalf("expected DeleteTaskList called once, got %d calls", fcr.deleteTaskListCalls)
	}
	if fcr.deleteTaskListCtx != ctx {
		t.Fatal("context mismatch")
	}
}

func TestCacheRepoMarkTaskFinished_Error_DoesNotCallCacheController(t *testing.T) {
	wantErr := errors.New("boom")
	fcr := &fakeCacheController{}
	cr := NewCachedRepository(
		&fakeRepo{markTaskFinishedErr: wantErr},
		fcr)

	cr.MarkTaskFinished(context.Background(), 1)
	if fcr.cacheTaskCalls != 0 {
		t.Fatalf("expected CacheTask not called, got %d calls", fcr.cacheTaskCalls)
	}
	if fcr.deleteTaskListCalls != 0 {
		t.Fatalf("expected DeleteTaskList not called, got %d calls", fcr.deleteTaskListCalls)
	}
}
