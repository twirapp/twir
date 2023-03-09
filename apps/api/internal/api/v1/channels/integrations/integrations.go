package integrations

import (
	"github.com/gofiber/fiber/v2"
	"github.com/satont/tsuwari/apps/api/internal/api/v1/channels/integrations/donatello"
	"github.com/satont/tsuwari/apps/api/internal/api/v1/channels/integrations/donatepay"
	"github.com/satont/tsuwari/apps/api/internal/api/v1/channels/integrations/donationalerts"
	"github.com/satont/tsuwari/apps/api/internal/api/v1/channels/integrations/faceit"
	"github.com/satont/tsuwari/apps/api/internal/api/v1/channels/integrations/lastfm"
	"github.com/satont/tsuwari/apps/api/internal/api/v1/channels/integrations/spotify"
	"github.com/satont/tsuwari/apps/api/internal/types"
)

func NewController(router fiber.Router, services *types.Services) fiber.Router {
	route := router.Group("integrations")

	donatello.NewController(route, services)
	donatepay.NewController(route, services)
	donationalerts.NewController(route, services)
	faceit.NewController(route, services)
	lastfm.NewController(route, services)
	spotify.NewController(route, services)

	return route
}
