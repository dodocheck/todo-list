package app

import (
	"context"
	"log"

	"github.com/dodocheck/go-pet-project-1/services/db/internal/models"
)

type CachedRepository struct {
	mainDBClient  TaskRepository
	cacheDBClient CacheController
}

func NewCachedRepository(ctx context.Context, mainDBClient TaskRepository, cacheDBClient CacheController) *CachedRepository {
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
		if cacheErr := cr.cacheDBClient.DeleteTaskList(ctx); cacheErr != nil {
			log.Printf("cache delete tasks err: %v", cacheErr)
		}
	}

	return createdTask, err
}

func (cr *CachedRepository) DeleteTask(ctx context.Context, id int) error {
	err := cr.mainDBClient.DeleteTask(ctx, id)

	if err == nil {
		if cacheErr := cr.cacheDBClient.DeleteTaskList(ctx); cacheErr != nil {
			log.Printf("cache delete tasks err: %v", cacheErr)
		}
	}

	return err
}

func (cr *CachedRepository) ListAllTasks(ctx context.Context) ([]models.TaskExportData, error) {
	cacheTasks, cacheErr := cr.cacheDBClient.GetTaskList(ctx)
	if cacheErr == nil {
		log.Println("cache hit! returning tasks from cache!")
		return cacheTasks, nil
	}

	tasks, err := cr.mainDBClient.ListAllTasks(ctx)

	if err == nil {
		if cacheErr := cr.cacheDBClient.CacheTaskList(ctx, tasks); cacheErr != nil {
			log.Printf("cache cache tasks err: %v", cacheErr)
		}
		log.Println("returning tasks from main DB")
		return tasks, nil
	}

	return nil, err
}

func (cr *CachedRepository) MarkTaskFinished(ctx context.Context, id int) (models.TaskExportData, error) {
	updatedTask, err := cr.mainDBClient.MarkTaskFinished(ctx, id)

	if err == nil {
		if cacheErr := cr.cacheDBClient.DeleteTaskList(ctx); cacheErr != nil {
			log.Printf("cache delete tasks err: %v", cacheErr)
		}
	}

	return updatedTask, err
}
