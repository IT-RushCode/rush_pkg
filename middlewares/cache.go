package middlewares

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

// CacheMiddleware содержит все middleware для кэширования и их конфигурации.
type CacheMiddleware struct{ *CacheConfig }

type CacheConfig struct {
	cache                *redis.Client
	ActiveCache          bool     // Активно ли кэширование
	CacheTime            int64    // Глобальное время кэширования в минутах
	CacheFiles           bool     // Кэшировать ли файлы
	MaxFileSize          uint     // Максимальный размер файла (возможный максимум 10 MB)
	ExcludedContentTypes []string // Исключить видео
	AllowedContentTypes  []string // Кэшировать только изображения
}

var (
	DefaultCacheConfig = &CacheConfig{
		ActiveCache:          true,
		CacheTime:            60,
		CacheFiles:           true,
		MaxFileSize:          2 * 1024 * 1024,
		ExcludedContentTypes: DefaultExcludedContentTypes,
		AllowedContentTypes:  DefaultAllowedContentTypes,
	}

	DefaultExcludedContentTypes = []string{
		"video/mp4",             // MP4 видео
		"video/mkv",             // Matroska видео
		"video/webm",            // WebM видео
		"video/avi",             // AVI видео
		"video/mpeg",            // MPEG видео
		"video/ogg",             // Ogg видео
		"video/quicktime",       // QuickTime видео (MOV)
		"video/x-msvideo",       // Windows AVI
		"video/x-flv",           // Flash видео
		"video/x-ms-wmv",        // Windows Media видео (WMV)
		"application/x-mpegURL", // HLS список воспроизведения
		"video/3gpp",            // 3GPP видео
		"video/3gpp2",           // 3GPP2 видео
		"video/x-matroska",      // Matroska (MKV) видео
		"video/x-ms-asf",        // ASF (Advanced Systems Format)
		"video/x-ogm+ogg",       // OGM видео
	}

	DefaultAllowedContentTypes = []string{
		// Изображения
		"image/png",     // PNG
		"image/jpeg",    // JPEG
		"image/jpg",     // JPG (альтернативный заголовок)
		"image/gif",     // GIF
		"image/bmp",     // BMP
		"image/webp",    // WebP
		"image/svg+xml", // SVG
		"image/tiff",    // TIFF

		// Документы
		"application/pdf",    // PDF
		"application/msword", // Microsoft Word (DOC)
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document", // Microsoft Word (DOCX)
		"application/vnd.ms-excel", // Microsoft Excel (XLS)
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",         // Microsoft Excel (XLSX)
		"application/vnd.ms-powerpoint",                                             // Microsoft PowerPoint (PPT)
		"application/vnd.openxmlformats-officedocument.presentationml.presentation", // Microsoft PowerPoint (PPTX)
		"text/plain",       // Текстовый файл (TXT)
		"text/csv",         // CSV
		"application/json", // JSON
		"application/xml",  // XML

		// Аудио
		"audio/mpeg", // MP3
		"audio/ogg",  // Ogg Audio
		"audio/wav",  // WAV
		"audio/flac", // FLAC

		// Прочее
		"application/zip",              // ZIP архивы
		"application/x-7z-compressed",  // 7z архивы
		"application/x-tar",            // TAR архивы
		"application/gzip",             // GZIP архивы
		"application/x-rar-compressed", // RAR архивы
	}
)

// NewCacheMiddleware создаёт новый экземпляр CacheMiddleware с Redis клиентом.
//
// activeCache - определяет кэшировать запросы или вызвать следующий обработчик (пропустить).
//
// Шаблон конифга:
//
//	var CacheConfig = &CacheConfig{
//		activeCache:          true,
//		cacheTime:            60,
//		cacheFiles:           true,
//		maxFileSize:          2 * 1024 * 1024,
//		excludedContentTypes: rpm.DefaultExcludedContentTypes,
//		allowedContentTypes:  rpm.DefaultAllowedContentTypes,
//	}
func NewCacheMiddleware(cacheStorage *redis.Client, cacheConfig *CacheConfig) *CacheMiddleware {
	if cacheConfig.MaxFileSize > 10*1024*1024 {
		log.Fatal("Максимальный размер кешируемого файла слишком велик")
	}
	cacheConfig.cache = cacheStorage
	return &CacheMiddleware{CacheConfig: cacheConfig}
}

