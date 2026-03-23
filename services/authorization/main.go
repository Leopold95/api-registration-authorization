package main

import (
	"api-registration-authorization/services/authorization/api"
	"api-registration-authorization/services/authorization/application"
	"api-registration-authorization/services/authorization/dataaccess"
	"api-registration-authorization/shared"
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	app := fiber.New()

	db, err := shared.DataBasePostgres()
	if err != nil {
		log.Fatal(err)
	}

	authRepository := dataaccess.NewAuthorizationRepository(db)
	authService := application.NewAuthorizationService(authRepository, shared.NewJwtService())

	api.NewAuthController(app, authService)

	log.Fatal(app.Listen(":3000"))
}
