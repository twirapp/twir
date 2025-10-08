package baseapp

import (
	"context"
	"log/slog"
	"time"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	buscore "github.com/twirapp/twir/libs/bus-core"
	config "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/logger"
	"github.com/twirapp/twir/libs/logger/audit"
	auditlogs "github.com/twirapp/twir/libs/pubsub/audit-logs"
	auditlogsrepository "github.com/twirapp/twir/libs/repositories/audit_logs"
	auditlogsrepositoryclickhouse "github.com/twirapp/twir/libs/repositories/audit_logs/datasources/clickhouse"
	twirsentry "github.com/twirapp/twir/libs/sentry"
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
			uptrace.NewFx(opts.AppName),
			newRedis,
			newPgxPool,
			newGorm,
			NewClickHouse(opts.AppName),
			buscore.NewNatsBusFx(opts.AppName),
			fx.Annotate(
				auditlogs.NewBusPubSubFx,
				fx.As(new(auditlogs.PubSub)),
			),
			fx.Annotate(
				auditlogsrepositoryclickhouse.NewFx,
				fx.As(new(auditlogsrepository.Repository)),
			),
			fx.Annotate(
				audit.NewDatabaseFx,
				fx.As(new(slog.Handler)),
				fx.ResultTags(`group:"slog-handlers"`),
			),
			fx.Annotate(
				audit.NewPubsubFx(),
				fx.As(new(slog.Handler)),
				fx.ResultTags(`group:"slog-handlers"`),
			),
			logger.NewFx(
				logger.Opts{
					Service: opts.AppName,
					Level:   slog.LevelInfo,
				},
			),
		),
		fx.Invoke(uptrace.NewFx(opts.AppName)),
	)
}

func newRedis(cfg config.Config) (*redis.Client, error) {
	redisOpts, err := redis.ParseURL(cfg.RedisUrl)
	if err != nil {
		return nil, err
	}
	redisClient := redis.NewClient(redisOpts)

	if err := redisotel.InstrumentTracing(redisClient); err != nil {
		return nil, err
	}

	if err := redisotel.InstrumentMetrics(redisClient); err != nil {
		return nil, err
	}

	return redisClient, nil
}

type PgxResult struct {
	fx.Out

	PgxPool   *pgxpool.Pool
	TrManager trm.Manager
}

func newPgxPool(cfg config.Config) (PgxResult, error) {
	connConfig, err := pgxpool.ParseConfig(cfg.DatabaseUrl)
	if err != nil {
		return PgxResult{}, err
	}

	connConfig.ConnConfig.Tracer = otelpgx.NewTracer()
	connConfig.MaxConnLifetime = 5 * time.Minute
	connConfig.MaxConnIdleTime = 2 * time.Minute
	connConfig.MaxConns = 100
	connConfig.MinConns = 1
	connConfig.HealthCheckPeriod = 30 * time.Second
	connConfig.ConnConfig.Config.ConnectTimeout = 5 * time.Second
	connConfig.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	pool, err := pgxpool.NewWithConfig(
		context.Background(),
		connConfig,
	)
	if err != nil {
		return PgxResult{}, err
	}

	trManager, err := manager.New(trmpgx.NewDefaultFactory(pool))
	if err != nil {
		return PgxResult{}, err
	}

	return PgxResult{
		PgxPool:   pool,
		TrManager: trManager,
	}, nil
}
