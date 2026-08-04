package events

import (
	"log/slog"

	"github.com/redis/go-redis/v9"
	"github.com/twirapp/twir/apps/events/internal/hydrator"
	bus_core "github.com/twirapp/twir/libs/bus-core"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	config "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/grpc/websockets"
	commandsrepository "github.com/twirapp/twir/libs/repositories/commands"
	"github.com/twirapp/twir/libs/repositories/greetings"
	"github.com/twirapp/twir/libs/repositories/overlays_tts"
	usersrepository "github.com/twirapp/twir/libs/repositories/users"
	"github.com/twirapp/twir/libs/repositories/variables"
	channelservice "github.com/twirapp/twir/libs/services/channels"
	"github.com/twirapp/twir/libs/types/types/api/modules"
	"gorm.io/gorm"
)

func New(
	db *gorm.DB,
	redisClient *redis.Client,
	cfg config.Config,
	websocketsGrpc websockets.WebsocketClient,
	hydrator *hydrator.Hydrator,
	bus *bus_core.Bus,
	channelService *channelservice.ChannelService,
	commandsRepository commandsrepository.Repository,
	greetingsRepository greetings.Repository,
	variablesRepository variables.Repository,
	ttsRepository overlays_tts.Repository,
	ttsCache *generic_cacher.GenericCacher[modules.TTSSettings],
	logger *slog.Logger,
	usersRepository usersrepository.Repository,
) *Activity {
	return &Activity{
		db:                  db,
		redis:               redisClient,
		cfg:                 cfg,
		websocketsGrpc:      websocketsGrpc,
		bus:                 bus,
		hydrator:            hydrator,
		channelService:      channelService,
		commandsRepo:        commandsRepository,
		greetingsRepository: greetingsRepository,
		variablesRepository: variablesRepository,
		ttsRepository:       ttsRepository,
		ttsCache:            ttsCache,
		logger:              logger,
		usersRepo:           usersRepository,
	}
}

type Activity struct {
	db                  *gorm.DB
	redis               *redis.Client
	cfg                 config.Config
	websocketsGrpc      websockets.WebsocketClient
	hydrator            *hydrator.Hydrator
	bus                 *bus_core.Bus
	channelService      *channelservice.ChannelService
	commandsRepo        commandsrepository.Repository
	greetingsRepository greetings.Repository
	variablesRepository variables.Repository
	ttsRepository       overlays_tts.Repository
	ttsCache            *generic_cacher.GenericCacher[modules.TTSSettings]
	logger              *slog.Logger
	usersRepo           usersrepository.Repository
	newTwitchBotClient  twitchBotClientFactory
}

type channelRuntimeInfo struct {
	ChannelID         string
	BroadcasterUserID string
	TwitchPlatformID  string
	BotID             string
	IsBotMod          bool
	IsTwitchBanned    bool
}
