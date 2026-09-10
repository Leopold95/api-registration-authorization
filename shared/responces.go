package shared

import "github.com/gofiber/fiber/v3"

type GlobalResponse struct {
	Success    bool    `json:"success"`
	SystemCode int     `json:"system_code"`
	HttpCode   int     `json:"http_code"`
	Message    *string `json:"message"`
	Data       any     `json:"data"`
}

func ResponseOk(ctx fiber.Ctx, data any) error {
	return ctx.Status(fiber.StatusOK).JSON(
		&GlobalResponse{
			Success: true,
			Message: nil,
			Data:    data,
		})
}

func ResponseBadRequest(ctx fiber.Ctx, message string) error {
	return ctx.Status(fiber.StatusBadRequest).JSON(GlobalResponse{
		Success: false,
		Message: &message,
		Data:    nil,
	})
}
