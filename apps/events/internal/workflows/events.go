package workflows

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	eventsActivity "github.com/satont/twir/apps/events/internal/activities/events"
	"github.com/satont/twir/apps/events/internal/hydrator"
	"github.com/satont/twir/apps/events/internal/shared"
	config "github.com/satont/twir/libs/config"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/logger"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/log"
	"go.temporal.io/sdk/temporal"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type EventsWorkflowOpts struct {
	fx.In

	Cfg            config.Config
	EventsActivity *eventsActivity.Activity
	Gorm           *gorm.DB
	Redis          *redis.Client
	Hydrator       *hydrator.Hydrador
	Logger         logger.Logger
}

func NewEventsWorkflow(opts EventsWorkflowOpts) (*EventWorkflow, error) {
	c, err := client.Dial(
		client.Options{
			HostPort: opts.Cfg.TemporalHost,
			Logger:   log.NewStructuredLogger(opts.Logger.GetSlog()),
		},
	)
	if err != nil {
		return nil, err
	}

	return &EventWorkflow{
		cfg:            opts.Cfg,
		cl:             c,
		eventsActivity: opts.EventsActivity,
		db:             opts.Gorm,
		redis:          opts.Redis,
		hydrator:       opts.Hydrator,
	}, nil
}

type EventWorkflow struct {
	cfg            config.Config
	cl             client.Client
	eventsActivity *eventsActivity.Activity
	db             *gorm.DB
	redis          *redis.Client
	hydrator       *hydrator.Hydrador
}

func (c *EventWorkflow) Execute(
	ctx context.Context,
	eventType model.EventType,
	data shared.EventData,
) error {
	options := client.StartWorkflowOptions{
		ID:        fmt.Sprintf("%s - %s", shared.EventsWorkflow, uuid.NewString()),
		TaskQueue: shared.EventsWorkerTaskQueueName,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:        time.Second,
			BackoffCoefficient:     2.0,
			MaximumInterval:        time.Second * 100,
			MaximumAttempts:        3,
			NonRetryableErrorTypes: []string{},
		},
	}

	we, err := c.cl.ExecuteWorkflow(
		ctx,
		options,
		c.Flow,
		eventType,
		data,
	)
	if err != nil {
		return err
	}

	return we.Get(ctx, nil)
}
