package processor

import (
	"errors"
	"fmt"
	"github.com/samber/lo"
	"github.com/satont/go-helix/v2"
	model "github.com/satont/tsuwari/libs/gomodels"
	"strings"
)

func (c *Processor) getChannelVips() ([]helix.ChannelVips, error) {
	if c.cache.channelVips != nil {
		return c.cache.channelVips, nil
	}

	vips, err := c.streamerApiClient.GetChannelVips(&helix.GetChannelVipsParams{
		BroadcasterID: c.channelId,
	})

	if err != nil {
		return nil, errors.New(vips.ErrorMessage)
	}

	if vips.ErrorMessage != "" {
		return nil, errors.New(vips.ErrorMessage)
	}

	if len(vips.Data.ChannelsVips) == 0 {
		return nil, errors.New("cannot get vips")
	}

	c.cache.channelVips = vips.Data.ChannelsVips

	return vips.Data.ChannelsVips, nil
}

func (c *Processor) VipOrUnvip(input string, operation model.EventOperationType) error {
	hydratedName, err := hydrateStringWithData(input, c.data)

	if err != nil || len(hydratedName) == 0 {
		return fmt.Errorf("cannot hydrate string %w", err)
	}

	hydratedName = strings.ReplaceAll(hydratedName, "@", "")

	user, err := c.streamerApiClient.GetUsers(&helix.UsersParams{
		Logins: []string{hydratedName},
	})

	if err != nil || len(user.Data.Users) == 0 {
		if err != nil {
			return err
		}
		return errors.New("cannot find user")
	}

	vips, err := c.getChannelVips()
	if err != nil {
		return err
	}

	dbChannel, err := c.getDbChannel()
	if err != nil {
		return err
	}

	if user.Data.Users[0].ID == dbChannel.BotID || user.Data.Users[0].ID == dbChannel.ID {
		return InternalError
	}

	mods, err := c.getChannelMods()
	if err != nil {
		return err
	}

	if lo.SomeBy(mods, func(item helix.Moderator) bool {
		return item.UserID == user.Data.Users[0].ID
	}) {
		return InternalError
	}

	isAlreadyVip := lo.SomeBy(vips, func(item helix.ChannelVips) bool {
		return item.UserID == user.Data.Users[0].ID
	})

	if operation == "VIP" {
		if isAlreadyVip {
			return InternalError
		}

		resp, err := c.streamerApiClient.AddChannelVip(&helix.AddChannelVipParams{
			BroadcasterID: c.channelId,
			UserID:        user.Data.Users[0].ID,
		})
		if resp.ErrorMessage != "" || err != nil {
			if err != nil {
				return err
			} else {
				return errors.New(resp.ErrorMessage)
			}
		}
	} else {
		if !isAlreadyVip {
			return InternalError
		}
		resp, err := c.streamerApiClient.RemoveChannelVip(&helix.RemoveChannelVipParams{
			BroadcasterID: c.channelId,
			UserID:        user.Data.Users[0].ID,
		})
		if resp.ErrorMessage != "" || err != nil {
			if err != nil {
				return err
			} else {
				return errors.New(resp.Error)
			}
		}
	}

	return nil
}

func (c *Processor) UnvipRandom() error {
	vips, err := c.getChannelVips()
	if err != nil {
		return err
	}

	dbChannel, err := c.getDbChannel()
	if err != nil {
		return err
	}

	// choose random vip, but filter out bot account
	randomVip := lo.Sample(lo.Filter(vips, func(item helix.ChannelVips, index int) bool {
		return item.UserID != dbChannel.BotID
	}))
	removeReq, err := c.streamerApiClient.RemoveChannelVip(&helix.RemoveChannelVipParams{
		BroadcasterID: c.channelId,
		UserID:        randomVip.UserID,
	})

	if err != nil {
		return err
	}

	if removeReq.ErrorMessage != "" {
		return errors.New(removeReq.ErrorMessage)
	}

	if len(c.data.PrevOperation.UnvipedUserName) > 0 {
		c.data.PrevOperation.UnvipedUserName += ", " + randomVip.UserName
	} else {
		c.data.PrevOperation.UnvipedUserName = randomVip.UserName
	}

	return nil
}
