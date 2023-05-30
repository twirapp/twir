package webhooks

import (
	"github.com/gofiber/fiber/v2"
	donate_stream "github.com/satont/tsuwari/apps/api/internal/api/webhooks/integrations/donate_stream"
	"github.com/satont/tsuwari/apps/api/internal/types"
)

func Setup(router fiber.Router, services types.Services) {
	group := router.Group("/webhooks")

	donate_stream.Setup(group, services)
}
