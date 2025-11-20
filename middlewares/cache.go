package middlewares

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

// CacheMiddleware содержит все middleware для кэширования и их конфигурации.
type CacheMiddleware struct{ *CacheConfig }

// CacheConfig описывает минимальный набор параметров для работы кэша.
type CacheConfig struct {
	cache       *redis.Client
	ServiceName string // имя сервиса, которое добавляется в начало каждого ключа
	ActiveCache bool   // включено ли кэширование
	CacheTime   int64  // глобальное время кэширования (минуты)
	// TagDependencies позволяет описать связи между тегами.
	TagDependencies map[string][]string
	// UserScopedPaths - список путей, к которым нужно привязывать кэш к UserID.
	// По умолчанию содержит /auth/me; можно переопределить при инициализации.
	UserScopedPaths []string
}

var (
	DefaultCacheConfig = &CacheConfig{
		ServiceName:     "RushApp",
		ActiveCache:     true,
		CacheTime:       60,
		TagDependencies: map[string][]string{},
		UserScopedPaths: []string{"/auth/me"},
	}
)

// Пример комплексной настройки:
//
//	func InitMiddlewares(...) *rpm.Middlewares {
//		cacheCfg := *rpm.DefaultCacheConfig
//		cacheCfg.ServiceName = "RushApp"
//		cacheCfg.ActiveCache = cfg.APP.CACHE_ACTIVE
//
//		cacheMd := rpm.NewCacheMiddleware(redisClient, &cacheCfg)
//
//		cacheMd.RegisterTagDependencies(map[string][]string{
//			"products":        {"restaurants:list", "restaurants:details"},
//			"restaurants":     {"restaurants:list", "restaurants:menu"},
//			"restaurants:list": {"restaurants:menu"},
//			"categories":      {"courses", "lessons"},
//		})
//
//		return &rpm.Middlewares{
//			Cache: cacheMd,
//		}
//	}
//
//	func registerCategoryRoutes(router fiber.Router, m *rpm.Middlewares, ctrl *controllers.Controllers) {
//		category := router.Group("/categories")
//
//		category.Get("/",
//			m.Permission.CheckPermission("view:categories"),
//			m.Cache.RouteCache(10, "categories:list"),
//			ctrl.Category.GetAllCategories,
//		)
//
//		category.Get("/:id",
//			m.Permission.CheckPermission("view:category_by_id"),
//			m.Cache.RouteCache(10, "categories:details"),
//			ctrl.Category.FindCategoryByID,
//		)
//
//		category.Post("/",
//			m.Permission.CheckPermission("create:category"),
//			m.Cache.InvalidateTags("categories:list"),
//			ctrl.Category.CreateCategory,
//		)
//
//		category.Put("/:id",
//			m.Permission.CheckPermission("update:category"),
//			m.Cache.InvalidateTags("categories:list", "categories:details"),
//			ctrl.Category.UpdateCategory,
//		)
//	}
//
//	func registerCourseRoutes(router fiber.Router, m *rpm.Middlewares, ctrl *controllers.Controllers) {
//		course := router.Group("/courses")
//
//		course.Get("/",
//			m.Permission.CheckPermission("view:courses"),
//			m.Cache.RouteCache(15, "courses:list"),
//			ctrl.Course.GetAllCourses,
//		)
//
//		course.Get("/:id",
//			m.Permission.CheckPermission("view:course_by_id"),
//			m.Cache.RouteCache(15, "courses:details"),
//			ctrl.Course.FindCourseByID,
//		)
//
//		course.Post("/",
//			m.Permission.CheckPermission("create:course"),
//			m.Cache.InvalidateTags("courses", "categories"),
//			ctrl.Course.CreateCourse,
//		)
//
//		course.Patch("/:id",
//			m.Permission.CheckPermission("change:course_status"),
//			m.Cache.InvalidateTags("courses:details", "courses:list"),
//			ctrl.Course.ChangeStatusCourse,
//		)
//	}
//
// NewCacheMiddleware создаёт новый экземпляр CacheMiddleware с Redis клиентом.
//
// Если передан nil конфиг, то будет использоваться DefaultCacheConfig.
func NewCacheMiddleware(cacheStorage *redis.Client, cacheConfig *CacheConfig) *CacheMiddleware {
	if cacheStorage == nil {
		log.Fatal("redis cache storage is required")
	}

	cfg := cacheConfig
	if cfg == nil {
		def := *DefaultCacheConfig
		cfg = &def
	}

	if cfg.CacheTime <= 0 {
		cfg.CacheTime = DefaultCacheConfig.CacheTime
	}

	if cfg.TagDependencies == nil {
		cfg.TagDependencies = map[string][]string{}
	}
	if len(cfg.UserScopedPaths) == 0 {
		cfg.UserScopedPaths = append([]string{}, DefaultCacheConfig.UserScopedPaths...)
	}

	cfg.ServiceName = sanitizeSegment(cfg.ServiceName)
	if cfg.ServiceName == "" {
		cfg.ServiceName = sanitizeSegment(DefaultCacheConfig.ServiceName)
	}

	cfg.cache = cacheStorage

	return &CacheMiddleware{CacheConfig: cfg}
}

