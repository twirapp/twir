package variables

import (
	"github.com/gofiber/fiber/v2"
	"github.com/satont/tsuwari/apps/api-go/internal/types"
)

func Setup(router fiber.Router, services types.Services) fiber.Router {
	middleware := router.Group("variables")
	middleware.Get("", get(services))
	middleware.Post("", post(services))
	middleware.Delete(":variableId", delete(services))
	middleware.Put(":variableId", put(services))
	middleware.Get("builtin", getBuiltIn(services))

	return middleware
}

func get(services types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		variables, err := handleGet(c.Params("channelId"), services)
		if err != nil {
			return err
		}
		return c.JSON(variables)
	}
}

func getBuiltIn(services types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return nil
	}
}

func post(services types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return nil
	}
}

func delete(services types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return nil
	}
}

func put(services types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return nil
	}
}
