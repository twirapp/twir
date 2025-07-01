package mod_task_queue

import (
	"context"
	"errors"
	"log/slog"

	"github.com/goccy/go-json"
	"github.com/hibiken/asynq"
	"github.com/nicklaw5/helix/v2"
	"github.com/satont/twir/libs/twitch"
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

	req, err := twitchClient.AddChannelModerator(
		&helix.AddChannelModeratorParams{
			BroadcasterID: payload.ChannelID,
			UserID:        payload.UserID,
		},
	)
	if req.ErrorMessage != "" {
		return errors.New(req.ErrorMessage)
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
