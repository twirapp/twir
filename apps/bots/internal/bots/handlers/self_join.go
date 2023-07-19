package handlers

import (
	"github.com/gempir/go-twitch-irc/v3"
)

func (c *Handlers) OnSelfJoin(msg twitch.UserJoinMessage) {
	c.logger.Sugar().
		Infow(
			"Joined channel",
			"botId",
			c.BotClient.Model.ID,
			"botName",
			c.BotClient.TwitchUser.Login,
			"channel",
			msg.Channel,
		)
}
