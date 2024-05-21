package database

import (
	"fmt"

	"github.com/redis/go-redis/v9"
	"gitlab.arvand.tj/conveyor/arvand_pkg/config"
)

func REDIS_CONNECT(cfg *config.RedisConfig) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.HOST, cfg.PORT),
		Password: cfg.PASS,
		DB:       cfg.DB,
		Protocol: 3,
	})
	return client
}

func REDIS_CLOSE(client *redis.Client) {
	client.Close()
}
