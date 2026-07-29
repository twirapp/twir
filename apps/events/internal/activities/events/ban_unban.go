package events

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/kvizyx/twitchy/helix"
	"github.com/twirapp/twir/apps/events/internal/shared"
	deprecatedgormmodel "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/repositories/events/model"
	"go.temporal.io/sdk/activity"
	"golang.org/x/sync/errgroup"
)

func computeBanReason(reason *string) string {
	if reason != nil && *reason != "" {
		return *reason
	}

	return "banned from twirapp"
}

func banDuration(duration int) *int {
	if duration == 0 {
		return nil
	}

	return &duration
}

func (c *Activity) Ban(
	ctx context.Context,
	operation model.EventOperation,
	data shared.EventData,
) error {
	activity.RecordHeartbeat(ctx, nil)

	if operation.Input == nil || *operation.Input == "" {
		return errors.New("input is required for ban operation")
	}

	hydratedName, hydrateErr := c.hydrator.HydrateStringWithData(
		data.ChannelID,
		data.ChannelTwitchUserID,
		data.ChannelDBID,
		*operation.Input,
		data,
	)
	if hydrateErr != nil || len(hydratedName) == 0 {
		return fmt.Errorf("cannot hydrate string %w", hydrateErr)
	}

	hydratedName = strings.TrimSpace(strings.ReplaceAll(hydratedName, "@", ""))

	dbChannel, err := c.getTwitchChannelDbEntity(ctx, data)
	if err != nil {
		return err
	}

	var errwg errgroup.Group
	botTwitchClient, twitchClientError := c.getHelixChannelBotApiClient(
		ctx,
		dbChannel.BotID,
		dbChannel.BroadcasterUserID,
	)
	if twitchClientError != nil {
		return twitchClientError
	}
	broadcasterTwitchClient, twitchBotClientError := c.getHelixChannelApiClient(ctx, dbChannel.BroadcasterUserID)
	if twitchBotClientError != nil {
		return twitchBotClientError
	}

	var targetUser helix.User

	errwg.Go(
		func() error {
			u, err := c.getHelixUserByLogin(ctx, broadcasterTwitchClient, hydratedName)
			if err != nil {
				return err
			}
			targetUser = u
			return nil
		},
	)

	var mods []helix.Moderator

	errwg.Go(
		func() error {
			m, err := c.getChannelMods(ctx, broadcasterTwitchClient, twitchBroadcasterID(data))
			if err != nil {
				return err
			}
			mods = m
			return nil
		},
	)

	if err := errwg.Wait(); err != nil {
		return err
	}

	if targetUser.ID == dbChannel.BotID || targetUser.ID == dbChannel.BroadcasterUserID {
		return errors.New("cannot ban bot or channel owner")
	}

	for _, mod := range mods {
		if mod.UserID == targetUser.ID {
			return errors.New("cannot ban moderator")
		}
	}

	_, err = botTwitchClient.Moderation.BanUser(
		ctx, helix.BanUserRequest{
			BroadcasterID: twitchBroadcasterID(data),
			ModeratorID:   dbChannel.BotID,
			Data: helix.BanUserBody{
				Duration: banDuration(operation.TimeoutTime),
				Reason:   computeBanReason(operation.TimeoutMessage),
				UserID:   targetUser.ID,
			},
		},
	)

	if err != nil {
		return fmt.Errorf("cannot ban targetUser: %w", err)
	}
	return nil
}

