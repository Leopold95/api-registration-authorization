package main

import (
	"api-registration-authorization/internal/logging"
	"api-registration-authorization/services/authorization/api"
	"api-registration-authorization/services/authorization/application"
	"api-registration-authorization/services/authorization/dataaccess"
	"api-registration-authorization/shared"
	"github.com/rs/zerolog/log"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	defer logging.Init("authorization")()
	logging.InstallFiber(log.Logger)
	if err != nil {
		log.Info().Msg("Error loading .env file. Using system enviroment")
	}

	app := fiber.New()
	app.Use(logging.HTTP())

	db, err := shared.DataBasePostgres()
	if err != nil {
		log.Fatal().Err(err).Msg("Service initialization failed")
	}

	authRepository := dataaccess.NewAuthorizationRepository(db)
	jwtService, err := shared.NewJwtService()
	if err != nil {
		log.Fatal().Err(err).Msg("Service initialization failed")
	}
	authService := application.NewAuthorizationService(authRepository, jwtService)

	api.NewAuthController(app, authService)

	if err := app.Listen("0.0.0.0"+os.Getenv("AUTHORIZATION_PORT"), fiber.ListenConfig{DisableStartupMessage: true}); err != nil {
		log.Fatal().Err(err).Msg("HTTP server failed")
	}
}
