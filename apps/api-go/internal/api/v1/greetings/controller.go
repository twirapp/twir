package greetings

import (
	"github.com/gofiber/fiber/v2"
	"github.com/satont/tsuwari/apps/api-go/internal/middlewares"
	"github.com/satont/tsuwari/apps/api-go/internal/types"
)

func Setup(router fiber.Router, services types.Services) fiber.Router {
	middleware := router.Group("greetings")

	middleware.Get("", get(services))
	middleware.Post("", post(services))
	middleware.Delete(":greetingId", delete(services))
	middleware.Put(":greetingId", put(services))

	return middleware
}

func get(services types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		c.JSON(handleGet(c.Params("channelId"), services))

		return nil
	}
}

func post(services types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		dto := &greetingsDto{}
		err := middlewares.ValidateBody(
			c,
			services.Validator,
			services.ValidatorTranslator,
			dto,
		)
		if err != nil {
			return err
		}

		greeting, err := handlePost(c.Params("channelId"), dto, services)
		if err == nil && greeting != nil {
			return c.JSON(greeting)
		}

		return err
	}
}

func delete(services types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		err := handleDelete(c.Params("greetingId"), services)
		if err != nil {
			return err
		}
		return c.SendStatus(fiber.StatusOK)
	}
}

func put(services types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		dto := &greetingsDto{}
		err := middlewares.ValidateBody(
			c,
			services.Validator,
			services.ValidatorTranslator,
			dto,
		)
		if err != nil {
			return err
		}
		greeting, err := handleUpdate(c.Params("greetingId"), dto, services)
		if err != nil {
			return err
		}

		return c.JSON(greeting)
	}
}
