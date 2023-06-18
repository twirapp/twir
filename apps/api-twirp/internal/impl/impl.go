package impl

import (
	"github.com/redis/go-redis/v9"
	"github.com/satont/tsuwari/apps/api-twirp/internal/impl/bot"
	"github.com/satont/tsuwari/apps/api-twirp/internal/impl/commands"
	"github.com/satont/tsuwari/apps/api-twirp/internal/impl/community"
	"github.com/satont/tsuwari/apps/api-twirp/internal/impl/deps"
	"github.com/satont/tsuwari/apps/api-twirp/internal/impl/events"
	"github.com/satont/tsuwari/apps/api-twirp/internal/impl/greetings"
	"github.com/satont/tsuwari/apps/api-twirp/internal/impl/integrations"
	"github.com/satont/tsuwari/apps/api-twirp/internal/impl/keywords"
	"github.com/satont/tsuwari/apps/api-twirp/internal/impl/modules"
	"gorm.io/gorm"
)

type Api struct {
	*integrations.Integrations
	*keywords.Keywords
	*modules.Modules
	*bot.Bot
	*commands.Commands
	*community.Community
	*events.Events
	*greetings.Greetings
}

type Opts struct {
	Redis *redis.Client
	DB    *gorm.DB
}

func NewApi(opts Opts) *Api {
	d := &deps.Deps{
		Redis: opts.Redis,
		Db:    opts.DB,
	}

	return &Api{
		Integrations: &integrations.Integrations{Deps: d},
		Keywords:     &keywords.Keywords{Deps: d},
		Modules:      &modules.Modules{Deps: d},
		Bot:          &bot.Bot{Deps: d},
		Commands:     &commands.Commands{Deps: d},
		Community:    &community.Community{Deps: d},
		Events:       &events.Events{Deps: d},
		Greetings:    &greetings.Greetings{Deps: d},
	}
}
