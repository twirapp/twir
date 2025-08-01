package twitchactions

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/hibiken/asynq"
	"github.com/nicklaw5/helix/v2"
	mod_task_queue "github.com/twirapp/twir/apps/bots/internal/mod-task-queue"
	"github.com/twirapp/twir/libs/twitch"
)

type BanOpts struct {
	Reason         string
	BroadcasterID  string
	UserID         string
	ModeratorID    string
	Duration       int
	IsModerator    bool
	AddModAfterBan bool
}

func (c *TwitchActions) Ban(ctx context.Context, opts BanOpts) error {
	broadcasterHelixClient, err := twitch.NewUserClientWithContext(
		ctx,
		opts.BroadcasterID,
		c.config,
		c.twirBus,
	)
	if err != nil {
		return fmt.Errorf("cannot create helix client: %w", err)
	}

	botHelixClient, err := twitch.NewBotClientWithContext(
		ctx,
		opts.ModeratorID,
		c.config,
		c.twirBus,
	)
	if err != nil {
		c.logger.Error("cannot create helix client", slog.Any("err", err))
		return fmt.Errorf("cannot create helix client: %w", err)
	}

	if opts.IsModerator && opts.AddModAfterBan {
		err := c.modTaskDistributor.DistributeModUser(
			ctx,
			&mod_task_queue.TaskModUserPayload{
				ChannelID: opts.BroadcasterID,
				UserID:    opts.UserID,
			}, asynq.ProcessIn(time.Duration(opts.Duration+2)*time.Second),
		)
		if err != nil {
			return fmt.Errorf("cannot distribute mod user: %w", err)
		}

		removeModeratorResponse, err := broadcasterHelixClient.RemoveChannelModerator(
			&helix.RemoveChannelModeratorParams{
				BroadcasterID: opts.BroadcasterID,
				UserID:        opts.UserID,
			},
		)
		if err != nil {
			return fmt.Errorf("cannot remove moderator: %w", err)
		}
		if removeModeratorResponse.ErrorMessage != "" {
			return errors.New(removeModeratorResponse.ErrorMessage)
		}
	}

	resp, err := botHelixClient.BanUser(
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
