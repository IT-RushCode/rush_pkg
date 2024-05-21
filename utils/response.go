package utils

import (
	"github.com/gofiber/fiber/v2"
)

type Response struct {
	Status  bool        `json:"status"`
	Body    interface{} `json:"body"`
	Message string      `json:"message"`
}

func SendResponse(ctx *fiber.Ctx, success bool, message string, body interface{}, statusCode int) error {
	response := Response{
		Status:  success,
		Body:    body,
		Message: message,
	}

	ctx.Status(statusCode)

	if err := ctx.JSON(response); err != nil {
		ctx.Status(fiber.StatusInternalServerError)
		return err
	}

	return nil
}

func SuccessResponse(ctx *fiber.Ctx, message string, body interface{}) error {
	return SendResponse(ctx, true, message, body, fiber.StatusOK)
}

func ErrorResponse(ctx *fiber.Ctx, message string, body interface{}) error {
	return SendResponse(ctx, false, message, body, fiber.StatusInternalServerError)
}
