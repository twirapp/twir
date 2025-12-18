package channel

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/twirapp/twir/apps/bots/internal/twitchactions"
	"github.com/twirapp/twir/libs/bus-core/bots"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/logger"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/sync/errgroup"
)

var ErrNoChannelId = errors.New("channel id is not provided")

func (s *Service) SendMessage(ctx context.Context, req bots.SendMessageRequest) error {
	return s.workersPool.SubmitErr(
		func() error {
			span := trace.SpanFromContext(ctx)
			defer span.End()

			span.SetAttributes(
				attribute.String("channel_id", req.ChannelId),
				attribute.String("message", req.Message),
				attribute.String("reply_to", req.ReplyTo),
			)

			if req.ChannelId == "" {
				return ErrNoChannelId
			}

			err := s.twitchActions.SendMessage(
				ctx,
				twitchactions.SendMessageOpts{
					BroadcasterID:        req.ChannelId,
					SenderID:             "",
					Message:              req.Message,
					ReplyParentMessageID: req.ReplyTo,
					IsAnnounce:           req.IsAnnounce,
					SkipToxicityCheck:    req.SkipToxicityCheck,
					SkipRateLimits:       req.SkipRateLimits,
					AnnounceColor:        req.AnnounceColor,
				},
			)
			if err != nil {
				s.logger.Error("cannot send message", logger.Error(err))
				return err
			}

			return nil
		},
	).Wait()
}

func (s *Service) DeleteMessage(ctx context.Context, req bots.DeleteMessageRequest) error {
	return s.workersPool.SubmitErr(
		func() error {
			var channel model.Channels
			err := s.gorm.WithContext(ctx).Where("id = ?", req.ChannelId).Find(&channel).Error
			if err != nil {
				s.logger.Error(
					"cannot find channel to delete messages from",
					slog.String("channelId", req.ChannelId),
				)
				return fmt.Errorf("find channel: %w", err)
			}

			if channel.ID == "" {
				return nil
			}

			wg, _ := errgroup.WithContext(ctx)

			for _, msgId := range req.MessageIds {
				wg.Go(
					func() error {
						deleteErr := s.twitchActions.DeleteMessage(
							ctx,
							twitchactions.DeleteMessageOpts{
								BroadcasterID: req.ChannelId,
								ModeratorID:   channel.BotID,
								MessageID:     msgId,
							},
						)
						if deleteErr != nil {
							s.logger.Error("cannot delete message", logger.Error(deleteErr))
							return fmt.Errorf("delete message: %w", deleteErr)
						}

						return nil
					},
				)
			}

			return wg.Wait()
		}).Wait()
}
