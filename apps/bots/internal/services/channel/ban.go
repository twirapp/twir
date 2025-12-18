package channel

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/twirapp/twir/apps/bots/internal/twitchactions"
	"github.com/twirapp/twir/libs/bus-core/bots"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/logger"
)

var ErrFriendlyFire = errors.New("friendly fire")

func (s *Service) Ban(ctx context.Context, req bots.BanRequest) error {
	return s.workersPool.SubmitErr(
		func() error {
			if req.ChannelID == req.UserID {
				return ErrFriendlyFire
			}

			var channelEntity model.Channels
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

func (s *Service) BanMany(ctx context.Context, reqs []bots.BanRequest) error {
	return s.workersPool.SubmitErr(
		func() error {
			channels, err := s.findUniqueBanChannels(ctx, reqs)
			if err != nil {
				return fmt.Errorf("find unique ban channels: %w", err)
			}

			if len(channels) == 0 {
				return nil
			}

			timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
			defer cancel()

			var (
				wg        sync.WaitGroup
				banErrors []error
			)

			for _, r := range reqs {
				var channelEntity *model.Channels
				for _, channel := range channels {
					if channel.ID == r.ChannelID {
						channelEntity = &channel
						break
					}
				}

				if channelEntity == nil {
					continue
				}

				if r.ChannelID == r.UserID || channelEntity.BotID == r.UserID {
					continue
				}

				wg.Add(1)
				go func() {
					defer wg.Done()

					if err := s.twitchActions.Ban(
						timeoutCtx,
						twitchactions.BanOpts{
							Duration:       r.BanTime,
							Reason:         r.Reason,
							BroadcasterID:  r.ChannelID,
							UserID:         r.UserID,
							ModeratorID:    channelEntity.BotID,
							IsModerator:    r.IsModerator,
							AddModAfterBan: r.AddModAfterBan,
						},
					); err != nil {
						s.logger.Error("cannot ban user", logger.Error(err))
						banErrors = append(banErrors, err)
					}
				}()
			}

			wg.Wait()
			return errors.Join(banErrors...)
		},
	).Wait()
}

func (s *Service) findUniqueBanChannels(ctx context.Context, reqs []bots.BanRequest) ([]model.Channels, error) {
	uniqueChannelsIdsMap := make(map[string]struct{}, len(reqs))
	for _, r := range reqs {
		if r.ChannelID == r.UserID {
			continue
		}

		uniqueChannelsIdsMap[r.ChannelID] = struct{}{}
	}

	if len(uniqueChannelsIdsMap) == 0 {
		return nil, nil
	}

	uniqueChannelsIds := make([]string, 0, len(uniqueChannelsIdsMap))
	for k := range uniqueChannelsIdsMap {
		uniqueChannelsIds = append(uniqueChannelsIds, k)
	}

	var channels []model.Channels
	if err := s.gorm.WithContext(ctx).Where(
		`"id" IN ?`,
		uniqueChannelsIds,
	).Find(&channels).Error; err != nil {
		s.logger.Error("cannot get channels entities", logger.Error(err))
		return nil, err
	}

	return channels, nil
}
