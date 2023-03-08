package youtube_sr

import (
	"github.com/gofiber/fiber/v2"
	"github.com/satont/tsuwari/apps/api/internal/middlewares"
	"github.com/satont/tsuwari/apps/api/internal/types"
	youtube "github.com/satont/tsuwari/libs/types/types/api/modules"
)

func Setup(router fiber.Router, services *types.Services) fiber.Router {
	middleware := router.Group("youtube-sr")
	middleware.Get("", get(services))
	middleware.Post("", post(services))
	middleware.Get("search", getSearch(services))

	return middleware
}

// YouTube godoc
// @Security ApiKeyAuth
// @Summary      Get YouTube settings
// @Tags         Modules|YouTube
// @Accept       json
// @Produce      json
// @Param        channelId   path      string  true  "ChannelId" default({{channelId}})
// @Success      200  {object}  youtube.YouTubeSettings
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/modules/youtube-sr [get]
func get(services *types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		settings, err := handleGet(c.Params("channelId"), services)
		if err != nil {
			return err
		}

		return c.JSON(settings)
	}
}

// YouTube godoc
// @Security ApiKeyAuth
// @Summary      Search channel or video
// @Tags         Modules|YouTube
// @Accept       json
// @Produce      json
// @Param        channelId   path      string  true  "ChannelId" default({{channelId}})
// @Param        query   query      string  true  "Input string"
// @Param        type   query      string  true  "channel or video"
// @Success      200  {array}  youtube.SearchResult
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/modules/youtube-sr/search [get]
func getSearch(service *types.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		results, err := handleSearch(c.Query("query"), c.Query("type"))
		if err != nil {
			return err
		}

		return c.JSON(results)
	}
}

func post(services *types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		dto := youtube.YouTubeSettings{}
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
