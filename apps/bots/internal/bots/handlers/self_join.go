package handlers

import (
	"log/slog"
)

func (c *Handlers) OnSelfJoin(channel string) {
	c.logger.Info(
		"Joined channel",
		slog.String("botId", c.BotClient.Model.ID),
		slog.String("botName", c.BotClient.TwitchUser.Login),
		slog.String("channel", channel),
	)
}
