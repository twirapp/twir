package impl_deps

import (
	"github.com/alexedwards/scs/v2"
	"github.com/redis/go-redis/v9"
	cfg "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/grpc/bots"
	"github.com/satont/twir/libs/grpc/discord"
	"github.com/satont/twir/libs/grpc/eventsub"
	"github.com/satont/twir/libs/grpc/integrations"
	"github.com/satont/twir/libs/grpc/parser"
	"github.com/satont/twir/libs/grpc/scheduler"
	"github.com/satont/twir/libs/grpc/timers"
	"github.com/satont/twir/libs/grpc/tokens"
	"github.com/satont/twir/libs/grpc/websockets"
	"github.com/satont/twir/libs/logger"
	"gorm.io/gorm"
)

type Grpc struct {
	Tokens       tokens.TokensClient
	Integrations integrations.IntegrationsClient
	Bots         bots.BotsClient
	Parser       parser.ParserClient
	Websockets   websockets.WebsocketClient
	Scheduler    scheduler.SchedulerClient
	Timers       timers.TimersClient
	EventSub     eventsub.EventSubClient
	Discord      discord.DiscordClient
}

type Deps struct {
	Config         cfg.Config
	Redis          *redis.Client
	Db             *gorm.DB
	Grpc           *Grpc
	SessionManager *scs.SessionManager
	Logger         logger.Logger
}
