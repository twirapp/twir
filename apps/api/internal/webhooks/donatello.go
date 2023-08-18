package webhooks

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/satont/twir/libs/logger"
	"log/slog"
	"net/http"
	"time"

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
	Logger     logger.Logger
}

type Donatello struct {
	redis      *redis.Client
	db         *gorm.DB
	eventsGrpc events.EventsClient
	l          logger.Logger
}

func NewDonatello(opts DonatelloOpts) handlers.IHandler {
	return &Donatello{
		redis:      opts.Redis,
		db:         opts.Db,
		eventsGrpc: opts.EventsGrpc,
		l:          opts.Logger,
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

			if err := c.db.Create(
				&model.ChannelsEventsListItem{
					ID:        uuid.New().String(),
					ChannelID: integration.ChannelID,
					UserID:    "",
					Type:      model.ChannelEventListItemTypeDonation,
					Data: &model.ChannelsEventsListItemData{
						DonationAmount:   body.Amount,
						DonationCurrency: body.Currency,
						DonationMessage:  body.Message,
						DonationUsername: body.ClientName,
					},
					CreatedAt: time.Now(),
				},
			).Error; err != nil {
				c.l.Error("cannot create event", slog.Any("err", err))
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
