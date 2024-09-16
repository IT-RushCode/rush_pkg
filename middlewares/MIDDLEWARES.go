package middlewares

type Middlewares struct {
	Permission *PermissionMiddleware
	Cache      *CacheMiddleware
}
