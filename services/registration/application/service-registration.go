package application

import (
	"api-registration-authorization/services/registration/dataaccess"
	"api-registration-authorization/services/registration/temporal"
	"api-registration-authorization/shared"
	"errors"
	"strings"
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

func (self *RegistrationService) Register(email, password string) (error, *shared.RegistrationResponse) {
	normalizedEmail := strings.ToLower(strings.TrimSpace(email))

	hashedPassword, err := self.hasher.HashPassword(password)
	if err != nil {
		return err, nil
	}

	workflowResult, err := self.executor.BeginUserRegistration(normalizedEmail, hashedPassword)

	if err != nil {
		return err, nil
	}

	if workflowResult == nil || !workflowResult.GetSuccess() {
		return errors.New("registration workflow returned success=false"), nil
	}

	return nil, &shared.RegistrationResponse{
		Status: true,
	}
}
