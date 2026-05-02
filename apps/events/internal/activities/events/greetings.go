package events

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/twirapp/twir/apps/events/internal/shared"
	"github.com/twirapp/twir/libs/entities/platform"
	"github.com/twirapp/twir/libs/repositories/events/model"
	"github.com/twirapp/twir/libs/repositories/greetings"
	"go.temporal.io/sdk/activity"
)

func (c *Activity) CreateGreeting(
	ctx context.Context,
	_ model.EventOperation,
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

	dbUser, err := c.usersRepo.GetByPlatformID(ctx, platform.PlatformTwitch, user.ID)
	if err != nil {
		return err
	}

	channelDBID, err := uuid.Parse(data.ChannelDBID)
	if err != nil {
		return err
	}

	_, err = c.greetingsRepository.Create(
		ctx, greetings.CreateInput{
			ChannelID: channelDBID,
			UserID:    dbUser.ID,
			Enabled:   true,
			Text:      *data.RewardInput,
			IsReply:   true,
			Processed: false,
		},
	)

	return err
}
