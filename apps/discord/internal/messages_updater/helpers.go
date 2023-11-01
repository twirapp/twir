package messages_updater

import (
	"context"
	"fmt"

	model "github.com/satont/twir/libs/gomodels"
)

func (c *MessagesUpdater) getChannelDiscordIntegration(
	ctx context.Context,
	channelId string,
) (*model.ChannelsIntegrations, error) {
	discordIntegration := model.Integrations{}
	err := c.db.WithContext(ctx).Where(
		`service = ?`,
		model.IntegrationServiceDiscord,
	).First(&discordIntegration).Error
	if err != nil {
		return nil, err
	}

	integration := &model.ChannelsIntegrations{}
	err = c.db.WithContext(ctx).
		Where(
			`"integrationId" = ? AND "channelId" = ?`,
			discordIntegration.ID,
			channelId,
		).
		Preload("Channel").
		Find(integration).Error

	if integration.ID != "" && integration.Channel.IsEnabled {
		return integration, nil
	}

	additionalUserQuery := fmt.Sprintf(
		`
SELECT *
FROM channels_integrations
WHERE EXISTS (
    SELECT 1
    FROM jsonb_array_elements(data->'discord'->'guilds') guild
    WHERE (guild->'additionalUsersIdsForLiveCheck') @> '["%s"]'::jsonb
)
`, channelId,
	)

	err = c.db.WithContext(ctx).Raw(additionalUserQuery).Scan(integration).Error
	if err != nil {
		return nil, err
	}

	if integration.ID == "" {
		channel := &model.Channels{}
		if err := c.db.WithContext(ctx).Where(`id = ?`, channelId).First(channel).Error; err != nil {
			return nil, err
		}

		if !channel.IsEnabled {
			return nil, fmt.Errorf("channel is not enabled")
		}

		integration.Channel = channel

		return nil, fmt.Errorf("integration not found")
	}

	return integration, err
}
