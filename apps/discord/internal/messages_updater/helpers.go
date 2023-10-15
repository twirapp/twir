package messages_updater

import (
	"context"

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
	err = c.db.WithContext(ctx).Where(
		`"channelId" = ? AND "integrationId" = ?`,
		channelId,
		discordIntegration.ID,
	).First(integration).Error
	return integration, err
}
