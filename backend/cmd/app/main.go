package main

import (
	"backend/internal/dbs/pgDB"
	"backend/internal/middleware/errorHandler"
	"backend/internal/routes/authRoutes"
	"backend/internal/routes/socialRoutes"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/lpernett/godotenv"
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
	app.Use(errorHandler.ErrorHandler())

	authRoutes.SetupRoutes(app)
	socialRoutes.SetupRoutes(app)

	log.Fatal(app.Listen(":" + port))
}