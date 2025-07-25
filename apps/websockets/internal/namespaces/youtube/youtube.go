package youtube

import (
	"net/http"

	"github.com/olahol/melody"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/redis/go-redis/v9"
	"github.com/twirapp/twir/apps/websockets/internal/namespaces/helpers"
	"github.com/twirapp/twir/libs/logger"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type YouTube struct {
	manager *melody.Melody

	gorm   *gorm.DB
	logger logger.Logger
	redis  *redis.Client

	counter prometheus.Gauge
	bus     *buscore.Bus
}

type Opts struct {
	fx.In

	Gorm   *gorm.DB
	Logger logger.Logger
	Redis  *redis.Client
	Bus    *buscore.Bus
}

func NewYouTube(opts Opts) *YouTube {
	m := melody.New()
	m.Config.MaxMessageSize = 1024 * 1024 * 10
	youTube := &YouTube{
		manager: m,
		gorm:    opts.Gorm,
		logger:  opts.Logger,
		redis:   opts.Redis,
		bus:     opts.Bus,
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
