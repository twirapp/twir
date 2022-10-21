package greetings

import (
	"github.com/gofiber/fiber/v2"
	"github.com/satont/tsuwari/apps/api-go/internal/middlewares"
	"github.com/satont/tsuwari/apps/api-go/internal/types"
)

func Setup(router fiber.Router, services types.Services) fiber.Router {
	middleware := router.Group("greetings")

	middleware.Get("", func(c *fiber.Ctx) error {
		c.JSON(HandleGet(c.Params("channelId"), services))

		return nil
	})
	middleware.Post("", func(c *fiber.Ctx) error {
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

		greeting, err := HandlePost(c.Params("channelId"), dto, services)
		if err == nil && greeting != nil {
			return c.JSON(greeting)
		}

		return err
	})
	middleware.Delete(":greetingId", func(c *fiber.Ctx) error {
		err := HandleDelete(c.Params("greetingId"), services)
		if err != nil {
			return err
		}
		return c.SendStatus(fiber.StatusOK)
	})

	return middleware
}
