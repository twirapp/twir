package donationalerts

import (
	"github.com/gofiber/fiber/v2"
	"github.com/satont/tsuwari/apps/api/internal/middlewares"
	"github.com/satont/tsuwari/apps/api/internal/types"
)

func Setup(router fiber.Router, services types.Services) fiber.Router {
	middleware := router.Group("donationalerts")
	middleware.Get("auth", getAuth(services))
	middleware.Get("", get(services))
	middleware.Post("logout", logout(services))
	middleware.Post("", post((services)))

	return middleware
}

// Integrations godoc
// @Security ApiKeyAuth
// @Summary      Get DonationAlerts integration
// @Tags         Integrations|DonationAlerts
// @Accept       json
// @Produce      json
// @Param        channelId   path      string  true  "ChannelId"
// @Success      200  {object}  model.ChannelsIntegrationsData
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/integrations/donationalerts [get]
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
// @Summary      Get DonationAlerts auth link
// @Tags         Integrations|DonationAlerts
// @Accept       json
// @Produce      plain
// @Param        channelId   path      string  true  "ChannelId"
// @Success 200 {string} string	"Auth link"
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/integrations/donationalerts/auth [get]
func getAuth(services types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		authLink, err := handleGetAuth(services)
		if err != nil {
			return err
		}

		return c.SendString(*authLink)
	}
}

type tokenDto struct {
	Code string `validate:"required" json:"code"`
}

// Integrations godoc
// @Security ApiKeyAuth
// @Summary      Update auth of DonationAlerts
// @Tags         Integrations|DonationAlerts
// @Accept       json
// @Produce      json
// @Param data body tokenDto true "Data"
// @Param        channelId   path      string  true  "ID of channel"
// @Success      200
// @Failure 400 {object} types.DOCApiValidationError
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/integrations/donationalerts [post]
func post(services types.Services) func(c *fiber.Ctx) error {
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

		err = handlePost(c.Params("channelId"), dto, services)
		if err != nil {
			return err
		}

		return c.SendStatus(200)
	}
}

// Integrations godoc
// @Security ApiKeyAuth
// @Summary      Logout
// @Tags         Integrations|DonationAlerts
// @Accept       json
// @Produce      json
// @Param        channelId   path      string  true  "ID of channel"
// @Success      200
// @Failure 400 {object} types.DOCApiValidationError
// @Failure 404 {object} types.DOCApiBadRequest
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/integrations/donationalerts/logout [post]
func logout(services types.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		channelId := c.Params("channelId")
		err := handleLogout(channelId, services)
		if err != nil {
			return err
		}

		return c.SendStatus(200)
	}
}
