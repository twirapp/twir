package chat_client

import (
	"log/slog"

	"github.com/gempir/go-twitch-irc/v3"
)

func (c *ChatClient) onSelfJoin(msg twitch.UserJoinMessage, shard string) {
	c.services.Logger.Info(
		"Joined channel",
		slog.String("botId", c.TwitchUser.ID),
		slog.String("botName", c.TwitchUser.Login),
		slog.String("channel", msg.Channel),
		slog.String("shard", shard),
	)
}
