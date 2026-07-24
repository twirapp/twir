package channel

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/twirapp/twir/apps/bots/internal/twitchactions"
	"github.com/twirapp/twir/libs/bus-core/bots"
	channelentity "github.com/twirapp/twir/libs/entities/channel"
	"github.com/twirapp/twir/libs/entities/platform"
	"github.com/twirapp/twir/libs/logger"
	channelsrepository "github.com/twirapp/twir/libs/repositories/channels"
)

var ErrFriendlyFire = errors.New("friendly fire")

func (s *Service) getChannelByIDOrTwitchID(
	ctx context.Context,
	id string,
) (channelentity.Channel, error) {
	if parsed, err := uuid.Parse(id); err == nil {
		return s.channelService.GetChannelByID(ctx, parsed)
	}
	return s.channelService.GetChannelByPlatformChannelID(ctx, platform.PlatformTwitch, id)
}

type twitchBanTarget struct {
	broadcasterID string
	botID         string
}

func getTwitchBanTarget(channel channelentity.Channel) (twitchBanTarget, bool, error) {
	twitchBinding, botConfig, found, err := channel.TwitchBinding()
	if err != nil {
		return twitchBanTarget{}, false, fmt.Errorf("parse Twitch bot config: %w", err)
	}
	if !found || !twitchBinding.Enabled || !botConfig.IsBotMod || botConfig.IsTwitchBanned ||
		twitchBinding.PlatformChannelID == "" || botConfig.BotID == "" {
		return twitchBanTarget{}, false, nil
	}

	return twitchBanTarget{
		broadcasterID: twitchBinding.PlatformChannelID,
		botID:         botConfig.BotID,
	}, true, nil
}

func (s *Service) Ban(ctx context.Context, req bots.BanRequest) error {
	return s.workersPool.SubmitErr(
		func() error {
			if req.ChannelID == req.UserID {
				return ErrFriendlyFire
			}

			channelEntity, err := s.getChannelByIDOrTwitchID(ctx, req.ChannelID)
			if err != nil {
				if errors.Is(err, channelsrepository.ErrNotFound) {
					s.logger.Error("channel not found", slog.String("channelId", req.ChannelID))
				} else {
					s.logger.Error("cannot get channel entity", logger.Error(err))
				}
				return fmt.Errorf("get channel entity: %w", err)
			}

			target, canBan, err := getTwitchBanTarget(channelEntity)
			if err != nil {
				return fmt.Errorf("get Twitch ban target: %w", err)
			}
			if !canBan {
				return nil
			}

			if err := s.twitchActions.Ban(
				ctx,
				twitchactions.BanOpts{
					Duration:       req.BanTime,
					Reason:         req.Reason,
					BroadcasterID:  target.broadcasterID,
					UserID:         req.UserID,
					ModeratorID:    target.botID,
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

func (s *Service) BanMany(ctx context.Context, reqs []bots.BanRequest) error {
	return s.workersPool.SubmitErr(
		func() error {
			uniqueChannelIDs := make(map[string]struct{})
			for _, r := range reqs {
				if r.ChannelID == r.UserID {
					continue
				}
				uniqueChannelIDs[r.ChannelID] = struct{}{}
			}

			if len(uniqueChannelIDs) == 0 {
				return nil
			}

			channelsByID := make(map[string]twitchBanTarget, len(uniqueChannelIDs))
			for id := range uniqueChannelIDs {
				ch, err := s.getChannelByIDOrTwitchID(ctx, id)
				if err != nil {
					if errors.Is(err, channelsrepository.ErrNotFound) {
						s.logger.Error("channel not found", slog.String("channelId", id))
						continue
					}
					s.logger.Error(
						"cannot get channel entity",
						logger.Error(err),
						slog.String("channelId", id),
					)
					continue
				}
				target, canBan, targetErr := getTwitchBanTarget(ch)
				if targetErr != nil {
					s.logger.Error(
						"cannot get Twitch ban target",
						logger.Error(targetErr),
						slog.String("channelId", id),
					)
					continue
				}
				if canBan {
					channelsByID[id] = target
				}
			}

			if len(channelsByID) == 0 {
				return nil
			}

			timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
			defer cancel()

			var (
				wg        sync.WaitGroup
				mu        sync.Mutex
				banErrors []error
			)

			for _, r := range reqs {
				target, ok := channelsByID[r.ChannelID]
				if !ok {
					continue
				}

				if r.ChannelID == r.UserID || target.botID == r.UserID {
					continue
				}

				wg.Add(1)
				go func(target twitchBanTarget, r bots.BanRequest) {
					defer wg.Done()

					if err := s.twitchActions.Ban(
						timeoutCtx,
						twitchactions.BanOpts{
							Duration:       r.BanTime,
							Reason:         r.Reason,
							BroadcasterID:  target.broadcasterID,
							UserID:         r.UserID,
							ModeratorID:    target.botID,
							IsModerator:    r.IsModerator,
							AddModAfterBan: r.AddModAfterBan,
						},
					); err != nil {
						s.logger.Error("cannot ban user", logger.Error(err))
						mu.Lock()
						banErrors = append(banErrors, err)
						mu.Unlock()
					}
				}(target, r)
			}

			wg.Wait()
			return errors.Join(banErrors...)
		},
	).Wait()
}
