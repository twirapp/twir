package public

import (
	config "github.com/satont/twir/libs/config"
	"github.com/twirapp/twir/apps/api-gql/internal/httpserver"
	cachedcommands "github.com/twirapp/twir/libs/cache/commands"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Opts struct {
	fx.In

	Server         *httpserver.Server
	Gorm           *gorm.DB
	Config         config.Config
	CachedCommands *cachedcommands.CachedCommandsClient
}

type Public struct {
	gorm           *gorm.DB
	config         config.Config
	cachedCommands *cachedcommands.CachedCommandsClient
}

func New(opts Opts) *Public {
	p := &Public{
		gorm:           opts.Gorm,
		config:         opts.Config,
		cachedCommands: opts.CachedCommands,
	}

	opts.Server.GET("/v1/public/badges", p.HandleBadgesGet)
	opts.Server.GET("/v1/public/channels/:channelId/commands", p.HandleChannelCommandsGet)

	return p
}
