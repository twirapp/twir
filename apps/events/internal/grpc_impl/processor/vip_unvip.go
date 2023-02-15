package processor

import (
	"errors"
	"github.com/samber/lo"
	"github.com/satont/go-helix/v2"
	model "github.com/satont/tsuwari/libs/gomodels"
)

func (c *Processor) VipOrUnvip(operation model.EventOperationType) error {
	user, err := c.streamerApiClient.GetUsers(&helix.UsersParams{
		Logins: []string{c.data.UserName},
	})

	if err != nil || len(user.Data.Users) == 0 {
		if err != nil {
			return err
		}
		return errors.New("cannot find user")
	}

	vips, err := c.streamerApiClient.GetChannelVips(&helix.GetChannelVipsParams{
		BroadcasterID: c.channelId,
	})

	if err != nil || vips.ErrorMessage != "" {
		if err != nil {
			return err
		}
		return errors.New(vips.ErrorMessage)
	}

	if len(vips.Data.ChannelsVips) == 0 {
		return errors.New("cannot get vips")
	}

	isAlreadyVip := lo.SomeBy(vips.Data.ChannelsVips, func(item helix.ChannelVips) bool {
		return item.UserID == user.Data.Users[0].ID
	})

	if operation == "VIP" {
		if isAlreadyVip {
			return errors.New("user already vip")
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
			return errors.New("not a vip")
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
	channel := model.Channels{}
	err := c.services.DB.Where(`"id" = ?`, c.channelId).Find(&channel).Error
	if channel.ID == "" {
		return err
	}

	vips, err := c.streamerApiClient.GetChannelVips(&helix.GetChannelVipsParams{
		BroadcasterID: c.channelId,
	})

	if err != nil || vips.ErrorMessage != "" {
		if err != nil {
			return err
		}
		return errors.New(vips.ErrorMessage)
	}

	if len(vips.Data.ChannelsVips) == 0 {
		return errors.New("cannot get vips")
	}

	// choose random vip, but filter out bot account
	randomVip := lo.Sample(lo.Filter(vips.Data.ChannelsVips, func(item helix.ChannelVips, index int) bool {
		return item.UserID != channel.BotID
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
