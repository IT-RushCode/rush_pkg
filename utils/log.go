package utils

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

func RequestConsoleLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {

		if c.Method() != "OPTIONS" {
			log.Printf("%s----------- REQUEST ----------->%s\n", Yellow, Reset)
			log.Printf("Client`s IP - %s%s%s\n", LightGreen, c.IP(), Reset)
			log.Printf("API Path - %s%s%s\n", Blue, c.Path(), Reset)
			log.Printf("Method - %s%s%s\n", Red, c.Method(), Reset)
			if string(c.Context().URI().QueryString()) != "" {
				fmt.Printf("Query: %s%s%s\n", Magenta, c.Context().URI().QueryString(), Reset)
			}
			if string(c.Request().Body()) != "" {
				fmt.Printf("Body: \n%s%s%s\n", Blue, string(c.Request().Body()), Reset)
			}
		}
		return c.Next()
	}
}
