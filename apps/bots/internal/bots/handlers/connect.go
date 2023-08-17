package handlers

import "log/slog"

func (c *Handlers) OnConnect() {
	c.logger.Info(
		"Connected to twitch servers",
		slog.String("botName", c.BotClient.TwitchUser.Login),
		slog.String("botId", c.BotClient.Model.ID),
	)
}
