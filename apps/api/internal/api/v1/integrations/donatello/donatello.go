package donatello

import (
	"github.com/gofiber/fiber/v2"
	"github.com/satont/twir/apps/api/internal/types"
)

func Setup(router fiber.Router, services types.Services) fiber.Router {
	middleware := router.Group("donatello")
	middleware.Get("", get(services))

	return middleware
}

// Integrations godoc
// @Security ApiKeyAuth
// @Summary      Get Donatello token
// @Tags         Integrations|Donatello
// @Accept       json
// @Produce      plain
// @Param        channelId   path      string  true  "ChannelId" default({{channelId}})
// @Success      200  {string} string
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/integrations/donatepay [get]
func get(services types.Services) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		data, err := handleGet(services, ctx.Params("channelId"))

		if err != nil {
			return err
		}

		return ctx.JSON(data)
	}
}
