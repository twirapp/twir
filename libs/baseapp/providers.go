package baseapp

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"github.com/getsentry/sentry-go"
	"github.com/goforj/wire"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/twirapp/kv"
	kvredis "github.com/twirapp/kv/stores/redis"
	"github.com/twirapp/twir/libs/audit"
	"github.com/twirapp/twir/libs/audit/recorder"
	"github.com/twirapp/twir/libs/baseapp/clickhouse"
	"github.com/twirapp/twir/libs/baseapp/lifecycle"
	buscore "github.com/twirapp/twir/libs/bus-core"
	config "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/logger"
	"github.com/twirapp/twir/libs/otel"
	auditpubsub "github.com/twirapp/twir/libs/pubsub/audit-logs"
	auditrepository "github.com/twirapp/twir/libs/repositories/audit_logs"
	auditrepositoryclickhouse "github.com/twirapp/twir/libs/repositories/audit_logs/datasources/clickhouse"
	twirsentry "github.com/twirapp/twir/libs/sentry"
	"github.com/twirapp/twir/libs/uptime"
	otelapi "go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

type serviceName string

// Base contains the common dependencies shared by Twir Go applications.
type Base struct {
	Lifecycle       *lifecycle.Lifecycle
	Config          config.Config
	Tracer          trace.Tracer
	Redis           *redis.Client
	PgxPool         *pgxpool.Pool
	TrManager       trm.Manager
	Gorm            *gorm.DB
	ClickHouse      *clickhouse.ClickhouseClient
	Bus             *buscore.Bus
	AuditPubSub     auditpubsub.PubSub
	AuditRepository auditrepository.Repository
	AuditRecorder   audit.Recorder
	KV              kv.KV
	Logger          *slog.Logger
	Sentry          *sentry.Client
	UptimeReporter  *uptime.Reporter
}

var providerSet = wire.NewSet(
	lifecycle.New,
	config.NewFx,
	provideServiceName,
	newLogger,
	newTracer,
	newRedis,
	createPgxPool,
	newTransactionManager,
	newGorm,
	newClickHouse,
	newBus,
	newAuditPubSub,
	newAuditRepository,
	newAuditRecorder,
	newKV,
	newSentry,
	newUptimeReporter,
)

func provideServiceName(opts Opts) serviceName {
	return serviceName(opts.AppName)
}

func newLogger(name serviceName) *slog.Logger {
	result := logger.New(logger.Options{AppName: string(name), Level: slog.LevelInfo})
	slog.SetDefault(result)
	return result
}

func newTracer(cfg config.Config, name serviceName, lc *lifecycle.Lifecycle) trace.Tracer {
	tracer, err := otel.New(cfg, string(name))
	if err != nil {
		slog.Error("Failed to initialize OpenTelemetry", logger.Error(err))
		tracer = otelapi.Tracer(string(name))
	}

	lc.Append(lifecycle.Hook{OnStop: otel.Shutdown})
	return tracer
}

func newClickHouse(cfg config.Config, name serviceName) (*clickhouse.ClickhouseClient, error) {
	return NewClickHouse(string(name))(cfg)
}

func newBus(cfg config.Config, name serviceName) (*buscore.Bus, error) {
	return buscore.NewNatsBusFx(string(name))(cfg)
}

func newAuditPubSub(bus *buscore.Bus, lc *lifecycle.Lifecycle) auditpubsub.PubSub {
	pubsub := auditpubsub.NewBusPubSub(bus)
	lc.Append(lifecycle.Hook{
		OnStart: func(context.Context) error {
			return pubsub.Start()
		},
		OnStop: func(context.Context) error {
			pubsub.Stop()
			return nil
		},
	})
	return pubsub
}

func newAuditRepository(client *clickhouse.ClickhouseClient) auditrepository.Repository {
	return auditrepositoryclickhouse.NewFx(client)
}

func newAuditRecorder(repository auditrepository.Repository, pubsub auditpubsub.PubSub) audit.Recorder {
	databaseRecorder := recorder.NewDatabase(repository)
	pubsubRecorder := recorder.NewPubSub(pubsub)
	return recorder.NewFanout(databaseRecorder, pubsubRecorder)
}

func newKV(client *redis.Client) kv.KV {
	return kvredis.New(client)
}

func newSentry(cfg config.Config, name serviceName, lc *lifecycle.Lifecycle) (*sentry.Client, error) {
	client, err := twirsentry.New(cfg.SentryDsn, string(name))
	if err != nil {
		return nil, err
	}

	lc.Append(lifecycle.Hook{OnStop: func(context.Context) error {
		sentry.Flush(2 * time.Second)
		return nil
	}})
	return client, nil
}

func newUptimeReporter(
	client *redis.Client,
	name serviceName,
	lc *lifecycle.Lifecycle,
) (*uptime.Reporter, error) {
	reporter, err := uptime.NewReporter(client, uptime.ReporterOpts{ServiceName: string(name)})
	if err != nil {
		return nil, fmt.Errorf("create uptime reporter: %w", err)
	}

	reporterContext, cancel := context.WithCancel(context.Background())
	lc.Append(lifecycle.Hook{
		OnStart: func(context.Context) error {
			reporter.Start(reporterContext)
			return nil
		},
		OnStop: func(context.Context) error {
			cancel()
			return nil
		},
	})

	return reporter, nil
}
