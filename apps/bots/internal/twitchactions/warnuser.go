package twitchactions

import (
	"context"
	"errors"

	"github.com/nicklaw5/helix/v2"
)

type WarnUserOpts struct {
	BroadcasterID string
	ModeratorID   string
	UserID        string
	Reason        string
}

func (c *TwitchActions) WarnUser(ctx context.Context, opts WarnUserOpts) error {
	return c.withBotClient(
		ctx,
		opts.ModeratorID,
		func(client *helix.Client) (int, error) {
			resp, err := client.SendModeratorWarnMessage(
				&helix.SendModeratorWarnChatMessageParams{
					BroadcasterID: opts.BroadcasterID,
					ModeratorID:   opts.ModeratorID,
					Body: helix.SendModeratorWarnMessageRequestBody{
						UserID: opts.UserID,
						Reason: opts.Reason,
					},
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
