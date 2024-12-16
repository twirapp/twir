package public

import (
	config "github.com/satont/twir/libs/config"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/apps/api-gql/internal/server"
	badges_with_users "github.com/twirapp/twir/apps/api-gql/internal/services/badges-with-users"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Opts struct {
	fx.In

	Server                 *server.Server
	Gorm                   *gorm.DB
	Config                 config.Config
	CachedCommands         *generic_cacher.GenericCacher[[]model.ChannelsCommands]
	BadgesWithUsersService *badges_with_users.Service
}

type Public struct {
	gorm                   *gorm.DB
	config                 config.Config
	cachedCommands         *generic_cacher.GenericCacher[[]model.ChannelsCommands]
	badgesWithUsersService *badges_with_users.Service
}

func New(opts Opts) *Public {
	p := &Public{
		gorm:                   opts.Gorm,
		config:                 opts.Config,
		cachedCommands:         opts.CachedCommands,
		badgesWithUsersService: opts.BadgesWithUsersService,
	}

	opts.Server.GET("/v1/public/badges", p.HandleBadgesGet)
	opts.Server.GET("/v1/public/channels/:channelId/commands", p.HandleChannelCommandsGet)

	return p
}
