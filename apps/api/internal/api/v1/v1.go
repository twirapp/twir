package apiv1

import (
	"github.com/gofiber/fiber/v2"
	admin_users "github.com/satont/tsuwari/apps/api/internal/api/v1/admin/users"
	"github.com/satont/tsuwari/apps/api/internal/api/v1/community"
	"github.com/satont/tsuwari/apps/api/internal/api/v1/events"
	"github.com/satont/tsuwari/apps/api/internal/api/v1/integrations/donatello"
	"github.com/satont/tsuwari/apps/api/internal/api/v1/integrations/donatepay"
	public_commands "github.com/satont/tsuwari/apps/api/internal/api/v1/public/commands"
	"github.com/satont/tsuwari/apps/api/internal/api/v1/public/song_requests"
	"github.com/satont/tsuwari/apps/api/internal/api/v1/rewards"
	"github.com/satont/tsuwari/apps/api/internal/api/v1/twitch/users"
	"github.com/satont/tsuwari/apps/api/internal/middlewares"

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
	"github.com/satont/tsuwari/apps/api/internal/api/v1/variables"
	"github.com/satont/tsuwari/apps/api/internal/types"
)

func Setup(router fiber.Router, services types.Services) fiber.Router {
	feedback.Setup(router, services)
	stats.Setup(router, services)

	adminGroup := router.Group("admin")
	adminGroup.Use(middlewares.CheckUserAuth(services))
	adminGroup.Use(middlewares.IsAdmin)
	admin_users.Setup(adminGroup, services)

	twitchGroup := router.Group("twitch")
	users.Setup(twitchGroup, services)

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
	rewards.Setup(channelsGroup, services)
	community.Setup(channelsGroup, services)
	events.Setup(channelsGroup, services)

	integrationsGroup := channelsGroup.Group("integrations")
	donationalerts.Setup(integrationsGroup, services)
	streamlabs.Setup(integrationsGroup, services)
	faceit.Setup(integrationsGroup, services)
	lastfm.Setup(integrationsGroup, services)
	vk.Setup(integrationsGroup, services)
	spotify.Setup(integrationsGroup, services)
	donatepay.Setup(integrationsGroup, services)
	donatello.Setup(integrationsGroup, services)

	modulesGroup := channelsGroup.Group("modules")
	youtube_sr.Setup(modulesGroup, services)

	publicGroup := router.Group("p")
	public_commands.Setup(publicGroup, services)
	song_requests.Setup(publicGroup, services)

	return router
}
