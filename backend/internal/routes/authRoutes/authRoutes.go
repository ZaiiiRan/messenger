package authRoutes

import (
	"github.com/gofiber/fiber/v2"
	controller "backend/internal/controllers/authController"
)

func register(c *fiber.Ctx) error {
	return controller.RegisterUser(c)
}

func SetupRoutes(app fiber.Router) {
	auth := app.Group("/auth")
	auth.Post("/register", register)
}