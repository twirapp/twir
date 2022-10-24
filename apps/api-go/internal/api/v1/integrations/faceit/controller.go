package faceit

import (
	"github.com/gofiber/fiber/v2"
	"github.com/satont/tsuwari/apps/api-go/internal/types"
)

func Setup(router fiber.Router, services types.Services) fiber.Router {
	middleware := router.Group("faceit")
	middleware.Get("", get(services))

	return middleware
}

func get(services types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		integration, err := handleGet(c.Params("channelId"), services)
		if err != nil {
			return err
		}
		return c.JSON(integration)
	}
}
