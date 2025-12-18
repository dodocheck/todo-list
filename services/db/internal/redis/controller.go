package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisController struct {
	redisClient *redis.Client
	taskListKey string
	ttlSeconds  int
}

func NewRedisController(ctx context.Context, address string, ttlSeconds int) (*RedisController, error) {
	redisClient, err := initRedis(ctx, address)
	if err != nil {
		return nil, err
	}

	return &RedisController{
		redisClient: redisClient,
		taskListKey: "tasks",
		ttlSeconds:  ttlSeconds,
	}, nil
}
