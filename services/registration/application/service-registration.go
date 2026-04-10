package application

import (
	"api-registration-authorization/services/registration/dataaccess"
	"api-registration-authorization/services/registration/temporal"
	"api-registration-authorization/shared"
)

type RegistrationService struct {
	jwtService *shared.JwtService
	repository *dataaccess.RegistrationRepository
	hasher     *PasswordHasherService
	executor   *temporal.TemporalExecutor
}

func NewRegistrationService(
	s *shared.JwtService,
	r *dataaccess.RegistrationRepository,
	e *temporal.TemporalExecutor,
) *RegistrationService {
	return &RegistrationService{
		repository: r,
		executor:   e,
		hasher:     NewPasswordHasherService(),
	}
}

func (self *RegistrationService) Register(email, password, username string) (error, *shared.RegistrationResponse) {
	_, err := self.executor.BeginUserRegistration(email, username, self.hasher.HashPassword(password))

	if err != nil {
		return err, nil
	}

	return nil, &shared.RegistrationResponse{
		Status: true,
	}
}
