package processor

import (
	"github.com/satont/go-helix/v2"
	model "github.com/satont/tsuwari/libs/gomodels"
)

func (c *Processor) BanRandom() {
	randomOnlineUser := &model.UsersOnline{}
	err := c.services.DB.Where(`"channelId" = ?`, c.channelId).Find(&randomOnlineUser).Error
	if err != nil {
		c.services.Logger.Sugar().Error(err)
		return
	}

	if randomOnlineUser == nil || !randomOnlineUser.UserId.Valid {
		return
	}

	c.streamerApiClient.BanUser(&helix.BanUserParams{
		BroadcasterID: c.channelId,
		ModeratorId:   c.channelId,
		Body: helix.BanUserRequestBody{
			Duration: 0,
			Reason:   "randomly banned from twirapp",
			UserId:   randomOnlineUser.UserId.String,
		},
	})
}
