package moderation

import (
	"github.com/gofiber/fiber/v2"
	"github.com/satont/tsuwari/apps/api/internal/middlewares"
	"github.com/satont/tsuwari/apps/api/internal/types"
)

type Moderation struct {
	services *types.Services
}

func NewController(router fiber.Router, services *types.Services) fiber.Router {
	moderation := &Moderation{
		services: services,
	}

	return router.Group("moderation").
		Get("", moderation.get).
		Post("", moderation.post)
}

// Moderation godoc
// @Security ApiKeyAuth
// @Summary      Get moderation settings
// @Tags         Moderation
// @Accept       json
// @Produce      json
// @Param        channelId   path      string  true  "ChannelId" default({{channelId}})
// @Success      200  {array}  model.ChannelsModerationSettings
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/moderation [get]
func (c *Moderation) get(ctx *fiber.Ctx) error {
	settings, err := c.getService(ctx.Params("channelId"))
	if err != nil {
		return err
	}

	return ctx.JSON(settings)
}

// Moderation godoc
// @Security ApiKeyAuth
// @Summary      Create command
// @Tags         Moderation
// @Accept       json
// @Produce      json
// @Param data body moderationDto true "Data"
// @Param        channelId   path      string  true  "ID of channel"
// @Success      200  {array}  model.ChannelsModerationSettings
// @Failure 400 {object} types.DOCApiValidationError
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/moderation [post]
func (c *Moderation) post(ctx *fiber.Ctx) error {
	dto := &moderationDto{}
	err := middlewares.ValidateBody(
		ctx,
		c.services.Validator,
		c.services.ValidatorTranslator,
		dto,
	)
	if err != nil {
		return err
	}
	settings, err := c.postService(ctx.Params("channelId"), dto)
	if err != nil {
		return err
	}

	return ctx.JSON(settings)
}
