package chat_alerts

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/samber/lo"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/bus-core/bots"
	"github.com/twirapp/twir/libs/bus-core/events"
)

func (c *ChatAlerts) ban(
	ctx context.Context,
	settings model.ChatAlertsSettings,
	req events.ChannelBanMessage,
) error {
	if !settings.Ban.Enabled {
		return nil
	}

	if len(settings.Ban.Messages) == 0 {
		return nil
	}

	for _, ignoredName := range settings.Ban.IgnoreTimeoutFrom {
		if strings.ToLower(ignoredName) == req.ModeratorUserLogin {
			return nil
		}
	}

	var sample string
	var time int

	if req.IsPermanent {
		permMessage, ok := lo.Find(
			settings.Ban.Messages,
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
			return err
		}

		if parsedEndsAt == 0 {
			parsedEndsAt = 1
		}

		sample = c.takeCountedSample(parsedEndsAt, settings.Ban.Messages)
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
		return nil
	}

	err := c.bus.Bots.SendMessage.Publish(
		bots.SendMessageRequest{
			ChannelId:      req.BaseInfo.ChannelID,
			Message:        sample,
			SkipRateLimits: true,
		},
	)

	return err
}
