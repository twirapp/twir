package gamesAliases

import (
	"github.com/gofiber/fiber"
	"github.com/satont/tsuwari/apps/api/internal/middlewares"
	"github.com/satont/tsuwari/apps/api/internal/types"
)

func Setup(router fiber.Router, services types.Services) fiber.Router {
	middleware := router.Group("games-aliases")

	middleware.Get("/", get(services))
	middleware.Post("/", post(services))
	middleware.Delete(":gameAliasId", delete(services))
	middleware.Put(":gameAliasId", put(services))

	return middleware
}

func get(services types.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		gamesAliases, err := handleGetGamesAliases(c.Params("channelId"), services)
		if err != nil {
			return err
		}

		return c.JSON(gamesAliases)
	}
}

func post(services types.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		dto := &gameAliasDto{}
		err := middlewares.ValidateBody(
			c,
			services.Validator,
			services.ValidatorTranslator,
			dto,
		)
		if err != nil {
			return err
		}

		return c.JSON(dto)
	}
}

func delete(services types.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {

	}
}

func put(services types.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {

	}
}
