package greetings

import (
	"github.com/gofiber/fiber/v2"
	"github.com/satont/tsuwari/apps/api/internal/middlewares"
	"github.com/satont/tsuwari/apps/api/internal/types"
)

type Greetings struct {
	services *types.Services
}

func NewGreetings(router fiber.Router, services *types.Services) fiber.Router {
	greetings := &Greetings{
		services: services,
	}

	return router.Group("greetings").
		Get("", greetings.get).
		Post("", greetings.post).
		Delete(":greetingId", greetings.delete).
		Put(":greetingId", greetings.put).
		Patch(":greetingId", greetings.patch)
}

// Greetings godoc
// @Security ApiKeyAuth
// @Summary      Get channel greeting list
// @Tags         Greetings
// @Accept       json
// @Produce      json
// @Param        channelId   path      string  true  "ChannelId" default({{channelId}})
// @Success      200  {array}  Greeting
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/greetings [get]
func (c *Greetings) get(ctx *fiber.Ctx) error {
	return ctx.JSON(c.getService(ctx.Params("channelId")))
}

// greetings godoc
// @Security ApiKeyAuth
// @Summary      Create greeting
// @Tags         Greetings
// @Accept       json
// @Produce      json
// @Param data body greetingsDto true "Data"
// @Param        channelId   path      string  true  "ID of channel"
// @Success      200  {object}  Greeting
// @Failure 400 {object} types.DOCApiValidationError
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/greetings [post]
func (c *Greetings) post(ctx *fiber.Ctx) error {
	dto := &greetingsDto{}
	err := middlewares.ValidateBody(
		ctx,
		c.services.Validator,
		c.services.ValidatorTranslator,
		dto,
	)
	if err != nil {
		return err
	}

	greeting, err := c.postService(ctx.Params("channelId"), dto)
	if err == nil && greeting != nil {
		return ctx.JSON(greeting)
	}

	return err
}

// greetings godoc
// @Security ApiKeyAuth
// @Summary      Delete greeting
// @Tags         Greetings
// @Accept       json
// @Produce      json
// @Param        channelId   path      string  true  "ID of channel"
// @Param        greetingId   path      string  true  "ID of greeting"
// @Success      200
// @Failure 400 {object} types.DOCApiValidationError
// @Failure 404
// @Router       /v1/channels/{channelId}/greetings/{greetingId} [delete]
func (c *Greetings) delete(ctx *fiber.Ctx) error {
	err := c.deleteService(ctx.Params("greetingId"))
	if err != nil {
		return err
	}
	return ctx.SendStatus(fiber.StatusOK)
}

// Greetings godoc
// @Security ApiKeyAuth
// @Summary      Update greeting
// @Tags         Greetings
// @Accept       json
// @Produce      json
// @Param data body greetingsDto true "Data"
// @Param        channelId   path      string  true  "ID of channel"
// @Param        channelId   path      string  true  "ID of greeting"
// @Success      200  {object}  Greeting
// @Failure 400 {object} types.DOCApiValidationError
// @Failure 404
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/greetings/{greetingId} [put]
func (c *Greetings) put(ctx *fiber.Ctx) error {
	dto := &greetingsDto{}
	err := middlewares.ValidateBody(
		ctx,
		c.services.Validator,
		c.services.ValidatorTranslator,
		dto,
	)
	if err != nil {
		return err
	}
	greeting, err := c.putService(ctx.Params("greetingId"), dto)
	if err != nil {
		return err
	}

	return ctx.JSON(greeting)
}

// Greetings godoc
// @Security ApiKeyAuth
// @Summary      Update greeting
// @Tags         Greetings
// @Accept       json
// @Produce      json
// @Param data body greetingsPatchDto true "Data"
// @Param        channelId   path      string  true  "ID of channel"
// @Param        channelId   path      string  true  "ID of greeting"
// @Success      200  {object}  Greeting
// @Failure 400 {object} types.DOCApiValidationError
// @Failure 404
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/greetings/{greetingId} [put]
func (c *Greetings) patch(ctx *fiber.Ctx) error {
	dto := &greetingsPatchDto{}
	err := middlewares.ValidateBody(
		ctx,
		c.services.Validator,
		c.services.ValidatorTranslator,
		dto,
	)
	if err != nil {
		return err
	}
	greeting, err := c.patchService(ctx.Params("channelId"), ctx.Params("greetingId"), dto)
	if err != nil {
		return err
	}

	return ctx.JSON(greeting)
}
