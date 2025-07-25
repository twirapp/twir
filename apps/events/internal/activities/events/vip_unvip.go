package events

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/nicklaw5/helix/v2"
	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/events/internal/shared"
	model "github.com/twirapp/twir/libs/gomodels"
	"go.temporal.io/sdk/activity"
	"golang.org/x/sync/errgroup"
)

func (c *Activity) VipOrUnvip(
	ctx context.Context,
	operation model.EventOperation,
	data shared.EventData,
) error {
	activity.RecordHeartbeat(ctx, nil)

	dbChannel, dbChannelErr := c.getChannelDbEntity(ctx, data.ChannelID)
	if dbChannelErr != nil {
		return dbChannelErr
	}

	hydratedName, hydrateErr := c.hydrator.HydrateStringWithData(
		data.ChannelID,
		operation.Input.String,
		data,
	)
	if hydrateErr != nil {
		return fmt.Errorf("cannot hydrate string %w", hydrateErr)
	}
	if hydratedName == "" {
		return nil
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

	var vips []helix.ChannelVips
	errWg.Go(
		func() error {
			v, err := c.getChannelVips(twitchClient, data.ChannelID)
			if err != nil {
				return err
			}

			vips = v
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

	if err := errWg.Wait(); err != nil {
		return err
	}

	if user.ID == dbChannel.BotID || user.ID == dbChannel.ID {
		return errors.New("cannot vip/unvip bot")
	}

	var isMod bool
	for _, mod := range mods {
		if mod.UserID == user.ID {
			isMod = true
			break
		}
	}

	if isMod {
		return errors.New("cannot vip/unvip mod")
	}

	var isVip bool
	for _, vip := range vips {
		if vip.UserID == user.ID {
			isVip = true
			break
		}
	}

	if operation.Type == model.OperationVip {
		if isVip {
			return nil
		}

		resp, err := twitchClient.AddChannelVip(
			&helix.AddChannelVipParams{
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
		if !isVip {
			return nil
		}
		resp, err := twitchClient.RemoveChannelVip(
			&helix.RemoveChannelVipParams{
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

func (c *Activity) UnvipRandom(
	ctx context.Context,
	operation model.EventOperation,
	data shared.EventData,
) error {
	activity.RecordHeartbeat(ctx, nil)

	slots := operation.Input.String
	if operation.Type == model.OperationUnvipRandomIfNoSlots && slots == "" {
		return errors.New("input is empty")
	}

	twitchClient, twitchClientErr := c.getHelixChannelApiClient(context.TODO(), data.ChannelID)
	if twitchClientErr != nil {
		return twitchClientErr
	}

	vips, vipsErr := c.getChannelVips(twitchClient, data.ChannelID)
	if vipsErr != nil {
		return vipsErr
	}

	// if there is still slots available, we should skip unvip
	if operation.Type == model.OperationUnvipRandomIfNoSlots {
		slotsInt, err := strconv.Atoi(slots)
		if err != nil {
			return err
		}

		if len(vips) < slotsInt {
			return nil
		}
	}

	dbChannel, dbChannelErr := c.getChannelDbEntity(ctx, data.ChannelID)
	if dbChannelErr != nil {
		return dbChannelErr
	}

	// choose random vip, but filter out bot account
	randomVip := lo.Sample(
		lo.Filter(
			vips,
			func(item helix.ChannelVips, index int) bool {
				return item.UserID != dbChannel.BotID
			},
		),
	)

	removeReq, err := twitchClient.RemoveChannelVip(
		&helix.RemoveChannelVipParams{
			BroadcasterID: data.ChannelID,
			UserID:        randomVip.UserID,
		},
	)
	if err != nil {
		return err
	}
	if removeReq.ErrorMessage != "" {
		return errors.New(removeReq.ErrorMessage)
	}

	// if len(c.data.PrevOperation.UnvipedUserName) > 0 {
	// 	c.data.PrevOperation.UnvipedUserName += ", " + randomVip.UserName
	// } else {
	// 	c.data.PrevOperation.UnvipedUserName = randomVip.UserName
	// }

	return nil
}
