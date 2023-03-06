package users

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/satont/tsuwari/apps/api/internal/types"
)

func Setup(router fiber.Router, services types.Services) fiber.Router {
	middleware := router.Group("users")
	middleware.Get("", get(services))
	middleware.Delete("clear", deleteAll(services))
	middleware.Delete(":userId", delete(services))

	return middleware
}

func get(services types.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		data, err := handleGet(c.Params("channelId"), services)
		if err != nil {
			return err
		}

		return c.JSON(data)
	}
}

func delete(services types.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := handleDelete(c.Params("channelId"), c.Params("userId"), services)
		if err != nil {
			return err
		}

		return c.SendStatus(http.StatusOK)
	}
}

func deleteAll(services types.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := handleDeleteAll(c.Params("channelId"), services)
		if err != nil {
			return err
		}

		return c.SendStatus(http.StatusOK)
	}
}
