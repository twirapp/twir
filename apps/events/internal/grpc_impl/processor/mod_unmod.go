package processor

import (
	"errors"
	"fmt"
	"github.com/samber/lo"
	"github.com/satont/go-helix/v2"
	model "github.com/satont/tsuwari/libs/gomodels"
	"strings"
)

func (c *Processor) getChannelMods() ([]helix.Moderator, error) {
	if c.cache.channelModerators != nil {
		return c.cache.channelModerators, nil
	}

	mods, err := c.streamerApiClient.GetModerators(&helix.GetModeratorsParams{
		BroadcasterID: c.channelId,
	})

	if err != nil {
		return nil, err
	}

	if mods.ErrorMessage != "" {
		return nil, errors.New(mods.ErrorMessage)
	}

	c.cache.channelModerators = mods.Data.Moderators

	cursor := ""
	if mods.Data.Pagination.Cursor != "" {
		for {
			mods, err = c.streamerApiClient.GetModerators(&helix.GetModeratorsParams{
				BroadcasterID: c.channelId,
				After:         cursor,
			})

			if err != nil {
				return nil, err
			}

			if mods.ErrorMessage != "" {
				return nil, errors.New(mods.ErrorMessage)
			}

			c.cache.channelModerators = append(c.cache.channelModerators, mods.Data.Moderators...)

			if mods.Data.Pagination.Cursor == "" {
				break
			}

			cursor = mods.Data.Pagination.Cursor
		}
	}

	return mods.Data.Moderators, nil
}

func (c *Processor) ModOrUnmod(input string, operation model.EventOperationType) error {
	hydratedName, err := c.HydrateStringWithData(input, c.data)

	if err != nil || len(hydratedName) == 0 {
		return fmt.Errorf("cannot hydrate string %w", err)
	}

	hydratedName = strings.TrimSpace(strings.ReplaceAll(hydratedName, "@", ""))

	user, err := c.streamerApiClient.GetUsers(&helix.UsersParams{
		Logins: []string{hydratedName},
	})

	if err != nil || len(user.Data.Users) == 0 {
		if err != nil {
			return err
		}
		return errors.New("cannot get user")
	}

	mods, err := c.getChannelMods()
	if err != nil {
		return err
	}

	dbChannel, err := c.getDbChannel()
	if err != nil {
		return err
	}

	if user.Data.Users[0].ID == dbChannel.BotID {
		return InternalError
	}

	isAlreadyMod := lo.SomeBy(mods, func(item helix.Moderator) bool {
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
		} else {
			c.cache.channelModerators = append(c.cache.channelModerators, helix.Moderator{
				UserID:    user.Data.Users[0].ID,
				UserLogin: user.Data.Users[0].Login,
				UserName:  user.Data.Users[0].DisplayName,
			})
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
		} else {
			c.cache.channelModerators = lo.Filter(c.cache.channelModerators, func(item helix.Moderator, index int) bool {
				return item.UserID != user.Data.Users[0].ID
			})
		}
	}

	return nil
}

func (c *Processor) UnmodRandom() error {
	dbChannel, err := c.getDbChannel()
	if err != nil {
		return err
	}

	mods, err := c.getChannelMods()
	if err != nil {
		return err
	}

	// choose random mod, but filter out bot account
	filteredMods := lo.Filter(mods, func(item helix.Moderator, index int) bool {
		return item.UserID != dbChannel.BotID
	})
	randomMod := lo.Sample(filteredMods)

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

	c.cache.channelModerators = lo.Filter(c.cache.channelModerators, func(item helix.Moderator, index int) bool {
		return item.UserID != randomMod.UserID
	})

	if len(c.data.PrevOperation.UnmodedUserName) > 0 {
		c.data.PrevOperation.UnmodedUserName += ", " + randomMod.UserName
	} else {
		c.data.PrevOperation.UnmodedUserName = randomMod.UserName
	}

	return nil
}
