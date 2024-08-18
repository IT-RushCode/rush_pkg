package middlewares

import (
	"strings"

	"github.com/IT-RushCode/rush_pkg/config"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func InitMiddlewares(app *fiber.App, apiVersion string, cfg *config.Config) {
	// Recover middleware - ловит panic ответы
	app.Use(recover.New(recover.Config{
		// Конфигурация по умолчанию
		// Next:             nil,
		// EnableStackTrace: false,
		// StackTraceHandler: defaultStackTraceHandler,
	}))

	// Pprof middleware - профилировщик
	app.Use(pprof.New())

	// Prometheus middleware
	app.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler())) // prometheus

	// RequestID middleware - создает уникальный ID запроса, для отлова метрики трассировки и логгирования
	app.Use(requestid.New(requestid.Config{
		// Header: "",
		// Generator: func() string {
		// },
		// ContextKey: nil,
	}))

	// Logger middleware
	// app.Use(middlewares.RequestResponseLogger()) - кастомный логгер в stdout
	app.Use(logger.New(
		logger.Config{
			// Конфигурация по умолчанию
			// Next:          nil,
			// Done:          nil,
			Format: "${time} | ${status} | ${latency} | ${ip} | ${method} | ${path} | ${queryParams} | ${error}\n",
			// TimeFormat:    "15:04:05",
			// TimeZone:      "Local",
			// TimeInterval:  500 * time.Millisecond,
			// Output:        os.Stdout,
			// DisableColors: false,
		},
	))

	// Cache middleware
	// app.Use(cache.New(cache.Config{
	// 	// Конфигурация по умолчанию
	// 	// Next:         nil,
	// 	// Expiration:   1 * time.Minute,
	// 	// CacheHeader:  "X-Cache",
	// 	// CacheControl: false,
	// 	// KeyGenerator: func(c *fiber.Ctx) string {
	// 	// 	return utils.CopyString(c.Path())
	// 	// },
	// 	// ExpirationGenerator:  nil,
	// 	// StoreResponseHeaders: false,
	// 	// Storage:              nil,
	// 	// MaxBytes:             0,
	// 	// Methods:              []string{fiber.MethodGet, fiber.MethodHead},
	// }))

	// HELMET middleware -
	app.Use(helmet.New(helmet.Config{
		HSTSMaxAge:                0,
		HSTSExcludeSubdomains:     false,
		ContentSecurityPolicy:     "default-src 'self'; script-src 'self'; style-src 'self'; img-src 'self' data:; font-src 'self';",
		CSPReportOnly:             false,
		HSTSPreloadEnabled:        false,
		PermissionPolicy:          "",
		XSSProtection:             "1",
		ContentTypeNosniff:        "nosniff",
		XFrameOptions:             "SAMEORIGIN",
		ReferrerPolicy:            "no-referrer",
		CrossOriginEmbedderPolicy: "require-corp",
		CrossOriginOpenerPolicy:   "same-origin",
		CrossOriginResourcePolicy: "same-origin",
		OriginAgentCluster:        "?1",
		XDNSPrefetchControl:       "off",
		XDownloadOptions:          "noopen",
		XPermittedCrossDomain:     "none",
	}))

	// CORS middleware
	app.Use(cors.New(cors.Config{
		AllowHeaders:     "Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin,Authorization",
		AllowCredentials: false, // false для тестового окружения
		AllowOrigins:     "*",   // не может быть wildcard - "*", если AllowCredentials == true
		AllowMethods:     strings.Join(fiber.DefaultMethods, ","),
		// ExposeHeaders: "",
		// MaxAge:        0,
	}))

	// Compressor (сжатие) middleware
	app.Use(compress.New(
		compress.Config{
			Level: compress.LevelDefault,
		}))

	// Auth middleware
	// authMiddleware := NewAuthMiddleware(cfg, nil)
	// app.Use(authMiddleware.VerifyToken)

	// Подключение Swagger для документации API
	app.Use(swagger.New(swagger.Config{
		BasePath: apiVersion,
		FilePath: "./docs/swagger.yaml",
		Path:     "docs",
		Title:    "Arvand API GATEWAY",
		// FileContent: ,
		CacheAge: 3600,
	}))
}
