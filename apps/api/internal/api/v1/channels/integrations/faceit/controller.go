package faceit

import (
	"github.com/gofiber/fiber/v2"
	"github.com/satont/tsuwari/apps/api/internal/middlewares"
	"github.com/satont/tsuwari/apps/api/internal/types"
)

func Setup(router fiber.Router, services *types.Services) fiber.Router {
	middleware := router.Group("faceit")
	middleware.Get("", get(services))
	middleware.Get("auth", getAuthLink(services))
	middleware.Post("", post(services))
	middleware.Post("logout", logout(services))

	return middleware
}

// Integrations godoc
// @Security ApiKeyAuth
// @Summary      Get Faceit profile
// @Tags         Integrations|Faceit
// @Accept       json
// @Produce      json
// @Param        channelId   path      string  true  "ChannelId" default({{channelId}})
// @Success 200 {object} model.ChannelsIntegrationsData
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/integrations/faceit [get]
func get(services *types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		profile, err := handleGet(c.Params("channelId"), services)
		if err != nil {
			return err
		}

		return c.JSON(profile)
	}
}

// Integrations godoc
// @Security ApiKeyAuth
// @Summary      Get Faceit auth link
// @Tags         Integrations|Faceit
// @Accept       json
// @Produce      plain
// @Param        channelId   path      string  true  "ChannelId" default({{channelId}})
// @Success 200 {string} string	"Auth link"
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/integrations/faceit/auth [get]
func getAuthLink(services *types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		authLink, err := handleGetAuth(services)
		if err != nil {
			return err
		}

		return c.SendString(*authLink)
	}
}

// Integrations godoc
// @Security ApiKeyAuth
// @Summary      Authorize Faceit
// @Tags         Integrations|Faceit
// @Accept       json
// @Produce      json
// @Param data body tokenDto true "Data"
// @Param        channelId   path      string  true  "ID of channel"
// @Success      200
// @Failure 400 {object} types.DOCApiValidationError
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/integrations/faceit [post]
func post(services *types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		dto := &tokenDto{}
		err := middlewares.ValidateBody(
			c,
			services.Validator,
			services.ValidatorTranslator,
			dto,
		)
		if err != nil {
			return err
		}

		channelId := c.Params("channelId")
		err = handlePost(channelId, dto, services)
		if err != nil {
			return err
		}

		return c.SendStatus(200)
	}
}

// Integrations godoc
// @Security ApiKeyAuth
// @Summary      Logout
// @Tags         Integrations|Faceit
// @Accept       json
// @Produce      json
// @Param        channelId   path      string  true  "ID of channel"
// @Success      200
// @Failure 400 {object} types.DOCApiValidationError
// @Failure 404 {object} types.DOCApiBadRequest
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/integrations/faceit/logout [post]
func logout(services *types.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		channelId := c.Params("channelId")
		err := handleLogout(channelId, services)
		if err != nil {
			return err
		}

		return c.SendStatus(200)
	}
}
