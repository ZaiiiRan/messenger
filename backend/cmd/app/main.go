package main

import (
	"backend/internal/dbs/pgDB"
	"backend/internal/middleware/errorHandler"
	"backend/internal/routes/authRoutes"
	"backend/internal/routes/socialRoutes"
	"backend/internal/routes/chatRoutes"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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
	clientUrl := os.Getenv("CLIENT_URL")

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowHeaders:     "Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin,Authorization",
		AllowOrigins:     clientUrl, 
		AllowCredentials: true,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	}))
	app.Use(errorHandler.ErrorHandler())

	authRoutes.SetupRoutes(app)
	socialRoutes.SetupRoutes(app)
	chatRoutes.SetupRoutes(app)

	log.Fatal(app.Listen(":" + port))
}