func (c *Activity) Unban(
	ctx context.Context,
	operation model.EventOperation,
	data shared.EventData,
) error {
	activity.RecordHeartbeat(ctx, nil)

	if operation.Input == nil || *operation.Input == "" {
		return errors.New("input is required for unban operation")
	}

	hydratedName, hydrateErr := c.hydrator.HydrateStringWithData(
		data.ChannelID,
		data.ChannelTwitchUserID,
		data.ChannelDBID,
		*operation.Input,
		data,
	)
	if hydrateErr != nil || len(hydratedName) == 0 {
		return fmt.Errorf("cannot hydrate string %w", hydrateErr)
	}

	hydratedName = strings.TrimSpace(strings.ReplaceAll(hydratedName, "@", ""))

	dbChannel, dbChannelErr := c.getTwitchChannelDbEntity(ctx, data)
	if dbChannelErr != nil {
		return dbChannelErr
	}

	botTwitchClient, twitchClientError := c.getHelixChannelBotApiClient(
		ctx,
		dbChannel.BotID,
		dbChannel.BroadcasterUserID,
	)
	if twitchClientError != nil {
		return twitchClientError
	}
	broadcasterTwitchClient, twitchBotClientError := c.getHelixChannelApiClient(ctx, dbChannel.BroadcasterUserID)
	if twitchBotClientError != nil {
		return twitchBotClientError
	}

	targetUser, userErr := c.getHelixUserByLogin(ctx, broadcasterTwitchClient, hydratedName)
	if userErr != nil {
		return userErr
	}

	_, err := botTwitchClient.Moderation.UnbanUser(
		ctx, helix.UnbanUserRequest{
			BroadcasterID: twitchBroadcasterID(data),
			ModeratorID:   dbChannel.BotID,
			UserID:        targetUser.ID,
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func (c *Activity) BanRandom(
	ctx context.Context,
	operation model.EventOperation,
	data shared.EventData,
) error {
	activity.RecordHeartbeat(ctx, nil)

	dbChannel, err := c.getTwitchChannelDbEntity(ctx, data)
	if err != nil {
		return err
	}

	botTwitchClient, twitchClientError := c.getHelixChannelBotApiClient(
		ctx,
		dbChannel.BotID,
		dbChannel.BroadcasterUserID,
	)
	if twitchClientError != nil {
		return twitchClientError
	}
	broadcasterTwitchClient, twitchBotClientError := c.getHelixChannelApiClient(ctx, dbChannel.BroadcasterUserID)
	if twitchBotClientError != nil {
		return twitchBotClientError
	}

	mods, err := c.getChannelMods(ctx, broadcasterTwitchClient, twitchBroadcasterID(data))
	if err != nil {
		return err
	}

	// exclude mods, channel and bot
	var excludedForBanUsers []string
	for _, mod := range mods {
		excludedForBanUsers = append(excludedForBanUsers, mod.UserID)
	}
	excludedForBanUsers = append(excludedForBanUsers, dbChannel.BroadcasterUserID, dbChannel.BotID)

	randomOnlineUser := &deprecatedgormmodel.UsersOnline{}
	err = c.db.
		Where(`"userId" not in ?`, excludedForBanUsers).
		Order("random()").
		Find(randomOnlineUser).
		Error
	if err != nil {
		return err
	}

	if !randomOnlineUser.UserId.Valid {
		return errors.New("cannot get random user")
	}

	timeoutTime := operation.TimeoutTime
	if operation.Type == model.EventOperationTypeBanRandom {
		timeoutTime = 0
	} else if timeoutTime == 0 {
		timeoutTime = 600
	}

	_, err = botTwitchClient.Moderation.BanUser(
		ctx, helix.BanUserRequest{
			BroadcasterID: twitchBroadcasterID(data),
			ModeratorID:   dbChannel.BotID,
			Data: helix.BanUserBody{
				Duration: banDuration(timeoutTime),
				Reason:   computeBanReason(operation.TimeoutMessage),
				UserID:   randomOnlineUser.UserId.String,
			},
		},
	)
	if err != nil {
		return err
	}
	// if len(c.data.PrevOperation.BannedUserName) > 0 {
	// 	c.data.PrevOperation.BannedUserName += ", " + randomOnlineUser.UserName.String
	// } else {
	// 	c.data.PrevOperation.BannedUserName = randomOnlineUser.UserName.String
	// }

	return nil
}
