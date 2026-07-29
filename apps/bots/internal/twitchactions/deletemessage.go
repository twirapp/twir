package twitchactions

import (
	"context"

	"github.com/kvizyx/twitchy/helix"
)

type DeleteMessageOpts struct {
	BroadcasterID string
	ModeratorID   string
	MessageID     string
}

func (c *TwitchActions) DeleteMessage(ctx context.Context, opts DeleteMessageOpts) error {
	twitchClient, err := c.createChannelBotClient(ctx, opts.ModeratorID, opts.BroadcasterID)
	if err != nil {
		return err
	}

	_, err = twitchClient.Moderation.DeleteChatMessages(
		ctx,
		helix.DeleteChatMessagesRequest{
			BroadcasterID: opts.BroadcasterID,
			ModeratorID:   opts.ModeratorID,
			MessageID:     opts.MessageID,
		},
	)
	if err != nil {
		return err
	}

	return nil
}
