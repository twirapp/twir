package discord_go

import (
	"context"

	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/twirapp/twir/libs/logger"
)

func (c *Discord) handleGuildDelete(e *gateway.GuildDeleteEvent) {
	if e.Unavailable {
		return
	}

	ctx := context.Background()
	guildID := e.ID.String()

	// Get all integrations for this guild
	integrations, err := c.discordRepo.GetByGuildID(ctx, guildID)
	if err != nil {
		c.logger.Error("failed to find channels integrations", logger.Error(err))
		return
	}

	// Delete all integrations for this guild
	for _, integration := range integrations {
		if err := c.discordRepo.Delete(ctx, integration.ID); err != nil {
			c.logger.Error("failed to delete channels integration", logger.Error(err))
			continue
		}
	}
}
