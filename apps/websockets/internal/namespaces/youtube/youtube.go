package youtube

import (
	"net/http"

	"github.com/olahol/melody"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/redis/go-redis/v9"
	"github.com/satont/twir/apps/websockets/internal/namespaces/helpers"
	"github.com/satont/twir/libs/grpc/bots"
	"github.com/satont/twir/libs/logger"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type YouTube struct {
	manager *melody.Melody

	gorm     *gorm.DB
	logger   logger.Logger
	redis    *redis.Client
	botsGrpc bots.BotsClient

	counter prometheus.Gauge
}

type Opts struct {
	fx.In

	Gorm     *gorm.DB
	Logger   logger.Logger
	Redis    *redis.Client
	BotsGrpc bots.BotsClient
}

func NewYouTube(opts Opts) *YouTube {
	m := melody.New()
	m.Config.MaxMessageSize = 1024 * 1024 * 10
	youTube := &YouTube{
		manager:  m,
		gorm:     opts.Gorm,
		logger:   opts.Logger,
		redis:    opts.Redis,
		botsGrpc: opts.BotsGrpc,
		counter: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name:        "websockets_connections_count",
				ConstLabels: prometheus.Labels{"overlay": "youtube"},
			},
		),
	}

	youTube.manager.HandleConnect(
		func(session *melody.Session) {
			err := helpers.CheckUserByApiKey(opts.Gorm, session)
			if err != nil {
				opts.Logger.Error(err.Error())
			} else {
				youTube.handleConnect(session)
			}
		},
	)

	youTube.manager.HandleMessage(
		func(session *melody.Session, msg []byte) {
			youTube.handleMessage(session, msg)
		},
	)

	http.HandleFunc("/youtube", youTube.HandleRequest)

	return youTube
}
