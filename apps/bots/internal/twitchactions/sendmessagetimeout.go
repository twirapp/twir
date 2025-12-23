package twitchactions

import (
	"context"
	"strconv"
	"strings"

	"github.com/twirapp/twir/libs/repositories/channels/model"
)

func (c *TwitchActions) timeoutFromMessage(ctx context.Context, channel model.Channel, opts SendMessageOpts) error {
	splittedMsg := strings.Fields(opts.Message)

	// /timeout [username] [duration in seconds] [reason]
	var (
		userName = splittedMsg[1]
		reason   string
		duration = 600
	)

	if len(splittedMsg) >= 3 {
		d, err := strconv.Atoi(splittedMsg[2])
		if err == nil {
			duration = d
		}
	}

	if len(splittedMsg) >= 4 {
		reason = strings.Join(splittedMsg[3:], " ")
	}

	if splittedMsg[0] == "/ban" {
		duration = 0
	}

	twitchUser, err := c.cachedTwitchClient.GetUserByName(ctx, userName)
	if err != nil {
		return err
	}

	if twitchUser == nil || twitchUser.ID == "" || twitchUser.ID == channel.BotID {
		return nil
	}

	return c.Ban(
		ctx,
		BanOpts{
			BroadcasterID:  channel.ID,
			UserID:         twitchUser.ID,
			ModeratorID:    channel.BotID,
			Duration:       duration,
			IsModerator:    false,
			AddModAfterBan: false,
			Reason:         reason,
		},
	)
}
