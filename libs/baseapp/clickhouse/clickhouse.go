package clickhouse

import (
	"context"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/column"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
	"go.opentelemetry.io/otel/trace"
)

type ClickhouseClient struct {
	conn driver.Conn
}

func New(conn driver.Conn) *ClickhouseClient {
	return &ClickhouseClient{conn: conn}
}

var clickhouseTracer = otel.Tracer("go-clickhouse")

func (c *ClickhouseClient) getCtx(
	ctx context.Context,
	operation,
	query string,
) (context.Context, trace.Span) {
	ctx, newSpan := clickhouseTracer.Start(ctx, operation, trace.WithSpanKind(trace.SpanKindClient))
	ctx = clickhouse.Context(ctx, clickhouse.WithSpan(newSpan.SpanContext()))

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

type clickhouseTracedBatch struct {
	batch driver.Batch
	span  trace.Span
}

// Append adds a row to the batch and records a trace event.
func (b *clickhouseTracedBatch) Append(args ...interface{}) error {
	b.span.AddEvent("batch.append")
	return b.batch.Append(args...)
}

// AppendStruct adds a struct to the batch. This is a clickhouse-go extension,
// so we must type-assert the underlying batch.
func (b *clickhouseTracedBatch) AppendStruct(v any) error {
	b.span.AddEvent("batch.appendStruct")
	return b.batch.AppendStruct(v)
}

// Column returns a specific column in the batch.
// This is a simple delegation as it doesn't affect the span's lifecycle.
func (b *clickhouseTracedBatch) Column(i int) driver.BatchColumn {
	b.span.AddEvent("batch.column")
	return b.batch.Column(i)
}

// Flush sends any buffered data to the server.
func (b *clickhouseTracedBatch) Flush() error {
	b.span.AddEvent("batch.flush")
	return b.batch.Flush()
}

// IsSent reports whether the batch has been sent.
// This is a simple delegation.
func (b *clickhouseTracedBatch) IsSent() bool {
	return b.batch.IsSent()
}

// Rows returns the number of rows in the batch.
// We delegate this to the underlying batch, as it's the source of truth.
func (b *clickhouseTracedBatch) Rows() int {
	return b.batch.Rows()
}

// Send executes the batch insert and ends the trace span.
func (b *clickhouseTracedBatch) Send() error {
	defer b.span.End()

	// Use the underlying batch's row count for the most accurate number.
	b.span.SetAttributes(attribute.Int("db.batch.rows_appended", b.batch.Rows()))

	err := b.batch.Send()
	if err != nil {
		b.span.RecordError(err)
		b.span.SetStatus(codes.Error, err.Error())
	}
	return err
}

// Abort cancels the batch operation and ends the trace span with an error status.
func (b *clickhouseTracedBatch) Abort() error {
	defer b.span.End()
	b.span.SetStatus(codes.Error, "batch aborted by user")

	// The standard driver.Batch interface has an Abort method.
	return b.batch.Abort()
}

func (b *clickhouseTracedBatch) Columns() []column.Interface {
	return b.Columns()
}

// Close is not a standard method on the driver.Batch interface.
// A batch is terminated by `Send` or `Abort`. We'll treat Close as an Abort.
func (b *clickhouseTracedBatch) Close() error {
	// Closing a batch is equivalent to aborting it.
	return b.Abort()
}

func (c *ClickhouseClient) PrepareBatch(ctx context.Context, query string) (driver.Batch, error) {
	// Note: We are NOT deferring span.End() here.
	ctx, span := c.getCtx(ctx, "prepareBatch", query)

	batch, err := c.conn.PrepareBatch(ctx, query)
	if err != nil {
		// If preparation fails, we must end the span here before returning.
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		span.End()
		return nil, err
	}

	// Return our wrapper, which now controls the span's lifecycle.
	return &clickhouseTracedBatch{
		batch: batch,
		span:  span,
	}, nil
}
