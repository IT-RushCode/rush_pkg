package storage

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"

	"github.com/IT-RushCode/rush_pkg/config"
)

func REDIS_CONNECT(ctx context.Context, cfg *config.RedisConfig) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.HOST, cfg.PORT),
		Password: cfg.PASS,
		DB:       cfg.DB,
		Protocol: 3,
	})

	pong, err := client.Ping(ctx).Result()
	if err != nil {
		panic("redis is not connected!")
	}
	fmt.Println(pong)
	return client
}

func REDIS_CLOSE(client *redis.Client) {
	client.Close()
}
