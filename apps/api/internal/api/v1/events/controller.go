package events

import (
	"github.com/gofiber/fiber/v2"
	"github.com/satont/tsuwari/apps/api/internal/middlewares"
	"github.com/satont/tsuwari/apps/api/internal/types"
	"net/http"
)

func Setup(router fiber.Router, services types.Services) fiber.Router {
	middleware := router.Group("events")
	middleware.Get("", get(services))
	middleware.Post("", create(services))
	middleware.Put("", update(services))

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

		if err = handlePost(ctx.Params("channelId"), dto); err != nil {
			return err
		}

		return ctx.SendStatus(http.StatusOK)
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

		if err = handleUpdate(ctx.Params("channelId"), ctx.Params("eventId"), dto); err != nil {
			return err
		}

		return ctx.SendStatus(http.StatusOK)
	}
}
