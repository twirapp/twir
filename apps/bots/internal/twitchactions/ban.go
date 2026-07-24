package twitchactions

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/nicklaw5/helix/v2"
	kvoptions "github.com/twirapp/kv/options"
	mod_task_queue "github.com/twirapp/twir/apps/bots/internal/mod-task-queue"
	"github.com/twirapp/twir/libs/logger"
	"github.com/twirapp/twir/libs/redis_keys"
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

type twitchUserClientFactory func(context.Context, uuid.UUID) (*helix.Client, error)

type twitchBotClientFactory func(context.Context, string) (*helix.Client, error)

func (c *TwitchActions) createUserClient(ctx context.Context, userID uuid.UUID) (*helix.Client, error) {
	if c.newUserClient != nil {
		return c.newUserClient(ctx, userID)
	}

	return twitch.NewUserClientWithContext(ctx, userID, c.config, c.twirBus)
}

func (c *TwitchActions) createBotClient(ctx context.Context, botID string) (*helix.Client, error) {
	if c.newBotClient != nil {
		return c.newBotClient(ctx, botID)
	}

	return twitch.NewBotClientWithContext(ctx, botID, c.config, c.twirBus)
}

func (c *TwitchActions) Ban(ctx context.Context, opts BanOpts) error {
	channel, err := c.channelsByTwitchIDCache.Get(ctx, opts.BroadcasterID)
	if err != nil {
		return fmt.Errorf("cannot get channel by twitch id: %w", err)
	}
	twitchBinding, botConfig, found, err := channel.TwitchBinding()
	if err != nil {
		return fmt.Errorf("cannot parse Twitch bot config: %w", err)
	}
	if !found || !twitchBinding.Enabled || twitchBinding.PlatformChannelID == "" {
		return fmt.Errorf("channel has no enabled Twitch binding for broadcaster %s", opts.BroadcasterID)
	}
	if twitchBinding.PlatformChannelID != opts.BroadcasterID {
		return fmt.Errorf("Twitch binding channel id does not match broadcaster %s", opts.BroadcasterID)
	}
	if !botConfig.IsBotMod || botConfig.IsTwitchBanned {
		return nil
	}
	if botConfig.BotID == "" {
		return fmt.Errorf("channel has no Twitch bot id for broadcaster %s", opts.BroadcasterID)
	}

	twitchUserID := twitchBinding.UserID
	moderatorID := botConfig.BotID

	broadcasterHelixClient, err := c.createUserClient(ctx, twitchUserID)
	if err != nil {
		return fmt.Errorf("cannot create helix client: %w", err)
	}

	botHelixClient, err := c.createBotClient(ctx, moderatorID)
	if err != nil {
		c.logger.Error("cannot create helix client", logger.Error(err))
		return fmt.Errorf("cannot create helix client: %w", err)
	}

	if opts.IsModerator && opts.AddModAfterBan {
		err := c.modTaskDistributor.DistributeModUser(
			ctx,
			&mod_task_queue.TaskModUserPayload{
				ChannelID:    opts.BroadcasterID,
				TwitchUserID: twitchUserID,
				UserID:       opts.UserID,
			}, asynq.ProcessIn(time.Duration(opts.Duration+1)*time.Second),
		)
		if err != nil {
			return fmt.Errorf("cannot distribute mod user: %w", err)
		}

		// we'll listen unban event via eventsub and track that key for faster processing of mod user
		if err := c.kv.Set(
			ctx,
			redis_keys.CreateDistributedModTaskKey(opts.BroadcasterID, opts.UserID),
			true,
			kvoptions.WithExpire(time.Duration(opts.Duration+5)*time.Second),
		); err != nil {
			return fmt.Errorf(
				"cannot prepare distributed mod task, so we better not to ban user: %w",
				err,
			)
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
			ModeratorId:   moderatorID,
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
