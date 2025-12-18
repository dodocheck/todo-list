package cachedcontroller

import (
	"context"
	"log"

	"github.com/dodocheck/go-pet-project-1/services/db/internal/dbcontroller/postgres"
	"github.com/dodocheck/go-pet-project-1/services/db/internal/models"
	"github.com/redis/go-redis/v9"
)

type CachedController struct {
	redisClient    *redis.Client
	taskListKey    string
	ttlSeconds     int
	postgresClient *postgres.PostgresController
}

func NewCachedController(ctx context.Context, address string, ttlSeconds int) (*CachedController, error) {
	redisClient, err := initRedis(ctx, address)
	if err != nil {
		return nil, err
	}

	postgresClient := postgres.NewPostgresController()

	return &CachedController{
		redisClient:    redisClient,
		taskListKey:    "tasklist",
		ttlSeconds:     ttlSeconds,
		postgresClient: postgresClient,
	}, nil
}

func (cc *CachedController) Close() {
	cc.redisClient.Close()
	cc.postgresClient.Close()
}

func (cc *CachedController) AddTask(ctx context.Context, task models.TaskImportData) (models.TaskExportData, error) {
	createdTask, err := cc.postgresClient.AddTask(ctx, task)

	if err == nil {
		if redisErr := cc.deleteTaskList(ctx); redisErr != nil {
			log.Printf("redis delete tasks err: %v", redisErr)
		}
	}

	return createdTask, err
}

func (cc *CachedController) DeleteTask(ctx context.Context, id int) error {
	err := cc.postgresClient.DeleteTask(ctx, id)

	if err == nil {
		if redisErr := cc.deleteTaskList(ctx); redisErr != nil {
			log.Printf("redis delete tasks err: %v", redisErr)
		}
	}

	return err
}

func (cc *CachedController) ListAllTasks(ctx context.Context) ([]models.TaskExportData, error) {
	cacheTasks, cacheErr := cc.getTaskList(ctx)
	if cacheErr == nil {
		log.Println("cache hit! returning tasks from redis!")
		return cacheTasks, nil
	}

	tasks, err := cc.postgresClient.ListAllTasks(ctx)

	if err == nil {
		if redisErr := cc.cacheTaskList(ctx, tasks); redisErr != nil {
			log.Printf("redis cache tasks err: %v", redisErr)
		}
		log.Println("returning tasks from postgres")
		return tasks, nil
	}

	return nil, err
}

func (cc *CachedController) MarkTaskFinished(ctx context.Context, id int) (models.TaskExportData, error) {
	updatedTask, err := cc.postgresClient.MarkTaskFinished(ctx, id)

	if err == nil {
		if redisErr := cc.deleteTaskList(ctx); redisErr != nil {
			log.Printf("redis delete tasks err: %v", redisErr)
		}
	}

	return updatedTask, err
}
