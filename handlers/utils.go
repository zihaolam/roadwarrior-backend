package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type IErrorResponse struct {
	Message string `json:"message"`
}

func ErrorResponse(c *fiber.Ctx, statusCode int, err interface{}) error {
	switch err := err.(type) {
	case string:
		return c.Status(statusCode).JSON(&IErrorResponse{Message: err})

	case error:
		return c.Status(statusCode).JSON(&IErrorResponse{Message: err.Error()})

	default:
		return c.Status(statusCode).JSON(&IErrorResponse{Message: fmt.Sprintf("%s", err)})
	}
}

func OKResponse(c *fiber.Ctx, body interface{}) error {
	return c.JSON(body)
}
