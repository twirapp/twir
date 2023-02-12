package processor

import (
	"github.com/samber/lo"
	"github.com/satont/go-helix/v2"
	model "github.com/satont/tsuwari/libs/gomodels"
)

func (c *Processor) ModOrUnmod(operation model.EventOperationType) {
	user, err := c.streamerApiClient.GetUsers(&helix.UsersParams{
		Logins: []string{c.data.UserName},
	})

	if err != nil || len(user.Data.Users) == 0 {
		if err != nil {
			c.services.Logger.Sugar().Error(err)
		}
		return
	}

	if operation == "MOD" {
		resp, err := c.streamerApiClient.AddChannelModerator(&helix.AddChannelModeratorParams{
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
		resp, err := c.streamerApiClient.RemoveChannelModerator(&helix.RemoveChannelModeratorParams{
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

func (c *Processor) UnmodRandom() {
	channel := model.Channels{}
	c.services.DB.Where(`"id" = ?`, c.channelId).Find(&channel)
	if channel.ID == "" {
		return
	}

	mods, err := c.streamerApiClient.GetModerators(&helix.GetModeratorsParams{
		BroadcasterID: c.channelId,
	})

	if err != nil || mods.ErrorMessage != "" {
		c.services.Logger.Sugar().Error(err, mods.ErrorMessage)
		return
	}

	if len(mods.Data.Moderators) == 0 {
		return
	}

	// choose random mod, but filter out bot account
	randomMod := lo.Sample(lo.Filter(mods.Data.Moderators, func(item helix.Moderator, index int) bool {
		return item.UserID != channel.BotID
	}))

	c.streamerApiClient.RemoveChannelModerator(&helix.RemoveChannelModeratorParams{
		BroadcasterID: c.channelId,
		UserID:        randomMod.UserID,
	})

	if len(c.data.PrevOperation.UnmodedUserName) > 0 {
		c.data.PrevOperation.UnmodedUserName += ", " + randomMod.UserName
	} else {
		c.data.PrevOperation.UnmodedUserName = randomMod.UserName
	}
}
