package events

import (
	"context"
	"errors"

	"github.com/nicklaw5/helix/v2"
	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/events/internal/shared"
	"github.com/twirapp/twir/libs/repositories/events/model"
	"go.temporal.io/sdk/activity"
)

func (c *Activity) SwitchEmoteOnly(
	ctx context.Context,
	operation model.EventOperation,
	data shared.EventData,
) error {
	activity.RecordHeartbeat(ctx, nil)

	dbEntity, dbEntityErr := c.getChannelDbEntity(ctx, data.ChannelID)
	if dbEntityErr != nil {
		return dbEntityErr
	}

	twitchClient, twitchClientErr := c.getHelixBotApiClient(ctx, dbEntity.BotID)
	if twitchClientErr != nil {
		return twitchClientErr
	}

	resp, err := twitchClient.UpdateChatSettings(
		&helix.UpdateChatSettingsParams{
			BroadcasterID: data.ChannelID,
			ModeratorID:   dbEntity.BotID,
			EmoteMode: lo.ToPtr(
				lo.
					If(operation.Type == model.EventOperationTypeEnableEmoteOnly, true).
					Else(false),
			),
		},
	)
	if err != nil {
		return err
	}
	if resp.ErrorMessage != "" {
		return errors.New(resp.ErrorMessage)
	}

	return nil
}

func (c *Activity) SwitchSubMode(
	ctx context.Context,
	operation model.EventOperation,
	data shared.EventData,
) error {
	activity.RecordHeartbeat(ctx, nil)

	dbEntity, dbEntityErr := c.getChannelDbEntity(ctx, data.ChannelID)
	if dbEntityErr != nil {
		return dbEntityErr
	}

	twitchClient, twitchClientErr := c.getHelixBotApiClient(ctx, dbEntity.BotID)
	if twitchClientErr != nil {
		return twitchClientErr
	}

	resp, err := twitchClient.UpdateChatSettings(
		&helix.UpdateChatSettingsParams{
			BroadcasterID: data.ChannelID,
			ModeratorID:   dbEntity.BotID,
			SubscriberMode: lo.ToPtr(
				lo.
					If(operation.Type == model.EventOperationTypeEnableSubmode, true).
					Else(false),
			),
		},
	)
	if err != nil {
		return err
	}
	if resp.ErrorMessage != "" {
		return errors.New(resp.ErrorMessage)
	}

	return nil
}
