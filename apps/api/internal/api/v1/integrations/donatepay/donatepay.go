package donatepay

import (
	"github.com/gofiber/fiber/v2"
	"github.com/satont/tsuwari/apps/api/internal/types"
)

func Setup(router fiber.Router, services types.Services) fiber.Router {
	middleware := router.Group("donatepay")
	middleware.Get("", get(services))
	middleware.Post("", post(services))

	return middleware
}

func get(services types.Services) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		data, err := handleGet(services, ctx.Params("channelId"))

		if err != nil {
			return err
		}

		return ctx.JSON(data)
	}
}

func post(services types.Services) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		return nil
	}
}
