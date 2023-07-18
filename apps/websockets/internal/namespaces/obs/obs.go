package obs

import (
	"encoding/json"
	"github.com/olahol/melody"
	"github.com/satont/twir/apps/websockets/internal/namespaces/helpers"
	"github.com/satont/twir/apps/websockets/types"
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

	obs.manager.HandleConnect(
		func(session *melody.Session) {
			err := helpers.CheckUserByApiKey(services.Gorm, session)
			if err != nil {
				services.Logger.Error(err)
				return
			}
		},
	)

	obs.manager.HandleMessage(
		func(session *melody.Session, msg []byte) {
			obs.handleMessage(session, msg)
		},
	)

	return obs
}

func (c *OBS) IsUserConnected(userId string) (bool, error) {
	sessions, err := c.manager.Sessions()
	if err != nil {
		return false, err
	}

	for _, s := range sessions {
		value, exists := s.Get("userId")
		if !exists {
			continue
		}
		castedValue, ok := value.(string)
		if !ok {
			continue
		}
		if castedValue == userId {
			return true, nil
		}
	}

	return false, nil
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

	err = c.manager.BroadcastFilter(
		bytes, func(session *melody.Session) bool {
			socketUserId, ok := session.Get("userId")
			return ok && socketUserId.(string) == userId
		},
	)

	if err != nil {
		c.services.Logger.Error(err)
		return err
	}

	return nil
}
