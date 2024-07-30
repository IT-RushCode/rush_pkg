package utils

import (
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// RequestResponseLogger логирует запросы и ответы
func RequestResponseLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Method() != "OPTIONS" {
			contentType := string(c.Request().Header.ContentType())
			if !strings.Contains(contentType, "multipart/form-data") && !strings.Contains(contentType, "application/octet-stream") {
				log.Printf("%s----------- REQUEST ----------->%s\n", Yellow, Reset)
				log.Printf("Client`s IP - %s%s%s\n", LightGreen, c.IP(), Reset)
				log.Printf("API Path - %s%s%s\n", Blue, c.Path(), Reset)
				log.Printf("Method - %s%s%s\n", Red, c.Method(), Reset)
				if query := string(c.Context().URI().QueryString()); query != "" {
					log.Printf("Query: %s%s%s\n", Magenta, query, Reset)
				}
				if body := string(c.Request().Body()); body != "" {
					log.Printf("Body: \n%s%s%s\n", Blue, body, Reset)
				}
			}

			err := c.Next()

			responseContentType := string(c.Response().Header.ContentType())
			if !strings.Contains(responseContentType, "multipart/form-data") && !strings.Contains(responseContentType, "application/octet-stream") {
				log.Printf("%s----------- RESPONSE ----------->%s\n", Yellow, Reset)
				log.Printf("Status Code - %s%d%s\n", Red, c.Response().StatusCode(), Reset)
				if body := string(c.Response().Body()); body != "" {
					log.Printf("Response Body: \n%s%s%s\n", Purple, body, Reset)
				}
			}

			return err
		}
		return c.Next()
	}
}
