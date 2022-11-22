package words_counters

import (
	"github.com/gofiber/fiber/v2"
	"github.com/satont/tsuwari/apps/api/internal/middlewares"
	"github.com/satont/tsuwari/apps/api/internal/types"
)

func Setup(router fiber.Router, services types.Services) fiber.Router {
	middleware := router.Group("words_counters")
	middleware.Get("", get(services))
	middleware.Post("", post(services))
	middleware.Put(":id", put(services))
	middleware.Delete(":id", delete(services))

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

func post(services types.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		dto := &wordsCountersDto{}
		err := middlewares.ValidateBody(
			c,
			services.Validator,
			services.ValidatorTranslator,
			dto,
		)
		if err != nil {
			return err
		}
		data, err := handlePost(c.Params("channelId"), dto, services)
		if err != nil {
			return err
		}

		return c.Status(201).JSON(data)
	}
}

func put(services types.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		dto := &wordsCountersDto{}
		err := middlewares.ValidateBody(
			c,
			services.Validator,
			services.ValidatorTranslator,
			dto,
		)
		if err != nil {
			return err
		}
		data, err := handlePut(c.Params("channelId"), c.Params("id"), dto, services)
		if err != nil {
			return err
		}

		return c.Status(201).JSON(data)
	}
}

func delete(services types.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := handleDelete(c.Params("channelId"), c.Params("id"), services)
		if err != nil {
			return nil
		}

		return c.SendStatus(200)
	}
}
