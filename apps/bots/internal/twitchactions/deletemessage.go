package twitchactions

import (
	"context"
	"errors"

	"github.com/nicklaw5/helix/v2"
)

type DeleteMessageOpts struct {
	BroadcasterID string
	ModeratorID   string
	MessageID     string
}

func (c *TwitchActions) DeleteMessage(ctx context.Context, opts DeleteMessageOpts) error {
	return c.withBotClient(
		ctx,
		opts.ModeratorID,
		func(client *helix.Client) (int, error) {
			resp, err := client.DeleteChatMessage(
				&helix.DeleteChatMessageParams{
					BroadcasterID: opts.BroadcasterID,
					ModeratorID:   opts.ModeratorID,
					MessageID:     opts.MessageID,
				},
			)
			if err != nil {
				return 0, err
			}

			if resp.ErrorMessage != "" {
				return resp.StatusCode, errors.New(resp.ErrorMessage)
			}

			return resp.StatusCode, nil
		},
	)
}
