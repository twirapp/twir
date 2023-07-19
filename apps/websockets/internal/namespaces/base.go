package namespaces

import (
	"encoding/json"
	"net/http"

	"github.com/olahol/melody"
	"github.com/satont/twir/apps/websockets/internal/namespaces/helpers"
	"github.com/satont/twir/apps/websockets/types"
)

type NameSpace struct {
	manager  *melody.Melody
	services *types.Services
}

func NewNameSpace(services *types.Services) *NameSpace {
	m := melody.New()
	namespace := &NameSpace{
		manager:  m,
		services: services,
	}

	namespace.manager.HandleConnect(
		func(session *melody.Session) {
			session.Write([]byte(`{"event":"connected"}`))
			helpers.CheckUserByApiKey(services.Gorm, session)
		},
	)

	return namespace
}

func (c *NameSpace) HandleRequest(w http.ResponseWriter, r *http.Request) {
	c.manager.HandleRequest(w, r)
}

func (c *NameSpace) SendEvent(userId, eventName string, data any) error {
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
