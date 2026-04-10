package main

import (
	"api-registration-authorization/services/proxy/api"
	"api-registration-authorization/services/proxy/application"
	"api-registration-authorization/shared"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	app := fiber.New(fiber.Config{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	})

	app.Use(recover.New())
	app.Use(logger.New())

	app.Use(cors.New())

	jwtService, err := shared.NewJwtService()
	if err != nil {
		log.Fatal(err)
	}

	authService := application.NewTestAuthServiceImpl(jwtService)

	api.NewGatewayController(app, authService)

	app.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	log.Fatal(app.Listen("0.0.0.0" + os.Getenv("PROXY_PORT")))
}
