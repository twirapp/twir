package grpc_impl

import (
	"context"
	"github.com/satont/tsuwari/apps/scheduler/internal/timers"
	"github.com/satont/tsuwari/apps/scheduler/internal/types"
	"github.com/satont/tsuwari/libs/grpc/generated/scheduler"
	"google.golang.org/protobuf/types/known/emptypb"
)

type schedulerGrpc struct {
	scheduler.UnimplementedSchedulerServer

	commandsService *timers.DefaultCommandsTimer
	services        *types.Services
}

func NewGrpcImpl(cmds *timers.DefaultCommandsTimer, services *types.Services) *schedulerGrpc {
	return &schedulerGrpc{
		commandsService: cmds,
		services:        services,
	}
}

func (c *schedulerGrpc) CreateDefaultCommandsAndRoles(
	ctx context.Context,
	req *scheduler.CreateDefaultCommandsAndRolesRequest,
) (*emptypb.Empty, error) {
	cmds, err := c.services.Grpc.Parser.GetDefaultCommands(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, err
	}

	err = c.commandsService.CreateCommandsAndRoles(cmds, req.UserId)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
