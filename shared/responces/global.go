package responces

import "github.com/gofiber/fiber/v3"

func JsonParseError(c fiber.Ctx) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"error": "cant parse request json",
	})
}
