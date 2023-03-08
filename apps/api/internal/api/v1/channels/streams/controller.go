package streams

import (
	"github.com/gofiber/fiber/v2"
	"github.com/satont/tsuwari/apps/api/internal/types"
)

func Setup(router fiber.Router, services *types.Services) fiber.Router {
	middleware := router.Group("streams")
	middleware.Get("", get(services))

	return middleware
}

// Streams godoc
// @Security ApiKeyAuth
// @Summary      Get channel stream
// @Tags         Streams
// @Produce      json
// @Param        channelId   path      string  true  "ChannelId" default({{channelId}})
// @Success      200  {object}  model.ChannelsStreams
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/streams [get]
func get(services *types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		stream, err := handleGet(c.Params("channelId"), services)
		if err != nil {
			return err
		}
		return c.JSON(stream)
	}
}
