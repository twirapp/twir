package obs_websocket

import (
	"github.com/gofiber/fiber/v2"
	"github.com/satont/twir/apps/api/internal/middlewares"
	"github.com/satont/twir/apps/api/internal/types"
	modules "github.com/satont/twir/libs/types/types/api/modules"
)

func Setup(router fiber.Router, services types.Services) fiber.Router {
	middleware := router.Group("obs-websocket")
	middleware.Get("", get(services))
	middleware.Post("", post(services))
	middleware.Get("data", getData(services))

	return middleware
}

func get(services types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		settings, err := handleGet(c.Params("channelId"), services)
		if err != nil {
			return err
		}

		return c.JSON(settings)
	}
}

func post(services types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		dto := modules.OBSWebSocketSettings{}
		err := middlewares.ValidateBody(
			c,
			services.Validator,
			services.ValidatorTranslator,
			&dto,
		)
		if err != nil {
			return err
		}

		err = handlePost(c.Params("channelId"), &dto, services)
		if err != nil {
			return err
		}

		return c.SendStatus(204)
	}
}

func getData(services types.Services) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		data, err := handleGetData(ctx.Params("channelId"))
		if err != nil {
			return err
		}

		return ctx.JSON(data)
	}
}
