package app

import (
	"context"
	"log"

	"github.com/dodocheck/go-pet-project-1/services/db/internal/models"
)

type TaskRepository interface {
	AddTask(ctx context.Context, task models.TaskImportData) (models.TaskExportData, error)
	DeleteTask(ctx context.Context, id int) error
	ListAllTasks(ctx context.Context) ([]models.TaskExportData, error)
	MarkTaskFinished(ctx context.Context, id int) (models.TaskExportData, error)
	Close() error
}

type CacheController interface {
	CacheTaskList(ctx context.Context, tasks []models.TaskExportData) error
	DeleteTaskList(ctx context.Context) error
	GetTaskList(ctx context.Context) ([]models.TaskExportData, error)
	CacheTask(ctx context.Context, task models.TaskExportData) error
	DeleteTaskById(ctx context.Context, id int) error
	GetTaskById(ctx context.Context, id int) (models.TaskExportData, error)
	FlushAllData(ctx context.Context) error
	Close() error
}

type CachedRepository struct {
	mainDBClient  TaskRepository
	cacheDBClient CacheController
}

func NewCachedRepository(mainDBClient TaskRepository, cacheDBClient CacheController) *CachedRepository {
	return &CachedRepository{
		mainDBClient:  mainDBClient,
		cacheDBClient: cacheDBClient,
	}
}

func (cr *CachedRepository) Close() error {
	_ = cr.cacheDBClient.Close()
	return cr.mainDBClient.Close()
}

func (cr *CachedRepository) AddTask(ctx context.Context, task models.TaskImportData) (models.TaskExportData, error) {
	createdTask, err := cr.mainDBClient.AddTask(ctx, task)

	if err == nil {
		if cacheTaskErr := cr.cacheDBClient.CacheTask(ctx, createdTask); cacheTaskErr != nil {
			log.Printf("cache add task err: %v\n", cacheTaskErr)
		}
		if cacheTaskListErr := cr.cacheDBClient.DeleteTaskList(ctx); cacheTaskListErr != nil {
			log.Printf("cache delete tasklist err: %v\n", cacheTaskListErr)
		}
	}

	return createdTask, err
}

func (cr *CachedRepository) DeleteTask(ctx context.Context, id int) error {
	err := cr.mainDBClient.DeleteTask(ctx, id)

	if err == nil {
		if cacheTaskErr := cr.cacheDBClient.DeleteTaskById(ctx, id); cacheTaskErr != nil {
			log.Printf("cache delete task err: %v\n", cacheTaskErr)
		}
		if cacheTaskListErr := cr.cacheDBClient.DeleteTaskList(ctx); cacheTaskListErr != nil {
			log.Printf("cache delete tasklist err: %v\n", cacheTaskListErr)
		}
	}

	return err
}

func (cr *CachedRepository) ListAllTasks(ctx context.Context) ([]models.TaskExportData, error) {
	cacheTasks, cacheErr := cr.cacheDBClient.GetTaskList(ctx)
	if cacheErr == nil {
		log.Println("cache hit! returning tasklist from cache!")
		return cacheTasks, nil
	}

	if cacheErr != ErrTaskNotFound {
		log.Printf("cache degraded: %v\n", cacheErr)
	}

	tasks, err := cr.mainDBClient.ListAllTasks(ctx)

	if err == nil {
		if cacheTaskListErr := cr.cacheDBClient.CacheTaskList(ctx, tasks); cacheTaskListErr != nil {
			log.Printf("cache tasklist err: %v\n", cacheTaskListErr)
		}
		for _, task := range tasks {
			if err := cr.cacheDBClient.CacheTask(ctx, task); err != nil {
				log.Printf("cache add task err: %v\n", err)
			}
		}
		log.Println("returning tasklist from main DB")
		return tasks, nil
	}

	return nil, err
}

func (cr *CachedRepository) MarkTaskFinished(ctx context.Context, id int) (models.TaskExportData, error) {
	updatedTask, err := cr.mainDBClient.MarkTaskFinished(ctx, id)

	if err == nil {
		if cacheTaskErr := cr.cacheDBClient.CacheTask(ctx, updatedTask); cacheTaskErr != nil {
			log.Printf("cache add task err: %v\n", cacheTaskErr)
		}
		if cacheTaskListErr := cr.cacheDBClient.DeleteTaskList(ctx); cacheTaskListErr != nil {
			log.Printf("cache delete tasklist err: %v\n", cacheTaskListErr)
		}
	}

	return updatedTask, err
}
