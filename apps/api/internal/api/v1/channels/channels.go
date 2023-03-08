package channels

import (
	"github.com/gofiber/fiber/v2"
	"github.com/satont/tsuwari/apps/api/internal/api/v1/channels/bot"
	"github.com/satont/tsuwari/apps/api/internal/middlewares"
	"github.com/satont/tsuwari/apps/api/internal/types"
)

func CreateChannelRouter(router fiber.Router, services *types.Services) fiber.Router {
	return router.
		Group("channels/:channelId").
		Use(middlewares.AttachUser(services)).
		Use(middlewares.CheckHasAccessToDashboard).
		Use(bot.NewBot(router, services))
}
