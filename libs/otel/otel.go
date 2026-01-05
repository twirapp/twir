package otel

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	cfg "github.com/twirapp/twir/libs/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"google.golang.org/grpc/credentials/insecure"
)

type shutdownFunc func(context.Context) error

// parseHeaders parses comma-separated key=value pairs into a map
// Example: "authorization=Bearer token,x-api-key=secret" -> map[string]string{"authorization": "Bearer token", "x-api-key": "secret"}
func parseHeaders(headersStr string) map[string]string {
	if headersStr == "" {
		return nil
	}

	headers := make(map[string]string)
	pairs := strings.Split(headersStr, ",")
	
	for _, pair := range pairs {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}
		
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			headers[key] = value
		}
	}
	
	return headers
}

// cleanEndpoint removes protocol prefix from endpoint if present
// Example: "http://localhost:4317" -> "localhost:4317"
func cleanEndpoint(endpoint string) string {
	// First trim spaces
	endpoint = strings.TrimSpace(endpoint)
	
	// Remove http://, https://, or grpc:// prefix
	endpoint = strings.TrimPrefix(endpoint, "http://")
	endpoint = strings.TrimPrefix(endpoint, "https://")
	endpoint = strings.TrimPrefix(endpoint, "grpc://")
	
	return endpoint
}

// setupResource creates OpenTelemetry resource with service info
func setupResource(serviceName, environment string) (*resource.Resource, error) {
	return resource.New(
		context.Background(),
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
			semconv.DeploymentEnvironment(environment),
		),
	)
}

// setupTraceProvider creates and configures trace provider with OTLP exporter
func setupTraceProvider(ctx context.Context, res *resource.Resource, endpoint string, headers map[string]string, insecureConn bool) (*sdktrace.TracerProvider, error) {
	opts := []otlptracegrpc.Option{
		otlptracegrpc.WithEndpoint(endpoint),
	}

	if insecureConn {
		opts = append(opts, otlptracegrpc.WithTLSCredentials(insecure.NewCredentials()))
	}

	if len(headers) > 0 {
		opts = append(opts, otlptracegrpc.WithHeaders(headers))
	}

	traceExporter, err := otlptracegrpc.New(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	traceProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)

	return traceProvider, nil
}

// setupMeterProvider creates and configures meter provider with OTLP exporter
func setupMeterProvider(ctx context.Context, res *resource.Resource, endpoint string, headers map[string]string, insecureConn bool) (*sdkmetric.MeterProvider, error) {
	opts := []otlpmetricgrpc.Option{
		otlpmetricgrpc.WithEndpoint(endpoint),
	}

	if insecureConn {
		opts = append(opts, otlpmetricgrpc.WithTLSCredentials(insecure.NewCredentials()))
	}

	if len(headers) > 0 {
		opts = append(opts, otlpmetricgrpc.WithHeaders(headers))
	}

	metricExporter, err := otlpmetricgrpc.New(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create metrics exporter: %w", err)
	}

	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExporter)),
		sdkmetric.WithResource(res),
	)

	return meterProvider, nil
}

// New initializes OpenTelemetry with OTLP exporters
// nolint:ireturn
func New(config cfg.Config, serviceName string) (trace.Tracer, error) {
	// Skip initialization if OTLP endpoint is not configured
	if config.OtelEndpoint == "" {
		slog.Info("OpenTelemetry is not configured (OTEL_ENDPOINT is empty), skipping initialization")
		return otel.Tracer(serviceName), nil
	}

	// Clean endpoint (remove protocol prefix if present)
	endpoint := cleanEndpoint(config.OtelEndpoint)
	if endpoint == "" {
		return nil, fmt.Errorf("OTEL_ENDPOINT is invalid after cleaning: %q", config.OtelEndpoint)
	}

	ctx := context.Background()

	// Parse headers
	headers := parseHeaders(config.OtelHeaders)

	// Create resource
	res, err := setupResource(serviceName, config.AppEnv)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Setup trace provider if tracing is enabled
	if config.OtelTracingEnabled {
		traceProvider, err := setupTraceProvider(ctx, res, endpoint, headers, config.OtelInsecure)
		if err != nil {
			return nil, err
		}
		otel.SetTracerProvider(traceProvider)
	}

	// Setup meter provider if metrics are enabled
	if config.OtelMetricsEnabled {
		meterProvider, err := setupMeterProvider(ctx, res, endpoint, headers, config.OtelInsecure)
		if err != nil {
			return nil, err
		}
		otel.SetMeterProvider(meterProvider)
	}

	// Setup propagators for distributed tracing
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	logMsg := "OpenTelemetry initialized"
	logAttrs := []any{
		"service", serviceName,
		"endpoint", endpoint,
		"tracing", config.OtelTracingEnabled,
		"metrics", config.OtelMetricsEnabled,
	}
	if len(headers) > 0 {
		logAttrs = append(logAttrs, "headers", len(headers))
	}
	slog.Info(logMsg, logAttrs...)

	return otel.Tracer(serviceName), nil
}

// Shutdown gracefully shuts down OpenTelemetry providers
func Shutdown(ctx context.Context) error {
	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var errs []error

	if tp, ok := otel.GetTracerProvider().(*sdktrace.TracerProvider); ok {
		if err := tp.Shutdown(shutdownCtx); err != nil {
			errs = append(errs, fmt.Errorf("trace provider shutdown: %w", err))
		}
	}

	if mp, ok := otel.GetMeterProvider().(*sdkmetric.MeterProvider); ok {
		if err := mp.Shutdown(shutdownCtx); err != nil {
			errs = append(errs, fmt.Errorf("meter provider shutdown: %w", err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("shutdown errors: %v", errs)
	}

	return nil
}

// NewFx creates a tracer with fx lifecycle management
func NewFx(service string) func(config cfg.Config, lc fx.Lifecycle) trace.Tracer {
	return func(config cfg.Config, lc fx.Lifecycle) trace.Tracer {
		tracer, err := New(config, service)
		if err != nil {
			slog.Error("Failed to initialize OpenTelemetry", "error", err)
			// Return no-op tracer on error
			return otel.Tracer(service)
		}

		lc.Append(
			fx.Hook{
				OnStop: Shutdown,
			},
		)

		return tracer
	}
}
