package handlers

import (
	"strings"

	irc "github.com/gempir/go-twitch-irc/v3"
	model "github.com/satont/tsuwari/libs/gomodels"
)

func (c *Handlers) handleWordsCounter(msg irc.PrivateMessage) {
	counters := []model.ChannelWordCounter{}
	err := c.db.Where(`"channelId" = ? AND enabled = ?`, msg.RoomID, true).Find(&counters).Error
	if err != nil {
		c.logger.Sugar().Error(err)
		return
	}

	lowerCasedMsg := strings.ToLower(msg.Message)

	for _, counter := range counters {
		go func(counter model.ChannelWordCounter) {
			count := strings.Count(lowerCasedMsg, counter.Phrase)
			counter.Counter += int32(count)
			c.db.Save(&counter)
		}(counter)
	}
}
