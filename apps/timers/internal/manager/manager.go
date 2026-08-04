package manager

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
	"github.com/twirapp/twir/libs/baseapp/lifecycle"
	buscore "github.com/twirapp/twir/libs/bus-core"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	cfg "github.com/twirapp/twir/libs/config"
	channelentity "github.com/twirapp/twir/libs/entities/channel"
	timersentity "github.com/twirapp/twir/libs/entities/timers"
	channelsrepository "github.com/twirapp/twir/libs/repositories/channels"
	"github.com/twirapp/twir/libs/repositories/timers"
	channelservice "github.com/twirapp/twir/libs/services/channels"
)

func New(
	lc *lifecycle.Lifecycle,
	repository timers.Repository,
	logger *slog.Logger,
	channelCachedRepo *generic_cacher.GenericCacher[channelentity.Channel],
	redis *redis.Client,
	twirBus *buscore.Bus,
	cfg cfg.Config,
	channelsRepo channelsrepository.Repository,
	channelsService *channelservice.ChannelService,
) *Manager {
	m := &Manager{
		timers:            make(map[TimerID]*Timer),
		repository:        repository,
		logger:            logger,
		stopChan:          make(chan struct{}, 1),
		channelCachedRepo: channelCachedRepo,
		redis:             redis,
		twirBus:           twirBus,
		config:            cfg,
		channelsRepo:      channelsRepo,
		channelservice:    channelsService,
	}

	lc.Append(
		lifecycle.Hook{
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
	channelCachedRepo *generic_cacher.GenericCacher[channelentity.Channel]
	redis             *redis.Client
	twirBus           *buscore.Bus
	config            cfg.Config
	channelsRepo      channelsrepository.Repository
	channelservice    *channelservice.ChannelService
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

	channelsByID := make(map[string]channelentity.Channel)
	for page := 0; ; page++ {
		channels, err := c.channelsRepo.GetMany(
			ctx,
			channelsrepository.GetManyInput{
				Enabled: new(true),
				PerPage: 100,
				Page:    page,
			},
		)
		if err != nil {
			return fmt.Errorf("cannot get channels: %w", err)
		}

		for _, ch := range channels {
			channelsByID[ch.ID.String()] = ch
		}

		if len(channels) < 100 {
			break
		}
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

		for _, t := range timersBatch {
			foundChannel, ok := channelsByID[t.ChannelID.String()]
			if !ok || !hasSupportedTimerBinding(foundChannel) {
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
		timer.ticker = time.NewTicker(5 * time.Second)

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
		slog.String("channelId", dbRow.ChannelID.String()),
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
		slog.String("channelId", t.dbRow.ChannelID.String()),
	)

	delete(c.timers, id)
}

func (c *Manager) OnChatMessage(channelId uuid.UUID) {
	for _, t := range c.timers {
		if t.dbRow.ChannelID != channelId {
			continue
		}

		go c.tryTick(t.id)
	}
}
