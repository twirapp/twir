package discord_go

import (
	"log/slog"

	"github.com/diamondburned/arikawa/v3/gateway"
)

func (c *Discord) handleShardReady(e *gateway.ReadyEvent) {
	c.logger.Info(
		"Discord shard is ready",
		slog.Group(
			"bot",
			slog.String("id", e.User.ID.String()),
			slog.String("name", e.User.Username),
		),
		slog.Int("guilds", len(e.Guilds)),
		slog.Int("shard_id", e.Shard.ShardID()),
	)
}
