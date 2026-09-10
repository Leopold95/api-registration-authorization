package main

import (
	"api-registration-authorization/internal/logging"
	"api-registration-authorization/services/registration/api"
	"api-registration-authorization/services/registration/application"
	"api-registration-authorization/services/registration/dataaccess"
	"api-registration-authorization/services/registration/temporal"
	"api-registration-authorization/shared"
	"github.com/rs/zerolog/log"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/converter"
)

func main() {
	err := godotenv.Load()
	defer logging.Init("registration")()
	logging.InstallFiber(log.Logger)
	if err != nil {
		log.Info().Msg("Error loading .env file. Using system enviroment")
	}

	db, err := shared.DataBasePostgres()
	if err != nil {
		log.Fatal().Err(err).Msg("Service initialization failed")
	}

	c, err := client.Dial(client.Options{
		Logger:        logging.NewTemporal(log.Logger),
		HostPort:      os.Getenv("TEMPORAL_ADDRESS"),
		DataConverter: converter.NewCompositeDataConverter(converter.NewProtoPayloadConverter()),
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Service initialization failed")
	}
	defer c.Close()

	registrationRepository := dataaccess.NewRegistrationRepository(db)
	temporalExecutor := temporal.NewTemporalExecutor(c)
	jwtService, err := shared.NewJwtService()
	if err != nil {
		log.Fatal().Err(err).Msg("Service initialization failed")
	}
	registrationService := application.NewRegistrationService(jwtService, registrationRepository, temporalExecutor)
	_ = temporal.NewTemporalActivities(c, registrationRepository)

	app := fiber.New()
	app.Use(logging.HTTP())

	api.NewRegisterController(app, registrationService)

	if err := app.Listen("0.0.0.0"+os.Getenv("REGISTRATION_PORT"), fiber.ListenConfig{DisableStartupMessage: true}); err != nil {
		log.Fatal().Err(err).Msg("HTTP server failed")
	}
}
