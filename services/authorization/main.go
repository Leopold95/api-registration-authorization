package main

import (
	"api-registration-authorization/services/authorization/api"
	"api-registration-authorization/services/authorization/application"
	"api-registration-authorization/services/authorization/dataaccess"
	"api-registration-authorization/shared"
	"log"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file. Using system enviroment")
	}

	app := fiber.New()

	db, err := shared.DataBasePostgres()
	if err != nil {
		log.Fatal(err)
	}

	authRepository := dataaccess.NewAuthorizationRepository(db)
	jwtService, err := shared.NewJwtService()
	if err != nil {
		log.Fatal(err)
	}
	authService := application.NewAuthorizationService(authRepository, jwtService)

	api.NewAuthController(app, authService)

	log.Fatal(app.Listen("0.0.0.0" + os.Getenv("AUTHORIZATION_PORT")))
}
