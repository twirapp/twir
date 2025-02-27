package impl_deps

import (
	"github.com/alexedwards/scs/v2"
	"github.com/redis/go-redis/v9"
	cfg "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/logger"
	"github.com/satont/twir/libs/types/types/api/modules"
	buscore "github.com/twirapp/twir/libs/bus-core"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	"github.com/twirapp/twir/libs/grpc/discord"
	"github.com/twirapp/twir/libs/grpc/integrations"
	"github.com/twirapp/twir/libs/grpc/parser"
	"github.com/twirapp/twir/libs/grpc/tokens"
	"github.com/twirapp/twir/libs/grpc/websockets"
	channelsintegrationsspotify "github.com/twirapp/twir/libs/repositories/channels_integrations_spotify"
	"gorm.io/gorm"
)

type Grpc struct {
	Tokens       tokens.TokensClient
	Integrations integrations.IntegrationsClient
	Parser       parser.ParserClient
	Websockets   websockets.WebsocketClient
	Discord      discord.DiscordClient
}

type Deps struct {
	Logger            logger.Logger
	SpotifyRepo       channelsintegrationsspotify.Repository
	Redis             *redis.Client
	Db                *gorm.DB
	Grpc              *Grpc
	SessionManager    *scs.SessionManager
	Bus               *buscore.Bus
	TTSSettingsCacher *generic_cacher.GenericCacher[modules.TTSSettings]
	Config            cfg.Config
}
