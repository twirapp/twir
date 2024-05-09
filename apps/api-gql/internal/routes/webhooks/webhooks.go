package webhooks

import (
	"github.com/redis/go-redis/v9"
	cfg "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/logger"
	"github.com/satont/twir/libs/pubsub"
	"github.com/twirapp/twir/apps/api-gql/internal/httpserver"
	"github.com/twirapp/twir/libs/grpc/events"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Opts struct {
	fx.In

	Server     *httpserver.Server
	Redis      *redis.Client
	Db         *gorm.DB
	EventsGrpc events.EventsClient
	Logger     logger.Logger
	Config     cfg.Config
}

type Webhooks struct {
	redis      *redis.Client
	db         *gorm.DB
	eventsGrpc events.EventsClient
	logger     logger.Logger
	config     cfg.Config
	pubSub     *pubsub.PubSub
}

func New(opts Opts) (*Webhooks, error) {
	pb, err := pubsub.NewPubSub(opts.Config.RedisUrl)
	if err != nil {
		return nil, err
	}

	p := &Webhooks{
		redis:      opts.Redis,
		db:         opts.Db,
		eventsGrpc: opts.EventsGrpc,
		logger:     opts.Logger,
		config:     opts.Config,
		pubSub:     pb,
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
