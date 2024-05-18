package bus_listener

import (
	"context"
	"log/slog"
	"time"

	"github.com/hibiken/asynq"
	"github.com/nicklaw5/helix/v2"
	mod_task_queue "github.com/satont/twir/apps/bots/internal/mod-task-queue"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/twitch"
	"github.com/twirapp/twir/libs/bus-core/bots"
)

func (c *BusListener) banUser(
	ctx context.Context,
	req bots.BanRequest,
) struct{} {
	if req.ChannelID == req.UserID {
		return struct{}{}
	}

	channelEntity := model.Channels{}
	if err := c.gorm.WithContext(ctx).Where(
		`"id" = ?`,
		req.ChannelID,
	).First(&channelEntity).Error; err != nil {
		c.logger.Error("cannot get channel entity", slog.Any("err", err))
		return struct{}{}
	}

	broadcasterHelixClient, err := twitch.NewUserClientWithContext(
		ctx,
		req.ChannelID,
		c.config,
		c.tokensGrpc,
	)
	if err != nil {
		c.logger.Error("cannot create helix client", slog.Any("err", err))
		return struct{}{}
	}

	botHelixClient, err := twitch.NewUserClientWithContext(
		ctx,
		channelEntity.BotID,
		c.config,
		c.tokensGrpc,
	)
	if err != nil {
		c.logger.Error("cannot create helix client", slog.Any("err", err))
		return struct{}{}
	}

	if req.IsModerator && req.AddModAfterBan {
		err := c.modTaskDistributor.DistributeModUser(
			ctx,
			&mod_task_queue.TaskModUserPayload{
				ChannelID: req.ChannelID,
				UserID:    req.UserID,
			}, asynq.ProcessIn(time.Duration(req.BanTime+2)*time.Second),
		)
		if err != nil {
			c.logger.Error("cannot distribute mod user", slog.Any("err", err))
			return struct{}{}
		}

		removeModeratorResponse, err := broadcasterHelixClient.RemoveChannelModerator(
			&helix.RemoveChannelModeratorParams{
				BroadcasterID: req.ChannelID,
				UserID:        req.UserID,
			},
		)
		if err != nil {
			c.logger.Error("cannot remove moderator", slog.Any("err", err))
			return struct{}{}
		}
		if removeModeratorResponse.ErrorMessage != "" {
			c.logger.Error(
				"cannot remove moderator",
				slog.Any("err", removeModeratorResponse.ErrorMessage),
			)
			return struct{}{}
		}
	}

	banUserResponse, err := botHelixClient.BanUser(
		&helix.BanUserParams{
			BroadcasterID: req.ChannelID,
			ModeratorId:   channelEntity.BotID,
			Body: helix.BanUserRequestBody{
				Duration: req.BanTime,
				Reason:   req.Reason,
				UserId:   req.UserID,
			},
		},
	)
	if err != nil {
		c.logger.Error("cannot ban user", slog.Any("err", err))
		return struct{}{}
	}
	if banUserResponse.ErrorMessage != "" {
		c.logger.Error(
			"cannot ban user",
			slog.Any("err", banUserResponse.ErrorMessage),
		)
		return struct{}{}
	}

	return struct{}{}
}
