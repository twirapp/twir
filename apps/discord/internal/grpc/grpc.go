package grpc

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/satont/twir/libs/grpc/generated/discord"
	"github.com/satont/twir/libs/grpc/servers"
	"github.com/satont/twir/libs/logger"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

type Opts struct {
	fx.In

	LC     fx.Lifecycle
	Logger logger.Logger
}

func New(opts Opts) (discord.DiscordServer, error) {
	service := &Impl{}

	grpcNetListener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", servers.DISCORD_SERVER_PORT))
	if err != nil {
		return nil, err
	}
	grpcServer := grpc.NewServer(
		grpc.KeepaliveParams(
			keepalive.ServerParameters{
				MaxConnectionAge: 1 * time.Minute,
			},
		),
	)

	discord.RegisterDiscordServer(grpcServer, service)

	opts.LC.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				go grpcServer.Serve(grpcNetListener)
				opts.Logger.Info("Grpc server is running")
				return nil
			},
			OnStop: func(ctx context.Context) error {
				grpcServer.GracefulStop()
				return nil
			},
		},
	)

	return service, nil
}

type Impl struct {
	discord.UnimplementedDiscordServer
}
