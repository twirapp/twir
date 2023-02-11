package events

import (
	"github.com/gofiber/fiber/v2"
	"github.com/satont/tsuwari/apps/api/internal/middlewares"
	"github.com/satont/tsuwari/apps/api/internal/types"
)

func Setup(router fiber.Router, services types.Services) fiber.Router {
	middleware := router.Group("events")
	middleware.Get("", get(services))
	middleware.Post("", create(services))
	middleware.Put(":eventId", update(services))
	middleware.Delete(":eventId", delete(services))

	return middleware
}

func get(services types.Services) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		events := handleGet(ctx.Params("channelId"), services)

		return ctx.JSON(events)
	}
}

func create(services types.Services) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		dto := &eventDto{}
		err := middlewares.ValidateBody(
			ctx,
			services.Validator,
			services.ValidatorTranslator,
			dto,
		)
		if err != nil {
			return err
		}

		event, err := handlePost(ctx.Params("channelId"), dto)
		if err != nil {
			return err
		}

		return ctx.JSON(event)
	}
}

func update(services types.Services) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		dto := &eventDto{}
		err := middlewares.ValidateBody(
			ctx,
			services.Validator,
			services.ValidatorTranslator,
			dto,
		)
		if err != nil {
			return err
		}

		event, err := handleUpdate(ctx.Params("channelId"), ctx.Params("eventId"), dto)
		if err != nil {
			return err
		}

		return ctx.JSON(event)
	}
}

func delete(services types.Services) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		err := handleDelete(ctx.Params("channelId"), ctx.Params("eventId"))
		if err != nil {
			return err
		}

		return ctx.SendStatus(200)
	}
}
