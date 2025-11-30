package messages_updater

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	discordmodel "github.com/twirapp/twir/libs/repositories/channels_integrations_discord/model"
)

var ErrIntegrationNotFound = fmt.Errorf("integration not found")

// getChannelDiscordIntegrations returns all Discord integrations for a channel
func (c *MessagesUpdater) getChannelDiscordIntegrations(
	ctx context.Context,
	channelId string,
) ([]discordmodel.ChannelIntegrationDiscord, error) {
	integrations, err := c.discordRepo.GetByChannelID(ctx, channelId)
	if err != nil {
		return nil, err
	}

	if len(integrations) == 0 {
		// Also check for additional users
		additionalIntegrations, err := c.discordRepo.GetByAdditionalUserID(ctx, channelId)
		if err != nil {
			return nil, err
		}

		if len(additionalIntegrations) == 0 {
			return nil, ErrIntegrationNotFound
		}

		return additionalIntegrations, nil
	}

	return integrations, nil
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
