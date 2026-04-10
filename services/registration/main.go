package main

import (
	"api-registration-authorization/services/registration/api"
	"api-registration-authorization/services/registration/application"
	"api-registration-authorization/services/registration/dataaccess"
	"api-registration-authorization/services/registration/temporal"
	"api-registration-authorization/shared"
	"log"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/converter"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file. Using system enviroment")
	}

	db, err := shared.DataBasePostgres()
	if err != nil {
		log.Fatal(err)
	}

	c, err := client.Dial(client.Options{
		HostPort:      os.Getenv("TEMPORAL_ADDRESS"),
		DataConverter: converter.NewCompositeDataConverter(converter.NewProtoPayloadConverter()),
	})
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	registrationRepository := dataaccess.NewRegistrationRepository(db)
	temporalExecutor := temporal.NewTemporalExecutor(c)
	jwtService, err := shared.NewJwtService()
	if err != nil {
		log.Fatal(err)
	}
	registrationService := application.NewRegistrationService(jwtService, registrationRepository, temporalExecutor)
	_ = temporal.NewTemporalActivities(c, registrationRepository)

	app := fiber.New()

	api.NewRegisterController(app, registrationService)

	log.Fatal(app.Listen("0.0.0.0" + os.Getenv("REGISTRATION_PORT")))
}
