package events

import (
	"context"

	"github.com/kvizyx/twitchy/helix"
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

	dbEntity, dbEntityErr := c.getTwitchChannelDbEntity(ctx, data)
	if dbEntityErr != nil {
		return dbEntityErr
	}

	twitchClient, twitchClientErr := c.getHelixChannelBotApiClient(
		ctx,
		dbEntity.BotID,
		dbEntity.BroadcasterUserID,
	)
	if twitchClientErr != nil {
		return twitchClientErr
	}

	_, err := twitchClient.Chat.UpdateChatSettings(
		ctx, helix.UpdateChatSettingsRequest{
			BroadcasterID: twitchBroadcasterID(data),
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
	return nil
}

func (c *Activity) SwitchSubMode(
	ctx context.Context,
	operation model.EventOperation,
	data shared.EventData,
) error {
	activity.RecordHeartbeat(ctx, nil)

	dbEntity, dbEntityErr := c.getTwitchChannelDbEntity(ctx, data)
	if dbEntityErr != nil {
		return dbEntityErr
	}

	twitchClient, twitchClientErr := c.getHelixChannelBotApiClient(
		ctx,
		dbEntity.BotID,
		dbEntity.BroadcasterUserID,
	)
	if twitchClientErr != nil {
		return twitchClientErr
	}

	_, err := twitchClient.Chat.UpdateChatSettings(
		ctx, helix.UpdateChatSettingsRequest{
			BroadcasterID: twitchBroadcasterID(data),
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
	return nil
}
