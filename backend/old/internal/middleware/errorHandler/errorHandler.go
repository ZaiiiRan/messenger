package errorHandler

import (
	"backend/internal/errors/appError"
	"errors"

	"github.com/gofiber/fiber/v2"
)

func ErrorHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := c.Next()

		if err != nil {
			var appErr *appError.AppError

			if errors.As(err, &appErr) {
				return c.Status(appErr.StatusCode).JSON(fiber.Map{
					"error": appErr.Error(),
				})
			}

			if fiberErr, ok := err.(*fiber.Error); ok {
				return c.Status(fiberErr.Code).JSON(fiber.Map{
					"error": fiberErr.Message,
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "internal server error",
			})
		}

		return nil
	}
}
