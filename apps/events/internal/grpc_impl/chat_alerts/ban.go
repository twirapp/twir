package chat_alerts

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/samber/lo"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/grpc/generated/bots"
	"github.com/satont/twir/libs/grpc/generated/events"
)

const banCooldownKey = "ban"

func (c *ChatAlerts) Ban(ctx context.Context, req *events.ChannelBanMessage) {
	if !c.settings.Ban.Enabled {
		return
	}

	if len(c.settings.Ban.Messages) == 0 {
		return
	}

	if cooldowned, err := c.IsOnCooldown(
		ctx,
		req.BaseInfo.ChannelId,
		banCooldownKey,
	); err != nil || cooldowned {
		return
	}

	for _, ignoredName := range c.settings.Ban.IgnoreTimeoutFrom {
		if strings.ToLower(ignoredName) == req.ModeratorUserLogin {
			return
		}
	}

	var sample string
	var time int

	if req.IsPermanent {
		permMessage, ok := lo.Find(
			c.settings.Ban.Messages,
			func(m model.ChatAlertsCountedMessage) bool {
				return m.Count == 0
			},
		)
		if ok {
			sample = permMessage.Text
		}
	} else {
		parsedEndsAt, err := strconv.Atoi(req.EndsAt)
		if err != nil {
			return
		}

		if parsedEndsAt == 0 {
			parsedEndsAt = 1
		}

		sample = c.takeCountedSample(parsedEndsAt, c.settings.Ban.Messages)
		time = parsedEndsAt
	}

	sample = strings.ReplaceAll(sample, "{userName}", req.UserName)
	sample = strings.ReplaceAll(sample, "{moderatorName}", req.ModeratorUserName)
	sample = strings.ReplaceAll(
		sample,
		"{time}",
		lo.If(req.IsPermanent, "permanent").Else(fmt.Sprintf("%v", time)),
	)
	sample = strings.ReplaceAll(sample, "{reason}", req.Reason)

	if sample == "" {
		return
	}

	_, err := c.services.BotsGrpc.SendMessage(
		ctx,
		&bots.SendMessageRequest{
			ChannelId:      req.BaseInfo.ChannelId,
			Message:        sample,
			SkipRateLimits: true,
		},
	)
	if err != nil {
		fmt.Println(err)
	}

	if c.settings.Ban.Cooldown != 0 {
		c.SetCooldown(ctx, req.BaseInfo.ChannelId, banCooldownKey, c.settings.Ban.Cooldown)
	}
}
