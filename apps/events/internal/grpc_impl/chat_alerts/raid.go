package chat_alerts

import (
	"context"
	"fmt"
	"strings"

	"github.com/satont/twir/libs/grpc/generated/bots"
	"github.com/satont/twir/libs/grpc/generated/events"
)

const raidCooldownKey = "raid"

func (c *ChatAlerts) Raid(ctx context.Context, req *events.RaidedMessage) {
	if !c.settings.Raids.Enabled {
		return
	}

	if len(c.settings.Raids.Messages) == 0 {
		return
	}

	if cooldowned, err := c.IsOnCooldown(
		ctx,
		req.BaseInfo.ChannelId,
		raidCooldownKey,
	); err != nil || cooldowned {
		return
	}

	sample := c.takeCountedSample(int(req.Viewers), c.settings.Raids.Messages)
	sample = strings.ReplaceAll(sample, "{count}", fmt.Sprint(req.Viewers))
	sample = strings.ReplaceAll(sample, "{user}", req.UserName)

	if sample == "" {
		return
	}

	c.services.BotsGrpc.SendMessage(
		ctx, &bots.SendMessageRequest{
			ChannelId:      req.BaseInfo.ChannelId,
			Message:        sample,
			IsAnnounce:     nil,
			SkipRateLimits: true,
		},
	)

	if c.settings.Raids.Cooldown != 0 {
		c.SetCooldown(ctx, req.BaseInfo.ChannelId, raidCooldownKey, c.settings.Raids.Cooldown)
	}
}
