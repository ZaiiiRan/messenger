package main

import (
	"backend/internal/dbs/pgDB"
	"backend/internal/logger"
	"backend/internal/middleware/errorHandler"
	"backend/internal/routes/authRoutes"
	"backend/internal/routes/chatRoutes"
	"backend/internal/routes/socialRoutes"
	"backend/internal/routes/wsRoutes"
	"backend/internal/webSocketManager"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/lpernett/godotenv"
)

func main() {
	logger := logger.GetInstance()

	err := godotenv.Load()
	if err != nil {
		logger.Fatal("Error loading .env file", "App starting")
	}

	pgDB.SetConnectionData(os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_PORT"))
	pgDB.GetDB()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	clientUrl := os.Getenv("CLIENT_URL")

	webSocketManager.GetInstance()

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
	wsRoutes.SetupRoutes(app)

	log.Fatal(app.Listen(":" + port))
}