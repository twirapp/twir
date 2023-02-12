package processor

import (
	"github.com/samber/lo"
	"github.com/satont/go-helix/v2"
	model "github.com/satont/tsuwari/libs/gomodels"
)

func (c *Processor) VipOrUnvip(operation model.EventOperationType) {
	user, err := c.streamerApiClient.GetUsers(&helix.UsersParams{
		Logins: []string{c.data.UserName},
	})

	if err != nil || len(user.Data.Users) == 0 {
		if err != nil {
			c.services.Logger.Sugar().Error(err)
		}
		return
	}

	if operation == "VIP" {
		resp, err := c.streamerApiClient.AddChannelVip(&helix.AddChannelVipParams{
			BroadcasterID: c.channelId,
			UserID:        user.Data.Users[0].ID,
		})
		if resp.ErrorMessage != "" || err != nil {
			if err != nil {
				c.services.Logger.Sugar().Error(err)
			} else {
				c.services.Logger.Sugar().Error(resp.ErrorMessage)
			}
		}
	} else {
		resp, err := c.streamerApiClient.RemoveChannelVip(&helix.RemoveChannelVipParams{
			BroadcasterID: c.channelId,
			UserID:        user.Data.Users[0].ID,
		})
		if resp.ErrorMessage != "" || err != nil {
			if err != nil {
				c.services.Logger.Sugar().Error(err)
			} else {
				c.services.Logger.Sugar().Error(resp.ErrorMessage)
			}
		}
	}
}

func (c *Processor) UnvipRandom() {
	channel := model.Channels{}
	c.services.DB.Where(`"id" = ?`, c.channelId).Find(&channel)
	if channel.ID == "" {
		return
	}

	vips, err := c.streamerApiClient.GetChannelVips(&helix.GetChannelVipsParams{
		BroadcasterID: c.channelId,
	})

	if err != nil || vips.ErrorMessage != "" || len(vips.Data.ChannelsVips) == 0 {
		return
	}

	// choose random vip, but filter out bot account
	randomVip := lo.Sample(lo.Filter(vips.Data.ChannelsVips, func(item helix.ChannelVips, index int) bool {
		return item.UserID != channel.BotID
	}))
	c.streamerApiClient.RemoveChannelVip(&helix.RemoveChannelVipParams{
		BroadcasterID: c.channelId,
		UserID:        randomVip.UserID,
	})

	if len(c.data.PrevOperation.UnmodedUserName) > 0 {
		c.data.PrevOperation.UnmodedUserName += ", " + randomVip.UserName
	} else {
		c.data.PrevOperation.UnmodedUserName = randomVip.UserName
	}
}
