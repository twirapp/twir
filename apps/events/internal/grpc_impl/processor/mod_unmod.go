package processor

import (
	"errors"
	"fmt"
	"github.com/samber/lo"
	"github.com/satont/go-helix/v2"
	model "github.com/satont/tsuwari/libs/gomodels"
)

func (c *Processor) ModOrUnmod(input string, operation model.EventOperationType) error {
	hydratedName, err := hydrateStringWithData(input, c.data)

	if err != nil || len(hydratedName) == 0 {
		return fmt.Errorf("cannot hydrate string %w", err)
	}

	user, err := c.streamerApiClient.GetUsers(&helix.UsersParams{
		Logins: []string{hydratedName},
	})

	if err != nil || len(user.Data.Users) == 0 {
		if err != nil {
			return err
		}
		return errors.New("cannot get user")
	}

	mods, err := c.streamerApiClient.GetModerators(&helix.GetModeratorsParams{
		BroadcasterID: c.channelId,
	})

	if err != nil {
		return err
	}

	if mods.ErrorMessage != "" {
		return errors.New(mods.ErrorMessage)
	}

	isAlreadyMod := lo.SomeBy(mods.Data.Moderators, func(item helix.Moderator) bool {
		return item.UserID == user.Data.Users[0].ID
	})

	if operation == "MOD" {
		if isAlreadyMod {
			return errors.New("already mod")
		}

		resp, err := c.streamerApiClient.AddChannelModerator(&helix.AddChannelModeratorParams{
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
		if !isAlreadyMod {
			return errors.New("not a mod")
		}

		resp, err := c.streamerApiClient.RemoveChannelModerator(&helix.RemoveChannelModeratorParams{
			BroadcasterID: c.channelId,
			UserID:        user.Data.Users[0].ID,
		})
		if resp.ErrorMessage != "" || err != nil {
			if err != nil {
				return err
			}

			return errors.New(resp.ErrorMessage)
		}
	}

	return nil
}

func (c *Processor) UnmodRandom() error {
	channel := model.Channels{}
	c.services.DB.Where(`"id" = ?`, c.channelId).Find(&channel)
	if channel.ID == "" {
		return errors.New("cannot get channel")
	}

	mods, err := c.streamerApiClient.GetModerators(&helix.GetModeratorsParams{
		BroadcasterID: c.channelId,
	})

	if err != nil {
		return err
	}

	if mods.ErrorMessage != "" {
		return errors.New(mods.ErrorMessage)
	}

	if len(mods.Data.Moderators) == 0 {
		return errors.New("cannot get mods")
	}

	// choose random mod, but filter out bot account
	randomMod := lo.Sample(lo.Filter(mods.Data.Moderators, func(item helix.Moderator, index int) bool {
		return item.UserID != channel.BotID
	}))

	removeReq, err := c.streamerApiClient.RemoveChannelModerator(&helix.RemoveChannelModeratorParams{
		BroadcasterID: c.channelId,
		UserID:        randomMod.UserID,
	})

	if err != nil {
		return err
	}

	if removeReq.ErrorMessage != "" {
		return errors.New(removeReq.ErrorMessage)
	}

	if len(c.data.PrevOperation.UnmodedUserName) > 0 {
		c.data.PrevOperation.UnmodedUserName += ", " + randomMod.UserName
	} else {
		c.data.PrevOperation.UnmodedUserName = randomMod.UserName
	}

	return nil
}
