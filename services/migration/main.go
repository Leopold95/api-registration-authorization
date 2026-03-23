package main

import (
	"api-registration-authorization/services/registration/dataaccess"
	"api-registration-authorization/shared"
	"log"

	"github.com/gofiber/fiber/v3"
)

func main() {
	app := fiber.New()

	db, err := shared.DataBasePostgres()
	if err != nil {
		log.Fatal(err)
	}

	app.Post("api/auth/migrate", func(c fiber.Ctx) error {
		log.Println("migrating data")
		dataaccess.Migrate(db)
		log.Println("data migrated")
		return c.SendStatus(fiber.StatusOK)
	})

	log.Fatal(app.Listen(":3000"))
}
