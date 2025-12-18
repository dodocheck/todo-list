package cachedcontroller

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/dodocheck/go-pet-project-1/services/db/internal/models"
	"github.com/redis/go-redis/v9"
)

func initRedis(ctx context.Context, address string) (*redis.Client, error) {
	redisClient := redis.NewClient(
		&redis.Options{
			Addr:     address,
			Password: "",
			DB:       0,
		},
	)

	if _, err := redisClient.Ping(ctx).Result(); err != nil {
		return nil, err
	}

	return redisClient, nil
}

func (cc *CachedController) cacheTaskList(ctx context.Context, tasks []models.TaskExportData) error {
	taskList, err := json.Marshal(tasks)
	if err != nil {
		return err
	}

	if err := cc.redisClient.Set(ctx, cc.taskListKey, taskList, time.Duration(cc.ttlSeconds)*time.Second).Err(); err != nil {
		return err
	}

	return nil
}

func (cc *CachedController) deleteTaskList(ctx context.Context) error {
	return cc.redisClient.Del(ctx, cc.taskListKey).Err()
}

func (cc *CachedController) getTaskList(ctx context.Context) ([]models.TaskExportData, error) {
	tasksStr, err := cc.redisClient.Get(ctx, cc.taskListKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, errors.New("cache miss")
		} else {
			return nil, err
		}
	}

	tasksToReturn := make([]models.TaskExportData, 0)
	if err := json.Unmarshal([]byte(tasksStr), &tasksToReturn); err != nil {
		return nil, err
	}

	return tasksToReturn, nil
}
