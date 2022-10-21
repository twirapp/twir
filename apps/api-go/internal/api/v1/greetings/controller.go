package greetings

import (
	"github.com/gofiber/fiber/v2"
	"github.com/satont/tsuwari/apps/api-go/internal/types"
)

func Setup(router fiber.Router, services types.Services) fiber.Router {
	middleware := router.Group("greetings")

	middleware.Get("", func(c *fiber.Ctx) error {
		c.JSON(HandleGet(c.Params("channelId"), services))

		return nil
	})

	return middleware
}
