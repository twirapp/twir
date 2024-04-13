package public

import (
	config "github.com/satont/twir/libs/config"
	"github.com/twirapp/twir/apps/api-gql/internal/httpserver"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Opts struct {
	fx.In

	Server *httpserver.Server
	Gorm   *gorm.DB
	Config config.Config
}

type Public struct {
	gorm   *gorm.DB
	config config.Config
}

func New(opts Opts) *Public {
	p := &Public{
		gorm:   opts.Gorm,
		config: opts.Config,
	}

	opts.Server.GET("/v1/public/badges", p.HandleBadgesGet)
	opts.Server.GET("/v1/public/channels/:channelId/commands", p.HandleChannelCommandsGet)

	return p
}
