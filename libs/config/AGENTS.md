# AGENTS.md — libs/config

Configuration management for Go services.

## OVERVIEW

Unified configuration loading for Go applications. Supports environment variables, config files, and defaults.

## STRUCTURE

```
libs/config/
├── *.go                     # Config types and loader
├── go.mod
└── ...
```

## USAGE

```go
import "libs/config"

cfg := config.Load()
// cfg.Database.URL
// cfg.Redis.Addr
```

## NOTES

- Environment-based config
- Sensible defaults
- Validates required fields
