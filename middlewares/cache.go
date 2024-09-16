package middlewares

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

// CacheMiddleware содержит все middleware для кэширования и их конфигурации.
type CacheMiddleware struct {
	cache       *redis.Client
	activeCache bool
	cacheTime   int64
}

// NewCacheMiddleware создаёт новый экземпляр CacheMiddleware с Redis клиентом.
//
// activeCache - определяет кэшировать запросы или вызвать следующий обработчик (пропустить).
func NewCacheMiddleware(redisClient *redis.Client, activeCache bool, cacheTime int64) *CacheMiddleware {
	return &CacheMiddleware{
		cache:       redisClient,
		activeCache: activeCache,
		cacheTime:   cacheTime,
	}
}

// RouteCache создает кэширование на уровне маршрутов.
//
// Если на уровне роута cacheTime > 0, то указанное время будет приоритетом, иначе будет используется глобальный cacheTime мидлвейра.
func (f *CacheMiddleware) RouteCache(cacheTime int64) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// Если кэширование не активно, пропускаем кэширование и выполняем следующий обработчик
		if !f.activeCache {
			return ctx.Next()
		}

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

		if cacheTime == 0 {
			cacheTime = f.cacheTime
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
					prefix := fmt.Sprintf("route_cache_%s", path)

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
