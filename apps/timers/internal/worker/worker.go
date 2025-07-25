package worker

import (
	"context"

	"github.com/twirapp/twir/apps/timers/internal/activity"
	"github.com/twirapp/twir/apps/timers/internal/shared"
	"github.com/twirapp/twir/apps/timers/internal/workflow"
	config "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/logger"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/log"
	"go.temporal.io/sdk/worker"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In
	Lc fx.Lifecycle

	Logger   logger.Logger
	Cfg      config.Config
	Workflow *workflow.Workflow
	Activity *activity.Activity
}

func New(opts Opts) error {
	c, err := client.Dial(
		client.Options{
			HostPort: opts.Cfg.TemporalHost,
			Logger:   log.NewStructuredLogger(opts.Logger.GetSlog()),
		},
	)
	if err != nil {
		return err
	}

	temporalWorker := worker.New(c, shared.TimersWorkerTaskQueueName, worker.Options{})

	temporalWorker.RegisterWorkflow(opts.Workflow.Flow)
	temporalWorker.RegisterActivity(opts.Activity.SendMessage)

	opts.Lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				return temporalWorker.Start()
			},
			OnStop: func(ctx context.Context) error {
				temporalWorker.Stop()
				return nil
			},
		},
	)

	return nil
}

type Worker struct {
	cl client.Client
}
