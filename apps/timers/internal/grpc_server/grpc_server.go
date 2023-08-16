package grpc_server

import (
	"context"
	"fmt"
	"github.com/satont/twir/apps/timers/internal/queue"
	"github.com/satont/twir/libs/grpc/generated/timers"
	"github.com/satont/twir/libs/grpc/servers"
	"github.com/satont/twir/libs/logger"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/protobuf/types/known/emptypb"
	"log/slog"
	"net"
	"time"
)

type server struct {
	timers.UnimplementedTimersServer

	queue *queue.Queue
}

func New(queue *queue.Queue, logger logger.Logger, lc fx.Lifecycle) error {
	server := &server{queue: queue}

	addr := fmt.Sprintf("0.0.0.0:%v", servers.TIMERS_SERVER_PORT)
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
	timers.RegisterTimersServer(grpcServer, server)

	lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				go func() {
					if err = grpcServer.Serve(lis); err != nil {
						panic(err)
					}
				}()
				logger.Info("Timers grpc server started", slog.String("addr", addr))
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

func (c *server) AddTimerToQueue(_ context.Context, t *timers.Request) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, c.queue.Add(t.TimerId)
}

func (c *server) RemoveTimerFromQueue(_ context.Context, t *timers.Request) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, c.queue.Remove(t.TimerId)
}
