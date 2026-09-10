package main

import (
	"api-registration-authorization/internal/logging"
	"api-registration-authorization/shared"
	"github.com/rs/zerolog/log"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	defer logging.Init("migration")()
	if err != nil {
		log.Info().Msg("Error loading .env file. Using system enviroment")
	}

	db, err := shared.DataBasePostgres()
	if err != nil {
		log.Fatal().Err(err).Msg("Service initialization failed")
	}

	Migrate(db)

}
