package chat_client

import "log/slog"

func (c *ChatClient) onConnect(t string) {
	c.services.Logger.Info(
		"Connected to twitch servers",
		slog.String("type", t),
		slog.String("botName", c.TwitchUser.Login),
		slog.String("botId", c.Model.ID),
	)
}
