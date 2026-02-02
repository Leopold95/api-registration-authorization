package dataaccess

import (
	"api-registration-authorization/shared"

	"gorm.io/gorm"
)

type AuthorizationRepository struct {
	db *gorm.DB
}

func NewAuthorizationRepository(db *gorm.DB) *AuthorizationRepository {
	return &AuthorizationRepository{
		db: db,
	}
}

func (this *AuthorizationRepository) SelectUser(email string) *shared.UserModel {
	entity := &shared.UserEntity{}
	this.db.Where(&shared.UserEntity{Email: email}).First(entity)
	return &shared.UserModel{
		Id:    entity.Id,
		Email: entity.Email,
		Name:  entity.Name,
	}
}
