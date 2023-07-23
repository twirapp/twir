package grpc_impl

import (
	"context"
	services2 "github.com/satont/twir/apps/scheduler/internal/services"
	"github.com/satont/twir/libs/grpc/generated/scheduler"
	"google.golang.org/protobuf/types/known/emptypb"
)

type SchedulerGrpc struct {
	scheduler.UnimplementedSchedulerServer

	Commands *services2.Commands
	Roles    *services2.Roles
}

func (c *SchedulerGrpc) CreateDefaultCommands(
	ctx context.Context,
	req *scheduler.CreateDefaultCommandsRequest,
) (*emptypb.Empty, error) {
	if err := c.Commands.CreateDefaultCommands(ctx, req.UsersIds); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (c *SchedulerGrpc) CreateDefaultRoles(
	ctx context.Context,
	req *scheduler.CreateDefaultRolesRequest,
) (*emptypb.Empty, error) {
	if err := c.Roles.CreateDefaultRoles(ctx, req.UsersIds); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
