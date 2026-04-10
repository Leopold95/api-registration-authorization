package dataaccess

import (
	"api-registration-authorization/services/registration/domain"
	"api-registration-authorization/shared"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type RegistrationRepository struct {
	db  *gorm.DB
	ctx context.Context
}

func NewRegistrationRepository(db *gorm.DB) *RegistrationRepository {
	return &RegistrationRepository{
		db:  db,
		ctx: context.Background(),
	}
}

func (this *RegistrationRepository) Insert(user *shared.UserModel) error {
	entity := &shared.UserEntity{
		Id:               user.Id,
		Email:            user.Email,
		PasswordHash:     user.Password,
		ProfileId:        user.ProfileId,
		RegistrationDate: time.Now(),
	}

	err := this.db.Create(entity).Error

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return domain.ErrorUserExists
	}

	return err
}

func (self *RegistrationRepository) Delete(id uuid.UUID) error {
	err := self.db.Where("id = ?", id).Delete(&shared.UserEntity{}).Error
	return err
}
