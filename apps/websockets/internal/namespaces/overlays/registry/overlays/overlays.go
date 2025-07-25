package overlays

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/olahol/melody"
	"github.com/redis/go-redis/v9"
	"github.com/twirapp/twir/apps/websockets/internal/namespaces/helpers"
	"github.com/twirapp/twir/apps/websockets/types"
	"github.com/twirapp/twir/libs/logger"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Registry struct {
	manager *melody.Melody

	gorm   *gorm.DB
	logger logger.Logger
	redis  *redis.Client
	bus    *buscore.Bus
}

type Opts struct {
	fx.In

	Gorm   *gorm.DB
	Logger logger.Logger
	Redis  *redis.Client
	Bus    *buscore.Bus
}

func New(opts Opts) *Registry {
	m := melody.New()
	m.Config.MaxMessageSize = 1024 * 1024 * 10
	overlaysRegistry := &Registry{
		manager: m,
		gorm:    opts.Gorm,
		logger:  opts.Logger,
		redis:   opts.Redis,
		bus:     opts.Bus,
	}

	overlaysRegistry.manager.HandleConnect(
		func(session *melody.Session) {
			err := helpers.CheckUserByApiKey(opts.Gorm, session)
			if err != nil {
				if !errors.Is(err, helpers.ErrUserNotFound) {
					opts.Logger.Error("cannot check user by api key", slog.Any("err", err))
				}
				return
			}
			session.Write([]byte(`{"eventName":"connected to overlays namespace"}`))
		},
	)

	overlaysRegistry.manager.HandleMessage(
		func(session *melody.Session, msg []byte) {
			overlaysRegistry.handleMessage(session, msg)
		},
	)

	http.HandleFunc("/overlays/registry/overlays", overlaysRegistry.HandleRequest)

	return overlaysRegistry
}
func (c *Registry) HandleRequest(w http.ResponseWriter, r *http.Request) {
	_ = c.manager.HandleRequest(w, r)
}

func (c *Registry) SendEvent(channelId, eventName string, data any) error {
	message := &types.WebSocketMessage{
		EventName: eventName,
		Data:      data,
		CreatedAt: time.Now().UTC().String(),
	}

	bytes, err := json.Marshal(message)
	if err != nil {
		c.logger.Error(err.Error())
		return err
	}

	err = c.manager.BroadcastFilter(
		bytes, func(session *melody.Session) bool {
			socketUserId, ok := session.Get("userId")
			return ok && socketUserId.(string) == channelId
		},
	)

	if err != nil {
		c.logger.Error(err.Error())
		return err
	}

	return nil
}
