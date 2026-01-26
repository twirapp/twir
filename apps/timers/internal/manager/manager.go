package manager

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
	buscore "github.com/twirapp/twir/libs/bus-core"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	cfg "github.com/twirapp/twir/libs/config"
	timersentity "github.com/twirapp/twir/libs/entities/timers"
	channelsrepository "github.com/twirapp/twir/libs/repositories/channels"
	channelmodel "github.com/twirapp/twir/libs/repositories/channels/model"
	"github.com/twirapp/twir/libs/repositories/timers"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In
	LC fx.Lifecycle

	Repository        timers.Repository
	Logger            *slog.Logger
	ChannelCachedRepo *generic_cacher.GenericCacher[channelmodel.Channel]
	Redis             *redis.Client
	TwirBus           *buscore.Bus
	Config            cfg.Config
	ChannelsRepo      channelsrepository.Repository
}

func New(opts Opts) *Manager {
	m := &Manager{
		timers:            make(map[TimerID]*Timer),
		repository:        opts.Repository,
		logger:            opts.Logger,
		stopChan:          make(chan struct{}, 1),
		channelCachedRepo: opts.ChannelCachedRepo,
		redis:             opts.Redis,
		twirBus:           opts.TwirBus,
		config:            opts.Config,
		channelsRepo:      opts.ChannelsRepo,
	}

	opts.LC.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				return m.initialize(ctx)
			},
			OnStop: func(ctx context.Context) error {
				for id := range m.timers {
					m.RemoveTimerById(id)
				}

				m.stopChan <- struct{}{}

				return nil
			},
		},
	)

	return m
}

type Manager struct {
	timers map[TimerID]*Timer

	repository        timers.Repository
	logger            *slog.Logger
	stopChan          chan struct{}
	channelCachedRepo *generic_cacher.GenericCacher[channelmodel.Channel]
	redis             *redis.Client
	twirBus           *buscore.Bus
	config            cfg.Config
	channelsRepo      channelsrepository.Repository
}

func (c *Manager) initialize(ctx context.Context) error {
	totalTimers, err := c.repository.Count(
		ctx, timers.CountInput{
			Enabled: lo.ToPtr(true),
		},
	)
	if err != nil {
		return fmt.Errorf("cannot get count of timers: %w", err)
	}

	if totalTimers == 0 {
		return nil
	}

	for offset := int64(0); offset < totalTimers; {
		batchSize := int64(100)
		if offset+batchSize > totalTimers {
			batchSize = totalTimers - offset
		}

		timersBatch, err := c.repository.GetMany(
			ctx,
			timers.GetManyInput{
				Enabled: lo.ToPtr(true),
				Limit:   int(batchSize),
				Offset:  int(offset),
			},
		)
		if err != nil {
			return fmt.Errorf("cannot initialize timers manager: %w", err)
		}

		channels, err := c.channelsRepo.GetMany(
			ctx,
			channelsrepository.GetManyInput{
				Enabled: lo.ToPtr(true),
			},
		)
		if err != nil {
			return fmt.Errorf("cannot get channels: %w", err)
		}

		for _, t := range timersBatch {
			var foundChannel channelmodel.Channel
			for _, ch := range channels {
				if ch.ID == t.ChannelID {
					foundChannel = ch
					break
				}
			}
			if foundChannel == channelmodel.Nil || foundChannel.ID == "" || !foundChannel.IsBotMod ||
				foundChannel.IsTwitchBanned ||
				!foundChannel.IsEnabled {
				continue
			}

			c.addTimer(t)
		}

		offset += batchSize
	}

	return nil
}

func (c *Manager) addTimer(dbRow timersentity.Timer) {
	timerId := TimerID(dbRow.ID)

	c.RemoveTimerById(timerId)

	timer := Timer{
		id:                   timerId,
		ticker:               nil,
		lastTriggerTimestamp: time.Now(),
		dbRow:                dbRow,
	}

	if dbRow.TimeInterval != 0 {
		timer.ticker = time.NewTicker(time.Duration(dbRow.TimeInterval) * time.Second)

		go func() {
			for {
				select {
				case <-c.stopChan:
					return
				case <-timer.ticker.C:
					c.tryTick(timer.id)
				}
			}
		}()
	}

	c.timers[timer.id] = &timer

	c.logger.Info(
		"[manager] added timer",
		slog.String("timerId", timerId.String()),
		slog.String("channelId", dbRow.ChannelID),
		slog.Int("timeInterval", dbRow.TimeInterval),
		slog.Int("messageInterval", dbRow.MessageInterval),
	)
}

func (c *Manager) AddTimerById(ctx context.Context, id TimerID) error {
	dbRow, err := c.repository.GetByID(ctx, uuid.UUID(id))
	if err != nil {
		return fmt.Errorf("cannot add timer: %w", err)
	}

	c.addTimer(dbRow)

	return nil
}

func (c *Manager) RemoveTimerById(id TimerID) {
	t, ok := c.timers[id]
	if !ok {
		return
	}

	if t.ticker != nil {
		t.ticker.Stop()
	}

	c.logger.Info(
		"[manager] removed timer",
		slog.String("timerId", id.String()),
		slog.String("channelId", t.dbRow.ChannelID),
	)

	delete(c.timers, id)
}

func (c *Manager) OnChatMessage(channelId string) {
	for _, t := range c.timers {
		if t.dbRow.ChannelID != channelId {
			continue
		}

		if t.dbRow.OfflineEnabled {
			t.offlineMessageNumber++
		}

		go c.tryTick(t.id)
	}
}
