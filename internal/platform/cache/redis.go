package cache

import (
	"context"
	"fmt"
	"go-war-ticket-service/configs"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(cfg configs.Config) (*redis.Client, error) {
	addr := fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort)

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB, // default 0
	})

	// Test connection
	ctx := context.Background()
	if _, err := rdb.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return rdb, nil
}
