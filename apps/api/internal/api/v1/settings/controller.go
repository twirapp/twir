package settings

import (
	"github.com/gofiber/fiber/v2"
	dashboardaccess "github.com/satont/tsuwari/apps/api/internal/api/v1/settings/dashboardAccess"
	"github.com/satont/tsuwari/apps/api/internal/types"
)

func Setup(router fiber.Router, services types.Services) fiber.Router {
	middleware := router.Group("settings")
	dashboardaccess.Setup(middleware, services)

	return middleware
}
