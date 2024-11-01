package errorHandler

import (
	"errors"
	"backend/internal/errors/appError"

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

			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "internal server error",
			})
		}

		return nil
	}
}
