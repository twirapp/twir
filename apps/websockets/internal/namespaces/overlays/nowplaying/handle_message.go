package nowplaying

import (
	"time"

	"github.com/goccy/go-json"
	"github.com/olahol/melody"
	"github.com/satont/twir/apps/websockets/types"
)

func (c *NowPlaying) handleMessage(session *melody.Session, msg []byte) {
	userId, ok := session.Get("userId")
	if userId == "" || !ok {
		return
	}

	data := &types.WebSocketMessage{
		CreatedAt: time.Now().UTC().String(),
	}
	err := json.Unmarshal(msg, data)
	if err != nil {
		c.logger.Error(err.Error())
		return
	}
}
