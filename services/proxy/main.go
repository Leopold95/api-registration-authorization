package main

import (
	"api-registration-authorization/internal/logging"
	"api-registration-authorization/services/proxy/api"
	"api-registration-authorization/services/proxy/application"
	"api-registration-authorization/shared"
	"github.com/rs/zerolog/log"
	"os"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"

	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	defer logging.Init("proxy")()
	logging.InstallFiber(log.Logger)
	if err != nil {
		log.Info().Msg("No .env file found, using system environment variables")
	}

	app := fiber.New(fiber.Config{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	})

	app.Use(logging.HTTP())
	app.Use(recover.New())

	app.Use(cors.New())

	jwtService, err := shared.NewJwtService()
	if err != nil {
		log.Fatal().Err(err).Msg("Service initialization failed")
	}

	authService := application.NewTestAuthServiceImpl(jwtService)

	api.NewGatewayController(app, authService)

	app.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	if err := app.Listen("0.0.0.0"+os.Getenv("PROXY_PORT"), fiber.ListenConfig{DisableStartupMessage: true}); err != nil {
		log.Fatal().Err(err).Msg("HTTP server failed")
	}
}
