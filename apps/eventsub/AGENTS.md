# AGENTS.md — apps/eventsub

Twitch EventSub webhook handler and WebSocket manager.

## OVERVIEW

Receives Twitch EventSub notifications via webhooks and WebSocket. Processes events (streams, followers, bans) and forwards to events service for workflow execution.

## STRUCTURE

```
apps/eventsub/
├── cmd/
│   └── main.go              # Entry point
├── internal/
│   ├── manager/             # EventSub management
│   │   ├── websocket.go     # WebSocket connections
│   │   └── on_start.go      # Startup logic
│   └── ...
├── go.mod
└── Dockerfile
```

## ENTRY POINTS

| Type | Path          | Purpose                 |
| ---- | ------------- | ----------------------- |
| Main | `cmd/main.go` | HTTP server bootstrap   |
| HTTP | `/webhook`    | Twitch webhook endpoint |

## KEY COMMANDS

```bash
# Run locally
go run ./cmd/main.go

# Build
bun cli build eventsub
```

## DEPENDENCIES

- Twitch API (EventSub subscriptions)
- Postgres (subscription state)
- events (workflow triggers)

## NOTES

- Handles webhook verification
- Supports WebSocket EventSub (beta)
- Manages subscription lifecycle
