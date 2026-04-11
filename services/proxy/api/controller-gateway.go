package api

import (
	"api-registration-authorization/services/proxy/application"
	"api-registration-authorization/services/proxy/domain"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/proxy"
)

var headerUserId = "X-User-Id"
var headerUserEmail = "X-User-Email"
var headerUserName = "X-User-Name"

type GatewayController struct {
	authService application.IAuthService
	gw          *application.GatewayService
	app         *fiber.App
}

func NewGatewayController(app *fiber.App, as application.IAuthService) *GatewayController {
	controller := &GatewayController{
		authService: as,
		app:         app,
		gw:          application.InitGateway(),
	}

	app.All("*", controller.gate)

	return controller
}

func (c *GatewayController) gate(ctx fiber.Ctx) error {
	path := ctx.Path()

	var endpoint *application.ServiceEndpoint
	for routePath := range c.gw.GetRoutes() {
		if strings.Contains(path, routePath) {
			ep, available := c.gw.GetEndpoint(routePath)
			if available {
				endpoint = ep
				break
			}
		}
	}

	if endpoint == nil {
		return ctx.Status(503).JSON(fiber.Map{
			"error":   "service unavailable",
			"message": "no available backend service for this route",
		})
	}

	if endpoint.IsPrivate {
		exists, value := c.isHeaderExists("Authorization", ctx)

		if !exists && strings.TrimSpace(value) == "" {
			return ctx.Status(403).JSON(fiber.Map{
				"error": "no authorization header found",
			})
		}

		//TODO сделать попытку авторизации нормальную авторизацию
		token, err := extractToken(value)
		user, err := c.authService.ParseUserToken(token)

		if err != nil {
			return ctx.Status(401).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		if user == nil {
			return ctx.Status(401).JSON(fiber.Map{
				"error": "error logging in. user is null for some reason",
			})
		}

		c.processHeaders(user, ctx)
	}

	// Build target URL
	targetURL := "" + endpoint.Host + ":" + endpoint.Port + path

	// Add query parameters
	if len(ctx.Request().URI().QueryString()) > 0 {
		targetURL += "?" + string(ctx.Request().URI().QueryString())
	}

	// Forward request
	return proxy.Do(ctx, targetURL)
}

func (c *GatewayController) processHeaders(user *domain.UserModel, ctx fiber.Ctx) {
	ctx.Request().Header.Del(headerUserId)
	ctx.Request().Header.Set(headerUserId, user.ProfileId)

	ctx.Request().Header.Del(headerUserEmail)
	ctx.Request().Header.Set(headerUserEmail, user.Email)

	ctx.Request().Header.Del(headerUserName)
	ctx.Request().Header.Set(headerUserName, user.Name)
}

func (c *GatewayController) isHeaderExists(header string, ctx fiber.Ctx) (bool, string) {
	value := ctx.Get(header)
	if value == "" {
		return false, ""
	}
	return true, value
}

func extractToken(authHeader string) (string, error) {
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", fmt.Errorf("invalid authorization header")
	}
	return parts[1], nil
}
