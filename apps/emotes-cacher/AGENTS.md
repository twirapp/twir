# AGENTS.md — apps/emotes-cacher

Background service caching emote data from 3rd party providers.

## OVERVIEW

Caches emote metadata from BTTV, FFZ, and 7TV. Provides fast emote lookups for the parser and overlays.

## STRUCTURE

```
apps/emotes-cacher/
├── cmd/
│   └── main.go              # Entry point
├── internal/
│   └── ...
├── go.mod
└── Dockerfile
```

## ENTRY POINTS

| Type | Path          | Purpose           |
| ---- | ------------- | ----------------- |
| Main | `cmd/main.go` | Service bootstrap |

## KEY COMMANDS

```bash
# Run locally
go run ./cmd/main.go

# Build
bun cli build emotes-cacher
```

## DEPENDENCIES

- Postgres (cached emote data)
- External APIs: BTTV, FFZ, 7TV

## NOTES

- Periodic refresh of emote caches
- Provides gRPC API for emote lookups
