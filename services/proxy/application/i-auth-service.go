package application

import "api-registration-authorization/services/proxy/domain"

type IAuthService interface {
	ParseUserToken(token string) (*domain.UserModel, error)
}
