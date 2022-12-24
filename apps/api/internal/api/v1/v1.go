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
	"github.com/satont/tsuwari/apps/api/internal/api/v1/settings"
	"github.com/satont/tsuwari/apps/api/internal/api/v1/stats"
	"github.com/satont/tsuwari/apps/api/internal/api/v1/streams"
	"github.com/satont/tsuwari/apps/api/internal/api/v1/timers"
	"github.com/satont/tsuwari/apps/api/internal/api/v1/variables"
	"github.com/satont/tsuwari/apps/api/internal/middlewares"
)

func Setup(router fiber.Router) fiber.Router {
	feedback.Setup(router)
	stats.Setup(router)

	channelsGroup := router.Group("channels/:channelId")
	channelsGroup.Use(middlewares.CheckUserAuth)
	channelsGroup.Use(middlewares.CheckHasAccessToDashboard)

	commands.Setup(channelsGroup)
	greetings.Setup(channelsGroup)
	keywords.Setup(channelsGroup)
	timers.Setup(channelsGroup)
	moderation.Setup(channelsGroup)
	bot.Setup(channelsGroup)
	streams.Setup(channelsGroup)
	variables.Setup(channelsGroup)
	settings.Setup(channelsGroup)

	integrationsGroup := channelsGroup.Group("integrations")
	donationalerts.Setup(integrationsGroup)
	streamlabs.Setup(integrationsGroup)
	faceit.Setup(integrationsGroup)
	lastfm.Setup(integrationsGroup)
	vk.Setup(integrationsGroup)
	spotify.Setup(integrationsGroup)

	return router
}
