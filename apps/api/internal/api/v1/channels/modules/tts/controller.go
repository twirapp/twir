package tts

import (
	"github.com/gofiber/fiber/v2"
	"github.com/satont/tsuwari/apps/api/internal/api/v1/modules/tts/users"
	"github.com/satont/tsuwari/apps/api/internal/middlewares"
	"github.com/satont/tsuwari/apps/api/internal/types"
	"github.com/satont/tsuwari/libs/types/types/api/modules"
)

func Setup(router fiber.Router, services *types.Services) fiber.Router {
	middleware := router.Group("tts")
	middleware.Get("", get(services))
	middleware.Post("", post(services))
	middleware.Get("info", getInfo(services))

	users.Setup(middleware, services)

	return middleware
}

// TTS Godoc
// @Security ApiKeyAuth
// @Summary Get TTS settings
// @Description Get TTS settings
// @Tags TTS
// @Accept json
// @Produce json
// @Param channelId path string true "Channel ID"
// @Success 200 {object} modules.TTSSettings
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/modules/tts [get]
func get(services *types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		settings, err := handleGet(c.Params("channelId"), services)
		if err != nil {
			return err
		}

		return c.JSON(settings)
	}
}

// TTS Godoc
// @Security ApiKeyAuth
// @Summary Set TTS settings
// @Description Set TTS settings
// @Tags TTS
// @Accept json
// @Produce json
// @Param channelId path string true "Channel ID"
// @Param settings body modules.TTSSettings true "TTS settings"
// @Success 204
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/modules/tts [post]
func post(services *types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		dto := modules.TTSSettings{}
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

func getInfo(services *types.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		info := handleGetInfo(services)

		return c.JSON(info)
	}
}
