package alerts

import (
	"encoding/json"
	"github.com/olahol/melody"
	"github.com/satont/twir/apps/websockets/internal/namespaces/helpers"
	"github.com/satont/twir/apps/websockets/types"
	"net/http"
)

type Alerts struct {
	manager  *melody.Melody
	services *types.Services
}

func NewAlerts(services *types.Services) *Alerts {
	m := melody.New()
	m.Config.MaxMessageSize = 1024 * 1024 * 10
	alerts := &Alerts{
		manager:  m,
		services: services,
	}

	alerts.manager.HandleConnect(
		func(session *melody.Session) {
			err := helpers.CheckUserByApiKey(services.Gorm, session)
			if err != nil {
				services.Logger.Error(err)
				return
			}
			session.Write([]byte(`{"eventName":"connected to alerts namespace"}`))
		},
	)

	return alerts
}

func (c *Alerts) HandleRequest(w http.ResponseWriter, r *http.Request) {
	_ = c.manager.HandleRequest(w, r)
}

func (c *Alerts) SendEvent(channelId, eventName string, data any) error {
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
			return ok && socketUserId.(string) == channelId
		},
	)

	if err != nil {
		c.services.Logger.Error(err)
		return err
	}

	return nil
}
