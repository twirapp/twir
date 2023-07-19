package webhooks

import (
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"github.com/satont/twir/apps/api/internal/handlers"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/grpc/generated/events"
	"go.uber.org/fx"
	"gorm.io/gorm"
	"net/http"
)

type DonateStreamOpts struct {
	fx.In

	Redis      *redis.Client
	Db         *gorm.DB
	EventsGrpc events.EventsClient
}

type DonateStream struct {
	redis      *redis.Client
	db         *gorm.DB
	eventsGrpc events.EventsClient
}

func NewDonateStream(opts DonateStreamOpts) handlers.IHandler {
	return &DonateStream{
		redis:      opts.Redis,
		db:         opts.Db,
		eventsGrpc: opts.EventsGrpc,
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
				value, err := c.redis.Get(r.Context(), "donate_stream_confirmation"+integration.ID).Result()
				if err != nil {
					http.Error(w, "Internal error", http.StatusInternalServerError)
					return
				}

				w.Write([]byte(value))
				return
			}

			_, err := c.eventsGrpc.Donate(
				r.Context(),
				&events.DonateMessage{
					BaseInfo: &events.BaseInfo{ChannelId: integration.ChannelID},
					UserName: body.Nickname,
					Amount:   body.Sum,
					Currency: "RUB",
					Message:  body.Message,
				},
			)

			if err != nil {
				http.Error(w, "Internal error", http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
		},
	)
}
