package workflows

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	eventsActivity "github.com/twirapp/twir/apps/events/internal/activities/events"
	"github.com/twirapp/twir/apps/events/internal/hydrator"
	"github.com/twirapp/twir/apps/events/internal/shared"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	config "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/entities/platform"
	channelsrepository "github.com/twirapp/twir/libs/repositories/channels"
	channelmodel "github.com/twirapp/twir/libs/repositories/channels/model"
	"github.com/twirapp/twir/libs/repositories/events/model"
	usersrepository "github.com/twirapp/twir/libs/repositories/users"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/log"
	"go.temporal.io/sdk/temporal"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type EventsWorkflowOpts struct {
	fx.In

	Cfg                               config.Config
	EventsActivity                    *eventsActivity.Activity
	Gorm                              *gorm.DB
	Redis                             *redis.Client
	Hydrator                          *hydrator.Hydrator
	Logger                            *slog.Logger
	ChannelsEventsWithOperationsCache *generic_cacher.GenericCacher[[]model.Event]
	ChannelsCache                     *generic_cacher.GenericCacher[channelmodel.Channel]
	ChannelsRepo                      channelsrepository.Repository
	UsersRepo                         usersrepository.Repository
}

func NewEventsWorkflow(opts EventsWorkflowOpts) (*EventWorkflow, error) {
	c, err := client.Dial(
		client.Options{
			HostPort: opts.Cfg.TemporalHost,
			Logger:   log.NewStructuredLogger(opts.Logger),
		},
	)
	if err != nil {
		return nil, err
	}

	return &EventWorkflow{
		cfg:                               opts.Cfg,
		cl:                                c,
		eventsActivity:                    opts.EventsActivity,
		db:                                opts.Gorm,
		redis:                             opts.Redis,
		hydrator:                          opts.Hydrator,
		channelsEventsWithOperationsCache: opts.ChannelsEventsWithOperationsCache,
		channelsCache:                     opts.ChannelsCache,
		channelsRepo:                      opts.ChannelsRepo,
		usersRepo:                         opts.UsersRepo,
	}, nil
}

type EventWorkflow struct {
	cfg                               config.Config
	cl                                client.Client
	eventsActivity                    *eventsActivity.Activity
	db                                *gorm.DB
	redis                             *redis.Client
	hydrator                          *hydrator.Hydrator
	channelsEventsWithOperationsCache *generic_cacher.GenericCacher[[]model.Event]
	channelsCache                     *generic_cacher.GenericCacher[channelmodel.Channel]
	channelsRepo                      channelsrepository.Repository
	usersRepo                         usersrepository.Repository
}

func (c *EventWorkflow) Execute(
	ctx context.Context,
	eventType model.EventType,
	data shared.EventData,
) error {
	if data.ChannelDBID == "" {
		plat := data.Platform
		if plat == "" {
			plat = platform.PlatformTwitch
		}

		user, err := c.usersRepo.GetByPlatformID(ctx, plat, data.ChannelID)
		if err != nil {
			return fmt.Errorf("resolve user by platform id: %w", err)
		}

		var channel channelmodel.Channel
		var channelErr error
		switch plat {
		case platform.PlatformKick:
			channel, channelErr = c.channelsRepo.GetByKickUserID(ctx, user.ID)
		default:
			channel, channelErr = c.channelsRepo.GetByTwitchUserID(ctx, user.ID)
		}
		if channelErr != nil {
			return fmt.Errorf("resolve channel: %w", channelErr)
		}

		data.ChannelDBID = channel.ID.String()
	}

	channel, err := c.channelsCache.Get(ctx, data.ChannelDBID)
	if err != nil {
		return err
	}
	if channel == channelmodel.Nil {
		return errors.New("channel not found")
	}

	channelEvents, err := c.channelsEventsWithOperationsCache.Get(ctx, data.ChannelDBID)
	if err != nil {
		return err
	}

	var eventTypeExists bool
	for _, entity := range channelEvents {
		if entity.Type == eventType {
			eventTypeExists = true
			break
		}
	}

	if !eventTypeExists {
		return nil
	}

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
