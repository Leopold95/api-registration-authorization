package shared

import (
	"api-registration-authorization/internal/logging"
	"fmt"
	"github.com/rs/zerolog/log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func DataBasePostgres() (*gorm.DB, error) {
	dsn := os.Getenv("DB_CONNECTION_STRING")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logging.NewGORM(log.Logger)})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %w", err)
	}

	return db, nil
}
