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

func (this *AuthorizationRepository) SelectUser(email string) (*shared.UserModel, error) {
	entity := &shared.UserEntity{}
	err := this.db.Where(&shared.UserEntity{Email: email}).First(entity).Error
	return &shared.UserModel{
		Id:        entity.Id,
		Email:     entity.Email,
		Password:  entity.PasswordHash,
		ProfileId: entity.ProfileId,
	}, err
}
