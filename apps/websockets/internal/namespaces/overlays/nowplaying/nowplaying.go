package nowplaying

import (
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/goccy/go-json"

	"github.com/olahol/melody"
	"github.com/redis/go-redis/v9"
	"github.com/satont/twir/apps/websockets/internal/namespaces/helpers"
	"github.com/satont/twir/apps/websockets/types"
	config "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/logger"
	"github.com/twirapp/twir/libs/grpc/tokens"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type NowPlaying struct {
	manager *melody.Melody

	gorm       *gorm.DB
	logger     logger.Logger
	redis      *redis.Client
	config     config.Config
	tokensGrpc tokens.TokensClient
}

type Opts struct {
	fx.In

	Gorm       *gorm.DB
	Logger     logger.Logger
	Redis      *redis.Client
	Config     config.Config
	TokensGrpc tokens.TokensClient
}

func New(opts Opts) *NowPlaying {
	m := melody.New()
	m.Config.MaxMessageSize = 1024 * 1024 * 10
	np := &NowPlaying{
		manager:    m,
		gorm:       opts.Gorm,
		logger:     opts.Logger,
		redis:      opts.Redis,
		config:     opts.Config,
		tokensGrpc: opts.TokensGrpc,
	}

	var fetcherMutex sync.Mutex

	np.manager.HandleConnect(
		func(session *melody.Session) {
			_ = session.Write([]byte(`{"eventName":"connected to nowplaying namespace"}`))
			err := helpers.CheckUserByApiKey(opts.Gorm, session)
			if err != nil {
				opts.Logger.Error("cannot check user by api key", slog.Any("err", err))
				return
			}

			userId, ok := session.Get("userId")
			if !ok {
				return
			}
			castedUserId, ok := userId.(string)
			if !ok {
				return
			}
			_ = np.SendSettings(castedUserId)

			fetcherMutex.Lock()
			defer fetcherMutex.Unlock()

			go func() {
				if err := np.fetcher(session.Request.Context(), castedUserId); err != nil {
					opts.Logger.Error("cannot run fetcher", slog.Any("err", err))
					session.CloseWithMsg([]byte(`{"eventName":"error","data":"cannot run fetcher"}`))
				}
			}()
		},
	)

	np.manager.HandleMessage(
		func(session *melody.Session, msg []byte) {
			np.handleMessage(session, msg)
		},
	)

	http.HandleFunc("/overlays/nowplaying", np.HandleRequest)

	return np
}

func (c *NowPlaying) HandleRequest(w http.ResponseWriter, r *http.Request) {
	_ = c.manager.HandleRequest(w, r)
}

func (c *NowPlaying) SendEvent(channelId, eventName string, data any) error {
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
