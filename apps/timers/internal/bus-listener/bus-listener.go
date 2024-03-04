package bus_listener

import (
	"context"

	"github.com/satont/twir/apps/timers/internal/workflow"
	"github.com/satont/twir/libs/logger"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/timers"
	"go.uber.org/fx"
)

type server struct {
	workflow *workflow.Workflow
}

type Opts struct {
	fx.In

	Lc       fx.Lifecycle
	Logger   logger.Logger
	Workflow *workflow.Workflow
	Bus      *buscore.Bus
}

func New(opts Opts) error {
	s := &server{
		workflow: opts.Workflow,
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

func (c *server) addTimerToQueue(ctx context.Context, t timers.AddOrRemoveTimerRequest) struct{} {
	c.workflow.AddTimer(ctx, t.TimerID)

	return struct{}{}
}

func (c *server) removeTimerFromQueue(
	ctx context.Context,
	t timers.AddOrRemoveTimerRequest,
) struct{} {
	c.workflow.RemoveTimer(ctx, t.TimerID)

	return struct{}{}
}
