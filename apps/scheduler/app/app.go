package app

import (
	"github.com/satont/twir/apps/scheduler/internal/gorm"
	"github.com/satont/twir/apps/scheduler/internal/grpc_impl"
	"github.com/satont/twir/apps/scheduler/internal/services"
	"github.com/satont/twir/apps/scheduler/internal/timers"
	config "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/grpc/clients"
	"github.com/satont/twir/libs/grpc/generated/emotes_cacher"
	"github.com/satont/twir/libs/grpc/generated/parser"
	"github.com/satont/twir/libs/grpc/generated/tokens"
	"github.com/satont/twir/libs/logger"
	"github.com/satont/twir/libs/pubsub"
	twirsentry "github.com/satont/twir/libs/sentry"
	"go.uber.org/fx"
)

const service = "scheduler"

var App = fx.Module(
	service,
	fx.NopLogger,
	fx.Provide(
		config.NewFx,
		twirsentry.NewFx(twirsentry.NewFxOpts{Service: service}),
		logger.NewFx(logger.Opts{Service: service}),
		func(c config.Config) parser.ParserClient {
			return clients.NewParser(c.AppEnv)
		},
		func(c config.Config) tokens.TokensClient {
			return clients.NewTokens(c.AppEnv)
		},
		func(c config.Config) emotes_cacher.EmotesCacherClient {
			return clients.NewEmotesCacher(c.AppEnv)
		},
		gorm.New,
		func(c config.Config) (*pubsub.PubSub, error) {
			return pubsub.NewPubSub(c.RedisUrl)
		},
		services.NewRoles,
		services.NewCommands,
	),
	fx.Invoke(
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
