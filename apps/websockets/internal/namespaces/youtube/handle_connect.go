package youtube

import (
	"encoding/json"
	"github.com/olahol/melody"
	"github.com/satont/twir/apps/websockets/types"
	model "github.com/satont/twir/libs/gomodels"
)

func (c *YouTube) handleConnect(session *melody.Session) {
	//const songs = await repository.find({
	//where: {
	//	channelId,
	//		deletedAt: IsNull(),
	//},
	//order: {
	//queuePosition: 'asc',
	//},
	//});

	userId, _ := session.Get("userId")

	var currentSongs []model.RequestedSong
	err := c.services.Gorm.
		Where(`"channelId" = ? AND "deletedAt" IS NULL`, userId.(string)).
		Find(&currentSongs).
		Error

	if err != nil {
		c.services.Logger.Error(err)
		return
	}

	outCome := &types.WebSocketMessage{
		EventName: "currentQueue",
		Data:      currentSongs,
	}

	bytes, err := json.Marshal(outCome)
	if err != nil {
		c.services.Logger.Error(err)
		return
	}

	c.manager.BroadcastFilter(
		bytes, func(session *melody.Session) bool {
			id, ok := session.Get("userId")
			return ok && id.(string) == userId.(string)
		},
	)
}
