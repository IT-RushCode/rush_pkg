package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

type Middlewares struct {
	Permission *PermissionMiddleware
	Cache      *CacheMiddleware
}

func NewMiddlewares(app *fiber.App, checker PermissionChecker, cache *redis.Client) *Middlewares {
	return &Middlewares{
		Permission: NewPermissionMiddleware(checker),
		Cache:      NewCacheMiddleware(cache),
	}
}
