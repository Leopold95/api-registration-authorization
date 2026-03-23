package api

import (
	"api-registration-authorization/shared"
	"api-registration-authorization/shared/consts"
	"api-registration-authorization/shared/responces"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

type AuthorizationRequest struct {
	Password string `json:"password" validate:"required"`
	Email    string `json:"email" validate:"required"`
}

func RequestAuthorization(c fiber.Ctx) error {
	var update AuthorizationRequest

	if err := c.Bind().Body(&update); err != nil {
		return responces.JsonParseError(c)
	}

	if err := shared.Validate.Struct(update); err != nil {
		errors := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			errors[err.Field()] = err.Tag()
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"validation_errors": errors,
		})
	}

	c.Locals(consts.RequestKey, update)

	return c.Next()
}
