package http_webhooks

import (
	"log/slog"

	"github.com/twirapp/kv"
	"github.com/twirapp/twir/apps/api-gql/internal/server"
	"github.com/twirapp/twir/apps/api-gql/internal/services/webhook_notifications"
	buscore "github.com/twirapp/twir/libs/bus-core"
	cfg "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/pubsub"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Opts struct {
	fx.In

	Server                      *server.Server
	KV                          kv.KV
	Db                          *gorm.DB
	Logger                      *slog.Logger
	Config                      cfg.Config
	TwirBus                     *buscore.Bus
	WebhookNotificationsService *webhook_notifications.Service
}

type Webhooks struct {
	kv                          kv.KV
	db                          *gorm.DB
	logger                      *slog.Logger
	config                      cfg.Config
	pubSub                      *pubsub.PubSub
	twirBus                     *buscore.Bus
	webhookNotificationsService *webhook_notifications.Service
}

func New(opts Opts) (*Webhooks, error) {
	pb, err := pubsub.NewPubSub(opts.Config.RedisUrl)
	if err != nil {
		return nil, err
	}

	p := &Webhooks{
		kv:                          opts.KV,
		db:                          opts.Db,
		logger:                      opts.Logger,
		config:                      opts.Config,
		pubSub:                      pb,
		twirBus:                     opts.TwirBus,
		webhookNotificationsService: opts.WebhookNotificationsService,
	}

	opts.Server.POST("/webhooks/integrations/donatestream/:id", p.donateStreamHandler)
	opts.Server.POST("/webhooks/integrations/donatello", p.donatelloHandler)
	opts.Server.POST("/webhooks/modules/github", p.githubWebhookHandler)

	return p, nil
}

type pbMessage struct {
	TwitchUserId string `json:"twitchUserId"`
	Amount       string `json:"amount"`
	Currency     string `json:"currency"`
	Message      string `json:"message"`
	UserName     string `json:"userName"`
}
