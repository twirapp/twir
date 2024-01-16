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
	"github.com/satont/twir/libs/grpc/events"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type DonateStreamOpts struct {
	fx.In

	Redis      *redis.Client
	Db         *gorm.DB
	EventsGrpc events.EventsClient
	Logger     logger.Logger
	Config     cfg.Config
}

type DonateStream struct {
	redis      *redis.Client
	db         *gorm.DB
	eventsGrpc events.EventsClient
	l          logger.Logger
	pb         *pubsub.PubSub
}

func NewDonateStream(opts DonateStreamOpts) handlers.IHandler {
	pb, err := pubsub.NewPubSub(opts.Config.RedisUrl)
	if err != nil {
		panic(err)
	}

	return &DonateStream{
		redis:      opts.Redis,
		db:         opts.Db,
		eventsGrpc: opts.EventsGrpc,
		l:          opts.Logger,
		pb:         pb,
	}
}

func (c *DonateStream) Pattern() string {
	return "/webhooks/integrations/donatestream/"
}

type donateStreamIncomingData struct {
	Type     string `json:"type,omitempty"`
	Uid      string `json:"uid"`
	Message  string `json:"message"`
	Sum      string `json:"sum"`
	Nickname string `json:"nickname"`
}

func (c *DonateStream) Handler() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}

			reqUrl := r.URL.Path
			id := reqUrl[len(c.Pattern()):]

			integration := model.ChannelsIntegrations{}

			if err := c.db.
				WithContext(r.Context()).
				Where("id = ?", id).First(&integration).Error; err != nil {
				http.Error(w, "Integration not found", http.StatusNotFound)
				return
			}

			body := &donateStreamIncomingData{}
			if err := json.NewDecoder(r.Body).Decode(body); err != nil {
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}

			if body.Type == "confirm" {
				value, err := c.redis.Get(
					r.Context(),
					"donate_stream_confirmation"+integration.IntegrationID,
				).Result()
				if err != nil {
					c.l.Error("cannot get confirmation from redis", slog.Any("err", err))
					http.Error(w, "Internal error", http.StatusInternalServerError)
					return
				}

				w.Write([]byte(value))
				return
			}

			integrationsMessage := &pbMessage{
				TwitchUserId: integration.ChannelID,
				Amount:       body.Sum,
				Currency:     "RUB",
				Message:      body.Message,
				UserName:     body.Nickname,
			}
			integrationsNameBytes, err := json.Marshal(integrationsMessage)
			if err != nil {
				c.l.Error("cannot marshal message", slog.Any("err", err))
			}

			c.pb.Publish("donations:new", integrationsNameBytes)

			if err != nil {
				http.Error(w, "Internal error", http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
		},
	)
}
