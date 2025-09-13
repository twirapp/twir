package admin_actions

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/eventsub"
	"github.com/twirapp/twir/libs/bus-core/timers"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/repositories/channels"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Opts struct {
	fx.In

	Redis              *redis.Client
	ChannelsRepository channels.Repository
	TwirBus            *buscore.Bus
	Gorm               *gorm.DB
}

func New(opts Opts) *Service {
	return &Service{
		redis:              opts.Redis,
		channelsRepository: opts.ChannelsRepository,
		twirbus:            opts.TwirBus,
		gorm:               opts.Gorm,
	}
}

type Service struct {
	redis              *redis.Client
	channelsRepository channels.Repository
	twirbus            *buscore.Bus
	gorm               *gorm.DB
}

func (c *Service) DropAllAuthSessions(ctx context.Context) error {
	iter := c.redis.Scan(ctx, 0, "scs:*", 0).Iterator()

	_, err := c.redis.Pipelined(
		ctx, func(pipe redis.Pipeliner) error {
			for iter.Next(ctx) {
				if err := pipe.Del(ctx, iter.Val()).Err(); err != nil {
					return err
				}
			}

			if err := iter.Err(); err != nil {
				return fmt.Errorf("scan err: %w", err)
			}

			return nil
		},
	)
	if err != nil {
		return fmt.Errorf("scan err: %w", err)
	}

	return nil
}

type EventSubSubscribeInput struct {
	Type    string `json:"type"`
	Version string `json:"version"`
}

func (c *Service) EventSubSubscribe(ctx context.Context, input EventSubSubscribeInput) error {
	ch, err := c.channelsRepository.GetMany(
		ctx,
		channels.GetManyInput{
			Enabled: lo.ToPtr(true),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to get channels: %w", err)
	}

	for _, channel := range ch {
		go func() {
			c.twirbus.EventSub.Subscribe.Publish(
				ctx,
				eventsub.EventsubSubscribeRequest{
					ChannelID: channel.ID,
					Topic:     input.Type,
					Version:   input.Version,
				},
			)
		}()
	}

	return nil
}

func (c *Service) RescheduleTimers(ctx context.Context) error {
	var entities []model.ChannelsTimers
	if err := c.gorm.Select("id", "enabled").Find(&entities).Error; err != nil {
		return fmt.Errorf("failed to get timers: %w", err)
	}

	for _, timer := range entities {
		c.twirbus.Timers.RemoveTimer.Publish(
			ctx,
			timers.AddOrRemoveTimerRequest{
				TimerID: timer.ID,
			},
		)

		if timer.Enabled {
			c.twirbus.Timers.AddTimer.Publish(
				ctx,
				timers.AddOrRemoveTimerRequest{
					TimerID: timer.ID,
				},
			)
		}
	}

	return nil
}

func (c *Service) EventsubReinitChannels(ctx context.Context) error {
	c.twirbus.EventSub.InitChannels.Publish(ctx, struct{}{})

	return nil
}

func (c *Service) BanUser(ctx context.Context, userId string) error {
	if err := c.gorm.
		WithContext(ctx).
		Model(&model.Users{}).
		Where("id = ?", userId).
		Update("is_banned", true).Error; err != nil {
		return err
	}

	if err := c.twirbus.EventSub.Unsubscribe.Publish(ctx, userId); err != nil {
		return err
	}

	return nil
}
