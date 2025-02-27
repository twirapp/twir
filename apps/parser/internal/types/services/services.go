package services

import (
	"github.com/go-redsync/redsync/v4"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"github.com/satont/twir/apps/parser/internal/task-queue"
	cfg "github.com/satont/twir/libs/config"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/types/types/api/modules"
	buscore "github.com/twirapp/twir/libs/bus-core"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	"github.com/twirapp/twir/libs/grpc/dota"
	"github.com/twirapp/twir/libs/grpc/events"
	"github.com/twirapp/twir/libs/grpc/tokens"
	"github.com/twirapp/twir/libs/grpc/websockets"
	"github.com/twirapp/twir/libs/grpc/ytsr"
	"github.com/twirapp/twir/libs/integrations/seventv"
	channelscommandsprefixrepository "github.com/twirapp/twir/libs/repositories/channels_commands_prefix"
	channelscommandsprefixmodel "github.com/twirapp/twir/libs/repositories/channels_commands_prefix/model"
	channelsintegrationsspotify "github.com/twirapp/twir/libs/repositories/channels_integrations_spotify"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Grpc struct {
	WebSockets websockets.WebsocketClient
	Dota       dota.DotaClient
	Tokens     tokens.TokensClient
	Events     events.EventsClient
	Ytsr       ytsr.YtsrClient
}

type Services struct {
	Config                   *cfg.Config
	Logger                   *zap.Logger
	Gorm                     *gorm.DB
	Sqlx                     *sqlx.DB
	Redis                    *redis.Client
	GrpcClients              *Grpc
	TaskDistributor          task_queue.TaskDistributor
	Bus                      *buscore.Bus
	CommandsCache            *generic_cacher.GenericCacher[[]model.ChannelsCommands]
	CommandsPrefixCache      *generic_cacher.GenericCacher[channelscommandsprefixmodel.ChannelsCommandsPrefix]
	SevenTvCache             *generic_cacher.GenericCacher[*seventv.ProfileResponse]
	SevenTvCacheBySevenTvID  *generic_cacher.GenericCacher[*seventv.ProfileResponse]
	RedSync                  *redsync.Redsync
	CommandsLock             *redsync.Mutex
	CommandsPrefixRepository channelscommandsprefixrepository.Repository
	TTSCache                 *generic_cacher.GenericCacher[modules.TTSSettings]
	SpotifyRepo              channelsintegrationsspotify.Repository
}
