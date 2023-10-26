package tts

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/olahol/melody"
	"github.com/redis/go-redis/v9"
	"github.com/satont/twir/apps/websockets/internal/namespaces/helpers"
	"github.com/satont/twir/apps/websockets/types"
	"github.com/satont/twir/libs/logger"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type TTS struct {
	manager *melody.Melody

	gorm   *gorm.DB
	logger logger.Logger
	redis  *redis.Client
}

type Opts struct {
	fx.In

	Gorm   *gorm.DB
	Logger logger.Logger
	Redis  *redis.Client
}

func NewTts(opts Opts) *TTS {
	m := melody.New()
	tts := &TTS{
		manager: m,
		gorm:    opts.Gorm,
		logger:  opts.Logger,
		redis:   opts.Redis,
	}

	tts.manager.HandleConnect(
		func(session *melody.Session) {
			helpers.CheckUserByApiKey(opts.Gorm, session)
		},
	)

	http.HandleFunc("/tts", tts.HandleRequest)

	return tts
}

func (c *TTS) HandleRequest(w http.ResponseWriter, r *http.Request) {
	c.manager.HandleRequest(w, r)
}

func (c *TTS) SendEvent(userId, eventName string, data any) error {
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
			return ok && socketUserId.(string) == userId
		},
	)

	if err != nil {
		c.logger.Error("cannot broadcast message", slog.Any("err", err))
		return err
	}

	return nil
}
