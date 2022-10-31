package timers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/satont/tsuwari/apps/api-go/internal/middlewares"
	"github.com/satont/tsuwari/apps/api-go/internal/types"
)

func Setup(router fiber.Router, services types.Services) fiber.Router {
	middleware := router.Group("timers")
	middleware.Get("", get(services))
	middleware.Post("", post(services))
	middleware.Delete(":timerId", delete(services))
	middleware.Put(":timerId", put(services))

	return middleware
}

func get(services types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		timers := handleGet(c.Params("channelId"), services)

		return c.JSON(timers)
	}
}

func post(services types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		dto := &timerDto{}
		err := middlewares.ValidateBody(
			c,
			services.Validator,
			services.ValidatorTranslator,
			dto,
		)
		if err != nil {
			return err
		}

		cmd, err := handlePost(c.Params("channelId"), dto, services)
		if err == nil {
			return c.JSON(cmd)
		}

		return err
	}
}

func delete(services types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		err := handleDelete(c.Params("timerId"), services)
		if err != nil {
			return err
		}
		return c.SendStatus(200)
	}
}

func put(services types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		dto := &timerDto{}
		err := middlewares.ValidateBody(
			c,
			services.Validator,
			services.ValidatorTranslator,
			dto,
		)
		if err != nil {
			return err
		}

		cmd, err := handlePut(c.Params("channelId"), c.Params("timerId"), dto, services)
		if err == nil {
			return c.JSON(cmd)
		}

		return err
	}
}
