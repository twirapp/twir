package overlays

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/olahol/melody"
	"github.com/satont/twir/apps/websockets/internal/namespaces/helpers"
	"github.com/satont/twir/apps/websockets/types"
)

type Registry struct {
	manager  *melody.Melody
	services *types.Services
}

func New(services *types.Services) *Registry {
	m := melody.New()
	m.Config.MaxMessageSize = 1024 * 1024 * 10
	overlaysRegistry := &Registry{
		manager:  m,
		services: services,
	}

	overlaysRegistry.manager.HandleConnect(
		func(session *melody.Session) {
			err := helpers.CheckUserByApiKey(services.Gorm, session)
			if err != nil {
				services.Logger.Error(err)
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
