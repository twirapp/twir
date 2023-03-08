package rewards

import (
	"github.com/gofiber/fiber/v2"
	"github.com/satont/tsuwari/apps/api/internal/types"
)

type Rewards struct {
	services *types.Services
}

func NewRewards(router fiber.Router, services *types.Services) fiber.Router {
	rewards := &Rewards{
		services: services,
	}

	return router.Group("rewards").
		Get("", rewards.get)
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
func (c *Rewards) get(ctx *fiber.Ctx) error {
	rewards, err := c.getService(ctx.Params("channelId"))
	if err != nil {
		return err
	}

	return ctx.JSON(rewards)
}
