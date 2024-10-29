package main

import (
	"os"
	"log"
	"github.com/gofiber/fiber/v2"
	"github.com/lpernett/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{ "message": "Hello World!" })
	})

	log.Fatal(app.Listen(":" + port))
}