# AGENTS.md — apps/scheduler

Cron-like scheduler for timed tasks.

## OVERVIEW

Schedules and executes periodic tasks across the platform. Manages timers, scheduled commands, and background jobs.

## STRUCTURE

```
apps/scheduler/
├── cmd/
│   └── main.go              # Entry point
├── internal/
│   └── ...
├── go.mod
└── Dockerfile
```

## ENTRY POINTS

| Type | Path          | Purpose             |
| ---- | ------------- | ------------------- |
| Main | `cmd/main.go` | Scheduler bootstrap |

## KEY COMMANDS

```bash
# Run locally
go run ./cmd/main.go

# Build
bun cli build scheduler
```

## DEPENDENCIES

- Postgres (scheduled tasks)
- Redis (distributed locking)
- Other apps (via gRPC)

## NOTES

- Distributed scheduling with Redis locks
- Supports cron expressions
- Reschedules missed tasks
