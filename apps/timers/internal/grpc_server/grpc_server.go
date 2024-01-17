package grpc_server

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/satont/twir/apps/timers/internal/workflow"
	"github.com/satont/twir/libs/logger"
	"github.com/twirapp/twir/libs/grpc/constants"
	"github.com/twirapp/twir/libs/grpc/timers"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/protobuf/types/known/emptypb"
)

type server struct {
	timers.UnimplementedTimersServer

	workflow *workflow.Workflow
}

type Opts struct {
	fx.In

	Lc       fx.Lifecycle
	Logger   logger.Logger
	Workflow *workflow.Workflow
}

func New(opts Opts) error {
	s := &server{
		workflow: opts.Workflow,
	}

	addr := fmt.Sprintf("0.0.0.0:%v", constants.TIMERS_SERVER_PORT)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer(
		grpc.KeepaliveParams(
			keepalive.ServerParameters{
				MaxConnectionAge: 1 * time.Minute,
			},
		),
	)
	timers.RegisterTimersServer(grpcServer, s)

	opts.Lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				go func() {
					if err = grpcServer.Serve(lis); err != nil {
						panic(err)
					}
				}()
				opts.Logger.Info("Timers grpc server started", slog.String("addr", addr))
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

func (c *server) AddTimerToQueue(ctx context.Context, t *timers.Request) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, c.workflow.AddTimer(ctx, t.TimerId)
}

func (c *server) RemoveTimerFromQueue(ctx context.Context, t *timers.Request) (
	*emptypb.Empty,
	error,
) {
	return &emptypb.Empty{}, c.workflow.RemoveTimer(ctx, t.TimerId)
}
