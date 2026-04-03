# AGENTS.md — apps/bots

Twitch bot service handling chat interactions and commands.

## OVERVIEW

Manages multiple Twitch bot instances across channels. Handles chat messages, commands, moderation actions, and translations. Connects to Twitch IRC and processes incoming messages.

## STRUCTURE

```
apps/bots/
├── cmd/
│   └── main.go              # Entry point
├── internal/
│   ├── messagehandler/      # IRC message processing
│   ├── services/            # Bot services
│   │   ├── chat_translations/
│   │   └── ...
│   └── ...
├── go.mod
└── Dockerfile
```

## ENTRY POINTS

| Type | Path          | Purpose               |
| ---- | ------------- | --------------------- |
| Main | `cmd/main.go` | Bot service bootstrap |

## KEY COMMANDS

```bash
# Run locally
go run ./cmd/main.go

# Build
bun cli build bots
```

## DEPENDENCIES

- Twitch IRC (via libs/twitch)
- Postgres (bot settings, channels)
- Redis (state management)
- api-gql (for commands data)

## ANTI-PATTERNS

- **NEVER** use `context.TODO()` — propagate context from message handlers
- **ALWAYS** handle IRC connection drops gracefully

## PORTS

| Service | Port    | Protocol |
| ------- | ------- | -------- |
| gRPC    | Dynamic | Internal |

## NOTES

- Manages OAuth tokens via tokens service
- Supports multiple bot identities
- Chat translations use external APIs
