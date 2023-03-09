package channels

import (
	"github.com/gofiber/fiber/v2"
	"github.com/satont/tsuwari/apps/api/internal/api/v1/channels/bot"
	"github.com/satont/tsuwari/apps/api/internal/api/v1/channels/commands"
	"github.com/satont/tsuwari/apps/api/internal/api/v1/channels/events"
	"github.com/satont/tsuwari/apps/api/internal/api/v1/channels/greetings"
	"github.com/satont/tsuwari/apps/api/internal/api/v1/channels/integrations"
	"github.com/satont/tsuwari/apps/api/internal/api/v1/channels/keywords"
	"github.com/satont/tsuwari/apps/api/internal/api/v1/channels/moderation"
	"github.com/satont/tsuwari/apps/api/internal/api/v1/channels/rewards"
	"github.com/satont/tsuwari/apps/api/internal/api/v1/channels/roles"
	"github.com/satont/tsuwari/apps/api/internal/api/v1/channels/streams"
	"github.com/satont/tsuwari/apps/api/internal/api/v1/channels/timers"
	"github.com/satont/tsuwari/apps/api/internal/api/v1/channels/variables"
	"github.com/satont/tsuwari/apps/api/internal/middlewares"
	"github.com/satont/tsuwari/apps/api/internal/types"
)

func CreateChannelsRouter(router fiber.Router, services *types.Services) fiber.Router {
	channel := router.
		Group("channels/:channelId").
		Use(middlewares.AttachUser(services)).
		Use(middlewares.CheckHasAccessToDashboard)

	bot.NewController(channel, services)
	commands.NewController(channel, services)
	events.NewController(channel, services)
	greetings.NewController(channel, services)
	integrations.NewController(channel, services)
	keywords.NewController(channel, services)
	moderation.NewController(channel, services)
	rewards.NewController(channel, services)
	roles.NewController(channel, services)
	streams.NewController(channel, services)
	timers.NewController(channel, services)
	variables.NewController(channel, services)

	return channel
}
