package public

import (
	config "github.com/satont/twir/libs/config"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/apps/api-gql/internal/httpserver"
	db_generic_cacher "github.com/twirapp/twir/libs/cache/db-generic-cacher"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Opts struct {
	fx.In

	Server         *httpserver.Server
	Gorm           *gorm.DB
	Config         config.Config
	CachedCommands *db_generic_cacher.GenericCacher[[]model.ChannelsCommands]
}

type Public struct {
	gorm           *gorm.DB
	config         config.Config
	cachedCommands *db_generic_cacher.GenericCacher[[]model.ChannelsCommands]
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
