package events

import (
	"github.com/gofiber/fiber/v2"
	"github.com/satont/tsuwari/apps/api/internal/middlewares"
	"github.com/satont/tsuwari/apps/api/internal/types"
)

type Events struct {
	services *types.Services
}

func NewEvents(router fiber.Router, services *types.Services) fiber.Router {
	events := &Events{
		services: services,
	}

	return router.Group("events").
		Get("", events.get).
		Post("", events.post).
		Patch(":eventId", events.patch).
		Put(":eventId", events.put).
		Delete(":eventId", events.delete)
}

func (c *Events) get(ctx *fiber.Ctx) error {
	events := c.getService(ctx.Params("channelId"))

	return ctx.JSON(events)
}

func (c *Events) post(ctx *fiber.Ctx) error {
	dto := &eventDto{}
	err := middlewares.ValidateBody(
		ctx,
		c.services.Validator,
		c.services.ValidatorTranslator,
		dto,
	)
	if err != nil {
		return err
	}

	event, err := c.postService(ctx.Params("channelId"), dto)
	if err != nil {
		return err
	}

	return ctx.JSON(event)
}

func (c *Events) put(ctx *fiber.Ctx) error {
	dto := &eventDto{}
	err := middlewares.ValidateBody(
		ctx,
		c.services.Validator,
		c.services.ValidatorTranslator,
		dto,
	)
	if err != nil {
		return err
	}

	event, err := c.putService(ctx.Params("channelId"), ctx.Params("eventId"), dto)
	if err != nil {
		return err
	}

	return ctx.JSON(event)
}

func (c *Events) delete(ctx *fiber.Ctx) error {
	err := c.deleteService(ctx.Params("channelId"), ctx.Params("eventId"))
	if err != nil {
		return err
	}

	return ctx.SendStatus(200)

}

func (c *Events) patch(ctx *fiber.Ctx) error {
	dto := &eventPatchDto{}
	err := middlewares.ValidateBody(
		ctx,
		c.services.Validator,
		c.services.ValidatorTranslator,
		dto,
	)
	if err != nil {
		return err
	}
	greeting, err := c.patchService(ctx.Params("channelId"), ctx.Params("eventId"), dto)
	if err != nil {
		return err
	}

	return ctx.JSON(greeting)
}
