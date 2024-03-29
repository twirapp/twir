package events

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/guregu/null"
	"github.com/nicklaw5/helix/v2"
	"github.com/satont/twir/apps/events/internal/shared"
	model "github.com/satont/twir/libs/gomodels"
	"go.temporal.io/sdk/activity"
	"golang.org/x/sync/errgroup"
)

func computeBanReason(reason null.String) string {
	if reason.Valid && reason.String != "" {
		return reason.String
	}

	return "banned from twirapp"
}

func (c *Activity) Ban(
	ctx context.Context,
	operation model.EventOperation,
	data shared.EventData,
) error {
	activity.RecordHeartbeat(ctx, nil)

	hydratedName, hydrateErr := c.hydrator.HydrateStringWithData(
		data.ChannelID,
		operation.Input.String,
		data,
	)
	if hydrateErr != nil || len(hydratedName) == 0 {
		return fmt.Errorf("cannot hydrate string %w", hydrateErr)
	}

	hydratedName = strings.TrimSpace(strings.ReplaceAll(hydratedName, "@", ""))

	dbChannel, err := c.getChannelDbEntity(ctx, data.ChannelID)
	if err != nil {
		return err
	}

	var errwg errgroup.Group
	botTwitchClient, twitchClientError := c.getHelixBotApiClient(ctx, dbChannel.BotID)
	if twitchClientError != nil {
		return twitchClientError
	}
	broadcasterTwitchClient, twitchBotClientError := c.getHelixChannelApiClient(ctx, dbChannel.ID)
	if twitchBotClientError != nil {
		return twitchBotClientError
	}

	var targetUser helix.User

	errwg.Go(
		func() error {
			u, err := c.getHelixUserByLogin(broadcasterTwitchClient, hydratedName)
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
			m, err := c.getChannelMods(broadcasterTwitchClient, data.ChannelID)
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

	if targetUser.ID == dbChannel.BotID || targetUser.ID == dbChannel.ID {
		return errors.New("cannot ban bot or channel owner")
	}

	for _, mod := range mods {
		if mod.UserID == targetUser.ID {
			return errors.New("cannot ban moderator")
		}
	}

	banReq, err := botTwitchClient.BanUser(
		&helix.BanUserParams{
			BroadcasterID: data.ChannelID,
			ModeratorId:   dbChannel.BotID,
			Body: helix.BanUserRequestBody{
				Duration: operation.TimeoutTime,
				Reason:   computeBanReason(operation.TimeoutMessage),
				UserId:   targetUser.ID,
			},
		},
	)

	if err != nil {
		return fmt.Errorf("cannot ban targetUser: %w", err)
	}
	if banReq.ErrorMessage != "" {
		return fmt.Errorf("cannot ban targetUser: %s", banReq.ErrorMessage)
	}

	return nil
}

func (c *Activity) Unban(
	ctx context.Context,
	operation model.EventOperation,
	data shared.EventData,
) error {
	activity.RecordHeartbeat(ctx, nil)

	hydratedName, hydrateErr := c.hydrator.HydrateStringWithData(
		data.ChannelID,
		operation.Input.String,
		data,
	)
	if hydrateErr != nil || len(hydratedName) == 0 {
		return fmt.Errorf("cannot hydrate string %w", hydrateErr)
	}

	hydratedName = strings.TrimSpace(strings.ReplaceAll(hydratedName, "@", ""))

	dbChannel, dbChannelErr := c.getChannelDbEntity(ctx, data.ChannelID)
	if dbChannelErr != nil {
		return dbChannelErr
	}

	botTwitchClient, twitchClientError := c.getHelixBotApiClient(ctx, dbChannel.BotID)
	if twitchClientError != nil {
		return twitchClientError
	}
	broadcasterTwitchClient, twitchBotClientError := c.getHelixChannelApiClient(ctx, dbChannel.ID)
	if twitchBotClientError != nil {
		return twitchBotClientError
	}

	targetUser, userErr := c.getHelixUserByLogin(broadcasterTwitchClient, hydratedName)
	if userErr != nil {
		return userErr
	}

	resp, err := botTwitchClient.UnbanUser(
		&helix.UnbanUserParams{
			BroadcasterID: data.ChannelID,
			ModeratorID:   dbChannel.BotID,
			UserID:        targetUser.ID,
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

func (c *Activity) BanRandom(
	ctx context.Context,
	operation model.EventOperation,
	data shared.EventData,
) error {
	activity.RecordHeartbeat(ctx, nil)

	dbChannel, err := c.getChannelDbEntity(ctx, data.ChannelID)
	if err != nil {
		return err
	}

	botTwitchClient, twitchClientError := c.getHelixBotApiClient(ctx, dbChannel.BotID)
	if twitchClientError != nil {
		return twitchClientError
	}
	broadcasterTwitchClient, twitchBotClientError := c.getHelixChannelApiClient(ctx, dbChannel.ID)
	if twitchBotClientError != nil {
		return twitchBotClientError
	}

	mods, err := c.getChannelMods(broadcasterTwitchClient, data.ChannelID)
	if err != nil {
		return err
	}

	// exclude mods, channel and bot
	var excludedForBanUsers []string
	for _, mod := range mods {
		excludedForBanUsers = append(excludedForBanUsers, mod.UserID)
	}
	excludedForBanUsers = append(excludedForBanUsers, dbChannel.ID, dbChannel.BotID)

	randomOnlineUser := &model.UsersOnline{}
	err = c.db.
		Where(`"userId" not in ?`, excludedForBanUsers).
		Order("random()").
		Find(randomOnlineUser).
		Error
	if err != nil {
		return err
	}

	if randomOnlineUser == nil || !randomOnlineUser.UserId.Valid {
		return errors.New("cannot get random user")
	}

	timeoutTime := operation.TimeoutTime
	if operation.Type == model.OperationBanRandom {
		timeoutTime = 0
	} else if timeoutTime == 0 {
		timeoutTime = 600
	}

	banReq, err := botTwitchClient.BanUser(
		&helix.BanUserParams{
			BroadcasterID: data.ChannelID,
			ModeratorId:   dbChannel.BotID,
			Body: helix.BanUserRequestBody{
				Duration: timeoutTime,
				Reason:   computeBanReason(operation.TimeoutMessage),
				UserId:   randomOnlineUser.UserId.String,
			},
		},
	)
	if err != nil {
		return err
	}
	if banReq.ErrorMessage != "" {
		return errors.New(banReq.ErrorMessage)
	}

	// if len(c.data.PrevOperation.BannedUserName) > 0 {
	// 	c.data.PrevOperation.BannedUserName += ", " + randomOnlineUser.UserName.String
	// } else {
	// 	c.data.PrevOperation.BannedUserName = randomOnlineUser.UserName.String
	// }

	return nil
}
