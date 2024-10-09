package utils

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// RequestResponseLogger логирует запросы, ответы и вычисляет время выполнения
func RequestResponseLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		fmt.Println()
		if c.Method() != "OPTIONS" {
			startTime := time.Now() // Засекаем время начала запроса
			contentType := string(c.Request().Header.ContentType())

			// Проверка на бинарный контент
			if !isBinaryContent(contentType) {
				// Логирование текстовых данных запроса
				log.Printf("%s----------- REQUEST ----------->%s\n", Yellow, Reset)
				log.Printf("IP | Path | Method &  —> %s%s%s | %s%s%s | %s%s%s\n", LightGreen, c.IP(), Reset, Blue, c.Path(), Reset, Red, c.Method(), Reset)
				if query := string(c.Context().URI().QueryString()); query != "" {
					log.Printf("Query: %s%s%s\n", Magenta, query, Reset)
				}
				if body := string(c.Request().Body()); body != "" {
					log.Printf("Body: \n%s%s%s\n", Blue, body, Reset)
				}
			} else {
				// Логирование метаданных файлового запроса
				log.Printf("%s----------- FILE UPLOAD REQUEST ----------->%s\n", Yellow, Reset)
				log.Printf("IP | Path | Method &  —> %s%s%s | %s%s%s | %s%s%s\n", LightGreen, c.IP(), Reset, Blue, c.Path(), Reset, Red, c.Method(), Reset)

				// Обрабатываем файлы без двойного цикла
				form, err := c.MultipartForm()
				if err == nil && form != nil {
					for fileKey, fileHeaders := range form.File {
						if len(fileHeaders) > 0 {
							fileHeader := fileHeaders[0] // Берём только первый файл
							log.Printf("Key: %s | File Name: %s | File Size: %d bytes | Content-Type: %s\n",
								fileKey, fileHeader.Filename, fileHeader.Size, fileHeader.Header.Get("Content-Type"))
						}
					}
				}
			}

			// Выполнение следующего обработчика (например, выполнение основного запроса)
			err := c.Next()

			// После выполнения запроса фиксируем время завершения
			endTime := time.Now()
			latency := endTime.Sub(startTime) // Вычисляем латентность (разницу во времени)

			responseContentType := string(c.Response().Header.ContentType())

			// Логирование ответа
			if !isBinaryContent(responseContentType) {
				log.Printf("%s----------- RESPONSE ----------->%s\n", Yellow, Reset)
				log.Printf("Status | Latency —> %s%d%s | %s%v%s\n", Red, c.Response().StatusCode(), Reset, Magenta, latency, Reset)
				if body := string(c.Response().Body()); body != "" {
					log.Printf("Response Body: %s%s%s\n", Purple, body, Reset)
				}
			} else {
				// Логирование метаданных файлового ответа
				log.Printf("%s----------- FILE RESPONSE ----------->%s\n", Yellow, Reset)
				log.Printf("Status | Latency —> %s%d%s | %s%v%s\n", Red, c.Response().StatusCode(), Reset, Magenta, latency, Reset)
				log.Printf("Response Content-Type: %s%s%s\n", Blue, responseContentType, Reset)
				log.Printf("Response Content-Length: %s%d bytes%s\n", Blue, c.Response().Header.ContentLength(), Reset)
			}

			return err
		}
		return c.Next()
	}
}

// isBinaryContent проверяет, является ли контент бинарным, чтобы исключить его логирование
func isBinaryContent(contentType string) bool {
	// Список распространённых бинарных типов контента
	binaryContentTypes := []string{
		"multipart/form-data",
		"application/octet-stream",
		"image/",
		"video/",
		"audio/",
		"application/zip",
		"application/pdf",
	}

	for _, binaryType := range binaryContentTypes {
		if strings.HasPrefix(contentType, binaryType) {
			return true
		}
	}
	return false
}
