package workers

import (
	"context"
	"log/slog"

	eventsActivity "github.com/twirapp/twir/apps/events/internal/activities/events"
	"github.com/twirapp/twir/apps/events/internal/shared"
	"github.com/twirapp/twir/apps/events/internal/workflows"
	"github.com/twirapp/twir/libs/baseapp/lifecycle"
	config "github.com/twirapp/twir/libs/config"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/log"
	"go.temporal.io/sdk/worker"
)

func NewEventsWorker(
	lc *lifecycle.Lifecycle,
	cfg config.Config,
	workflow *workflows.EventWorkflow,
	logger *slog.Logger,
	activities *eventsActivity.Activity,
) error {
	c, err := client.Dial(
		client.Options{
			HostPort: cfg.TemporalHost,
			Logger:   log.NewStructuredLogger(logger),
		},
	)
	if err != nil {
		return err
	}

	temporalWorker := worker.New(c, shared.EventsWorkerTaskQueueName, worker.Options{})

	temporalWorker.RegisterWorkflow(workflow.Flow)

	temporalWorker.RegisterActivity(activities.SendMessage)
	temporalWorker.RegisterActivity(activities.Ban)
	temporalWorker.RegisterActivity(activities.Unban)
	temporalWorker.RegisterActivity(activities.BanRandom)
	temporalWorker.RegisterActivity(activities.ChangeTitle)
	temporalWorker.RegisterActivity(activities.ChangeCategory)
	temporalWorker.RegisterActivity(activities.CommandAllowOrRemoveUserPermission)
	temporalWorker.RegisterActivity(activities.CommandDenyOrRemoveUserPermission)
	temporalWorker.RegisterActivity(activities.CreateGreeting)
	temporalWorker.RegisterActivity(activities.SwitchEmoteOnly)
	temporalWorker.RegisterActivity(activities.SwitchSubMode)
	temporalWorker.RegisterActivity(activities.ModOrUnmod)
	temporalWorker.RegisterActivity(activities.UnmodRandom)
	temporalWorker.RegisterActivity(activities.ObsSetScene)
	temporalWorker.RegisterActivity(activities.ObsToggleSource)
	temporalWorker.RegisterActivity(activities.ObsToggleAudio)
	temporalWorker.RegisterActivity(activities.ObsAudioChangeVolume)
	temporalWorker.RegisterActivity(activities.ObsAudioSetVolume)
	temporalWorker.RegisterActivity(activities.ObsEnableOrDisableAudio)
	temporalWorker.RegisterActivity(activities.ObsStartOrStopStream)
	temporalWorker.RegisterActivity(activities.TtsSay)
	temporalWorker.RegisterActivity(activities.TtsSkip)
	temporalWorker.RegisterActivity(activities.TtsChangeState)
	temporalWorker.RegisterActivity(activities.TtsChangeAutoReadState)
	temporalWorker.RegisterActivity(activities.ChangeVariableValue)
	temporalWorker.RegisterActivity(activities.IncrementORDecrementVariable)
	temporalWorker.RegisterActivity(activities.VipOrUnvip)
	temporalWorker.RegisterActivity(activities.UnvipRandom)
	temporalWorker.RegisterActivity(activities.SevenTvEmoteManage)
	temporalWorker.RegisterActivity(activities.RaidChannel)
	temporalWorker.RegisterActivity(activities.TriggerAlert)
	temporalWorker.RegisterActivity(activities.ShoutoutChannel)
	temporalWorker.RegisterActivity(activities.MessageDelete)
	temporalWorker.RegisterActivity(activities.SendHttpRequest)

	lc.Append(
		lifecycle.Hook{
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
