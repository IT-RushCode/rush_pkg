package middlewares

import (
	"fmt"
	"strings"
	"time"

	"github.com/IT-RushCode/rush_pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

// CacheMiddleware кэширует запросы и возвращает кэшированные ответы, если они имеются.
//
// Параметры:
// - cache: клиент Redis для хранения кэшированных данных.
// - cacheTime: время хранения кэшированных данных в минутах.
// - noCachePaths: список путей, которые не должны кэшироваться (например, []string{"/api/v1/example", "/example"}).
// - cacheMethods: список HTTP-методов, которые должны кэшироваться (например, []string{"GET", "HEAD"}).
//
// Пример использования:
// app.Use(CacheMiddleware(redisClient, 30, []string{"/api/v1/example"}, []string{"GET", "HEAD"}))
//
// Поведение:
// - Если запрос попадает в список noCachePaths или использует метод, который не нужно кэшировать, то кэширование пропускается.
// - При попадании запроса в кэшируемые методы и пути, ответ сохраняется в Redis с заданным временем хранения.
// - При повторном запросе, если данные уже есть в кэше, они возвращаются без выполнения контроллера.
func CacheMiddleware(cache *redis.Client, cacheTime uint, noCachePaths []string, cacheMethods []string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// Проверка на URL, которые не должны кэшироваться
		for _, noCachePath := range noCachePaths {
			if utils.IsRouteMatch(ctx.Path(), noCachePath) {
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
			// Возвращаем кэшированный ответ
			ctx.Set("Content-Type", "application/json")
			return ctx.SendString(cached)
		} else if err != redis.Nil {
			// Если произошла ошибка, отличная от "ключ не найден"
			return err
		}

		// Выполняем следующий обработчик (контроллер)
		err = ctx.Next()
		if err != nil {
			return err
		}

		// После выполнения обработчика проверяем статус код ответа
		resStatus := ctx.Response().StatusCode()
		if resStatus == fiber.StatusUnauthorized || resStatus == fiber.StatusForbidden {
			// Не кэшируем ответы с ошибками 401 или 403
			return nil
		}

		// Сохраняем ответ в кеш на указанное время
		if cacheTime == 0 {
			cacheTime = 5
		}
		expiration := time.Duration(cacheTime) * time.Minute

		// Сохраняем данные в кэше в виде строки
		if err := cache.Set(ctx.Context(), cacheKey, string(ctx.Response().Body()), expiration).Err(); err != nil {
			return err
		}

		return nil
	}
}

// CacheInvalidationMiddleware удаляет кэш при выполнении операций, изменяющих данные.
//
// Параметры:
// - cache: клиент Redis для удаления кэшированных данных.
// - cacheablePaths: список путей, которые могут быть кэшированы (например, []string{"/api/v1/example", "/example"}).
//
// Пример использования:
// app.Use(CacheInvalidationMiddleware(redisClient, []string{"/api/v1/example"}))
//
// Поведение:
// - Middleware отслеживает методы, которые изменяют данные (POST, PUT, PATCH, DELETE).
// - При обнаружении одного из этих методов, middleware проверяет, затронут ли запросом один из кэшируемых путей.
// - Если текущий путь совпадает с кэшируемым, middleware использует команду SCAN для поиска всех кэшированных ключей, связанных с этим путём, и удаляет их.
// - SCAN используется для безопасного удаления ключей без блокировки Redis, особенно если существует много ключей, соответствующих указанному префиксу.
func CacheInvalidationMiddleware(cache *redis.Client, cacheablePaths []string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// Определяем методы, которые изменяют данные и требуют очистки кэша
		invalidatingMethods := map[string]struct{}{
			"POST":   {},
			"PUT":    {},
			"PATCH":  {},
			"DELETE": {},
		}

		// Проверяем, если метод является изменяющим данные
		if _, shouldInvalidate := invalidatingMethods[ctx.Method()]; shouldInvalidate {
			// Проходим по списку кэшируемых путей и проверяем совпадение с текущим путем
			for _, path := range cacheablePaths {
				// Проверяем, начинается ли текущий путь с кэшируемого пути
				if strings.HasPrefix(ctx.Path(), path) {
					// Создаем префикс для поиска ключей в Redis
					prefix := fmt.Sprintf("cache_%s", path)

					// Используем SCAN для поиска всех ключей, начинающихся с указанного префикса
					var cursor uint64
					for {
						keys, nextCursor, err := cache.Scan(ctx.Context(), cursor, prefix+"*", 10).Result()
						if err != nil {
							return err
						}

						// Удаляем найденные ключи
						if len(keys) > 0 {
							if err := cache.Del(ctx.Context(), keys...).Err(); err != nil {
								return err
							}
						}

						// Если дошли до конца, выходим из цикла
						if nextCursor == 0 {
							break
						}
						cursor = nextCursor
					}
				}
			}
		}

		// Выполняем основной обработчик (контроллер)
		return ctx.Next()
	}
}
