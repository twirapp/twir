package webhooks

import (
	"encoding/json"
	"log/slog"
	"net/http"

	cfg "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/logger"
	"github.com/satont/twir/libs/pubsub"

	"github.com/redis/go-redis/v9"
	"github.com/satont/twir/apps/api/internal/handlers"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/grpc/events"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type DonatelloOpts struct {
	fx.In

	Redis      *redis.Client
	Db         *gorm.DB
	EventsGrpc events.EventsClient
	Logger     logger.Logger
	Config     cfg.Config
}

type Donatello struct {
	redis      *redis.Client
	db         *gorm.DB
	eventsGrpc events.EventsClient
	l          logger.Logger
	pb         *pubsub.PubSub
}

func NewDonatello(opts DonatelloOpts) handlers.IHandler {
	pb, err := pubsub.NewPubSub(opts.Config.RedisUrl)
	if err != nil {
		panic(err)
	}

	return &Donatello{
		redis:      opts.Redis,
		db:         opts.Db,
		eventsGrpc: opts.EventsGrpc,
		l:          opts.Logger,
		pb:         pb,
	}
}

func (c *Donatello) Pattern() string {
	return "/webhooks/integrations/donatello"
}

type donatelloBody struct {
	PubId       string `json:"pubId"`
	ClientName  string `json:"clientName"`
	Message     string `json:"message"`
	Amount      string `json:"amount"`
	Currency    string `json:"currency"`
	Source      string `json:"source"`
	Goal        string `json:"goal"`
	IsPublished bool   `json:"isPublished"`
	CreatedAt   string `json:"createdAt"`
}

func (c *Donatello) Handler() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}

			apiKey := r.Header.Get("X-Key")
			if apiKey == "" {
				http.Error(w, "X-Key header is required", http.StatusBadRequest)
				return
			}

			integration := &model.ChannelsIntegrations{}
			if err := c.db.
				WithContext(r.Context()).
				Where(`"id" = ?`, apiKey).
				First(integration).
				Error; err != nil {
				http.Error(w, "Integration not found", http.StatusNotFound)
				return
			}

			body := &donatelloBody{}
			if err := json.NewDecoder(r.Body).Decode(body); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			integrationsMessage := &pbMessage{
				TwitchUserId: integration.ChannelID,
				Amount:       body.Amount,
				Currency:     body.Currency,
				Message:      body.Message,
				UserName:     body.ClientName,
			}
			integrationsNameBytes, err := json.Marshal(integrationsMessage)
			if err != nil {
				c.l.Error("cannot marshal message", slog.Any("err", err))
			}

			c.pb.Publish("donations:new", integrationsNameBytes)

			w.WriteHeader(http.StatusOK)
		},
	)
}
