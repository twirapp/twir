package handlers

import (
	"github.com/gempir/go-twitch-irc/v3"
	"log/slog"
)

func (c *Handlers) OnSelfJoin(msg twitch.UserJoinMessage) {
	c.logger.Info(
		"Joined channel",
		slog.String("botId", c.BotClient.Model.ID),
		slog.String("botName", c.BotClient.TwitchUser.Login),
		slog.String("channel", msg.Channel),
	)
}
