package giveaways

import (
	"github.com/gofiber/fiber/v2"
	"github.com/satont/tsuwari/apps/api/internal/middlewares"
	"github.com/satont/tsuwari/apps/api/internal/types"
)

func Setup(router fiber.Router, services types.Services) fiber.Router {
	middleware := router.Group("giveaways")

	middleware.Get("", getMany(services))
	middleware.Get(":giveawayId", get(services))
	middleware.Post("", post(services))
	middleware.Delete(":giveawayId", delete(services))
	middleware.Patch(":greetingId", patch(services))

	return middleware
}

func post(services types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		dto := &postGiveawayDto{}
		err := middlewares.ValidateBody(
			c,
			services.Validator,
			services.ValidatorTranslator,
			dto,
		)
		if err != nil {
			return err
		}

		giveaway, err := handlePost(c.Params("channelId"), services)
		if err == nil && giveaway != nil {
			return c.JSON(giveaway)
		}
		return err
	}
}

func get(services types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return c.JSON(handleGetAll(c.Params("channelId"), services))
	}
}

func getMany(services types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return c.JSON(handleGetAll(c.Params("channelId"), services))
	}
}
