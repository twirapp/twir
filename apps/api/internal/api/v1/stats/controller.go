package stats

import (
	"github.com/gofiber/fiber/v2"
	"github.com/satont/tsuwari/apps/api/internal/types"
)

func Setup(router fiber.Router, services types.Services) fiber.Router {
	middleware := router.Group("stats")
	middleware.Get("", get(services))

	return middleware
}

func get(services types.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		stats, err := handleGet(services)
		if err != nil {
			return err
		}
		return c.JSON(stats)
	}
}
