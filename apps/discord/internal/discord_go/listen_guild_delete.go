package discord_go

import (
	"fmt"
	"log/slog"

	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/samber/lo"
	model "github.com/twirapp/twir/libs/gomodels"
)

func (c *Discord) handleGuildDelete(e *gateway.GuildDeleteEvent) {
	if e.Unavailable {
		return
	}

	var integrations []model.ChannelsIntegrations
	if err := c.db.
		Where(`data->'discord'->'guilds' @> ?::jsonb`, fmt.Sprintf(`[{"id": "%s"}]`, e.ID)).
		Find(&integrations).Error; err != nil {
		c.logger.Error("failed to find channels integrations", slog.Any("error", err))
		return
	}

	for _, integration := range integrations {
		integration.Data.Discord.Guilds = lo.Filter(
			integration.Data.Discord.Guilds,
			func(guild model.ChannelIntegrationDataDiscordGuild, _ int) bool {
				return guild.ID != e.ID.String()
			},
		)

		if err := c.db.Save(&integration).Error; err != nil {
			c.logger.Error("failed to save channels integrations", slog.Any("error", err))
			return
		}
	}
}
