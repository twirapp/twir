package webhooks

import (
	"encoding/json"
	"net/http"

	"github.com/redis/go-redis/v9"
	"github.com/satont/twir/apps/api/internal/handlers"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/grpc/generated/events"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type DonatelloOpts struct {
	fx.In

	Redis      *redis.Client
	Db         *gorm.DB
	EventsGrpc events.EventsClient
}

type Donatello struct {
	redis      *redis.Client
	db         *gorm.DB
	eventsGrpc events.EventsClient
}

func NewDonatello(opts DonatelloOpts) handlers.IHandler {
	return &Donatello{
		redis:      opts.Redis,
		db:         opts.Db,
		eventsGrpc: opts.EventsGrpc,
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

			_, err := c.eventsGrpc.Donate(
				r.Context(),
				&events.DonateMessage{
					BaseInfo: &events.BaseInfo{ChannelId: integration.ChannelID},
					UserName: body.ClientName,
					Amount:   body.Amount,
					Currency: body.Currency,
					Message:  body.Message,
				},
			)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
		},
	)
}
