package streams

import (
	"github.com/gofiber/fiber/v2"
	"github.com/satont/tsuwari/apps/api/internal/types"
)

type Streams struct {
	services *types.Services
}

func NewController(router fiber.Router, services *types.Services) fiber.Router {
	streams := &Streams{
		services: services,
	}

	return router.Group("streams").
		Get("", streams.get)
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
func (c *Streams) get(ctx *fiber.Ctx) error {
	stream, err := c.getService(ctx.Params("channelId"))
	if err != nil {
		return err
	}
	return ctx.JSON(stream)
}
