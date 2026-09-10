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
	var existing shared.UserEntity
	findErr := this.db.Where("id = ?", user.Id).First(&existing).Error
	if findErr == nil {
		if existing.Email == user.Email && existing.ProfileId == user.ProfileId {
			return nil
		}

		return domain.ErrorRegistrationConflict
	}

	if !errors.Is(findErr, gorm.ErrRecordNotFound) {
		return findErr
	}

	entity := &shared.UserEntity{
		Id:               user.Id,
		Email:            user.Email,
		PasswordHash:     user.Password,
		ProfileId:        user.ProfileId,
		RegistrationDate: time.Now(),
	}

	err := this.db.Create(entity).Error
	if err == nil {
		return nil
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		findErr = this.db.Where("id = ?", user.Id).First(&existing).Error
		if findErr == nil && existing.Email == user.Email && existing.ProfileId == user.ProfileId {
			return nil
		}

		return domain.ErrorUserExists
	}

	return err
}

func (self *RegistrationRepository) Delete(id uuid.UUID) error {
	err := self.db.Where("id = ?", id).Delete(&shared.UserEntity{}).Error
	return err
}
