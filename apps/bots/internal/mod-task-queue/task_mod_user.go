package mod_task_queue

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/goccy/go-json"
	"github.com/hibiken/asynq"
	"github.com/kvizyx/twitchy/helix"
	"github.com/twirapp/twir/libs/twitch"
)

func (p *RedisTaskProcessor) ProcessDistributeMod(
	ctx context.Context,
	task *asynq.Task,
) error {
	var payload TaskModUserPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return err
	}

	twitchClient, err := twitch.NewUserClientWithContext(
		ctx,
		payload.TwitchUserID,
		p.config,
		p.twirBus,
	)
	if err != nil {
		return err
	}

	checkModReq, err := twitchClient.Moderation.GetModerators(
		ctx,
		helix.GetModeratorsRequest{
			BroadcasterID: payload.ChannelID,
			UserIDs:       []string{payload.UserID},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to check existing moderator: %w", err)
	}

	if len(checkModReq.Data) > 0 {
		p.logger.Warn(
			"user is already a moderator",
			slog.String("channelId", payload.ChannelID),
			slog.String("userId", payload.UserID),
		)
		return nil
	}

	_, err = twitchClient.Moderation.AddChannelModerator(
		ctx,
		helix.AddChannelModeratorRequest{
			BroadcasterID: payload.ChannelID,
			UserID:        payload.UserID,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to add moderator: %w", err)
	}

	return nil
}

func (d *ModTaskDistributor) DistributeModUser(
	ctx context.Context,
	payload *TaskModUserPayload,
	opts ...asynq.Option,
) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	task := asynq.NewTask(TaskModUser, jsonPayload, opts...)
	info, err := d.client.EnqueueContext(ctx, task)
	if err != nil {
		return err
	}

	d.logger.Info("task sent", slog.String("id", info.ID))

	return nil
}
