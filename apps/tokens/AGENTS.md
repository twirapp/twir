# AGENTS.md — apps/tokens

OAuth token management service.

## OVERVIEW

Manages OAuth tokens for Twitch and other integrations. Handles token refresh, validation, and secure storage. Provides token access to other services.

## STRUCTURE

```
apps/tokens/
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
| Main | `cmd/main.go` | Token service bootstrap |
| gRPC | —             | Token access API        |

## KEY COMMANDS

```bash
# Run locally
go run ./cmd/main.go

# Build
bun cli build tokens
```

## DEPENDENCIES

- Postgres (encrypted tokens)
- Twitch OAuth API

## SECURITY

- Tokens encrypted at rest
- Automatic refresh before expiry
- Scoped access control

## NOTES

- Centralized token management
- Other services request tokens via gRPC
- Handles multiple OAuth providers
