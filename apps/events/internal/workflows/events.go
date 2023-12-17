package workflows

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	eventsActivity "github.com/satont/twir/apps/events/internal/activities/events"
	"github.com/satont/twir/apps/events/internal/hydrator"
	"github.com/satont/twir/apps/events/internal/shared"
	config "github.com/satont/twir/libs/config"
	model "github.com/satont/twir/libs/gomodels"
	"go.temporal.io/sdk/client"
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
}

func NewEventsWorkflow(opts EventsWorkflowOpts) (*EventWorkflow, error) {
	c, err := client.Dial(
		client.Options{
			HostPort: opts.Cfg.TemporalHost,
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
	data shared.EvenData,
) error {
	options := client.StartWorkflowOptions{
		ID:        fmt.Sprintf("%s - %s", shared.EventsWorkflow, uuid.NewString()),
		TaskQueue: shared.EventsWorkerTaskQueueName,
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
