package alerts

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/olahol/melody"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/redis/go-redis/v9"
	"github.com/satont/twir/apps/websockets/internal/namespaces/helpers"
	"github.com/satont/twir/apps/websockets/types"
	"github.com/satont/twir/libs/logger"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Alerts struct {
	manager *melody.Melody

	gorm   *gorm.DB
	logger logger.Logger
	redis  *redis.Client

	counter prometheus.Gauge
}

type Opts struct {
	fx.In

	Gorm   *gorm.DB
	Logger logger.Logger
	Redis  *redis.Client
}

func NewAlerts(opts Opts) *Alerts {
	m := melody.New()
	m.Config.MaxMessageSize = 1024 * 1024 * 10
	alerts := &Alerts{
		manager: m,
		gorm:    opts.Gorm,
		logger:  opts.Logger,
		redis:   opts.Redis,
		counter: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name:        "websockets_connections_count",
				ConstLabels: prometheus.Labels{"overlay": "alerts"},
			},
		),
	}

	alerts.manager.HandleConnect(
		func(session *melody.Session) {
			err := helpers.CheckUserByApiKey(opts.Gorm, session)
			if err != nil {
				if !errors.Is(err, helpers.ErrUserNotFound) {
					opts.Logger.Error("cannot check user by api key", slog.Any("err", err))
				}
				return
			}
			alerts.counter.Inc()
			session.Write([]byte(`{"eventName":"connected to alerts namespace"}`))
		},
	)

	alerts.manager.HandleDisconnect(
		func(session *melody.Session) {
			alerts.counter.Dec()
		},
	)

	http.HandleFunc("/overlays/alerts", alerts.HandleRequest)

	return alerts
}

func (c *Alerts) HandleRequest(w http.ResponseWriter, r *http.Request) {
	_ = c.manager.HandleRequest(w, r)
}

func (c *Alerts) SendEvent(channelId, eventName string, data any) error {
	message := &types.WebSocketMessage{
		EventName: eventName,
		Data:      data,
		CreatedAt: time.Now().UTC().String(),
	}

	bytes, err := json.Marshal(message)
	if err != nil {
		c.logger.Error("cannot process message", slog.Any("err", err))
		return err
	}

	err = c.manager.BroadcastFilter(
		bytes, func(session *melody.Session) bool {
			socketUserId, ok := session.Get("userId")
			return ok && socketUserId.(string) == channelId
		},
	)

	if err != nil {
		c.logger.Error("cannot broadcast message", slog.Any("err", err))
		return err
	}

	return nil
}
