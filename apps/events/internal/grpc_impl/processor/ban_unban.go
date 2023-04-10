package processor

import (
	"errors"
	"fmt"
	"strings"

	"github.com/nicklaw5/helix/v2"
	"github.com/samber/lo"
	model "github.com/satont/tsuwari/libs/gomodels"
)

func (c *Processor) Timeout(input string, timeoutTime int) error {
	hydratedName, err := c.HydrateStringWithData(input, c.data)

	if err != nil || len(hydratedName) == 0 {
		return fmt.Errorf("cannot hydrate string %w", err)
	}

	hydratedName = strings.TrimSpace(strings.ReplaceAll(hydratedName, "@", ""))

	user, err := c.streamerApiClient.GetUsers(&helix.UsersParams{
		Logins: []string{hydratedName},
	})

	if err != nil || user.ErrorMessage != "" || len(user.Data.Users) == 0 {
		if err != nil {
			return err
		}
		return fmt.Errorf("user not found")
	}

	mods, err := c.getChannelMods()
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

	if lo.SomeBy(mods, func(item helix.Moderator) bool {
		return item.UserID == user.Data.Users[0].ID
	}) {
		return InternalError
	}

	banReq, err := c.streamerApiClient.BanUser(&helix.BanUserParams{
		BroadcasterID: c.channelId,
		ModeratorId:   c.channelId,
		Body: helix.BanUserRequestBody{
			Duration: timeoutTime,
			Reason:   "banned from twirapp",
			UserId:   user.Data.Users[0].ID,
		},
	})

	if err != nil || banReq.ErrorMessage != "" {
		return fmt.Errorf("cannot ban user %w", err)
	}

	return nil
}

func (c *Processor) BanOrUnban(input string, operation model.EventOperationType) error {
	hydratedName, err := c.HydrateStringWithData(input, c.data)

	if err != nil || len(hydratedName) == 0 {
		return fmt.Errorf("cannot hydrate string %w", err)
	}

	hydratedName = strings.TrimSpace(strings.ReplaceAll(hydratedName, "@", ""))

	user, err := c.streamerApiClient.GetUsers(&helix.UsersParams{
		Logins: []string{hydratedName},
	})

	if err != nil || user.ErrorMessage != "" || len(user.Data.Users) == 0 {
		if err != nil {
			return err
		}
		return fmt.Errorf("cannot find user")
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

	if lo.SomeBy(mods, func(item helix.Moderator) bool {
		return item.UserID == user.Data.Users[0].ID
	}) {
		return InternalError
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
				return err
			} else {
				return errors.New(resp.ErrorMessage)
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
				return err
			} else {
				return errors.New(resp.ErrorMessage)
			}
		}
	}

	return nil
}

func (c *Processor) BanRandom(timeoutTime int) error {
	mods, err := c.getChannelMods()
	if err != nil {
		return err
	}

	dbChannel, err := c.getDbChannel()
	if err != nil {
		return err
	}

	excludedUsers := lo.Map(mods, func(item helix.Moderator, _ int) string {
		return item.UserID
	})

	excludedUsers = append(excludedUsers, dbChannel.ID, dbChannel.BotID)

	randomOnlineUser := &model.UsersOnline{}
	err = c.services.DB.
		Where(`"userId" not in ?`, excludedUsers).
		Order("random()").
		Find(randomOnlineUser).
		Error

	if err != nil {
		return err
	}

	if randomOnlineUser == nil || !randomOnlineUser.UserId.Valid {
		return errors.New("cannot get random user")
	}

	banReq, err := c.streamerApiClient.BanUser(&helix.BanUserParams{
		BroadcasterID: c.channelId,
		ModeratorId:   c.channelId,
		Body: helix.BanUserRequestBody{
			Duration: timeoutTime,
			Reason:   "randomly banned from twirapp",
			UserId:   randomOnlineUser.UserId.String,
		},
	})
	if err != nil {
		return err
	}

	if banReq.ErrorMessage != "" {
		return errors.New(banReq.ErrorMessage)
	}

	if len(c.data.PrevOperation.BannedUserName) > 0 {
		c.data.PrevOperation.BannedUserName += ", " + randomOnlineUser.UserName.String
	} else {
		c.data.PrevOperation.BannedUserName = randomOnlineUser.UserName.String
	}

	return nil
}
