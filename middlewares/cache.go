package middlewares

import (
	"fmt"
	"strings"
	"time"

	"github.com/IT-RushCode/rush_pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

// CacheMiddleware содержит все middleware для кэширования и их конфигурации.
type CacheMiddleware struct {
	cache *redis.Client
}

// NewCacheMiddleware создаёт новый экземпляр CacheMiddleware с Redis клиентом.
func NewCacheMiddleware(redisClient *redis.Client) *CacheMiddleware {
	return &CacheMiddleware{
		cache: redisClient,
	}
}

// GlobalCache глобально кэширует запросы и возвращает кэшированные ответы, если они имеются.
func (f *CacheMiddleware) GlobalCache(
	cacheTime uint,
	noCachePaths []string,
	cacheMethods []string,
) fiber.Handler {
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

		// Создаем ключ для кэша на основе URL запроса
		cacheKey := fmt.Sprintf("cache_%s", ctx.OriginalURL())

		// Попробуем получить значение из кеша
		cached, err := f.cache.Get(ctx.Context(), cacheKey).Result()
		if err == nil && cached != "" {
			// Возвращаем кэшированный ответ
			ctx.Set("Content-Type", "application/json")
			return ctx.SendString(cached)
		} else if err != redis.Nil {
			// Если произошла ошибка, отличная от "ключ не найден"
			return err
		}

		// Выполняем следующие middleware и контроллер, включая PermissionMiddleware
		err = ctx.Next()
		if err != nil {
			return err
		}

		// Проверяем статус ответа после выполнения PermissionMiddleware
		resStatus := ctx.Response().StatusCode()
		if resStatus < 200 || resStatus >= 300 {
			// Не кэшируем неуспешные ответы, такие как 401 или 403
			return nil
		}

		// Сохраняем ответ в кеш на указанное время
		if cacheTime == 0 {
			cacheTime = 5
		}
		expiration := time.Duration(cacheTime) * time.Minute

		// Сохраняем данные в кэше в виде строки
		if err := f.cache.Set(ctx.Context(), cacheKey, string(ctx.Response().Body()), expiration).Err(); err != nil {
			return err
		}

		return nil
	}
}

// RouteCache создает кэширование на уровне маршрутов.
func (f *CacheMiddleware) RouteCache(cacheTime uint) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// Создаем ключ кэша на основе URL запроса
		cacheKey := fmt.Sprintf("route_cache_%s", ctx.OriginalURL())

		// Попытка получения кэшированных данных по ключу
		cached, err := f.cache.Get(ctx.Context(), cacheKey).Result()
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

		// Проверяем статус ответа после выполнения контроллера
		resStatus := ctx.Response().StatusCode()
		if resStatus < 200 || resStatus >= 300 {
			// Не кэшируем неуспешные ответы, такие как 401 или 403
			return nil
		}

		// Сохраняем ответ в кэш на указанное время
		if cacheTime == 0 {
			cacheTime = 5
		}
		expiration := time.Duration(cacheTime) * time.Minute

		// Сохраняем данные в кэше в виде строки
		if err := f.cache.Set(ctx.Context(), cacheKey, string(ctx.Response().Body()), expiration).Err(); err != nil {
			return err
		}

		return nil
	}
}

// CacheInvalidationMiddleware удаляет кэш при выполнении операций, изменяющих данные.
func (f *CacheMiddleware) CacheInvalidation(cacheablePaths []string) fiber.Handler {
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
						keys, nextCursor, err := f.cache.Scan(ctx.Context(), cursor, prefix+"*", 10).Result()
						if err != nil {
							return err
						}

						// Удаляем найденные ключи
						if len(keys) > 0 {
							if err := f.cache.Del(ctx.Context(), keys...).Err(); err != nil {
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
