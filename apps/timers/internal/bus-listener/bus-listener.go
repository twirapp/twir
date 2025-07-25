package bus_listener

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/twirapp/twir/apps/timers/internal/workflow"
	"github.com/twirapp/twir/libs/logger"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/timers"
	"github.com/twirapp/twir/libs/redis_keys"
	"go.uber.org/fx"
)

type server struct {
	workflow *workflow.Workflow
	redis    *redis.Client
}

type Opts struct {
	fx.In

	Lc       fx.Lifecycle
	Logger   logger.Logger
	Workflow *workflow.Workflow
	Bus      *buscore.Bus
	Redis    *redis.Client
}

func New(opts Opts) error {
	s := &server{
		workflow: opts.Workflow,
		redis:    opts.Redis,
	}

	opts.Lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				opts.Bus.Timers.AddTimer.SubscribeGroup("timers", s.addTimerToQueue)
				opts.Bus.Timers.RemoveTimer.SubscribeGroup("timers", s.removeTimerFromQueue)

				opts.Logger.Info("Timers grpc server started")
				return nil
			},
			OnStop: func(ctx context.Context) error {
				opts.Bus.Timers.AddTimer.Unsubscribe()
				opts.Bus.Timers.RemoveTimer.Unsubscribe()
				return nil
			},
		},
	)

	return nil
}

func (c *server) addTimerToQueue(ctx context.Context, t timers.AddOrRemoveTimerRequest) (
	struct{},
	error,
) {
	if err := c.redis.Del(ctx, redis_keys.TimersCurrentResponse(t.TimerID)).Err(); err != nil {
		return struct{}{}, err
	}

	c.workflow.RemoveTimer(ctx, t.TimerID)
	if err := c.workflow.AddTimer(ctx, t.TimerID); err != nil {
		return struct{}{}, err
	}

	return struct{}{}, nil
}

func (c *server) removeTimerFromQueue(
	ctx context.Context,
	t timers.AddOrRemoveTimerRequest,
) (struct{}, error) {
	c.workflow.RemoveTimer(ctx, t.TimerID)

	return struct{}{}, nil
}
