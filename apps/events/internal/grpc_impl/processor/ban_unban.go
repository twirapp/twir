package processor

import (
	"github.com/satont/go-helix/v2"
	model "github.com/satont/tsuwari/libs/gomodels"
)

func (c *Processor) BanOrUnban(operation model.EventOperationType) {
	user, err := c.streamerApiClient.GetUsers(&helix.UsersParams{
		Logins: []string{c.data.UserName},
	})

	if err != nil || len(user.Data.Users) == 0 {
		if err != nil {
			c.services.Logger.Sugar().Error(err)
		}
		return
	}

	if operation == "BAN" {
		resp, err := c.streamerApiClient.BanUser(&helix.BanUserParams{
			BroadcasterID: c.channelId,
			ModeratorId:   c.channelId,
			Body: helix.BanUserRequestBody{
				Duration: 0,
				Reason:   "banned from twirapp",
				UserId:   user.Data.Users[0].ID,
			},
		})
		if resp.ErrorMessage != "" || err != nil {
			if err != nil {
				c.services.Logger.Sugar().Error(err)
			} else {
				c.services.Logger.Sugar().Error(resp.ErrorMessage)
			}
		}
	} else {
		resp, err := c.streamerApiClient.UnbanUser(&helix.UnbanUserParams{
			BroadcasterID: c.channelId,
			ModeratorID:   c.channelId,
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
