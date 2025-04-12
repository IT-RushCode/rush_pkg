package middlewares

import (
	"time"

	"github.com/IT-RushCode/rush_pkg/config"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func ZapLoggerMiddleware(logger *zap.Logger, cfg *config.LogConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Создаем логгер для запроса без caller
		reqLogger := logger.With(
			zap.String("request_id", c.GetRespHeader(fiber.HeaderXRequestID)),
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.String("client_ip", c.IP()),
		)

		// Логируем начало запроса
		reqLogger.Info("-- REQUEST --")

		// Логируем тело запроса, если нужно
		if cfg.LogRequestBody && len(c.Body()) > 0 {
			reqLogger.Info("request body", zap.String("body", string(c.Body())))
		}

		// Засекаем время выполнения
		start := time.Now()
		err := c.Next()
		latency := time.Since(start)

		// Логируем завершение запроса
		reqLogger.Info("-- RESPONSE --",
			zap.Int("status", c.Response().StatusCode()),
			zap.String("latency", latency.String()),
		)

		return err
	}
}
