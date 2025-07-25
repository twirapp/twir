package http_webhooks

import (
	"github.com/redis/go-redis/v9"
	cfg "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/logger"
	"github.com/twirapp/twir/libs/pubsub"
	"github.com/twirapp/twir/apps/api-gql/internal/server"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Opts struct {
	fx.In

	Server  *server.Server
	Redis   *redis.Client
	Db      *gorm.DB
	Logger  logger.Logger
	Config  cfg.Config
	TwirBus *buscore.Bus
}

type Webhooks struct {
	redis   *redis.Client
	db      *gorm.DB
	logger  logger.Logger
	config  cfg.Config
	pubSub  *pubsub.PubSub
	twirBus *buscore.Bus
}

func New(opts Opts) (*Webhooks, error) {
	pb, err := pubsub.NewPubSub(opts.Config.RedisUrl)
	if err != nil {
		return nil, err
	}

	p := &Webhooks{
		redis:   opts.Redis,
		db:      opts.Db,
		logger:  opts.Logger,
		config:  opts.Config,
		pubSub:  pb,
		twirBus: opts.TwirBus,
	}

	opts.Server.POST("/webhooks/integrations/donatestream/:id", p.donateStreamHandler)
	opts.Server.POST("/webhooks/integrations/donatello", p.donatelloHandler)

	return p, nil
}

type pbMessage struct {
	TwitchUserId string `json:"twitchUserId"`
	Amount       string `json:"amount"`
	Currency     string `json:"currency"`
	Message      string `json:"message"`
	UserName     string `json:"userName"`
}
