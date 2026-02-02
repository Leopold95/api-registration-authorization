package main

import (
	"api-registration-authorization/services/registration/api"
	"api-registration-authorization/services/registration/application"
	"api-registration-authorization/services/registration/dataaccess"
	"api-registration-authorization/shared"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
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
