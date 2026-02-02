package shared

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func DataBasePostgres() (*gorm.DB, error) {
	dsn := "host=localhost user=user password=password dbname=auth.db port=6121 sslmode=disable TimeZone=Europe/Kiev"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %w", err)
	}

	return db, nil
}
