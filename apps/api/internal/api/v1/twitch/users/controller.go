package users

import (
	"github.com/gofiber/fiber/v2"
	"github.com/satont/tsuwari/apps/api/internal/types"
	"github.com/satont/tsuwari/libs/twitch"
)

func Setup(router fiber.Router, services types.Services) fiber.Router {
	categoriesMiddleware := router.Group("categories")
	categoriesMiddleware.Get("", getCategories(services))

	usersMiddleware := router.Group("users")
	usersMiddleware.Get("", get(services))

	return router
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

func getCategories(services types.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		category := c.Query("category", "")

		categories, err := handleGetCategories(category)
		if err != nil {
			return err
		}

		return c.JSON(categories)
	}
}
