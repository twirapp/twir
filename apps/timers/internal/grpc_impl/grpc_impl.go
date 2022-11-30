package grpc_impl

import (
	"context"

	"github.com/satont/tsuwari/apps/timers/internal/scheduler"
	"github.com/satont/tsuwari/apps/timers/internal/types"
	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/grpc/generated/timers"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

type TimersGrpcServerOpts struct {
	Db        *gorm.DB
	Logger    *zap.Logger
	Scheduler *scheduler.Scheduler
}

type TimersGrpcServer struct {
	timers.UnimplementedTimersServer

	db        *gorm.DB
	logger    *zap.Logger
	scheduler *scheduler.Scheduler
}

func New(opts *TimersGrpcServerOpts) *TimersGrpcServer {
	return &TimersGrpcServer{
		db:        opts.Db,
		logger:    opts.Logger,
		scheduler: opts.Scheduler,
	}
}

func (c *TimersGrpcServer) AddTimerToQueue(
	ctx context.Context,
	data *timers.Request,
) (*emptypb.Empty, error) {
	timer := &model.ChannelsTimers{}
	if err := c.db.Where(`"id" = ?`, data.TimerId).Preload("Responses").Take(timer).Error; err != nil {
		c.logger.Sugar().Error(err)
		return &emptypb.Empty{}, nil
	}
	c.scheduler.AddTimer(&types.Timer{
		Model:     timer,
		SendIndex: 0,
	})

	return &emptypb.Empty{}, nil
}

func (c *TimersGrpcServer) RemoveTimerFromQueue(
	ctx context.Context,
	data *timers.Request,
) (*emptypb.Empty, error) {
	c.scheduler.RemoveTimer(data.TimerId)

	return &emptypb.Empty{}, nil
}
