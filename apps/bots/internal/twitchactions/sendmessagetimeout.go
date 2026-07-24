package twitchactions

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	channelentity "github.com/twirapp/twir/libs/entities/channel"
)

func (c *TwitchActions) timeoutFromMessage(ctx context.Context, channel channelentity.Channel, opts SendMessageOpts) error {
	twitchBinding, botConfig, found, err := channel.TwitchBinding()
	if err != nil {
		return fmt.Errorf("cannot parse Twitch bot config: %w", err)
	}
	if !found || !twitchBinding.Enabled || twitchBinding.PlatformChannelID == "" ||
		!botConfig.IsBotMod || botConfig.IsTwitchBanned || botConfig.BotID == "" {
		return nil
	}
	if twitchBinding.PlatformChannelID != opts.BroadcasterID {
		return fmt.Errorf("Twitch binding channel id does not match broadcaster %s", opts.BroadcasterID)
	}

	splittedMsg := strings.Fields(opts.Message)

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

	if twitchUser == nil || twitchUser.ID == "" || twitchUser.ID == botConfig.BotID {
		return nil
	}

	return c.Ban(
		ctx,
		BanOpts{
			BroadcasterID:  twitchBinding.PlatformChannelID,
			UserID:         twitchUser.ID,
			ModeratorID:    botConfig.BotID,
			Duration:       duration,
			IsModerator:    false,
			AddModAfterBan: false,
			Reason:         reason,
		},
	)
}
