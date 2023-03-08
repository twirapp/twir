package channels

import (
	"github.com/gofiber/fiber/v2"
	"github.com/satont/tsuwari/apps/api/internal/api/v1/channels/bot"
	"github.com/satont/tsuwari/apps/api/internal/api/v1/channels/commands"
	"github.com/satont/tsuwari/apps/api/internal/api/v1/channels/events"
	"github.com/satont/tsuwari/apps/api/internal/api/v1/channels/greetings"
	"github.com/satont/tsuwari/apps/api/internal/api/v1/channels/keywords"
	"github.com/satont/tsuwari/apps/api/internal/middlewares"
	"github.com/satont/tsuwari/apps/api/internal/types"
)

func CreateChannelsRouter(router fiber.Router, services *types.Services) fiber.Router {
	channel := router.
		Group("channels/:channelId").
		Use(middlewares.AttachUser(services)).
		Use(middlewares.CheckHasAccessToDashboard)

	bot.NewBot(channel, services)
	commands.NewCommands(channel, services)
	events.NewEvents(channel, services)
	greetings.NewGreetings(channel, services)
	keywords.NewKeywords(channel, services)

	return channel
}
