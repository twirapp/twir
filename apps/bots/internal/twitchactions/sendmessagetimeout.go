package twitchactions

import (
	"context"
	"strconv"
	"strings"

	"github.com/twirapp/twir/libs/repositories/channels/model"
)

func (c *TwitchActions) timeoutFromMessage(ctx context.Context, channel model.Channel, opts SendMessageOpts) error {
	splittedMsg := strings.Fields(opts.Message)
	var (
		durationStr string
		reason      string
		userName    = splittedMsg[1]
	)

	if len(splittedMsg) >= 2 {
		durationStr = splittedMsg[1]
	} else {
		durationStr = "600"
	}

	if len(splittedMsg) >= 3 {
		reason = strings.Join(splittedMsg[2:], " ")
	}

	duration, err := strconv.Atoi(durationStr)
	if err != nil {
		duration = 600
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
