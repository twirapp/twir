package song_requests

import (
	"github.com/gofiber/fiber/v2"
	"github.com/satont/twir/apps/api/internal/types"
)

func Setup(router fiber.Router, services types.Services) fiber.Router {
	middleware := router.Group("song-requests")

	middleware.Get(":channelId", get(services))

	return router
}

func get(services types.Services) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		commands, err := handleGet(ctx.Params("channelId"), services)

		if err != nil {
			return err
		}

		return ctx.JSON(commands)
	}
}
