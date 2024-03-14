package events

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/satont/twir/apps/events/internal/shared"
	model "github.com/satont/twir/libs/gomodels"
	"go.temporal.io/sdk/activity"
)

func (c *Activity) CreateGreeting(
	ctx context.Context,
	operation model.EventOperation,
	data shared.EventData,
) error {
	activity.RecordHeartbeat(ctx, nil)

	if data.RewardInput == nil {
		return errors.New("reward input is empty")
	}

	dbChannel, dbChannelErr := c.getChannelDbEntity(ctx, data.ChannelID)
	if dbChannelErr != nil {
		return dbChannelErr
	}

	twitchClient, twitchClientErr := c.getHelixBotApiClient(ctx, dbChannel.BotID)
	if twitchClientErr != nil {
		return twitchClientErr
	}

	user, userErr := c.getHelixUserById(twitchClient, data.UserID)
	if userErr != nil {
		return userErr
	}

	newGreeting := model.ChannelsGreetings{
		ID:        uuid.NewString(),
		ChannelID: data.ChannelID,
		UserID:    user.ID,
		Enabled:   true,
		Text:      *data.RewardInput,
		IsReply:   true,
		Processed: false,
	}

	return c.db.Create(&newGreeting).Error
}
