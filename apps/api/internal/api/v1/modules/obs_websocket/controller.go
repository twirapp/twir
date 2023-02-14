package obs_websocket

import (
	"github.com/gofiber/fiber/v2"
	"github.com/satont/tsuwari/apps/api/internal/middlewares"
	"github.com/satont/tsuwari/apps/api/internal/types"
	modules "github.com/satont/tsuwari/libs/types/types/api/modules"
)

func Setup(router fiber.Router, services types.Services) fiber.Router {
	middleware := router.Group("obs-websocket")
	middleware.Get("", get(services))
	middleware.Post("", post(services))

	return middleware
}

// OBSWebSocket godoc
// @Security ApiKeyAuth
// @Summary      Get OBSWebSocket settings
// @Tags         Modules|OBSWebSocket
// @Accept       json
// @Produce      json
// @Param        channelId   path      string  true  "ChannelId"
// @Success      200  {object}  youtube.YouTubeSettings
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/modules/obs-websocket [get]
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
