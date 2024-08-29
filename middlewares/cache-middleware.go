package middlewares

import (
	"fmt"
	"time"

	"github.com/IT-RushCode/rush_pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

func CacheMiddleware(cache *redis.Client, cacheTime uint, noCachePaths []string, cacheMethods []string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// Проверка на URL, которые не должны кэшироваться
		for _, noCachePath := range noCachePaths {
			if ctx.Path() == noCachePath {
				return ctx.Next()
			}
		}

		// Проверка на методы, которые должны кэшироваться
		methodInCacheList := make(map[string]struct{}, len(cacheMethods))
		for _, method := range cacheMethods {
			methodInCacheList[method] = struct{}{}
		}

		// Если метод не должен кэшироваться, пропустить кэширование
		if _, shouldCacheMethod := methodInCacheList[ctx.Method()]; !shouldCacheMethod {
			return ctx.Next()
		}

		// Создаем ключ для кеша на основе URL запроса
		cacheKey := fmt.Sprintf("cache_%s", ctx.OriginalURL())

		// Попробуем получить значение из кеша
		cached, err := cache.Get(ctx.Context(), cacheKey).Result()
		if err == nil && cached != "" {
			// Если данные есть в кеше, возвращаем их
			return utils.SuccessResponse(ctx, utils.Success, string(cached))
		} else if err != redis.Nil {
			// Если произошла ошибка, отличная от "ключ не найден"
			return err
		}

		// Выполняем следующий обработчик (контроллер)
		err = ctx.Next()
		if err != nil {
			return err
		}

		// Сохраняем ответ в кеш на указанное время
		if cacheTime == 0 {
			cacheTime = 5
		}

		responseData := ctx.Response().Body()
		expiration := time.Duration(cacheTime) * time.Minute

		// Сохраняем данные в кэше в формате строки
		if err := cache.Set(ctx.Context(), cacheKey, responseData, expiration).Err(); err != nil {
			return err
		}

		return nil
	}
}
