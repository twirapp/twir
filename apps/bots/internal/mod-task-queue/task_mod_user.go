package mod_task_queue

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/goccy/go-json"
	"github.com/hibiken/asynq"
	"github.com/nicklaw5/helix/v2"
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
		payload.ChannelID,
		p.config,
		p.twirBus,
	)
	if err != nil {
		return err
	}

	checkModReq, err := twitchClient.GetModerators(
		&helix.GetModeratorsParams{
			BroadcasterID: payload.ChannelID,
			UserIDs:       []string{payload.UserID},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to check existing moderator: %w", err)
	}
	if checkModReq.ErrorMessage != "" {
		return errors.New(checkModReq.ErrorMessage)
	}

	if len(checkModReq.Data.Moderators) > 0 {
		p.logger.Warn(
			"user is already a moderator",
			slog.String("channelId", payload.ChannelID),
			slog.String("userId", payload.UserID),
		)
		return nil
	}

	addModReq, err := twitchClient.AddChannelModerator(
		&helix.AddChannelModeratorParams{
			BroadcasterID: payload.ChannelID,
			UserID:        payload.UserID,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to add moderator: %w", err)
	}
	if addModReq.ErrorMessage != "" {
		return errors.New(addModReq.ErrorMessage)
	}

	return err
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
