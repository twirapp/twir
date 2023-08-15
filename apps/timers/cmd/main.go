package main

import (
	"fmt"
	"github.com/satont/twir/libs/grpc/generated/bots"
	"github.com/satont/twir/libs/grpc/generated/parser"
	"go.uber.org/fx"
	"golang.org/x/exp/slog"
	"net"
	"time"

	"github.com/satont/twir/libs/grpc/generated/tokens"

	"github.com/satont/twir/apps/timers/internal/scheduler"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	"github.com/satont/twir/libs/grpc/clients"
	"github.com/satont/twir/libs/grpc/servers"

	cfg "github.com/satont/twir/libs/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"

	"github.com/satont/twir/apps/timers/internal/grpc_impl"
	timersgrpc "github.com/satont/twir/libs/grpc/generated/timers"
)

func main() {
	fx.New(
		fx.Provide(
			func() (*cfg.Config, error) {
				return cfg.New()
			},
			func(config *cfg.Config) (*gorm.DB, error) {
				return gorm.Open(
					postgres.Open(config.DatabaseUrl), &gorm.Config{
						Logger: gormLogger.Default.LogMode(gormLogger.Silent),
					},
				)
			},
			func(config *cfg.Config) parser.ParserClient {
				return clients.NewParser(config.AppEnv)
			},
			func(config *cfg.Config) bots.BotsClient {
				return clients.NewBots(config.AppEnv)
			},
			func(config *cfg.Config) tokens.TokensClient {
				return clients.NewTokens(config.AppEnv)
			},
			fx.Annotate(
				grpc_impl.New,
				fx.As(new(timersgrpc.TimersServer)),
			),
			scheduler.New,
		),
		fx.Invoke(
			func(g timersgrpc.TimersServer) error {
				lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", servers.TIMERS_SERVER_PORT))
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
				timersgrpc.RegisterTimersServer(grpcServer, g)
				slog.Info("Grpc server started")
				return grpcServer.Serve(lis)
			},
		),
	).Run()
}
