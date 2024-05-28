package middlewares

import (
	"log"

	"github.com/IT-RushCode/rush_pkg/utils"
	"github.com/gofiber/fiber/v2"
)

// RequestResponseLogger логирует запросы и ответы
func RequestResponseLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Method() != "OPTIONS" {
			log.Printf("%s----------- REQUEST ----------->%s\n", utils.Yellow, utils.Reset)
			log.Printf("Client`s IP - %s%s%s\n", utils.LightGreen, c.IP(), utils.Reset)
			log.Printf("API Path - %s%s%s\n", utils.Blue, c.Path(), utils.Reset)
			log.Printf("Method - %s%s%s\n", utils.Red, c.Method(), utils.Reset)
			if query := string(c.Context().URI().QueryString()); query != "" {
				log.Printf("Query: %s%s%s\n", utils.Magenta, query, utils.Reset)
			}
			if body := string(c.Request().Body()); body != "" {
				log.Printf("Body: \n%s%s%s\n", utils.Blue, body, utils.Reset)
			}

			err := c.Next()

			log.Printf("%s----------- RESPONSE ----------->%s\n", utils.Yellow, utils.Reset)
			log.Printf("Status Code - %s%d%s\n", utils.Red, c.Response().StatusCode(), utils.Reset)
			if body := string(c.Response().Body()); body != "" {
				log.Printf("Response Body: \n%s%s%s\n", utils.Purple, body, utils.Reset)
			}

			return err
		}
		return c.Next()
	}
}
