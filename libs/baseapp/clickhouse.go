package baseapp

import (
	"context"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	config "github.com/satont/twir/libs/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
	"go.opentelemetry.io/otel/trace"
)

type ClickhouseClient struct {
	conn driver.Conn
}

func NewClickHouse(appName string) func(config.Config) (*ClickhouseClient, error) {
	return func(cfg config.Config) (*ClickhouseClient, error) {
		options, err := clickhouse.ParseDSN(cfg.ClickhouseUrl)
		if err != nil {
			return nil, err
		}

		conn, err := clickhouse.Open(
			&clickhouse.Options{
				Addr: options.Addr,
				Auth: options.Auth,
				Settings: clickhouse.Settings{
					"max_execution_time":    30,
					"async_insert":          "1",
					"wait_for_async_insert": "1",
				},
				DialTimeout: time.Second * 5,
				Compression: &clickhouse.Compression{
					Method: clickhouse.CompressionLZ4,
				},
				Debug:                false,
				BlockBufferSize:      10,
				MaxCompressionBuffer: 10240,
				ClientInfo: clickhouse.ClientInfo{
					Products: []struct {
						Name    string
						Version string
					}{
						{Name: appName},
					},
				},
			},
		)

		pingCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		if err := conn.Ping(pingCtx); err != nil {
			return nil, err
		}

		return &ClickhouseClient{conn}, nil
	}
}

var clickhouseTracer = otel.Tracer("go-clickhouse")

func (c *ClickhouseClient) getCtx(
	ctx context.Context,
	operation,
	query string,
) (context.Context, trace.Span) {
	span := trace.SpanFromContext(ctx)

	if !span.IsRecording() {
		return ctx, span
	}

	ctx, newSpan := clickhouseTracer.Start(ctx, operation, trace.WithSpanKind(trace.SpanKindClient))
	ctx = clickhouse.Context(ctx, clickhouse.WithSpan(span.SpanContext()))

	newSpan.SetAttributes(
		attribute.String("db.system", "clickhouse"),
		semconv.DBStatement(query),
	)

	return ctx, newSpan
}

func (c *ClickhouseClient) Query(
	ctx context.Context,
	query string,
	args ...interface{},
) (driver.Rows, error) {
	ctx, span := c.getCtx(ctx, "query", query)
	defer span.End()

	return c.conn.Query(ctx, query, args...)
}

func (c *ClickhouseClient) QueryRow(
	ctx context.Context,
	query string,
	args ...interface{},
) driver.Row {
	ctx, span := c.getCtx(ctx, "queryRow", query)
	defer span.End()

	return c.conn.QueryRow(ctx, query, args...)
}

func (c *ClickhouseClient) Exec(ctx context.Context, query string, args ...interface{}) error {
	ctx, span := c.getCtx(ctx, "exec", query)
	defer span.End()

	return c.conn.Exec(ctx, query, args...)
}

func (c *ClickhouseClient) PrepareBatch(ctx context.Context, query string) (driver.Batch, error) {
	ctx, span := c.getCtx(ctx, "prepareBatch", query)
	defer span.End()

	return c.conn.PrepareBatch(ctx, query)
}
