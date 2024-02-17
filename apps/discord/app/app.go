package app

import (
	"github.com/satont/twir/apps/discord/internal/discord_go"
	"github.com/satont/twir/apps/discord/internal/gorm"
	"github.com/satont/twir/apps/discord/internal/grpc"
	"github.com/satont/twir/apps/discord/internal/messages_updater"
	"github.com/satont/twir/apps/discord/internal/redis"
	"github.com/satont/twir/apps/discord/internal/sended_messages_store"
	cfg "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/logger"
	twirsentry "github.com/satont/twir/libs/sentry"
	"github.com/twirapp/twir/libs/grpc/clients"
	"github.com/twirapp/twir/libs/grpc/tokens"
	"github.com/twirapp/twir/libs/uptrace"
	"go.uber.org/fx"
)

var App = fx.Module(
	"discord",
	fx.Provide(
		cfg.NewFx,
		gorm.New,
		redis.New,
		sended_messages_store.New,
		messages_updater.New,
		discord_go.New,
		twirsentry.NewFx(twirsentry.NewFxOpts{Service: "discord"}),
		uptrace.NewFx("discord"),
		logger.NewFx(
			logger.Opts{
				Service: "discord",
			},
		),
		func(config cfg.Config) tokens.TokensClient {
			return clients.NewTokens(config.AppEnv)
		},
		grpc.New,
	),
	fx.Invoke(
		uptrace.NewFx("discord"),
		redis.New,
		gorm.New,
		// discord_go.New,
		messages_updater.New,
		grpc.New,
	),
)