// RouteCache создает кэширование на уровне маршрутов.
//
//	cacheTime := 60 // cacheTime принимает только минуты
//	tags := []string{"products", "products:list"} // теги для последующей инвалидации
//
//	Если cacheTime на уровне роута > 0, то он перекрывает глобальное значение.
//	Теги используются для точечного сброса кэшей. Если они не переданы, в качестве тега используется путь роута.
//
//	// Пример использования в роутере:
//	api.Get("/restaurants",
//		m.Cache.RouteCache(5, "restaurants:list"),
//		ctrl.Restaurant.List,
//	)
func (f *CacheMiddleware) RouteCache(cacheTime int64, tags ...string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		if !f.ActiveCache {
			return ctx.Next()
		}

		cacheKey := f.buildCacheKey(ctx)
		normalizedTags := f.normalizeTags(ctx, tags)

		// Попытка получить данные из кэша
		cached, err := f.cache.Get(ctx.Context(), cacheKey).Bytes()
		if err == nil && len(cached) > 0 {
			var payload struct {
				ContentType string `json:"ct,omitempty"`
				Body        []byte `json:"b"`
			}
			if unmarshalErr := json.Unmarshal(cached, &payload); unmarshalErr == nil && len(payload.Body) > 0 {
				if payload.ContentType != "" {
					ctx.Set(fiber.HeaderContentType, payload.ContentType)
				}
				return ctx.Send(payload.Body)
			}

			// Старые записи без упаковки просто отдаем как есть.
			return ctx.Send(cached)
		} else if err != nil && err != redis.Nil {
			return err
		}

		if err := ctx.Next(); err != nil {
			return err
		}

		if !f.shouldCacheResponse(ctx) {
			return nil
		}

		expiration := f.resolveExpiration(cacheTime)
		payload := struct {
			ContentType string `json:"ct,omitempty"`
			Body        []byte `json:"b"`
		}{
			ContentType: string(ctx.Response().Header.ContentType()),
			Body:        append([]byte(nil), ctx.Response().Body()...),
		}
		data, marshalErr := json.Marshal(payload)
		if marshalErr != nil {
			return marshalErr
		}
		if err := f.cache.Set(ctx.Context(), cacheKey, data, expiration).Err(); err != nil {
			return err
		}
		if err := f.bindTags(ctx, normalizedTags, cacheKey, expiration); err != nil {
			return err
		}

		return nil
	}
}

// CacheInvalidation оставлена для обратной совместимости и теперь работает через теги.
// Эквивалентно InvalidateTags(cacheablePaths...).
func (f *CacheMiddleware) CacheInvalidation(cacheablePaths []string) fiber.Handler {
	return f.InvalidateTags(cacheablePaths...)
}

// InvalidateTags очищает кэши по тегам после успешного изменения данных (POST/PUT/PATCH/DELETE).
// Теги можно переиспользовать между разными маршрутами.
//
//	// Пример: при изменении ресторана очищаем список и карточку.
//	api.Put("/restaurants/:id",
//		m.Cache.InvalidateTags("restaurants:list", "restaurants:details"),
//		ctrl.Restaurant.Update,
//	)
func (f *CacheMiddleware) InvalidateTags(tags ...string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		if !f.ActiveCache {
			return ctx.Next()
		}

		if err := ctx.Next(); err != nil {
			return err
		}

		if !isMutatingMethod(ctx.Method()) || ctx.Response().StatusCode() >= 400 {
			return nil
		}

		normalizedTags := f.normalizeTags(ctx, tags)
		return f.dropTags(ctx, normalizedTags)
	}
}

// RegisterTagDependencies позволяет передать карту всех зависимостей единым блоком.
// Удобно объявить её отдельно (например, в routes/cache_tags.go) и применять при инициализации middleware.
//
//	// Пример:
//	var CacheTagDependencies = map[string][]string{
//		"products":    {"restaurants:list", "restaurants:details"},
//		"restaurants": {"restaurants:list", "restaurants:menu"},
//	}
//
//	cacheMd.RegisterTagDependencies(CacheTagDependencies)
func (f *CacheMiddleware) RegisterTagDependencies(deps map[string][]string) {
	if len(deps) == 0 {
		return
	}
	for tag, chained := range deps {
		tag = strings.TrimSpace(tag)
		if tag == "" {
			continue
		}
		if f.TagDependencies == nil {
			f.TagDependencies = map[string][]string{}
		}
		f.TagDependencies[tag] = append(f.TagDependencies[tag], chained...)
	}
}

func (f *CacheMiddleware) shouldCacheResponse(ctx *fiber.Ctx) bool {
	status := ctx.Response().StatusCode()
	if status < 200 || status >= 300 {
		return false
	}

	return len(ctx.Response().Body()) > 0
}

