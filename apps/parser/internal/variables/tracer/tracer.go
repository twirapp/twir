package tracer

import (
	"go.opentelemetry.io/otel"
)

var VariablesTracer = otel.Tracer("message-handler")
