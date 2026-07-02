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
	"github.com/twirapp/twir/libs/logger"
	channelsrepository "github.com/twirapp/twir/libs/repositories/channels"
	channelmodel "github.com/twirapp/twir/libs/repositories/channels/model"
)

var ErrFriendlyFire = errors.New("friendly fire")

func (s *Service) getChannelByIDOrTwitchID(
	ctx context.Context,
	id string,
) (channelmodel.Channel, error) {
	if parsed, err := uuid.Parse(id); err == nil {
		return s.channelsRepo.GetByID(ctx, parsed)
	}
	return s.channelsRepo.GetByTwitchPlatformID(ctx, id)
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

			broadcasterID := req.ChannelID
			if channelEntity.TwitchPlatformID != nil {
				broadcasterID = *channelEntity.TwitchPlatformID
			}

			if err := s.twitchActions.Ban(
				ctx,
				twitchactions.BanOpts{
					Duration:       req.BanTime,
					Reason:         req.Reason,
					BroadcasterID:  broadcasterID,
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

			channelsByID := make(map[string]channelmodel.Channel, len(uniqueChannelIDs))
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
				channelsByID[id] = ch
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
				channelEntity, ok := channelsByID[r.ChannelID]
				if !ok {
					continue
				}

				if r.ChannelID == r.UserID || channelEntity.BotID == r.UserID {
					continue
				}

				broadcasterID := r.ChannelID
				if channelEntity.TwitchPlatformID != nil {
					broadcasterID = *channelEntity.TwitchPlatformID
				}

				wg.Add(1)
				go func(broadcasterID, botID string, r bots.BanRequest) {
					defer wg.Done()

					if err := s.twitchActions.Ban(
						timeoutCtx,
						twitchactions.BanOpts{
							Duration:       r.BanTime,
							Reason:         r.Reason,
							BroadcasterID:  broadcasterID,
							UserID:         r.UserID,
							ModeratorID:    botID,
							IsModerator:    r.IsModerator,
							AddModAfterBan: r.AddModAfterBan,
						},
					); err != nil {
						s.logger.Error("cannot ban user", logger.Error(err))
						mu.Lock()
						banErrors = append(banErrors, err)
						mu.Unlock()
					}
				}(broadcasterID, channelEntity.BotID, r)
			}

			wg.Wait()
			return errors.Join(banErrors...)
		},
	).Wait()
}
