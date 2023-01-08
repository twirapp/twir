package users

import (
	"github.com/gofiber/fiber/v2"
	"github.com/satont/tsuwari/apps/api/internal/types"
)

func Setup(router fiber.Router, services types.Services) fiber.Router {
	middleware := router.Group("users")
	middleware.Get("", get(services))

	return middleware
}

func get(services types.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ids := c.Query("ids", "")
		names := c.Query("names", "")

		if ids == "" && names == "" {
			return c.JSON(fiber.Map{"message": "ids or names not provided"})
		}

		users, err := handleGet(ids, names, services)
		if err != nil {
			return err
		}

		return c.JSON(users)
	}
}