// RouteCache создает кэширование на уровне маршрутов.
//
//	cacheTime := 60 // cacheTime принимает только минуты
//
//	Если на уровне роута (т.е. m.RouteCache(60)) cacheTime > 0, то указанное время будет приоритетом, иначе используется глобальный cacheTime мидлвейра.
func (f *CacheMiddleware) RouteCache(cacheTime int64) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// Если кэширование не активно, пропускаем обработку
		if !f.ActiveCache {
			return ctx.Next()
		}

		// Создаем ключ кэша
		cacheKey := fmt.Sprintf("route_cache_%s", ctx.OriginalURL())
		metaKey := fmt.Sprintf("%s_meta", cacheKey)

		// Попытка получить данные из кэша
		cached, err := f.cache.Get(ctx.Context(), cacheKey).Bytes()
		if err == nil && len(cached) > 0 {
			// Попытка получить мета-информацию о Content-Type
			contentType, _ := f.cache.Get(ctx.Context(), metaKey).Result()
			if contentType != "" {
				// Если Content-Type найден, возвращаем кэшированные данные
				ctx.Set("Content-Type", contentType)
				return ctx.Send(cached)
			}

			// Если Content-Type отсутствует, удаляем запись из кэша
			_ = f.cache.Del(ctx.Context(), cacheKey).Err() // Удаляем данные
			_ = f.cache.Del(ctx.Context(), metaKey).Err()  // Удаляем мета-информацию

			// Перенаправляем запрос на обработчик
			return ctx.Next()
		} else if err != redis.Nil {
			// Если ошибка не связана с отсутствием ключа
			return err
		}

		// Выполняем обработчик
		err = ctx.Next()
		if err != nil {
			return err
		}

		// Проверяем статус ответа
		resStatus := ctx.Response().StatusCode()
		if resStatus < 200 || resStatus >= 300 {
			return nil // Не кэшируем ошибки
		}

		// Получаем Content-Type и тело ответа
		contentType := string(ctx.Response().Header.ContentType())
		if contentType == "" {
			fmt.Printf("Skipping cache: %s (Content-Type is empty)\n", ctx.OriginalURL())
			return nil
		}
		body := ctx.Response().Body()

		// Проверяем Content-Type на исключенные
		for _, excludedType := range f.ExcludedContentTypes {
			if strings.HasPrefix(contentType, excludedType) {
				return nil // Не кэшируем запрещенные типы
			}
		}

		// Проверяем размер файла
		if f.CacheFiles && uint(len(body)) > f.MaxFileSize {
			return nil // Не кэшируем слишком большие файлы
		}

		// Проверяем, допустим ли тип файла
		isAllowed := false
		for _, allowedType := range f.AllowedContentTypes {
			if strings.HasPrefix(contentType, allowedType) {
				isAllowed = true
				break
			}
		}
		if !isAllowed {
			return nil // Не кэшируем недопустимые типы
		}

		// Устанавливаем время кэширования
		if cacheTime == 0 {
			cacheTime = f.CacheTime
		}
		expiration := time.Duration(cacheTime) * time.Minute

		// Сохраняем Content-Type и данные в кэше
		if err := f.cache.Set(ctx.Context(), metaKey, contentType, expiration).Err(); err != nil {
			return err
		}
		if err := f.cache.Set(ctx.Context(), cacheKey, body, expiration).Err(); err != nil {
			return err
		}

		return nil
	}
}

// CacheInvalidation удаляет все кэши определенного роута при выполнении запросов, изменяющих данные (POST, PUT, PATCH, DELETE).
//
// Пример передачи роутов кэши которых надо сбрасывать:
//
//	package middlewares
//
//	// Роуты кэши которых нужно очищать
//	var (
//		CachedRoutes = []string{
//			"/api/v1/examples",
//			"/public",
//			"и т.д.",
//		}
//	)
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
