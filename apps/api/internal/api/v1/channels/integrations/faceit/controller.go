package faceit

import (
	"github.com/gofiber/fiber/v2"
	"github.com/satont/tsuwari/apps/api/internal/middlewares"
	"github.com/satont/tsuwari/apps/api/internal/types"
)

type Faceit struct {
	services *types.Services
}

func NewController(router fiber.Router, services *types.Services) fiber.Router {
	faceit := &Faceit{
		services,
	}

	return router.Group("faceit").
		Get("", faceit.get).
		Get("auth", faceit.getAuthLink).
		Post("", faceit.post).
		Post("logout", faceit.logout)
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
func (c *Faceit) get(ctx *fiber.Ctx) error {
	profile, err := c.getService(ctx.Params("channelId"))
	if err != nil {
		return err
	}

	return ctx.JSON(profile)
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
func (c *Faceit) getAuthLink(ctx *fiber.Ctx) error {
	authLink, err := c.getAuthLinkService()
	if err != nil {
		return err
	}

	return ctx.SendString(*authLink)
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
func (c *Faceit) post(ctx *fiber.Ctx) error {
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

	channelId := ctx.Params("channelId")
	err = c.postService(channelId, dto)
	if err != nil {
		return err
	}

	return ctx.SendStatus(200)
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
func (c *Faceit) logout(ctx *fiber.Ctx) error {
	channelId := ctx.Params("channelId")
	err := c.logoutService(channelId)
	if err != nil {
		return err
	}

	return ctx.SendStatus(200)
}
