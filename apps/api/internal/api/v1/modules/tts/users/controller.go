package users

import (
	"github.com/satont/twir/apps/api/internal/middlewares"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/satont/twir/apps/api/internal/types"
)

func Setup(router fiber.Router, services types.Services) fiber.Router {
	middleware := router.Group("users")
	middleware.Get("", get(services))
	middleware.Delete("", delete(services))

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
		dto := &deleteDto{}
		err := middlewares.ValidateBody(
			c,
			services.Validator,
			services.ValidatorTranslator,
			dto,
		)
		if err != nil {
			return err
		}

		err = handleDelete(c.Params("channelId"), dto, services)

		return c.SendStatus(http.StatusOK)
	}
}
