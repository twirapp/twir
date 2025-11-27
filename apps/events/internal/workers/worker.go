package workers

import (
	"context"
	"log/slog"

	eventsActivity "github.com/twirapp/twir/apps/events/internal/activities/events"
	"github.com/twirapp/twir/apps/events/internal/shared"
	"github.com/twirapp/twir/apps/events/internal/workflows"
	config "github.com/twirapp/twir/libs/config"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/log"
	"go.temporal.io/sdk/worker"
	"go.uber.org/fx"
)

type EventsWorkerOpts struct {
	fx.In
	Lc fx.Lifecycle

	Cfg        config.Config
	Workflow   *workflows.EventWorkflow
	Logger     *slog.Logger
	Activities *eventsActivity.Activity
}

func NewEventsWorker(opts EventsWorkerOpts) error {
	c, err := client.Dial(
		client.Options{
			HostPort: opts.Cfg.TemporalHost,
			Logger:   log.NewStructuredLogger(opts.Logger),
		},
	)
	if err != nil {
		return err
	}

	temporalWorker := worker.New(c, shared.EventsWorkerTaskQueueName, worker.Options{})

	temporalWorker.RegisterWorkflow(opts.Workflow.Flow)

	temporalWorker.RegisterActivity(opts.Activities.SendMessage)
	temporalWorker.RegisterActivity(opts.Activities.Ban)
	temporalWorker.RegisterActivity(opts.Activities.Unban)
	temporalWorker.RegisterActivity(opts.Activities.BanRandom)
	temporalWorker.RegisterActivity(opts.Activities.ChangeTitle)
	temporalWorker.RegisterActivity(opts.Activities.ChangeCategory)
	temporalWorker.RegisterActivity(opts.Activities.CommandAllowOrRemoveUserPermission)
	temporalWorker.RegisterActivity(opts.Activities.CommandDenyOrRemoveUserPermission)
	temporalWorker.RegisterActivity(opts.Activities.CreateGreeting)
	temporalWorker.RegisterActivity(opts.Activities.SwitchEmoteOnly)
	temporalWorker.RegisterActivity(opts.Activities.SwitchSubMode)
	temporalWorker.RegisterActivity(opts.Activities.ModOrUnmod)
	temporalWorker.RegisterActivity(opts.Activities.UnmodRandom)
	temporalWorker.RegisterActivity(opts.Activities.ObsSetScene)
	temporalWorker.RegisterActivity(opts.Activities.ObsToggleSource)
	temporalWorker.RegisterActivity(opts.Activities.ObsToggleAudio)
	temporalWorker.RegisterActivity(opts.Activities.ObsAudioChangeVolume)
	temporalWorker.RegisterActivity(opts.Activities.ObsAudioSetVolume)
	temporalWorker.RegisterActivity(opts.Activities.ObsEnableOrDisableAudio)
	temporalWorker.RegisterActivity(opts.Activities.ObsStartOrStopStream)
	temporalWorker.RegisterActivity(opts.Activities.TtsSay)
	temporalWorker.RegisterActivity(opts.Activities.TtsSkip)
	temporalWorker.RegisterActivity(opts.Activities.TtsChangeState)
	temporalWorker.RegisterActivity(opts.Activities.TtsChangeAutoReadState)
	temporalWorker.RegisterActivity(opts.Activities.ChangeVariableValue)
	temporalWorker.RegisterActivity(opts.Activities.IncrementORDecrementVariable)
	temporalWorker.RegisterActivity(opts.Activities.VipOrUnvip)
	temporalWorker.RegisterActivity(opts.Activities.UnvipRandom)
	temporalWorker.RegisterActivity(opts.Activities.SevenTvEmoteManage)
	temporalWorker.RegisterActivity(opts.Activities.RaidChannel)
	temporalWorker.RegisterActivity(opts.Activities.TriggerAlert)
	temporalWorker.RegisterActivity(opts.Activities.ShoutoutChannel)
	temporalWorker.RegisterActivity(opts.Activities.MessageDelete)

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
