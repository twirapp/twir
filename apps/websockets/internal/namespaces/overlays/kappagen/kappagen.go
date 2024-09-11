package kappagen

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/goccy/go-json"

	"github.com/olahol/melody"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/redis/go-redis/v9"
	"github.com/satont/twir/apps/websockets/internal/namespaces/helpers"
	"github.com/satont/twir/apps/websockets/types"
	config "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/logger"
	"github.com/twirapp/twir/libs/grpc/tokens"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Kappagen struct {
	manager *melody.Melody

	gorm       *gorm.DB
	logger     logger.Logger
	redis      *redis.Client
	config     config.Config
	tokensGrpc tokens.TokensClient
	counter    prometheus.Gauge
}

type Opts struct {
	fx.In

	Gorm       *gorm.DB
	Logger     logger.Logger
	Redis      *redis.Client
	Config     config.Config
	TokensGrpc tokens.TokensClient
}

func New(opts Opts) *Kappagen {
	m := melody.New()
	m.Config.MaxMessageSize = 1024 * 1024 * 10
	kappagen := &Kappagen{
		manager:    m,
		gorm:       opts.Gorm,
		logger:     opts.Logger,
		redis:      opts.Redis,
		config:     opts.Config,
		tokensGrpc: opts.TokensGrpc,
		counter: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name:        "websockets_connections_count",
				ConstLabels: prometheus.Labels{"overlay": "kappagen"},
			},
		),
	}

	kappagen.manager.HandleConnect(
		func(session *melody.Session) {
			err := helpers.CheckUserByApiKey(opts.Gorm, session)
			if err != nil {
				if !errors.Is(err, helpers.ErrUserNotFound) {
					opts.Logger.Error("cannot check user by api key", slog.Any("err", err))
				}
				return
			}

			kappagen.counter.Inc()
			session.Write([]byte(`{"eventName":"connected to kappagen namespace"}`))
		},
	)

	kappagen.manager.HandleMessage(
		func(session *melody.Session, msg []byte) {
			kappagen.handleMessage(session, msg)
		},
	)

	kappagen.manager.HandleDisconnect(
		func(session *melody.Session) {
			kappagen.counter.Dec()
		},
	)

	http.HandleFunc("/overlays/kappagen", kappagen.HandleRequest)

	return kappagen
}

func (c *Kappagen) HandleRequest(w http.ResponseWriter, r *http.Request) {
	_ = c.manager.HandleRequest(w, r)
}

func (c *Kappagen) SendEvent(channelId, eventName string, data any) error {
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
