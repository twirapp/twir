package events

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/nicklaw5/helix/v2"
	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/events/internal/shared"
	deprecatedgormmodel "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/repositories/events/model"
	"go.temporal.io/sdk/activity"
	"golang.org/x/sync/errgroup"
)

func (c *Activity) ModOrUnmod(
	ctx context.Context,
	operation model.EventOperation,
	data shared.EventData,
) error {
	activity.RecordHeartbeat(ctx, nil)

	if operation.Input == nil || *operation.Input == "" {
		return fmt.Errorf("input is required for operation %s", operation.Type)
	}

	hydratedName, hydrationErr := c.hydrator.HydrateStringWithData(
		data.ChannelID,
		*operation.Input,
		data,
	)
	if hydrationErr != nil || len(hydratedName) == 0 {
		return fmt.Errorf("cannot hydrate string %w", hydrationErr)
	}

	hydratedName = strings.TrimSpace(strings.ReplaceAll(hydratedName, "@", ""))

	var errWg errgroup.Group
	twitchClient, twitchClientErr := c.getHelixChannelApiClient(ctx, data.ChannelID)
	if twitchClientErr != nil {
		return twitchClientErr
	}

	var user helix.User
	errWg.Go(
		func() error {
			u, err := c.getHelixUserByLogin(twitchClient, hydratedName)
			if err != nil {
				return err
			}
			user = u
			return nil
		},
	)

	var mods []helix.Moderator

	errWg.Go(
		func() error {
			m, err := c.getChannelMods(twitchClient, data.ChannelID)
			if err != nil {
				return err
			}
			mods = m
			return nil
		},
	)

	var dbChannel deprecatedgormmodel.Channels

	errWg.Go(
		func() error {
			ch, err := c.getChannelDbEntity(ctx, data.ChannelID)
			if err != nil {
				return err
			}
			dbChannel = ch
			return nil
		},
	)

	if err := errWg.Wait(); err != nil {
		return err
	}

	if user.ID == dbChannel.BotID || user.ID == dbChannel.ID {
		return errors.New("cannot mod/unmod bot")
	}

	var isAlreadyMod bool
	for _, item := range mods {
		if item.UserID == user.ID {
			isAlreadyMod = true
			break
		}
	}

	if operation.Type == model.EventOperationTypeMod {
		if isAlreadyMod {
			return nil
		}

		resp, err := twitchClient.AddChannelModerator(
			&helix.AddChannelModeratorParams{
				BroadcasterID: data.ChannelID,
				UserID:        user.ID,
			},
		)
		if err != nil {
			return err
		}
		if resp.ErrorMessage != "" {
			return errors.New(resp.ErrorMessage)
		}
	} else {
		if !isAlreadyMod {
			return errors.New("not a mod")
		}

		resp, err := twitchClient.RemoveChannelModerator(
			&helix.RemoveChannelModeratorParams{
				BroadcasterID: data.ChannelID,
				UserID:        user.ID,
			},
		)
		if err != nil {
			return err
		}
		if resp.ErrorMessage != "" {
			return errors.New(resp.ErrorMessage)
		}
	}

	return nil
}

func (c *Activity) UnmodRandom(
	ctx context.Context,
	_ model.EventOperation,
	data shared.EventData,
) error {
	activity.RecordHeartbeat(ctx, nil)

	dbChannel, dbChannelErr := c.getChannelDbEntity(ctx, data.ChannelID)
	if dbChannelErr != nil {
		return dbChannelErr
	}

	twitchClient, twitchClientErr := c.getHelixChannelApiClient(ctx, data.ChannelID)
	if twitchClientErr != nil {
		return twitchClientErr
	}

	mods, modsErr := c.getChannelMods(twitchClient, data.ChannelID)
	if modsErr != nil {
		return modsErr
	}

	// choose random mod, but filter out bot account
	filteredMods := lo.Filter(
		mods, func(item helix.Moderator, index int) bool {
			return item.UserID != dbChannel.BotID
		},
	)
	randomMod := lo.Sample(filteredMods)

	removeReq, err := twitchClient.RemoveChannelModerator(
		&helix.RemoveChannelModeratorParams{
			BroadcasterID: data.ChannelID,
			UserID:        randomMod.UserID,
		},
	)
	if err != nil {
		return err
	}
	if removeReq.ErrorMessage != "" {
		return errors.New(removeReq.ErrorMessage)
	}

	// if len(c.data.PrevOperation.UnmodedUserName) > 0 {
	// 	c.data.PrevOperation.UnmodedUserName += ", " + randomMod.UserName
	// } else {
	// 	c.data.PrevOperation.UnmodedUserName = randomMod.UserName
	// }

	return nil
}
