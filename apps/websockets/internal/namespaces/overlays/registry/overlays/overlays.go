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
	buscore "github.com/twirapp/twir/libs/bus-core"
	twirlogger "github.com/twirapp/twir/libs/logger"
	"github.com/twirapp/twir/libs/repositories/channels"
	"github.com/twirapp/twir/libs/repositories/channels_overlays"
	"github.com/twirapp/twir/libs/repositories/users"
	"github.com/twirapp/twir/libs/wsrouter"
	"golang.org/x/time/rate"
	"gorm.io/gorm"
)

// maxMessagesPerSecond caps one connection to 60 events ("frames") per second,
// so overlays cannot spam expensive handlers like parseLayerVariables (DB + parser call per message).
const maxMessagesPerSecond = 60

const rateLimiterSessionKey = "messagesRateLimiter"

type Registry struct {
	manager  *melody.Melody
	wsRouter wsrouter.WsRouter

	gorm                       *gorm.DB
	logger                     *slog.Logger
	redis                      *redis.Client
	bus                        *buscore.Bus
	channelsOverlaysRepository channels_overlays.Repository
}

func New(
	gorm *gorm.DB,
	logger *slog.Logger,
	redis *redis.Client,
	bus *buscore.Bus,
	wsRouter wsrouter.WsRouter,
	channelsOverlaysRepository channels_overlays.Repository,
	channelsRepository channels.Repository,
	usersRepository users.Repository,
) *Registry {
	m := melody.New()
	m.Config.MaxMessageSize = 1024 * 1024 * 10
	overlaysRegistry := &Registry{
		manager:                    m,
		wsRouter:                   wsRouter,
		gorm:                       gorm,
		logger:                     logger,
		redis:                      redis,
		bus:                        bus,
		channelsOverlaysRepository: channelsOverlaysRepository,
	}

	overlaysRegistry.manager.HandleConnect(
		func(session *melody.Session) {
			session.Set(
				rateLimiterSessionKey,
				rate.NewLimiter(maxMessagesPerSecond, maxMessagesPerSecond),
			)

			err := helpers.CheckChannelByApiKey(session, channelsRepository, usersRepository)
			if err != nil {
				if !errors.Is(err, helpers.ErrUserNotFound) {
					logger.Error("cannot check user by api key", twirlogger.Error(err))
				}
				return
			}
			session.Write([]byte(`{"eventName":"connected to overlays namespace"}`))
		},
	)

	overlaysRegistry.manager.HandleMessage(
		func(session *melody.Session, msg []byte) {
			if value, ok := session.Get(rateLimiterSessionKey); ok {
				if limiter, ok := value.(*rate.Limiter); ok && !limiter.Allow() {
					return
				}
			}

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
