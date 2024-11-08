package authRoutes

import (
	"github.com/gofiber/fiber/v2"
	controller "backend/internal/controllers/authController"
	"backend/internal/middleware/authMiddleware"
)

func register(c *fiber.Ctx) error {
	return controller.RegisterUser(c)
}

func activate(c *fiber.Ctx) error {
	return controller.ActivateAccount(c)
}

func resend(c *fiber.Ctx) error {
	return controller.ResendActivationCode(c)
}

func login(c *fiber.Ctx) error {
	return controller.Login(c)
}

func logout(c *fiber.Ctx) error {
	return controller.Logout(c)
}

func refresh(c *fiber.Ctx) error {
	return controller.Refresh(c)
}

func SetupRoutes(app fiber.Router) {
	auth := app.Group("/auth")
	auth.Post("/register", register)
	auth.Post("/activate", authMiddleware.AuthMiddleware, activate)
	auth.Get("/resend", authMiddleware.AuthMiddleware, resend)
	auth.Post("/login", login)
	auth.Get("/logout", logout)
	auth.Get("/refresh", refresh)
}