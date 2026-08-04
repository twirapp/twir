package workers

import (
	"context"
	"log/slog"
	"time"

	"github.com/alitto/pond/v2"
	"github.com/twirapp/twir/libs/baseapp/lifecycle"
	channelentity "github.com/twirapp/twir/libs/entities/channel"
	"github.com/twirapp/twir/libs/entities/platform"
	"github.com/twirapp/twir/libs/logger"
	channelsrepository "github.com/twirapp/twir/libs/repositories/channels"
	"go.uber.org/atomic"
)

type Opts struct {
	LC *lifecycle.Lifecycle

	ChannelsRepository channelsrepository.Repository
	Logger             *slog.Logger
}

func New(opts Opts) *Pool {
	w := &Pool{
		Pool:               pond.NewPool(1),
		channelsRepository: opts.ChannelsRepository,
		logger:             opts.Logger,
	}

	workersResizerCtx, workersResizerCtxCancel := context.WithCancel(context.Background())

	opts.LC.Append(
		lifecycle.Hook{
			OnStart: func(ctx context.Context) error {
				w.setSize(workersResizerCtx)

				go func() {
					ticker := time.NewTicker(1 * time.Minute)
					defer ticker.Stop()

					for {
						select {
						case <-workersResizerCtx.Done():
							return
						case <-ticker.C:
							w.setSize(workersResizerCtx)
						}
					}
				}()
				return nil
			},
			OnStop: func(ctx context.Context) error {
				workersResizerCtxCancel()
				w.StopAndWait()
				return nil
			},
		},
	)

	return w
}

type Pool struct {
	pond.Pool

	channelsRepository channelsrepository.Repository
	logger             *slog.Logger
	lastSize           atomic.Int64
}

const (
	proposedMessageHandlers = 15
	proposedTwitchActions   = 6
)

func (c *Pool) setSize(ctx context.Context) {
	channels, err := c.channelsRepository.GetAllByBindingPlatform(ctx, platform.PlatformTwitch)
	if err != nil {
		c.logger.Error("cannot get Twitch channels", logger.Error(err))
		return
	}

	channelsCount := countEnabledTwitchChannels(channels)
	newSize := channelsCount * proposedMessageHandlers * proposedTwitchActions
	if c.lastSize.Swap(int64(newSize)) == int64(newSize) {
		return
	}

	c.Resize(newSize)
	c.logger.Info(
		"workers pool resized",
		slog.Int("size", newSize),
	)
}

func countEnabledTwitchChannels(channels []channelentity.Channel) int {
	count := 0
	for _, channel := range channels {
		binding, ok := channel.Binding(platform.PlatformTwitch)
		if ok && binding.Enabled {
			count++
		}
	}

	return count
}
