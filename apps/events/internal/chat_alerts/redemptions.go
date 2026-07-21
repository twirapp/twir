package chat_alerts

import (
	"context"
	"slices"
	"strings"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/twirapp/twir/libs/bus-core/events"
	model "github.com/twirapp/twir/libs/gomodels"
)

func (c *ChatAlerts) redemption(
	ctx context.Context,
	settings model.ChatAlertsSettings,
	channelID uuid.UUID,
	req events.RedemptionCreatedMessage,
) error {
	if !settings.Redemptions.Enabled {
		return nil
	}

	if len(settings.Redemptions.Messages) == 0 {
		return nil
	}

	if slices.Contains(settings.Redemptions.IgnoredRewardsIDS, req.ID) {
		return nil
	}

	sample := lo.Sample(settings.Redemptions.Messages)

	text := sample.Text
	text = strings.ReplaceAll(text, "{user}", req.UserName)
	text = strings.ReplaceAll(text, "{reward}", req.RewardName)

	if text == "" {
		return nil
	}

	return c.sendMessage(ctx, channelID, req.BaseInfo.Platform, text)
}
