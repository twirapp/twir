package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/satont/twir/apps/eventsub/internal/manager"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/grpc/constants"
	"github.com/twirapp/twir/libs/grpc/eventsub"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

type EventSubGrpcImpl struct {
	eventsub.UnimplementedEventSubServer

	eventSubClient *manager.Manager
	gorm           *gorm.DB
}

type Opts struct {
	fx.In
	Lc fx.Lifecycle

	Manager *manager.Manager
	Gorm    *gorm.DB
}

func New(opts Opts) (*EventSubGrpcImpl, error) {
	impl := &EventSubGrpcImpl{
		eventSubClient: opts.Manager,
		gorm:           opts.Gorm,
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", constants.EVENTSUB_SERVER_PORT))
	if err != nil {
		return nil, err
	}
	grpcServer := grpc.NewServer(grpc.StatsHandler(otelgrpc.NewServerHandler()))
	eventsub.RegisterEventSubServer(grpcServer, impl)

	opts.Lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				go func() {
					if err := grpcServer.Serve(lis); err != nil {
						panic(err)
					}
				}()

				return nil
			},
			OnStop: func(ctx context.Context) error {
				grpcServer.Stop()
				return nil
			},
		},
	)

	return impl, nil
}

func (c *EventSubGrpcImpl) SubscribeToEvents(
	ctx context.Context, msg *eventsub.SubscribeToEventsRequest,
) (*emptypb.Empty, error) {
	channel := model.Channels{}
	err := c.gorm.
		WithContext(ctx).
		Where(
			`"id" = ?`,
			msg.GetChannelId(),
		).First(&channel).Error
	if err != nil {
		return nil, status.Error(codes.NotFound, "channel not found")
	}

	if err := c.eventSubClient.SubscribeToNeededEvents(
		ctx,
		msg.GetChannelId(),
		channel.BotID,
	); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
