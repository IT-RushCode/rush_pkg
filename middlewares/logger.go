package middlewares

import (
	"os"
	"time"

	"github.com/IT-RushCode/rush_pkg/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type LoggerMiddleware struct {
	logDir    string
	logFile   string
	debug     bool
	logOutput *os.File
}

// NewLoggerMiddleware создает новый экземпляр LoggerMiddleware.
func NewLoggerMiddleware(cfg *config.AppConfig, logDir, logFile string) (*LoggerMiddleware, error) {
	// Создаем директорию для логов, если её нет
	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		return nil, err
	}

	// Открываем файл для логов
	file, err := os.OpenFile(logDir+"/"+logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	return &LoggerMiddleware{
		logDir:    logDir,
		logFile:   logFile,
		debug:     cfg.DEBUG,
		logOutput: file,
	}, nil
}

// Logger middleware - логирование запросов и ответов
//
// // Пример использования мидлвейра:
//
//	loggerMiddleware, err := rpm.NewLoggerMiddleware(&cfg.APP, "./logs", "server.log")
//
//	if err != nil {
//			log.Fatalf("Ошибка создания логгера: %v", err)
//	}
//
// // Используем мидлвейр логгера
//
//	app.Use(loggerMiddleware.Middleware())
//
// // Закрываем файл логов при завершении приложения
//
//	app.Hooks().OnShutdown(func() error {
//		return loggerMiddleware.Close()
//	})
func (m *LoggerMiddleware) Middleware() fiber.Handler {
	return logger.New(logger.Config{
		Format:        "${time} | ${status} | ${latency} | ${ip} | ${method} | ${path} | ${queryParams} | ${body} | ${response}\n",
		TimeFormat:    time.RFC3339,
		TimeZone:      "Local",
		Output:        m.logOutput,
		DisableColors: !m.debug, // Отключение цветов, если debug отключен
	})
}

// Close закрывает файл логов.
func (m *LoggerMiddleware) Close() error {
	return m.logOutput.Close()
}
