package timers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/satont/tsuwari/apps/api/internal/middlewares"
	"github.com/satont/tsuwari/apps/api/internal/types"
)

type Timers struct {
	services *types.Services
}

func NewController(router fiber.Router, services *types.Services) fiber.Router {
	timers := &Timers{
		services: services,
	}

	return router.Group("timers").
		Get("", timers.get).
		Post("", timers.post).
		Delete(":timerId", timers.delete).
		Put(":timerId", timers.put).
		Patch(":timerId", timers.patch)

}

// Timers godoc
// @Security ApiKeyAuth
// @Summary      Get channel timers list
// @Tags         Timers
// @Accept       json
// @Produce      json
// @Param        channelId   path      string  true  "ChannelId" default({{channelId}})
// @Success      200  {array}  model.ChannelsTimers
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/timers [get]
func (c *Timers) get(ctx *fiber.Ctx) error {
	timers := c.getService(ctx.Params("channelId"))

	return ctx.JSON(timers)
}

// Timers godoc
// @Security ApiKeyAuth
// @Summary      Create timer
// @Tags         Timers
// @Accept       json
// @Produce      json
// @Param data body timerDto true "Data"
// @Param        channelId   path      string  true  "ID of channel"
// @Success      200  {object}  model.ChannelsTimers
// @Failure 400 {object} types.DOCApiValidationError
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/timers [post]
func (c *Timers) post(ctx *fiber.Ctx) error {
	dto := &timerDto{}
	err := middlewares.ValidateBody(
		ctx,
		c.services.Validator,
		c.services.ValidatorTranslator,
		dto,
	)
	if err != nil {
		return err
	}

	cmd, err := c.postService(ctx.Params("channelId"), dto)
	if err == nil {
		return ctx.JSON(cmd)
	}

	return err
}

// Timers godoc
// @Security ApiKeyAuth
// @Summary      Delete timer
// @Tags         Timers
// @Accept       json
// @Produce      json
// @Param        channelId   path      string  true  "ID of channel"
// @Param        timerId   path      string  true  "ID of timer"
// @Success      200
// @Failure 400 {object} types.DOCApiValidationError
// @Failure 404
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/timers/{timerId} [delete]
func (c *Timers) delete(ctx *fiber.Ctx) error {
	err := c.deleteService(ctx.Params("timerId"))
	if err != nil {
		return err
	}
	return ctx.SendStatus(200)
}

// Timers godoc
// @Security ApiKeyAuth
// @Summary      Update timer
// @Tags         Timers
// @Accept       json
// @Produce      json
// @Param data body timerDto true "Data"
// @Param        channelId   path      string  true  "ID of channel"
// @Param        timerId   path      string  true  "ID of timer"
// @Success      200  {object}  model.ChannelsTimers
// @Failure 400 {object} types.DOCApiValidationError
// @Failure 404
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/timers/{timerId} [put]
func (c *Timers) put(ctx *fiber.Ctx) error {
	dto := &timerDto{}
	err := middlewares.ValidateBody(
		ctx,
		c.services.Validator,
		c.services.ValidatorTranslator,
		dto,
	)
	if err != nil {
		return err
	}

	cmd, err := c.putService(ctx.Params("timerId"), dto)
	if err == nil {
		return ctx.JSON(cmd)
	}

	return err
}

// Timers godoc
// @Security ApiKeyAuth
// @Summary      Partially update timer
// @Tags         Timers
// @Accept       json
// @Produce      json
// @Param data body timerPatchDto true "Data"
// @Param        channelId   path      string  true  "ID of channel"
// @Param        timerId   path      string  true  "ID of timer"
// @Success      200  {object}  model.ChannelsTimers
// @Failure 400 {object} types.DOCApiValidationError
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/timers/{timerId} [patch]
func (c *Timers) patch(ctx *fiber.Ctx) error {
	dto := &timerPatchDto{}
	err := middlewares.ValidateBody(
		ctx,
		c.services.Validator,
		c.services.ValidatorTranslator,
		dto,
	)
	if err != nil {
		return err
	}

	cmd, err := c.patchService(ctx.Params("timerId"), dto)
	if err == nil {
		return ctx.JSON(cmd)
	}

	return err

}
