package rewards

import (
	"github.com/gofiber/fiber/v2"
	"github.com/satont/twir/apps/api/internal/types"
)

func Setup(router fiber.Router, services types.Services) fiber.Router {
	middleware := router.Group("rewards")
	middleware.Get("", get(services))

	return middleware
}

// Rewards godoc
// @Security ApiKeyAuth
// @Summary      Get channel rewards list
// @Tags         Rewards
// @Accept       json
// @Produce      json
// @Param        channelId   path      string  true  "ChannelId" default({{channelId}})
// @Success      200  {array}  helix.ChannelCustomReward
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/rewards [get]
func get(services types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		rewards, err := handleGet(c.Params("channelId"))
		if err != nil {
			return err
		}

		return c.JSON(rewards)
	}
}
