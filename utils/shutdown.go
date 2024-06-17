package utils

import (
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
)

// GracefulShutdown плавно отключает Fiber сервер.
func GracefulShutdown(app *fiber.App, signalCh chan os.Signal, timeout time.Duration) {
	<-signalCh

	if err := app.Shutdown(); err != nil {
		log.Printf("Ошибка при завершении работы Fiber-сервиса: %v", err)
	}

	log.Println("Сервер Fiber корректно завершен.")
}
