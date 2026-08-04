package http_webhooks

import (
	"log/slog"

	"github.com/twirapp/kv"
	"github.com/twirapp/twir/apps/api-gql/internal/server"
	"github.com/twirapp/twir/apps/api-gql/internal/services/webhook_notifications"
	buscore "github.com/twirapp/twir/libs/bus-core"
	cfg "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/pubsub"
	"gorm.io/gorm"
)

type Webhooks struct {
	kv                          kv.KV
	db                          *gorm.DB
	logger                      *slog.Logger
	config                      cfg.Config
	pubSub                      *pubsub.PubSub
	twirBus                     *buscore.Bus
	webhookNotificationsService *webhook_notifications.Service
}

func New(srv *server.Server, kvClient kv.KV, db *gorm.DB, logger *slog.Logger, config cfg.Config, twirBus *buscore.Bus, webhookNotificationsService *webhook_notifications.Service) (*Webhooks, error) {
	pb, err := pubsub.NewPubSub(config.RedisUrl)
	if err != nil {
		return nil, err
	}

	p := &Webhooks{
		kv:                          kvClient,
		db:                          db,
		logger:                      logger,
		config:                      config,
		pubSub:                      pb,
		twirBus:                     twirBus,
		webhookNotificationsService: webhookNotificationsService,
	}

	srv.POST("/webhooks/integrations/donatestream/:id", p.donateStreamHandler)
	srv.POST("/webhooks/integrations/donatello", p.donatelloHandler)
	srv.POST("/webhooks/modules/github", p.githubWebhookHandler)

	return p, nil
}

type pbMessage struct {
	TwitchUserId string `json:"twitchUserId"`
	Amount       string `json:"amount"`
	Currency     string `json:"currency"`
	Message      string `json:"message"`
	UserName     string `json:"userName"`
}
