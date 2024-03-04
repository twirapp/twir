package app

import (
	"github.com/satont/twir/apps/scheduler/internal/gorm"
	"github.com/satont/twir/apps/scheduler/internal/grpc_impl"
	"github.com/satont/twir/apps/scheduler/internal/services"
	"github.com/satont/twir/apps/scheduler/internal/timers"
	config "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/logger"
	twirsentry "github.com/satont/twir/libs/sentry"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/grpc/clients"
	"github.com/twirapp/twir/libs/grpc/parser"
	"github.com/twirapp/twir/libs/grpc/tokens"
	"github.com/twirapp/twir/libs/uptrace"
	"go.uber.org/fx"
)

const service = "scheduler"

var App = fx.Module(
	service,
	fx.Provide(
		config.NewFx,
		twirsentry.NewFx(twirsentry.NewFxOpts{Service: service}),
		logger.NewFx(logger.Opts{Service: service}),
		uptrace.NewFx(service),
		func(c config.Config) parser.ParserClient {
			return clients.NewParser(c.AppEnv)
		},
		func(c config.Config) tokens.TokensClient {
			return clients.NewTokens(c.AppEnv)
		},
		gorm.New,
		services.NewRoles,
		services.NewCommands,
		buscore.NewNatsBusFx(service),
	),
	fx.Invoke(
		uptrace.NewFx(service),
		grpc_impl.New,
		timers.NewEmotes,
		timers.NewOnlineUsers,
		timers.NewStreams,
		timers.NewCommandsAndRoles,
		timers.NewBannedChannels,
		timers.NewWatched,
		func(l logger.Logger) {
			l.Info("Started")
		},
	),
)
