package baseapp

import (
	"github.com/redis/go-redis/v9"
	config "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/logger"
	auditlogs "github.com/satont/twir/libs/pubsub/audit-logs"
	twirsentry "github.com/satont/twir/libs/sentry"
	"github.com/twirapp/twir/libs/uptrace"
	"go.uber.org/fx"
)

type Opts struct {
	AppName string

	WithAudit bool
}

func CreateBaseApp(opts Opts) fx.Option {
	return fx.Options(
		fx.Provide(
			config.NewFx,
			twirsentry.NewFx(twirsentry.NewFxOpts{Service: opts.AppName}),
			logger.NewFx(
				logger.Opts{
					Service: opts.AppName,
				},
			),
			fx.Annotate(
				auditlogs.NewBusPubSub,
				fx.As(new(auditlogs.PubSub)),
			),
			newRedis,
			newGorm(opts.WithAudit),
			uptrace.NewFx(opts.AppName),
		),
	)
}

func newRedis(cfg config.Config) (*redis.Client, error) {
	redisOpts, err := redis.ParseURL(cfg.RedisUrl)
	if err != nil {
		return nil, err
	}
	redisClient := redis.NewClient(redisOpts)
	return redisClient, nil
}
