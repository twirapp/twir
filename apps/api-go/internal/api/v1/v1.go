package apiv1

import (
	"github.com/gofiber/fiber/v2"

	commands "github.com/satont/tsuwari/apps/api-go/internal/api/v1/commands"
	greetings "github.com/satont/tsuwari/apps/api-go/internal/api/v1/greetings"
	keywords "github.com/satont/tsuwari/apps/api-go/internal/api/v1/keywords"
	"github.com/satont/tsuwari/apps/api-go/internal/api/v1/timers"
	"github.com/satont/tsuwari/apps/api-go/internal/types"
)

func Setup(router fiber.Router, services types.Services) fiber.Router {
	channelsGroup := router.Group("channels/:channelId")
	commands.Setup(channelsGroup, services)
	greetings.Setup(channelsGroup, services)
	keywords.Setup(channelsGroup, services)
	timers.Setup(channelsGroup, services)

	return router
}
