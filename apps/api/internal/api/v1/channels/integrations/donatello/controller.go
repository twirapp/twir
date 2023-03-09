package donatello

import (
	"github.com/gofiber/fiber/v2"
	"github.com/satont/tsuwari/apps/api/internal/middlewares"
	"github.com/satont/tsuwari/apps/api/internal/types"
)

type Donatello struct {
	services *types.Services
}

func NewController(router fiber.Router, services *types.Services) fiber.Router {
	donatello := &Donatello{services: services}

	return router.Group("donatello").
		Get("", donatello.get).
		Post("", donatello.post)
}

// Integrations godoc
// @Security ApiKeyAuth
// @Summary      Get Donatello token
// @Tags         Integrations|Donatello
// @Accept       json
// @Produce      plain
// @Param        channelId   path      string  true  "ChannelId" default({{channelId}})
// @Success      200  {string} string
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/integrations/donatepay [get]
func (c *Donatello) get(ctx *fiber.Ctx) error {
	data, err := c.getService(ctx.Params("channelId"))
	if err != nil {
		return err
	}

	return ctx.JSON(data)
}

// Integrations godoc
// @Security ApiKeyAuth
// @Summary      Authorize Donatello
// @Tags         Integrations|Donatello
// @Accept       json
// @Param data body createOrUpdateDTO true "Data"
// @Param        channelId   path      string  true  "ID of channel"
// @Success      200
// @Failure 400 {object} types.DOCApiValidationError
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/integrations/donatepay [post]
func (c *Donatello) post(ctx *fiber.Ctx) error {
	dto := &createOrUpdateDTO{}
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
