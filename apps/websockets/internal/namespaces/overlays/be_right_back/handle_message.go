package be_right_back

import (
	"time"

	"github.com/goccy/go-json"
	"github.com/olahol/melody"
	"github.com/twirapp/twir/apps/websockets/types"
)

func (c *BeRightBack) handleMessage(session *melody.Session, msg []byte) {
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

	if data.EventName == "getSettings" {
		err := c.SendSettings(userId.(string))
		if err != nil {
			c.logger.Error(err.Error())
		}
	}
}
