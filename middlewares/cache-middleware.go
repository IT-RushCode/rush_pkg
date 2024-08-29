package middlewares

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

func CacheMiddleware(cache *redis.Client, cacheTime uint) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// Создаем ключ для кеша на основе URL запроса
		cacheKey := fmt.Sprintf("cache_%s", ctx.OriginalURL())

		// Попробуем получить значение из кеша
		cached, err := cache.Get(ctx.Context(), cacheKey).Result()
		if err == nil && cached != "" {
			// Если данные есть в кеше, десериализуем их и возвращаем
			var responseBody map[string]interface{}
			if err := json.Unmarshal([]byte(cached), &responseBody); err == nil {
				return ctx.Status(fiber.StatusOK).JSON(responseBody)
			} else {
				return err
			}
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

		// Преобразуем тело ответа в JSON и сохраняем в Redis
		var responseBody interface{}
		if err := json.Unmarshal(responseData, &responseBody); err == nil {
			cachedData, err := json.Marshal(responseBody)
			if err != nil {
				return err
			}
			if err := cache.Set(ctx.Context(), cacheKey, cachedData, time.Duration(cacheTime)*time.Minute).Err(); err != nil {
				return err
			}
		}

		return nil
	}
}
