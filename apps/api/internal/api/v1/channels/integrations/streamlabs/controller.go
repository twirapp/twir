package streamlabs

import (
	"github.com/gofiber/fiber/v2"
	"github.com/satont/tsuwari/apps/api/internal/middlewares"
	"github.com/satont/tsuwari/apps/api/internal/types"
)

type StreamLabs struct {
	services *types.Services
}

func NewController(router fiber.Router, services *types.Services) fiber.Router {
	streamLabs := &StreamLabs{
		services,
	}

	return router.Group("streamlabs").
		Get("", streamLabs.get).
		Get("auth", streamLabs.getAuthLink).
		Post("", streamLabs.post).
		Post("logout", streamLabs.logout)
}

// Integrations godoc
// @Security ApiKeyAuth
// @Summary      Get Streamlabs integration
// @Tags         Integrations|Streamlabs
// @Accept       json
// @Produce      json
// @Param        channelId   path      string  true  "ChannelId" default({{channelId}})
// @Success      200  {object}  model.ChannelsIntegrations
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/integrations/streamlabs [get]
func (c *StreamLabs) get(ctx *fiber.Ctx) error {
	integration, err := c.getService(ctx.Params("channelId"))
	if err != nil {
		return err
	}
	return ctx.JSON(integration)
}

// Integrations godoc
// @Security ApiKeyAuth
// @Summary      Get DonationAlerts auth link
// @Tags         Integrations|Streamlabs
// @Accept       json
// @Produce      plain
// @Param        channelId   path      string  true  "ChannelId" default({{channelId}})
// @Success 200 {string} string	"Auth link"
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/integrations/streamlabs/auth [get]
func (c *StreamLabs) getAuthLink(ctx *fiber.Ctx) error {
	authLink, err := c.getAuthLinkService()
	if err != nil {
		return err
	}

	return ctx.SendString(*authLink)
}

type tokenDto struct {
	Code string `validate:"required" json:"code"`
}

// Integrations godoc
// @Security ApiKeyAuth
// @Summary      Update auth of Streamlabs
// @Tags         Integrations|Streamlabs
// @Accept       json
// @Produce      json
// @Param data body tokenDto true "Data"
// @Param        channelId   path      string  true  "ID of channel"
// @Success      200
// @Failure 400 {object} types.DOCApiValidationError
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/integrations/streamlabs [post]
func (c *StreamLabs) post(ctx *fiber.Ctx) error {
	dto := &tokenDto{}
	err := middlewares.ValidateBody(
		ctx,
		c.services.Validator,
		c.services.ValidatorTranslator,
		dto,
	)
	if err != nil {
		return err
	}

	err = c.postService(ctx.Params("channelId"), dto)
	if err != nil {
		return err
	}

	return ctx.SendStatus(200)
}

// Integrations godoc
// @Security ApiKeyAuth
// @Summary      Logout
// @Tags         Integrations|Streamlabs
// @Accept       json
// @Produce      json
// @Param        channelId   path      string  true  "ID of channel"
// @Success      200
// @Failure 400 {object} types.DOCApiValidationError
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/integrations/streamlabs/logout [post]
func (c *StreamLabs) logout(ctx *fiber.Ctx) error {
	err := c.logoutService(ctx.Params("channelId"))
	if err != nil {
		return err
	}

	return ctx.SendStatus(200)
}
