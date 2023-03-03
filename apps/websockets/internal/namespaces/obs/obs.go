package obs

import (
	"encoding/json"
	"github.com/olahol/melody"
	"github.com/satont/tsuwari/apps/websockets/internal/namespaces/helpers"
	"github.com/satont/tsuwari/apps/websockets/types"
)

type OBS struct {
	manager  *melody.Melody
	services *types.Services
}

func NewObs(services *types.Services) *OBS {
	m := melody.New()
	m.Config.MaxMessageSize = 1024 * 1024 * 10
	obs := &OBS{
		manager:  m,
		services: services,
	}

	obs.manager.HandleConnect(func(session *melody.Session) {
		err := helpers.CheckUserByApiKey(services.Gorm, session)
		if err != nil {
			services.Logger.Error(err)
		}
	})

	obs.manager.HandleMessage(func(session *melody.Session, msg []byte) {
		obs.handleMessage(session, msg)
	})

	return obs
}

func (c *OBS) SendEvent(userId, eventName string, data any) error {
	message := &types.WebSocketMessage{
		EventName: eventName,
		Data:      data,
	}

	bytes, err := json.Marshal(message)
	if err != nil {
		c.services.Logger.Error(err)
		return err
	}

	err = c.manager.BroadcastFilter(bytes, func(session *melody.Session) bool {
		socketUserId, ok := session.Get("userId")
		return ok && socketUserId.(string) == userId
	})

	if err != nil {
		c.services.Logger.Error(err)
		return err
	}

	return nil
}
