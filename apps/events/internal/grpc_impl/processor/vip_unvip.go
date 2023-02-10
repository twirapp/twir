package processor

import (
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
