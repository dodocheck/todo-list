package redis

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/dodocheck/go-pet-project-1/services/db/internal/app"
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

func (rc *RedisController) Close() error {
	return rc.redisClient.Close()
}

func (rc *RedisController) CacheTaskList(ctx context.Context, tasks []models.TaskExportData) error {
	taskList, err := json.Marshal(tasks)
	if err != nil {
		return err
	}

	return rc.redisClient.Set(ctx, rc.taskListKey, taskList, time.Duration(rc.ttlSeconds)*time.Second).Err()
}

func (rc *RedisController) DeleteTaskList(ctx context.Context) error {
	return rc.redisClient.Del(ctx, rc.taskListKey).Err()
}

func (rc *RedisController) GetTaskList(ctx context.Context) ([]models.TaskExportData, error) {
	tasksStr, err := rc.redisClient.Get(ctx, rc.taskListKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, app.ErrTaskNotFound
		}
		return nil, err
	}

	var tasksToReturn []models.TaskExportData
	if err := json.Unmarshal([]byte(tasksStr), &tasksToReturn); err != nil {
		return nil, err
	}

	return tasksToReturn, nil
}

func (rc *RedisController) CacheTask(ctx context.Context, task models.TaskExportData) error {
	key := "task:" + strconv.Itoa(task.Id)
	taskStr, err := json.Marshal(task)
	if err != nil {
		return err
	}

	return rc.redisClient.Set(ctx, key, taskStr, time.Duration(rc.ttlSeconds)*time.Second).Err()
}

func (rc *RedisController) DeleteTaskById(ctx context.Context, id int) error {
	key := "task:" + strconv.Itoa(id)
	return rc.redisClient.Del(ctx, key).Err()
}

func (rc *RedisController) GetTaskById(ctx context.Context, id int) (models.TaskExportData, error) {
	key := "task:" + strconv.Itoa(id)
	taskStr, err := rc.redisClient.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return models.TaskExportData{}, app.ErrTaskNotFound
		}
		return models.TaskExportData{}, err
	}

	var taskToReturn models.TaskExportData
	if err := json.Unmarshal([]byte(taskStr), &taskToReturn); err != nil {
		return models.TaskExportData{}, err
	}

	return taskToReturn, nil
}

func (rc *RedisController) FlushAllData(ctx context.Context) error {
	return rc.redisClient.FlushAll(ctx).Err()
}
