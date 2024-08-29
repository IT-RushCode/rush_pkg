package utils

import (
	"context"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

func ClearCache(clearCache *bool, cache *redis.Client) {
	if *clearCache {
		err := cache.FlushAll(context.Background()).Err()
		if err != nil {
			log.Fatalf("Ошибка очистки кеша Redis: %v", err)
		}
		log.Println("Кеш Redis успешно очищен")
		os.Exit(0)
	}
}
