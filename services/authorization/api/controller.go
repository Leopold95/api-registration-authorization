package api

import (
	"api-registration-authorization/services/authorization/application"
	"api-registration-authorization/shared"
	"api-registration-authorization/shared/consts"

	"github.com/gofiber/fiber/v3"
)

type AuthController struct {
	service *application.AuthorizationService
}

func NewAuthController(app *fiber.App, s *application.AuthorizationService) *AuthController {
	controller := &AuthController{
		service: s,
	}

	app.Post("/api/v2/auth/authorize", RequestAuthorization, controller.auth)

	return controller
}

func (this *AuthController) auth(ctx fiber.Ctx) error {
	request := ctx.Locals(consts.RequestKey).(AuthorizationRequest)
	data, err := this.service.TryAuthenticateUser(request.Email, request.Password)
	if err != nil {
		return shared.ResponseBadRequest(ctx, "error")
	}
	return shared.ResponseOk(ctx, &AuthorizeResponse{Token: data.Token, RefreshToken: data.RefreshToken})
}
