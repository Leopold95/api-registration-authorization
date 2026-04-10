package application

import (
	"api-registration-authorization/services/authorization/dataaccess"
	"api-registration-authorization/shared"
)

type AuthorizationService struct {
	repository   *dataaccess.AuthorizationRepository
	tokenService *shared.JwtService
}

func NewAuthorizationService(repository *dataaccess.AuthorizationRepository, tokenService *shared.JwtService) *AuthorizationService {
	return &AuthorizationService{
		repository:   repository,
		tokenService: tokenService,
	}
}

func (this *AuthorizationService) TryAuthenticateUser(email, password string) (*shared.TokenModel, error) {
	userModel, err := this.repository.SelectUser(email)

	if err != nil {
		return nil, err
	}

	token, err := this.tokenService.CreateToken(userModel)
	if err != nil {
		return nil, err
	}
	return &shared.TokenModel{
		Token:        token,
		RefreshToken: "",
	}, nil
}
