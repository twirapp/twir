package workflow

import (
	"time"

	"github.com/satont/twir/apps/timers/internal/activity"
	"github.com/satont/twir/apps/timers/internal/repositories/channels"
	"github.com/satont/twir/apps/timers/internal/repositories/streams"
	"github.com/satont/twir/apps/timers/internal/shared"
	cfg "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/logger"
	"github.com/twirapp/twir/libs/grpc/parser"
	timersrepository "github.com/twirapp/twir/libs/repositories/timers"
	timersmodel "github.com/twirapp/twir/libs/repositories/timers/model"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/log"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Opts struct {
	fx.In

	Lc     fx.Lifecycle
	Logger logger.Logger
	Cfg    cfg.Config

	TimersRepository   timersrepository.Repository
	ChannelsRepository channels.Repository
	StreamsRepository  streams.Repository

	ParserGrpc parser.ParserClient
	Activity   *activity.Activity
}

func New(opts Opts) (*Workflow, error) {
	cl, err := client.Dial(
		client.Options{
			HostPort: opts.Cfg.TemporalHost,
			Logger:   log.NewStructuredLogger(opts.Logger.GetSlog()),
		},
	)
	if err != nil {
		return nil, err
	}

	w := &Workflow{
		logger:             opts.Logger,
		config:             opts.Cfg,
		cl:                 cl,
		timersRepository:   opts.TimersRepository,
		channelsRepository: opts.ChannelsRepository,
		streamsRepository:  opts.StreamsRepository,
		parserGrpc:         opts.ParserGrpc,
		activity:           opts.Activity,
	}

	return w, nil
}

type Workflow struct {
	logger logger.Logger
	config cfg.Config
	db     *gorm.DB
	cl     client.Client

	timersRepository   timersrepository.Repository
	channelsRepository channels.Repository
	streamsRepository  streams.Repository

	parserGrpc parser.ParserClient
	activity   *activity.Activity
}

func (c *Workflow) Flow(ctx workflow.Context, timer timersmodel.Timer) error {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 5,
		HeartbeatTimeout:    time.Second * 10,
		TaskQueue:           shared.TimersWorkerTaskQueueName,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:        time.Second,
			BackoffCoefficient:     2.0,
			MaximumInterval:        time.Second * 100,
			MaximumAttempts:        3,
			NonRetryableErrorTypes: []string{},
		},
	}
	ctx = workflow.WithActivityOptions(ctx, options)

	var newResponse int
	err := workflow.ExecuteActivity(
		ctx,
		c.activity.SendMessage,
		timer.ID,
	).Get(ctx, &newResponse)
	if err != nil {
		return err
	}

	return nil
}
