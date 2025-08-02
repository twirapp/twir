package impl_deps

import (
	"github.com/alexedwards/scs/v2"
	"github.com/redis/go-redis/v9"
	buscore "github.com/twirapp/twir/libs/bus-core"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	cfg "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/grpc/discord"
	"github.com/twirapp/twir/libs/grpc/websockets"
	"github.com/twirapp/twir/libs/logger"
	channelsintegrationsspotify "github.com/twirapp/twir/libs/repositories/channels_integrations_spotify"
	eventsmodel "github.com/twirapp/twir/libs/repositories/events/model"
	"github.com/twirapp/twir/libs/types/types/api/modules"
	"gorm.io/gorm"
)

type Grpc struct {
	Websockets websockets.WebsocketClient
	Discord    discord.DiscordClient
}

type Deps struct {
	Logger                            logger.Logger
	SpotifyRepo                       channelsintegrationsspotify.Repository
	Redis                             *redis.Client
	Db                                *gorm.DB
	Grpc                              *Grpc
	SessionManager                    *scs.SessionManager
	Bus                               *buscore.Bus
	TTSSettingsCacher                 *generic_cacher.GenericCacher[modules.TTSSettings]
	Config                            cfg.Config
	ChannelsEventsWithOperationsCache *generic_cacher.GenericCacher[[]eventsmodel.Event]
}
