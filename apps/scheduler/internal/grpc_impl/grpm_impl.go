package grpc_impl

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"github.com/satont/twir/apps/scheduler/internal/services"
	"github.com/satont/twir/libs/logger"
	"github.com/twirapp/twir/libs/grpc/constants"
	"github.com/twirapp/twir/libs/grpc/scheduler"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type schedulerGrpc struct {
	scheduler.UnimplementedSchedulerServer

	commandsService *services.Commands
	rolesService    *services.Roles
}

type Opts struct {
	fx.In
	Lc fx.Lifecycle

	Logger logger.Logger

	CommandsService *services.Commands
	RolesService    *services.Roles
}

func New(opts Opts) error {
	impl := &schedulerGrpc{
		commandsService: opts.CommandsService,
		rolesService:    opts.RolesService,
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", constants.SCHEDULER_SERVER_PORT))
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer(grpc.StatsHandler(otelgrpc.NewServerHandler()))
	scheduler.RegisterSchedulerServer(
		grpcServer,
		impl,
	)

	opts.Lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				go grpcServer.Serve(lis)
				opts.Logger.Info("Grpc service started", slog.Int("port", constants.SCHEDULER_SERVER_PORT))
				return nil
			},
			OnStop: func(ctx context.Context) error {
				grpcServer.GracefulStop()
				return nil
			},
		},
	)

	return nil
}

func (c *schedulerGrpc) CreateDefaultCommands(
	ctx context.Context,
	req *scheduler.CreateDefaultCommandsRequest,
) (*emptypb.Empty, error) {
	if err := c.commandsService.CreateDefaultCommands(ctx, req.UsersIds); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (c *schedulerGrpc) CreateDefaultRoles(
	ctx context.Context,
	req *scheduler.CreateDefaultRolesRequest,
) (*emptypb.Empty, error) {
	if err := c.rolesService.CreateDefaultRoles(ctx, req.UsersIds); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
