package application

import (
	"api-registration-authorization/services/registration/dataaccess"
	"api-registration-authorization/shared"
)

type RegistrationService struct {
	repository *dataaccess.RegistrationRepository
	hasher     *PasswordHasherService
}

func NewRegistrationService(
	r *dataaccess.RegistrationRepository) *RegistrationService {
	return &RegistrationService{
		repository: r,
		hasher:     NewPasswordHasherService(),
	}
}

func (self *RegistrationService) Register(email, password, username string) (error, *shared.TokenModel) {
	err := self.repository.Insert(&shared.UserModel{
		Email:    email,
		Name:     username,
		Password: self.hasher.HashPassword(password),
	})

	if err != nil {
		return err, nil
	}

	return nil, &shared.TokenModel{
		Token:        "ok",
		RefreshToken: "ok",
	}
}
