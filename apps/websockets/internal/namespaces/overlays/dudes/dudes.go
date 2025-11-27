package dudes

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
	"github.com/twirapp/twir/apps/websockets/internal/namespaces/helpers"
	"github.com/twirapp/twir/apps/websockets/types"
	buscore "github.com/twirapp/twir/libs/bus-core"
	config "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/logger"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Dudes struct {
	manager *melody.Melody

	gorm    *gorm.DB
	logger  *slog.Logger
	redis   *redis.Client
	config  config.Config
	counter prometheus.Gauge
	twirBus *buscore.Bus
}

type Opts struct {
	fx.In

	Gorm    *gorm.DB
	Logger  *slog.Logger
	Redis   *redis.Client
	Config  config.Config
	TwirBus *buscore.Bus
}

func New(opts Opts) *Dudes {
	m := melody.New()
	m.Config.MaxMessageSize = 1024 * 1024 * 10
	dudes := &Dudes{
		manager: m,
		gorm:    opts.Gorm,
		logger:  opts.Logger,
		redis:   opts.Redis,
		config:  opts.Config,
		counter: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name:        "websockets_connections_count",
				ConstLabels: prometheus.Labels{"overlay": "dudes"},
			},
		),
		twirBus: opts.TwirBus,
	}

	dudes.manager.HandleConnect(
		func(session *melody.Session) {
			err := helpers.CheckUserByApiKey(opts.Gorm, session)
			if err != nil {
				if !errors.Is(err, helpers.ErrUserNotFound) {
					opts.Logger.Error("cannot check user by api key", logger.Error(err))
				}
				return
			}

			session.Set("id", session.Request.URL.Query().Get("id"))

			dudes.counter.Inc()
			session.Write([]byte(`{"eventName":"connected to dudes namespace"}`))
		},
	)

	dudes.manager.HandleMessage(
		func(session *melody.Session, msg []byte) {
			dudes.handleMessage(session, msg)
		},
	)

	dudes.manager.HandleDisconnect(
		func(session *melody.Session) {
			dudes.counter.Dec()
		},
	)

	http.HandleFunc("/overlays/dudes", dudes.HandleRequest)

	return dudes
}

func (c *Dudes) HandleRequest(w http.ResponseWriter, r *http.Request) {
	_ = c.manager.HandleRequest(w, r)
}

func (c *Dudes) SendEvent(channelId, eventName string, data any) error {
	message := &types.WebSocketMessage{
		EventName: eventName,
		Data:      data,
		CreatedAt: time.Now().UTC().String(),
	}

	bytes, err := json.Marshal(message)
	if err != nil {
		c.logger.Error("cannot process message", logger.Error(err))
		return err
	}

	err = c.manager.BroadcastFilter(
		bytes,
		func(session *melody.Session) bool {
			socketUserId, ok := session.Get("userId")
			return ok && socketUserId.(string) == channelId
		},
	)

	if err != nil {
		c.logger.Error("cannot broadcast message", logger.Error(err))
		return err
	}

	return nil
}
