package workers

import (
	"context"
	"log/slog"
	"time"

	"github.com/alitto/pond/v2"
	"github.com/satont/twir/libs/logger"
	channelsrepository "github.com/twirapp/twir/libs/repositories/channels"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In
	LC fx.Lifecycle

	ChannelsRepository channelsrepository.Repository
	Logger             logger.Logger
}

func New(opts Opts) *Pool {
	w := &Pool{
		Pool:               pond.NewPool(1),
		channelsRepository: opts.ChannelsRepository,
		logger:             opts.Logger,
	}

	workersResizerCtx, workersResizerCtxCancel := context.WithCancel(context.Background())

	opts.LC.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				w.setSize(workersResizerCtx)

				go func() {
					for {
						select {
						case <-ctx.Done():
							return
						case <-time.After(1 * time.Minute):
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
	logger             logger.Logger
}

const (
	proposedMessageHandlers = 15
	proposedTwitchActions   = 6
)

func (c *Pool) setSize(ctx context.Context) {
	channelsCount, err := c.channelsRepository.GetCount(
		ctx, channelsrepository.GetCountInput{
			OnlyEnabled: true,
		},
	)
	if err != nil {
		c.logger.Error("cannot get channels count", slog.Any("err", err))
		return
	}
	newSize := channelsCount * proposedMessageHandlers * proposedTwitchActions
	c.Resize(newSize)
	c.logger.Info(
		"workers pool resized",
		slog.Int("size", newSize),
	)
}
