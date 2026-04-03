# AGENTS.md — libs/logger

Structured logging for Go services.

## OVERVIEW

Centralized logging with structured output. Provides consistent log formatting across all Go applications.

## STRUCTURE

```
libs/logger/
├── *.go                     # Logger implementation
├── go.mod
└── ...
```

## USAGE

```go
import "libs/logger"

// Error logging
logger.Error(err)

// With context
logger.InfoContext(ctx, "message", logger.String("key", "value"))
```

## CONVENTIONS

- Use `*Context` methods when ctx available
- Use `logger.String()`, `logger.Int()` for fields
- Don't use inline arguments

## NOTES

- Structured JSON output in production
- Colored output in development
- Integrated with OpenTelemetry
