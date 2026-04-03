# AGENTS.md — apps/timers

Chat timer message dispatcher.

## OVERVIEW

Sends periodic messages to chat channels based on configured timers. Integrates with parser for command execution in timers.

## STRUCTURE

```
apps/timers/
├── cmd/
│   └── main.go              # Entry point
├── internal/
│   └── ...
├── go.mod
└── Dockerfile
```

## ENTRY POINTS

| Type | Path          | Purpose                 |
| ---- | ------------- | ----------------------- |
| Main | `cmd/main.go` | Timer service bootstrap |

## KEY COMMANDS

```bash
# Run locally
go run ./cmd/main.go

# Build
bun cli build timers
```

## DEPENDENCIES

- Postgres (timer configs)
- bots (chat sending)
- parser (command execution)

## NOTES

- Timers configured per channel
- Supports interval and line-based triggers
- Respects chat activity
