# AGENTS.md — apps/events

Workflow engine for channel events and alerts.

## OVERVIEW

Processes Twitch events (follows, subs, raids, etc.) and triggers configured workflows. Handles chat alerts and notification routing.

## STRUCTURE

```
apps/events/
├── cmd/
│   └── main.go              # Entry point
├── internal/
│   ├── activities/          # Workflow activities
│   │   └── events/          # Event handlers
│   ├── chat_alerts/         # Alert processing
│   └── ...
├── go.mod
└── Dockerfile
```

## ENTRY POINTS

| Type | Path          | Purpose                   |
| ---- | ------------- | ------------------------- |
| Main | `cmd/main.go` | Workflow worker bootstrap |

## KEY COMMANDS

```bash
# Run locally
go run ./cmd/main.go

# Build
bun cli build events
```

## DEPENDENCIES

- Temporal (workflow orchestration)
- Postgres (event configs)
- Redis (state)
- bots (for chat alerts)

## NOTES

- Uses Temporal for durable workflows
- Event configurations from dashboard
