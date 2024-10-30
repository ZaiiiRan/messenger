package main

import (
	"os"
	"log"
	"github.com/gofiber/fiber/v2"
	"github.com/lpernett/godotenv"
	"backend/internal/dbs/pgDB"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	pgDB.Connect(os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_PORT"))

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