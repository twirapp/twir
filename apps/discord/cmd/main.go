package main

import (
	"github.com/getsentry/sentry-go"
	"github.com/satont/twir/apps/discord/internal/discord_go"
	"github.com/satont/twir/apps/discord/internal/gorm"
	"github.com/satont/twir/apps/discord/internal/messages_updater"
	"github.com/satont/twir/apps/discord/internal/sended_messages_store"
	"github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/grpc/clients"
	"github.com/satont/twir/libs/grpc/generated/tokens"
	"github.com/satont/twir/libs/logger"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.NopLogger,
		fx.Provide(
			cfg.NewFx,
			gorm.New,
			sended_messages_store.New,
			messages_updater.New,
			discord_go.New,
			func(config cfg.Config) (*sentry.Client, error) {
				if config.SentryDsn == "" {
					return nil, nil
				}

				s, err := sentry.NewClient(
					sentry.ClientOptions{
						Dsn:              config.SentryDsn,
						AttachStacktrace: true,
					},
				)

				return s, err
			},
			func(config cfg.Config, s *sentry.Client) logger.Logger {
				l := logger.New(
					logger.Opts{
						Env:     config.AppEnv,
						Service: "discord",
						Sentry:  s,
					},
				)

				return l
			},
			func(config cfg.Config) tokens.TokensClient {
				return clients.NewTokens(config.AppEnv)
			},
		),
		fx.Invoke(
			gorm.New,
			discord_go.New,
			messages_updater.New,
		),
	).Run()
}
