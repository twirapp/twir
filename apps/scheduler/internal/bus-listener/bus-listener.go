package bus_listener

import (
	"context"
	"log/slog"

	"github.com/twirapp/twir/apps/scheduler/internal/services"
	"github.com/twirapp/twir/libs/logger"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/scheduler"
	"go.uber.org/fx"
)

type schedulerListener struct {
	commandsService *services.Commands
	rolesService    *services.Roles
	bus             *buscore.Bus
	logger          logger.Logger
}

type Opts struct {
	fx.In
	Lc fx.Lifecycle

	Logger logger.Logger

	CommandsService *services.Commands
	RolesService    *services.Roles
	Bus             *buscore.Bus
}

func New(opts Opts) error {
	impl := &schedulerListener{
		commandsService: opts.CommandsService,
		rolesService:    opts.RolesService,
		bus:             opts.Bus,
		logger:          opts.Logger,
	}

	opts.Lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				impl.bus.Scheduler.CreateDefaultCommands.SubscribeGroup(
					"scheduler",
					impl.createDefaultCommands,
				)

				impl.bus.Scheduler.CreateDefaultRoles.SubscribeGroup(
					"scheduler",
					impl.createDefaultRoles,
				)

				return nil
			},
			OnStop: func(ctx context.Context) error {
				impl.bus.Scheduler.CreateDefaultCommands.Unsubscribe()
				impl.bus.Scheduler.CreateDefaultRoles.Unsubscribe()
				return nil
			},
		},
	)

	return nil
}

func (c *schedulerListener) createDefaultCommands(
	ctx context.Context,
	req scheduler.CreateDefaultCommandsRequest,
) (struct{}, error) {
	if err := c.commandsService.CreateDefaultCommands(ctx); err != nil {
		c.logger.Error("failed to create default commands", slog.Any("err", err))
		return struct{}{}, err
	}

	return struct{}{}, nil
}

func (c *schedulerListener) createDefaultRoles(
	ctx context.Context,
	req scheduler.CreateDefaultRolesRequest,
) (struct{}, error) {
	if err := c.rolesService.CreateDefaultRoles(ctx); err != nil {
		c.logger.Error("failed to create default roles", slog.Any("err", err))
		return struct{}{}, err
	}

	return struct{}{}, nil
}
