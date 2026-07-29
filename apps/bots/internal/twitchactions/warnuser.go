package twitchactions

import (
	"context"

	"github.com/kvizyx/twitchy/helix"
)

type WarnUserOpts struct {
	BroadcasterID string
	ModeratorID   string
	UserID        string
	Reason        string
}

func (c *TwitchActions) WarnUser(ctx context.Context, opts WarnUserOpts) error {
	twitchClient, err := c.createChannelBotClient(ctx, opts.ModeratorID, opts.BroadcasterID)
	if err != nil {
		return err
	}

	_, err = twitchClient.Moderation.WarnChatUser(
		ctx,
		helix.WarnChatUserRequest{
			BroadcasterID: opts.BroadcasterID,
			ModeratorID:   opts.ModeratorID,
			Data: helix.WarnChatUserBody{
				UserID: opts.UserID,
				Reason: opts.Reason,
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}
