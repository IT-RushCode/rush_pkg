package utils

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

var (
	fileLogger *log.Logger
)

// InitFileLogger инициализирует логгер для записи в файл.
func InitFileLogger(logDir, logFile string) error {
	// Создаем директорию для логов, если её нет
	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		return err
	}

	// Открываем файл для логов
	file, err := os.OpenFile(fmt.Sprintf("%s/%s", logDir, logFile), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	// Создаем новый логгер для файла
	fileLogger = log.New(file, "", log.LstdFlags)
	return nil
}

// RequestResponseLogger логирует запросы, ответы и вычисляет время выполнения.
//
//	logToFile: записывать ли логи в файл.
//
//	logRequestBody: логировать ли тело запроса.
//
//	logResponseBody: логировать ли тело ответа.
func RequestResponseLogger(logToFile, logRequestBody, logResponseBody bool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Method() != "OPTIONS" {
			startTime := time.Now() // Засекаем время начала запроса
			contentType := string(c.Request().Header.ContentType())

			var logMessage, consoleMessage string

			// Проверка на бинарный контент
			if !isBinaryContent(contentType) {
				// Логирование текстовых данных запроса
				consoleMessage = fmt.Sprintf("\033[33m----------- REQUEST ----------->\033[0m\n")
				consoleMessage += fmt.Sprintf("IP | Path | Method —> \033[92m%s\033[0m | \033[34m%s\033[0m | \033[31m%s\033[0m\n", c.IP(), c.Path(), c.Method())

				logMessage = "----------- REQUEST ----------->\n"
				logMessage += fmt.Sprintf("IP | Path | Method —> %s | %s | %s\n", c.IP(), c.Path(), c.Method())

				if query := string(c.Context().URI().QueryString()); query != "" {
					consoleMessage += fmt.Sprintf("Query: \033[35m%s\033[0m\n", query)
					logMessage += fmt.Sprintf("Query: %s\n", query)
				}
				if logRequestBody {
					if body := string(c.Request().Body()); body != "" {
						consoleMessage += fmt.Sprintf("Body: \n\033[34m%s\033[0m\n", body)
						logMessage += fmt.Sprintf("Body: \n%s\n", body)
					}
				}
			} else {
				// Логирование метаданных файлового запроса
				consoleMessage = fmt.Sprintf("\033[33m----------- FILE UPLOAD REQUEST ----------->\033[0m\n")
				consoleMessage += fmt.Sprintf("IP | Path | Method —> \033[92m%s\033[0m | \033[34m%s\033[0m | \033[31m%s\033[0m\n", c.IP(), c.Path(), c.Method())

				logMessage = "----------- FILE UPLOAD REQUEST ----------->\n"
				logMessage += fmt.Sprintf("IP | Path | Method —> %s | %s | %s\n", c.IP(), c.Path(), c.Method())

				form, err := c.MultipartForm()
				if err == nil && form != nil {
					for fileKey, fileHeaders := range form.File {
						if len(fileHeaders) > 0 {
							fileHeader := fileHeaders[0]
							consoleMessage += fmt.Sprintf("Key: \033[35m%s\033[0m | File Name: \033[34m%s\033[0m | File Size: \033[92m%d bytes\033[0m | Content-Type: \033[31m%s\033[0m\n",
								fileKey, fileHeader.Filename, fileHeader.Size, fileHeader.Header.Get("Content-Type"))

							logMessage += fmt.Sprintf("Key: %s | File Name: %s | File Size: %d bytes | Content-Type: %s\n",
								fileKey, fileHeader.Filename, fileHeader.Size, fileHeader.Header.Get("Content-Type"))
						}
					}
				}
			}

			// Логируем в консоль
			log.Print(consoleMessage)
			// Логируем в файл
			if logToFile && fileLogger != nil {
				fileLogger.Print(stripANSI(logMessage))
			}

			// Выполнение следующего обработчика
			err := c.Next()

			// После выполнения запроса фиксируем время завершения
			endTime := time.Now()
			latency := endTime.Sub(startTime) // Вычисляем латентность

			responseContentType := string(c.Response().Header.ContentType())

			// Логирование ответа
			if !isBinaryContent(responseContentType) {
				consoleMessage = fmt.Sprintf("\033[33m----------- RESPONSE ----------->\033[0m\n")
				consoleMessage += fmt.Sprintf("Status | Latency —> \033[31m%d\033[0m | \033[35m%v\033[0m\n", c.Response().StatusCode(), latency)

				logMessage = "----------- RESPONSE ----------->\n"
				logMessage += fmt.Sprintf("Status | Latency —> %d | %v\n", c.Response().StatusCode(), latency)

				if logResponseBody {
					if body := string(c.Response().Body()); body != "" {
						consoleMessage += fmt.Sprintf("Response Body: \033[34m%s\033[0m\n", body)
						logMessage += fmt.Sprintf("Response Body: %s\n", body)
					}
				}
			} else {
				consoleMessage = fmt.Sprintf("\033[33m----------- FILE RESPONSE ----------->\033[0m\n")
				consoleMessage += fmt.Sprintf("Status | Latency —> \033[31m%d\033[0m | \033[35m%v\033[0m\n", c.Response().StatusCode(), latency)
				consoleMessage += fmt.Sprintf("Response Content-Type: \033[34m%s\033[0m\n", responseContentType)
				consoleMessage += fmt.Sprintf("Response Content-Length: \033[34m%d bytes\033[0m\n", c.Response().Header.ContentLength())

				logMessage = "----------- FILE RESPONSE ----------->\n"
				logMessage += fmt.Sprintf("Status | Latency —> %d | %v\n", c.Response().StatusCode(), latency)
				logMessage += fmt.Sprintf("Response Content-Type: %s\n", responseContentType)
				logMessage += fmt.Sprintf("Response Content-Length: %d bytes\n", c.Response().Header.ContentLength())
			}

			// Логируем в консоль
			log.Print(consoleMessage)
			// Логируем в файл
			if logToFile && fileLogger != nil {
				fileLogger.Print(stripANSI(logMessage))
			}

			return err
		}
		return c.Next()
	}
}

// isBinaryContent проверяет, является ли контент бинарным.
func isBinaryContent(contentType string) bool {
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

// stripANSI удаляет ANSI-коды для цвета из строки.
func stripANSI(input string) string {
	replacer := strings.NewReplacer(
		"\033[33m", "", // Жёлтый
		"\033[0m", "", // Сброс
		"\033[92m", "", // Светло-зелёный
		"\033[34m", "", // Синий
		"\033[31m", "", // Красный
		"\033[35m", "", // Фиолетовый
	)
	return replacer.Replace(input)
}
