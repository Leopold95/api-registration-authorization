package api

import (
	"api-registration-authorization/services/registration/application"
	"api-registration-authorization/services/registration/domain"
	"api-registration-authorization/shared"
	"api-registration-authorization/shared/consts"
	"errors"

	"github.com/gofiber/fiber/v2"
)

type RegisterController struct {
	service *application.RegistrationService
}

func NewRegisterController(app *fiber.App, s *application.RegistrationService) *RegisterController {
	controller := &RegisterController{
		service: s,
	}

	app.Post("api/v2/auth/register", RequestRegistration, controller.register)

	return controller
}

func (this *RegisterController) register(ctx *fiber.Ctx) error {
	request := ctx.Locals(consts.RequestKey).(RegistrationRequest)
	err, result := this.service.Register(request.Email, request.Password, request.Username)
	if err != nil {
		if errors.Is(err, domain.ErrorUserExists) {
			return shared.ResponseBadRequest(ctx, domain.ErrorUserExists.Error())
		}

		return shared.ResponseBadRequest(ctx, err.Error())
	}

	if result == nil {
		return shared.ResponseBadRequest(ctx, "user data is is null for some reasons")
	}

	return shared.ResponseOk(ctx, result)
}