func (f *CacheMiddleware) resolveExpiration(cacheTime int64) time.Duration {
	if cacheTime <= 0 {
		cacheTime = f.CacheTime
	}
	if cacheTime <= 0 {
		cacheTime = 1
	}
	return time.Duration(cacheTime) * time.Minute
}

// buildCacheKey формирует ключ кэша.
// Публичные и одинаковые для всех ответы получают единый ключ,
// но user-specific ручки (например, auth/me) персонализируются по UserID,
// чтобы не было утечек между пользователями при кэшировании.
func (f *CacheMiddleware) buildCacheKey(ctx *fiber.Ctx) string {
	pathSegment := sanitizeSegment(ctx.Path())
	hash := hashURL(ctx.OriginalURL())
	userSegment := ""
	if isPublic, _ := ctx.Locals("IsPublic").(bool); !isPublic {
		// Персонализируем кэш только для user-specific ручек.
		if f.isUserScopedPath(ctx.Path()) {
			if userID, ok := ctx.Locals("UserID").(uint); ok && userID > 0 {
				userSegment = fmt.Sprintf("user_%d:", userID)
			}
		}
	}
	if userSegment == "" {
		return fmt.Sprintf("%s:%s:%s", f.ServiceName, pathSegment, hash)
	}
	return fmt.Sprintf("%s:%s%s:%s", f.ServiceName, userSegment, pathSegment, hash)
}

func (f *CacheMiddleware) tagKey(tag string) string {
	return fmt.Sprintf("%s:tag:%s", f.ServiceName, sanitizeSegment(tag))
}

func (f *CacheMiddleware) isUserScopedPath(path string) bool {
	path = strings.TrimRight(path, "/")
	for _, scoped := range f.UserScopedPaths {
		scoped = strings.TrimRight(strings.TrimSpace(scoped), "/")
		if scoped == "" {
			continue
		}

		if path == scoped || strings.HasSuffix(path, scoped) {
			return true
		}
	}
	return false
}

func (f *CacheMiddleware) normalizeTags(ctx *fiber.Ctx, tags []string) []string {
	var normalized []string
	for _, tag := range tags {
		tag = strings.TrimSpace(tag)
		if tag == "" {
			continue
		}
		normalized = append(normalized, tag)
	}

	if len(normalized) == 0 {
		normalized = append(normalized, ctx.Path())
	}

	return normalized
}

func (f *CacheMiddleware) bindTags(ctx *fiber.Ctx, tags []string, cacheKey string, expiration time.Duration) error {
	for _, tag := range tags {
		tagKey := f.tagKey(tag)
		if err := f.cache.SAdd(ctx.Context(), tagKey, cacheKey).Err(); err != nil {
			return err
		}
		if expiration > 0 {
			if err := f.cache.Expire(ctx.Context(), tagKey, expiration).Err(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (f *CacheMiddleware) dropTags(ctx *fiber.Ctx, tags []string) error {
	if len(tags) == 0 {
		return nil
	}

	cascadeTags := f.expandTagCascade(tags)
	for _, tag := range cascadeTags {
		tagKey := f.tagKey(tag)
		cacheKeys, err := f.cache.SMembers(ctx.Context(), tagKey).Result()
		if err != nil && err != redis.Nil {
			return err
		}

		if len(cacheKeys) > 0 {
			if err := f.cache.Del(ctx.Context(), cacheKeys...).Err(); err != nil {
				return err
			}
		}

		if err := f.cache.Del(ctx.Context(), tagKey).Err(); err != nil {
			return err
		}
	}

	return nil
}

func hashURL(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	// Берем первые 8 байт, чтобы ключи оставались короткими и читаемыми.
	return hex.EncodeToString(sum[:8])
}

func sanitizeSegment(value string) string {
	value = strings.TrimSpace(strings.ToLower(value))
	if value == "" {
		return ""
	}

	var b strings.Builder
	b.Grow(len(value))
	lastUnderscore := false

	for _, r := range value {
		if (r >= '0' && r <= '9') || (r >= 'a' && r <= 'z') {
			b.WriteRune(r)
			lastUnderscore = false
			continue
		}

		if !lastUnderscore {
			b.WriteByte('_')
			lastUnderscore = true
		}
	}

	res := strings.Trim(b.String(), "_")
	if res == "" {
		return ""
	}

	return res
}

func isMutatingMethod(method string) bool {
	switch method {
	case fiber.MethodPost, fiber.MethodPut, fiber.MethodPatch, fiber.MethodDelete:
		return true
	default:
		return false
	}
}

func (f *CacheMiddleware) expandTagCascade(tags []string) []string {
	visited := make(map[string]struct{})
	var ordered []string

	var visit func(string)
	visit = func(tag string) {
		tag = strings.TrimSpace(tag)
		if tag == "" {
			return
		}
		if _, ok := visited[tag]; ok {
			return
		}
		visited[tag] = struct{}{}
		ordered = append(ordered, tag)
		for _, dependant := range f.TagDependencies[tag] {
			visit(dependant)
		}
	}

	for _, tag := range tags {
		visit(tag)
	}

	return ordered
}
