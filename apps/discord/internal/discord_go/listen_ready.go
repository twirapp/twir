package discord_go

import (
	"log/slog"

	"github.com/switchupcb/disgo"
)

func (c *Discord) handleReady(e *disgo.Ready) {
	c.logger.Info(
		"Discord shard is ready",
		slog.Group(
			"bot",
			slog.String("id", e.User.ID),
			slog.String("name", e.User.Username),
			slog.Bool("verified", *e.User.Verified),
		),
		slog.Int("guilds", len(e.Guilds)),
	)
}
