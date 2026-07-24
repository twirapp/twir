# AGENTS.md — libs/integrations

Third-party API integrations.

## OVERVIEW

Clients for external streaming/creator services. Handles API authentication, rate limiting, and error handling.

## STRUCTURE

```
libs/integrations/
├── seventv/                 # 7TV (genqlient-generated GraphQL client in api/generated.go)
├── spotify/                 # Spotify API
├── lastfm/                  # LastFM API
├── streamelements/          # StreamElements API
├── valorant/                # Valorant (HenrikDev) API
├── vk/                      # VK API
├── Makefile                 # genqlient generation
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
