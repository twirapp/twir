package apiv1

import (
	"github.com/gofiber/fiber/v2"

	commands "github.com/satont/tsuwari/apps/api-go/internal/api/v1/commands"
	"github.com/satont/tsuwari/apps/api-go/internal/types"
)

func Setup(router fiber.Router, services types.Services) fiber.Router {
	channelsGroup := router.Group("channels")
	commands.Setup(channelsGroup, services)

	return router
}
