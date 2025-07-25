package messages_updater

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	model "github.com/twirapp/twir/libs/gomodels"
)

var ErrIntegrationNotFound = fmt.Errorf("integration not found")

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

		return nil, ErrIntegrationNotFound
	}

	return integration, err
}

type replaceMessageVarsOpts struct {
	UserName    string `json:"userName"`
	DisplayName string `json:"displayName"`

	CategoryName string `json:"categoryName"`
	Title        string `json:"title"`
}

func (c *MessagesUpdater) replaceMessageVars(text string, vars replaceMessageVarsOpts) string {
	varsMap := map[string]string{}

	bytes, err := json.Marshal(vars)
	if err != nil {
		return text
	}

	err = json.Unmarshal(bytes, &varsMap)
	if err != nil {
		return text
	}

	for key, value := range varsMap {
		text = strings.ReplaceAll(text, fmt.Sprintf("{%s}", key), value)
	}

	return text
}
