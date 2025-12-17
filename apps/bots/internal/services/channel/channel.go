package channel

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/twirapp/twir/apps/bots/internal/twitchactions"
	"github.com/twirapp/twir/apps/bots/internal/workers"
	"github.com/twirapp/twir/libs/bus-core/bots"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/logger"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

var (
	ErrFriendlyFire = errors.New("friendly fire")
	ErrNoChannelId  = errors.New("channel id is not provided")
)

type Opts struct {
	fx.In

	Logger        *slog.Logger
	Gorm          *gorm.DB
	TwitchActions *twitchactions.TwitchActions
	WorkersPool   *workers.Pool
}

func New(opts Opts) *Service {
	return &Service{
		gorm:          opts.Gorm,
		logger:        opts.Logger,
		twitchActions: opts.TwitchActions,
		workersPool:   opts.WorkersPool,
	}
}

type Service struct {
	logger        *slog.Logger
	gorm          *gorm.DB
	twitchActions *twitchactions.TwitchActions
	workersPool   *workers.Pool
}

func (s *Service) Ban(ctx context.Context, req bots.BanRequest) error {
	return s.workersPool.SubmitErr(
		func() error {
			if req.ChannelID == req.UserID {
				return ErrFriendlyFire
			}

			channelEntity := model.Channels{}
			if err := s.gorm.WithContext(ctx).Where(
				`"id" = ?`,
				req.ChannelID,
			).First(&channelEntity).Error; err != nil {
				s.logger.Error("cannot get channel entity", logger.Error(err))
				return fmt.Errorf("get channel entity: %w", err)
			}

			if err := s.twitchActions.Ban(
				ctx,
				twitchactions.BanOpts{
					Duration:       req.BanTime,
					Reason:         req.Reason,
					BroadcasterID:  req.ChannelID,
					UserID:         req.UserID,
					ModeratorID:    channelEntity.BotID,
					IsModerator:    req.IsModerator,
					AddModAfterBan: req.AddModAfterBan,
				},
			); err != nil {
				s.logger.Error("cannot ban user", logger.Error(err))
				return err
			}

			return nil
		},
	).Wait()
}

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
