package tracer

import (
	"go.opentelemetry.io/otel"
)

var CommandsTracer = otel.Tracer("message-handler")
