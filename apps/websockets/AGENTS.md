# AGENTS.md — apps/websockets

Real-time WebSocket server for overlays.

## OVERVIEW

Manages WebSocket connections for browser overlays. Handles overlay events, dudes (avatars), and real-time data streaming from backend to browser.

## STRUCTURE

```
apps/websockets/
├── cmd/
│   └── main.go              # Entry point
├── internal/
│   ├── namespaces/          # Socket.io namespaces
│   │   └── overlays/        # Overlay handlers
│   │       └── registry/    # Overlay registry
│   └── ...
├── go.mod
└── Dockerfile
```

## ENTRY POINTS

| Type      | Path          | Purpose                    |
| --------- | ------------- | -------------------------- |
| Main      | `cmd/main.go` | WebSocket server bootstrap |
| Socket.io | `/`           | Client connections         |

## KEY COMMANDS

```bash
# Run locally
go run ./cmd/main.go

# Build
bun cli build websockets
```

## DEPENDENCIES

- Redis (pub/sub for events)
- events (event streaming)
- overlays frontend (browser clients)

## NOTES

- Socket.io for browser compatibility
- Namespace per overlay type
- Handles thousands of concurrent connections
