package events

import (
	"github.com/redis/go-redis/v9"
	"github.com/twirapp/twir/apps/events/internal/hydrator"
	bus_core "github.com/twirapp/twir/libs/bus-core"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	config "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/grpc/websockets"
	"github.com/twirapp/twir/libs/repositories/greetings"
	"github.com/twirapp/twir/libs/repositories/overlays_tts"
	"github.com/twirapp/twir/libs/repositories/variables"
	"github.com/twirapp/twir/libs/types/types/api/modules"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Opts struct {
	fx.In

	Gorm                *gorm.DB
	Redis               *redis.Client
	Cfg                 config.Config
	WebsocketsGrpc      websockets.WebsocketClient
	Hydrator            *hydrator.Hydrator
	Bus                 *bus_core.Bus
	GreetingsRepository greetings.Repository
	VariablesRepository variables.Repository
	TTSRepository       overlays_tts.Repository
	TTSCache            *generic_cacher.GenericCacher[modules.TTSSettings]
}

func New(opts Opts) *Activity {
	return &Activity{
		db:                  opts.Gorm,
		redis:               opts.Redis,
		cfg:                 opts.Cfg,
		websocketsGrpc:      opts.WebsocketsGrpc,
		bus:                 opts.Bus,
		hydrator:            opts.Hydrator,
		greetingsRepository: opts.GreetingsRepository,
		variablesRepository: opts.VariablesRepository,
		ttsRepository:       opts.TTSRepository,
		ttsCache:            opts.TTSCache,
	}
}

type Activity struct {
	db                  *gorm.DB
	redis               *redis.Client
	cfg                 config.Config
	websocketsGrpc      websockets.WebsocketClient
	hydrator            *hydrator.Hydrator
	bus                 *bus_core.Bus
	greetingsRepository greetings.Repository
	variablesRepository variables.Repository
	ttsRepository       overlays_tts.Repository
	ttsCache            *generic_cacher.GenericCacher[modules.TTSSettings]
}
