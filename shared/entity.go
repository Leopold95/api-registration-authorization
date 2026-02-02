package shared

import (
	"time"

	"github.com/google/uuid"
)

type UserEntity struct {
	Id               uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ProfileId        uuid.UUID `gorm:"type:uuid;not null"`
	Name             string    `gorm:"type:varchar(255);not null"`
	Email            string    `gorm:"type:varchar(255);not null;unique"`
	PasswordHash     string    `gorm:"type:text;not null"`
	RegistrationDate time.Time `gorm:"type:date;not null"`
}
