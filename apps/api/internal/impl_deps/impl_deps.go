package impl_deps

import (
	"log/slog"

	"github.com/alexedwards/scs/v2"
	"github.com/redis/go-redis/v9"
	buscore "github.com/twirapp/twir/libs/bus-core"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	cfg "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/grpc/discord"
	"github.com/twirapp/twir/libs/grpc/websockets"
	channelsintegrationsspotify "github.com/twirapp/twir/libs/repositories/channels_integrations_spotify"
	commandwithgroupandresponsesmodel "github.com/twirapp/twir/libs/repositories/commands_with_groups_and_responses/model"
	eventsmodel "github.com/twirapp/twir/libs/repositories/events/model"
	"gorm.io/gorm"
)

type Grpc struct {
	Websockets websockets.WebsocketClient
	Discord    discord.DiscordClient
}

type Deps struct {
	Logger                            *slog.Logger
	SpotifyRepo                       channelsintegrationsspotify.Repository
	Redis                             *redis.Client
	Db                                *gorm.DB
	Grpc                              *Grpc
	SessionManager                    *scs.SessionManager
	Bus                               *buscore.Bus
	Config                            cfg.Config
	ChannelsEventsWithOperationsCache *generic_cacher.GenericCacher[[]eventsmodel.Event]
	ChannelsCommandsCache             *generic_cacher.GenericCacher[[]commandwithgroupandresponsesmodel.CommandWithGroupAndResponses]
}
