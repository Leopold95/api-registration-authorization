package main

import (
	"api-registration-authorization/shared"
	"log"

	"github.com/joho/godotenv"
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

	Migrate(db)

	//app.Post("api/auth/migrate", func(c fiber.Ctx) error {
	//	log.Println("migrating data")
	//	Migrate(db)
	//	log.Println("data migrated")
	//	return c.SendStatus(fiber.StatusOK)
	//})
	//
	//log.Fatal(app.Listen(":3000"))
}
