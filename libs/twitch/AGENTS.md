# AGENTS.md — libs/twitch

Twitch API and IRC clients.

## OVERVIEW

Twitch Helix API and IRC chat clients. Provides authenticated access to Twitch services for all bot functionality.

## STRUCTURE

```
libs/twitch/
├── *.go                     # API and IRC clients
├── go.mod
└── ...
```

## USAGE

```go
import "libs/twitch"

client := twitch.NewClient(token)
// Use Helix API
```

## DEPENDENCIES

- tokens service (OAuth tokens)

## NOTES

- Helix API wrapper
- IRC connection management
- Token refresh handling
