package channel

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	botplatforms "github.com/twirapp/twir/apps/bots/internal/platforms"
	"github.com/twirapp/twir/apps/bots/internal/twitchactions"
	"github.com/twirapp/twir/libs/bus-core/bots"
	"github.com/twirapp/twir/libs/entities/platform"
	"github.com/twirapp/twir/libs/logger"
	channelsrepository "github.com/twirapp/twir/libs/repositories/channels"
	usersmodel "github.com/twirapp/twir/libs/repositories/users/model"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/sync/errgroup"
)

var ErrNoChannelId = errors.New("channel id is not provided")

type deleteMessageChannel struct {
	BotID         string
	BroadcasterID string
}

func (s *Service) SendMessage(ctx context.Context, req bots.SendMessageRequest) error {
	return s.workersPool.SubmitErr(
		func() error {
			span := trace.SpanFromContext(ctx)
			defer span.End()

			span.SetAttributes(
				attribute.String("channel_id", req.ChannelID.String()),
				attribute.String("message", req.Message),
				attribute.String("reply_to", req.ReplyTo),
			)

			if req.ChannelID == uuid.Nil {
				return ErrNoChannelId
			}

			platformsSliceAttribute := make([]string, len(req.Platforms))
			for i, p := range req.Platforms {
				platformsSliceAttribute[i] = p.String()
			}
			span.SetAttributes(attribute.StringSlice("platforms", platformsSliceAttribute))

			channel, err := s.channelService.GetChannelByID(ctx, req.ChannelID)
			if err != nil {
				return err
			}

			err = botplatforms.Dispatch(
				ctx,
				s.chatRegistry,
				channel.Bindings,
				req.Platforms,
				req.Message,
				req.ReplyTo,
				botplatforms.ChatOptions{
					IsAnnounce:        req.IsAnnounce,
					SkipToxicityCheck: req.SkipToxicityCheck,
					SkipRateLimits:    req.SkipRateLimits,
					AnnounceColor:     req.AnnounceColor,
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
			channel, found, err := s.getDeleteMessageChannel(ctx, req.ChannelId)
			if err != nil {
				s.logger.Error(
					"cannot find channel to delete messages from",
					slog.String("channelId", req.ChannelId),
				)
				return fmt.Errorf("find channel: %w", err)
			}

			if !found {
				return nil
			}

			wg, _ := errgroup.WithContext(ctx)

			for _, msgId := range req.MessageIds {
				wg.Go(
					func() error {
						deleteErr := s.twitchActions.DeleteMessage(
							ctx,
							twitchactions.DeleteMessageOpts{
								BroadcasterID: channel.BroadcasterID,
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
		},
	).Wait()
}

func (s *Service) getDeleteMessageChannel(ctx context.Context, twitchUserID string) (deleteMessageChannel, bool, error) {
	user, err := s.usersRepo.GetByPlatformID(ctx, platform.PlatformTwitch, twitchUserID)
	if err != nil {
		if errors.Is(err, usersmodel.ErrNotFound) {
			return deleteMessageChannel{}, false, nil
		}

		return deleteMessageChannel{}, false, err
	}

	channel, err := s.channelService.GetChannelByBindingUserID(ctx, platform.PlatformTwitch, user.ID)
	if err != nil {
		if errors.Is(err, channelsrepository.ErrNotFound) {
			return deleteMessageChannel{}, false, nil
		}

		return deleteMessageChannel{}, false, err
	}

	twitchBinding, botConfig, found, err := channel.TwitchBinding()
	if err != nil {
		return deleteMessageChannel{}, false, fmt.Errorf("parse Twitch bot config: %w", err)
	}
	if !found || !twitchBinding.Enabled || twitchBinding.PlatformChannelID == "" || botConfig.BotID == "" {
		return deleteMessageChannel{}, false, nil
	}

	return deleteMessageChannel{
		BotID:         botConfig.BotID,
		BroadcasterID: twitchBinding.PlatformChannelID,
	}, true, nil
}
