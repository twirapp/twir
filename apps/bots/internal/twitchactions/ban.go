package twitchactions

import (
	"context"
	"errors"

	"github.com/nicklaw5/helix/v2"
	"github.com/satont/twir/libs/twitch"
)

type BanOpts struct {
	Duration      int
	Reason        string
	BroadcasterID string
	UserID        string
	ModeratorID   string
}

func (c *TwitchActions) Ban(ctx context.Context, opts BanOpts) error {
	twitchClient, err := twitch.NewBotClientWithContext(ctx, opts.ModeratorID, c.config, c.tokensGrpc)
	if err != nil {
		return err
	}

	resp, err := twitchClient.BanUser(
		&helix.BanUserParams{
			BroadcasterID: opts.BroadcasterID,
			ModeratorId:   opts.ModeratorID,
			Body: helix.BanUserRequestBody{
				Duration: opts.Duration,
				Reason:   opts.Reason,
				UserId:   opts.UserID,
			},
		},
	)
	if err != nil {
		return err
	}

	if resp.ErrorMessage != "" {
		return errors.New(resp.ErrorMessage)
	}

	return nil
}
