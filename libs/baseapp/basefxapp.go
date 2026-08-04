package baseapp

import (
	"context"
	"fmt"
	"time"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	config "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/entities/platform"
)

type Opts struct {
	AppName string

	WithAudit bool
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

func createPgxPool(cfg config.Config) (*pgxpool.Pool, error) {
	connConfig, err := pgxpool.ParseConfig(cfg.DatabaseUrl)
	if err != nil {
		return nil, err
	}

	connConfig.ConnConfig.Tracer = otelpgx.NewTracer()
	connConfig.MaxConnLifetime = 5 * time.Minute
	connConfig.MaxConnIdleTime = 2 * time.Minute
	connConfig.MaxConns = 100
	connConfig.MinConns = 1
	connConfig.HealthCheckPeriod = 30 * time.Second
	connConfig.ConnConfig.Config.ConnectTimeout = 5 * time.Second
	// connConfig.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	connConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		platformType, err := conn.LoadType(ctx, "platform")
		if err != nil {
			return fmt.Errorf("load platform type: %w", err)
		}
		conn.TypeMap().RegisterType(platformType)

		platformArrayType, err := conn.LoadType(ctx, "_platform")
		if err != nil {
			return fmt.Errorf("load _platform type: %w", err)
		}
		conn.TypeMap().RegisterType(platformArrayType)

		conn.TypeMap().RegisterDefaultPgType(platform.Platform(""), "platform")
		conn.TypeMap().RegisterDefaultPgType([]platform.Platform{}, "_platform")
		return nil
	}

	pool, err := pgxpool.NewWithConfig(
		context.Background(),
		connConfig,
	)
	if err != nil {
		return nil, err
	}

	return pool, nil
}

func newTransactionManager(pool *pgxpool.Pool) (trm.Manager, error) {
	return manager.New(trmpgx.NewDefaultFactory(pool))
}
