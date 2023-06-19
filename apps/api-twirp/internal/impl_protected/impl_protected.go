package impl_protected

import (
	"github.com/alexedwards/scs/v2"
	"github.com/redis/go-redis/v9"
	"github.com/satont/tsuwari/apps/api-twirp/internal/impl_deps"
	"github.com/satont/tsuwari/apps/api-twirp/internal/impl_protected/bot"
	"github.com/satont/tsuwari/apps/api-twirp/internal/impl_protected/commands"
	"github.com/satont/tsuwari/apps/api-twirp/internal/impl_protected/community"
	"github.com/satont/tsuwari/apps/api-twirp/internal/impl_protected/events"
	"github.com/satont/tsuwari/apps/api-twirp/internal/impl_protected/greetings"
	"github.com/satont/tsuwari/apps/api-twirp/internal/impl_protected/integrations"
	"github.com/satont/tsuwari/apps/api-twirp/internal/impl_protected/keywords"
	"github.com/satont/tsuwari/apps/api-twirp/internal/impl_protected/modules"
	"github.com/satont/tsuwari/apps/api-twirp/internal/impl_protected/rewards"
	"github.com/satont/tsuwari/apps/api-twirp/internal/impl_protected/roles"
	"github.com/satont/tsuwari/apps/api-twirp/internal/impl_protected/timers"
	config "github.com/satont/tsuwari/libs/config"
	"github.com/satont/tsuwari/libs/grpc/generated/bots"
	"github.com/satont/tsuwari/libs/grpc/generated/tokens"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Protected struct {
	*integrations.Integrations
	*keywords.Keywords
	*modules.Modules
	*bot.Bot
	*commands.Commands
	*community.Community
	*events.Events
	*greetings.Greetings
	*rewards.Rewards
	*roles.Roles
	*timers.Timers
}

type Opts struct {
	fx.In

	Redis          *redis.Client
	DB             *gorm.DB
	Config         *config.Config
	SessionManager *scs.SessionManager
	TokensGrpc     tokens.TokensClient
	BotsGrpc       bots.BotsClient
}

func New(opts Opts) *Protected {
	d := &impl_deps.Deps{
		Redis:          opts.Redis,
		Db:             opts.DB,
		Config:         opts.Config,
		SessionManager: opts.SessionManager,
		Grpc: &impl_deps.Grpc{
			Tokens: opts.TokensGrpc,
			Bots:   opts.BotsGrpc,
		},
	}

	return &Protected{
		Integrations: &integrations.Integrations{Deps: d},
		Keywords:     &keywords.Keywords{Deps: d},
		Modules:      &modules.Modules{Deps: d},
		Bot:          &bot.Bot{Deps: d},
		Commands:     &commands.Commands{Deps: d},
		Community:    &community.Community{Deps: d},
		Events:       &events.Events{Deps: d},
		Greetings:    &greetings.Greetings{Deps: d},
		Rewards:      &rewards.Rewards{Deps: d},
		Roles:        &roles.Roles{Deps: d},
		Timers:       &timers.Timers{Deps: d},
	}
}
