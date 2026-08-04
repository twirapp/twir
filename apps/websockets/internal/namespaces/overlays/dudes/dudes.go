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
	twirlogger "github.com/twirapp/twir/libs/logger"
	"github.com/twirapp/twir/libs/repositories/channels"
	"github.com/twirapp/twir/libs/repositories/users"
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

func New(
	gorm *gorm.DB,
	logger *slog.Logger,
	redis *redis.Client,
	cfg config.Config,
	twirBus *buscore.Bus,
	channelsRepository channels.Repository,
	usersRepository users.Repository,
) *Dudes {
	m := melody.New()
	m.Config.MaxMessageSize = 1024 * 1024 * 10
	dudes := &Dudes{
		manager: m,
		gorm:    gorm,
		logger:  logger,
		redis:   redis,
		config:  cfg,
		counter: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name:        "websockets_connections_count",
				ConstLabels: prometheus.Labels{"overlay": "dudes"},
			},
		),
		twirBus: twirBus,
	}

	dudes.manager.HandleConnect(
		func(session *melody.Session) {
			err := helpers.CheckChannelByApiKey(session, channelsRepository, usersRepository)
			if err != nil {
				if !errors.Is(err, helpers.ErrUserNotFound) {
					logger.Error("cannot check user by api key", twirlogger.Error(err))
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
		c.logger.Error("cannot process message", twirlogger.Error(err))
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
		c.logger.Error("cannot broadcast message", twirlogger.Error(err))
		return err
	}

	return nil
}
