package middlewares

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"strings"
	"time"

	"github.com/IT-RushCode/rush_pkg/config"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func ZapLoggerMiddleware(logger *zap.Logger, cfg *config.LogConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Path() == "/metrics" {
			// Пропускаем логирование метрик
			return c.Next()
		}

		// Создаем логгер для запроса без caller
		reqLogger := logger.With(
			zap.String("request_id", c.GetRespHeader(fiber.HeaderXRequestID)),
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.String("client_ip", c.IP()),
		)

		// Логируем начало запроса
		reqLogger.Info("[ REQUEST ]")

		// Логируем тело запроса, если нужно
		if cfg.LogRequestBody && len(c.Body()) > 0 {
			contentType := c.Get("Content-Type")
			if strings.HasPrefix(contentType, "multipart/form-data") {
				bodyStr := logMultipartForm(c)
				reqLogger.Info("request multipart form", zap.String("fields", bodyStr))
			} else {
				reqLogger.Info("request body", zap.String("body", string(c.Body())))
			}
		}

		// Засекаем время выполнения
		start := time.Now()
		err := c.Next()
		latency := time.Since(start)

		// Логируем завершение запроса
		reqLogger.Info("[ RESPONSE ]",
			zap.Int("status", c.Response().StatusCode()),
			zap.String("latency", latency.String()),
		)

		return err
	}
}

// logMultipartForm фильтрует multipart/form-data, логируя только текстовые поля, файлы заменяет на плейсхолдер.
func logMultipartForm(c *fiber.Ctx) string {
	contentType := c.Get("Content-Type")
	boundary := ""
	// Ищем boundary
	for _, part := range strings.Split(contentType, ";") {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "boundary=") {
			boundary = strings.TrimPrefix(part, "boundary=")
			break
		}
	}
	if boundary == "" {
		return "[invalid multipart: no boundary]"
	}
	body := c.Body()
	reader := multipart.NewReader(bytes.NewReader(body), boundary)
	result := make(map[string]interface{})
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Sprintf("[multipart parse error: %v]", err)
		}
		fieldName := part.FormName()
		if fieldName == "" {
			continue
		}
		fileName := part.FileName()
		if fileName != "" {
			result[fieldName] = "[binary file omitted]"
			continue
		}
		val, err := io.ReadAll(part)
		if err != nil {
			result[fieldName] = "[error reading value]"
			continue
		}
		result[fieldName] = string(val)
	}
	// Форматируем результат как строку
	var sb strings.Builder
	sb.WriteString("{")
	i := 0
	for k, v := range result {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("%q: %q", k, v))
		i++
	}
	sb.WriteString("}")
	return sb.String()
}
