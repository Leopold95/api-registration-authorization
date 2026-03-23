package main

import (
	"api-registration-authorization/services/registration/api"
	"api-registration-authorization/services/registration/application"
	"api-registration-authorization/services/registration/dataaccess"
	"api-registration-authorization/shared"
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file")
	}

	app := fiber.New()

	db, err := shared.DataBasePostgres()
	if err != nil {
		log.Fatal(err)
	}

	registrationRepository := dataaccess.NewRegistrationRepository(db)
	registrationService := application.NewRegistrationService(registrationRepository)

	api.NewRegisterController(app, registrationService)

	log.Fatal(app.Listen(":3001"))
}
