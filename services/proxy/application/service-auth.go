package application

import (
	"api-registration-authorization/services/proxy/domain"
	"api-registration-authorization/shared"
)

type TestAuthServiceImpl struct {
	tokenService *shared.JwtService
}

func NewTestAuthServiceImpl(
	t *shared.JwtService,
) IAuthService {
	return &TestAuthServiceImpl{
		tokenService: t,
	}
}

func (t *TestAuthServiceImpl) ParseUserToken(token string) (*domain.UserModel, error) {

	claims, err := t.tokenService.ParseToken(token)

	if err != nil {
		return nil, err
	}

	model := &domain.UserModel{
		Id:        claims.Subject,
		Email:     claims.Email,
		Name:      claims.Name,
		ProfileId: claims.ProfileId,
	}

	return model, nil
}
