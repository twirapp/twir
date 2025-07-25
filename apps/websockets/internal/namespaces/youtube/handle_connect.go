package youtube

import (
	"encoding/json"
	"time"

	"github.com/olahol/melody"
	"github.com/twirapp/twir/apps/websockets/types"
	model "github.com/twirapp/twir/libs/gomodels"
)

func (c *YouTube) handleConnect(session *melody.Session) {
	// const songs = await repository.find({
	// where: {
	//	channelId,
	//		deletedAt: IsNull(),
	// },
	// order: {
	// queuePosition: 'asc',
	// },
	// });

	userId, _ := session.Get("userId")

	var currentSongs []model.RequestedSong
	err := c.gorm.
		Where(`"channelId" = ? AND "deletedAt" IS NULL`, userId.(string)).
		Order(`"queuePosition" ASC`).
		Find(&currentSongs).
		Error

	if err != nil {
		c.logger.Error(err.Error())
		return
	}

	outCome := &types.WebSocketMessage{
		EventName: "currentQueue",
		Data:      currentSongs,
		CreatedAt: time.Now().UTC().String(),
	}

	bytes, err := json.Marshal(outCome)
	if err != nil {
		c.logger.Error(err.Error())
		return
	}

	c.manager.BroadcastFilter(
		bytes, func(session *melody.Session) bool {
			id, ok := session.Get("userId")
			return ok && id.(string) == userId.(string)
		},
	)
}
