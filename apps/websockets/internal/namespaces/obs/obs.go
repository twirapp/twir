package obs

import (
	"encoding/json"
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

type OBS struct {
	manager *melody.Melody
	gorm    *gorm.DB
	logger  logger.Logger
	redis   *redis.Client
}

type Opts struct {
	fx.In

	Gorm   *gorm.DB
	Logger logger.Logger
	Redis  *redis.Client
}

func NewObs(opts Opts) *OBS {
	m := melody.New()
	m.Config.MaxMessageSize = 1024 * 1024 * 10
	obs := &OBS{
		manager: m,
		gorm:    opts.Gorm,
		logger:  opts.Logger,
		redis:   opts.Redis,
	}

	obs.manager.HandleConnect(
		func(session *melody.Session) {
			err := helpers.CheckUserByApiKey(opts.Gorm, session)
			if err != nil {
				opts.Logger.Error(err.Error())
				return
			}
			session.Write([]byte(`{"eventName":"connected"}`))
		},
	)

	obs.manager.HandleMessage(
		func(session *melody.Session, msg []byte) {
			obs.handleMessage(session, msg)
		},
	)

	http.HandleFunc("/obs", obs.HandleRequest)

	return obs
}

func (c *OBS) IsUserConnected(userId string) (bool, error) {
	sessions, err := c.manager.Sessions()
	if err != nil {
		return false, err
	}

	for _, s := range sessions {
		userIdValue, isUserIdExists := s.Get("userId")
		isConnectedValue, isConnectedExists := s.Get("obsConnected")
		if !isUserIdExists || !isConnectedExists {
			continue
		}
		castedUserId, isUserCastOk := userIdValue.(string)
		castedIsConnected, isConnectCastOk := isConnectedValue.(bool)
		if !isUserCastOk || !isConnectCastOk {
			continue
		}
		if castedUserId == userId {
			return castedIsConnected, nil
		}
	}

	return false, nil
}

func (c *OBS) SendEvent(userId, eventName string, data any) error {
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
			return ok && socketUserId.(string) == userId
		},
	)

	if err != nil {
		c.logger.Error(err.Error())
		return err
	}

	return nil
}
