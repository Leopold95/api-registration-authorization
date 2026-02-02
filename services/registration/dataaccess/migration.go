package dataaccess

import (
	"api-registration-authorization/shared"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&shared.UserEntity{})
}
