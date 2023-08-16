package grpc_impl

import (
	"context"
	"go.uber.org/fx"
	"log/slog"

	"github.com/satont/twir/apps/timers/internal/scheduler"
	"github.com/satont/twir/apps/timers/internal/types"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/grpc/generated/timers"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

type TimersGrpcServer struct {
	timers.UnimplementedTimersServer

	db        *gorm.DB
	scheduler *scheduler.Scheduler
}

type TimersGrpcServerOpts struct {
	fx.In

	Db        *gorm.DB
	Scheduler *scheduler.Scheduler
}

func New(opts TimersGrpcServerOpts) timers.TimersServer {
	return &TimersGrpcServer{
		db:        opts.Db,
		scheduler: opts.Scheduler,
	}
}

func (c *TimersGrpcServer) AddTimerToQueue(
	ctx context.Context,
	data *timers.Request,
) (*emptypb.Empty, error) {
	timer := model.ChannelsTimers{}
	if err := c.db.
		WithContext(ctx).
		Where(`"id" = ?`, data.TimerId).
		Preload("Responses").
		Take(&timer).Error; err != nil {
		slog.Error(err.Error())
		return &emptypb.Empty{}, nil
	}
	c.scheduler.AddTimer(
		types.Timer{
			Model:     timer,
			SendIndex: 0,
		},
	)

	return &emptypb.Empty{}, nil
}

func (c *TimersGrpcServer) RemoveTimerFromQueue(
	ctx context.Context,
	data *timers.Request,
) (*emptypb.Empty, error) {
	c.scheduler.RemoveTimer(data.TimerId)

	return &emptypb.Empty{}, nil
}
