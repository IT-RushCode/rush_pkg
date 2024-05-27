package middlewares

import (
	"log"

	color "github.com/IT-RushCode/rush_pkg/utils"
	"github.com/gofiber/fiber/v2"
)

// RequestResponseLogger логирует запросы и ответы
func RequestResponseLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Method() != "OPTIONS" {
			log.Printf("%s----------- REQUEST ----------->%s\n", color.Yellow, color.Reset)
			log.Printf("Client`s IP - %s%s%s\n", color.LightGreen, c.IP(), color.Reset)
			log.Printf("API Path - %s%s%s\n", color.Blue, c.Path(), color.Reset)
			log.Printf("Method - %s%s%s\n", color.Red, c.Method(), color.Reset)
			if query := string(c.Context().URI().QueryString()); query != "" {
				log.Printf("Query: %s%s%s\n", color.Magenta, query, color.Reset)
			}
			if body := string(c.Request().Body()); body != "" {
				log.Printf("Body: \n%s%s%s\n", color.Blue, body, color.Reset)
			}

			err := c.Next()

			log.Printf("%s----------- RESPONSE ----------->%s\n", color.Yellow, color.Reset)
			log.Printf("Status Code - %s%d%s\n", color.Red, c.Response().StatusCode(), color.Reset)
			if body := string(c.Response().Body()); body != "" {
				log.Printf("Response Body: \n%s%s%s\n", color.Blue, body, color.Reset)
			}

			return err
		}
		return c.Next()
	}
}
