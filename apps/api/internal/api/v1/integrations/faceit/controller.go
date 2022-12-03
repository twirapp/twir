package faceit

import (
	"github.com/gofiber/fiber/v2"
	"github.com/satont/tsuwari/apps/api/internal/middlewares"
	"github.com/satont/tsuwari/apps/api/internal/types"
)

func Setup(router fiber.Router, services types.Services) fiber.Router {
	middleware := router.Group("faceit")
	middleware.Get("", get(services))
	middleware.Post("", post(services))

	return middleware
}

// Integrations godoc
// @Security ApiKeyAuth
// @Summary      Get Faceit integration
// @Tags         Integrations
// @Accept       json
// @Produce      json
// @Param        channelId   path      string  true  "ChannelId"
// @Success      200  {object}  model.ChannelsIntegrations
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/integrations/faceit [get]
func get(services types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		integration, err := handleGet(c.Params("channelId"), services)
		if err != nil {
			return err
		}
		return c.JSON(integration)
	}
}

// Integrations godoc
// @Security ApiKeyAuth
// @Summary      Update Faceit
// @Tags         Integrations
// @Accept       json
// @Produce      json
// @Param data body faceitUpdateDto true "Data"
// @Param        channelId   path      string  true  "ID of channel"
// @Success      200  {object} model.ChannelsIntegrations
// @Failure 400 {object} types.DOCApiValidationError
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/integrations/faceit [post]
func post(services types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		dto := &faceitUpdateDto{}
		err := middlewares.ValidateBody(
			c,
			services.Validator,
			services.ValidatorTranslator,
			dto,
		)
		if err != nil {
			return err
		}

		err = handlePost(c.Params("channelId"), dto, services)

		if err != nil {
			return err
		}

		return c.SendStatus(200)
	}
}
