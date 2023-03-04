package tts

import (
	"encoding/json"
	"github.com/olahol/melody"
	"github.com/satont/tsuwari/apps/websockets/internal/namespaces/helpers"
	"github.com/satont/tsuwari/apps/websockets/types"
	"net/http"
)

type TTS struct {
	manager  *melody.Melody
	services *types.Services
}

func NewTts(services *types.Services) *TTS {
	m := melody.New()
	tts := &TTS{
		manager:  m,
		services: services,
	}

	tts.manager.HandleConnect(func(session *melody.Session) {
		helpers.CheckUserByApiKey(services.Gorm, session)
	})

	return tts
}

func (c *TTS) HandleRequest(w http.ResponseWriter, r *http.Request) {
	c.manager.HandleRequest(w, r)
}

func (c *TTS) SendEvent(userId, eventName string, data any) error {
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
