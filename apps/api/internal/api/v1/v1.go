package apiv1

import (
	"github.com/gofiber/fiber/v2"

	"github.com/satont/tsuwari/apps/api/internal/api/v1/bot"
	"github.com/satont/tsuwari/apps/api/internal/api/v1/commands"
	"github.com/satont/tsuwari/apps/api/internal/api/v1/feedback"
	"github.com/satont/tsuwari/apps/api/internal/api/v1/greetings"
	"github.com/satont/tsuwari/apps/api/internal/api/v1/integrations/donationalerts"
	"github.com/satont/tsuwari/apps/api/internal/api/v1/integrations/faceit"
	"github.com/satont/tsuwari/apps/api/internal/api/v1/integrations/lastfm"
	"github.com/satont/tsuwari/apps/api/internal/api/v1/integrations/spotify"
	"github.com/satont/tsuwari/apps/api/internal/api/v1/integrations/streamlabs"
	"github.com/satont/tsuwari/apps/api/internal/api/v1/integrations/vk"
	"github.com/satont/tsuwari/apps/api/internal/api/v1/keywords"
	"github.com/satont/tsuwari/apps/api/internal/api/v1/moderation"
	"github.com/satont/tsuwari/apps/api/internal/api/v1/modules/youtube_sr"
	"github.com/satont/tsuwari/apps/api/internal/api/v1/settings"
	"github.com/satont/tsuwari/apps/api/internal/api/v1/stats"
	"github.com/satont/tsuwari/apps/api/internal/api/v1/streams"
	"github.com/satont/tsuwari/apps/api/internal/api/v1/timers"
	"github.com/satont/tsuwari/apps/api/internal/api/v1/twitch/users"
	"github.com/satont/tsuwari/apps/api/internal/api/v1/variables"
	"github.com/satont/tsuwari/apps/api/internal/middlewares"
	"github.com/satont/tsuwari/apps/api/internal/types"
)

func Setup(router fiber.Router, services types.Services) fiber.Router {
	feedback.Setup(router, services)
	stats.Setup(router, services)

	channelsGroup := router.Group("channels/:channelId")
	channelsGroup.Use(middlewares.CheckUserAuth(services))
	channelsGroup.Use(middlewares.CheckHasAccessToDashboard)

	commands.Setup(channelsGroup, services)
	greetings.Setup(channelsGroup, services)
	keywords.Setup(channelsGroup, services)
	timers.Setup(channelsGroup, services)
	moderation.Setup(channelsGroup, services)
	bot.Setup(channelsGroup, services)
	streams.Setup(channelsGroup, services)
	variables.Setup(channelsGroup, services)
	settings.Setup(channelsGroup, services)

	integrationsGroup := channelsGroup.Group("integrations")
	donationalerts.Setup(integrationsGroup, services)
	streamlabs.Setup(integrationsGroup, services)
	faceit.Setup(integrationsGroup, services)
	lastfm.Setup(integrationsGroup, services)
	vk.Setup(integrationsGroup, services)
	spotify.Setup(integrationsGroup, services)

	modulesGroup := channelsGroup.Group("modules")
	youtube_sr.Setup(modulesGroup, services)

	twitchGroup := router.Group("twitch")
	twitchGroup.Use(middlewares.CheckUserAuth(services))
	users.Setup(twitchGroup, services)

	return router
}
