# AGENTS.md — libs/integrations

Third-party API integrations.

## OVERVIEW

Clients for external services (7TV, Spotify, Discord, etc.). Handles API authentication, rate limiting, and error handling.

## STRUCTURE

```
libs/integrations/
├── seventv/                 # 7TV integration
│   └── api/
│       └── generated.go    # Generated GraphQL client
├── spotify/                 # Spotify API
├── discord/                 # Discord API
├── Makefile                # genqlient generation
├── go.mod
└── ...
```

## KEY COMMANDS

```bash
# Regenerate 7TV GraphQL client
make generate
# or
cd seventv && go generate ./...
```

## CONVENTIONS

- Generated clients in `api/generated.go`
- Use genqlient for GraphQL APIs
- Handle rate limits gracefully

## NOTES

- Some clients auto-generated
- OAuth tokens from tokens service
- Rate limit handling built-in
