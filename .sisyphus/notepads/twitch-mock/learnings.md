## [2026-04-12] Session: twitch-mock event duplication fix

### Architecture

- Mock OAuth2: `localhost:7777/oauth2/*`
- Mock Helix API: `localhost:7777/helix/*`
- Mock Admin UI: `localhost:3333/admin`
- Mock EventSub WS: `localhost:8081/ws`
- Mock broadcaster: ID=`12345`, login=`mockstreamer`
- Mock bot: ID=`67890`, login=`mockbot`

### Config fields

- TwitchMockApiUrl = "http://localhost:7777/helix" // needs /helix for helix Go lib
- TwitchMockAuthUrl = "http://localhost:7777" // without /helix
- TwitchMockWsUrl = "ws://localhost:8081/ws"

### Key files

- `apps/twitch-mock/internal/websocket/eventsub.go` - WS server, Broadcast(), clients map
- `apps/twitch-mock/internal/state/state.go` - conduit/subscription state
- `apps/twitch-mock/internal/admin/admin.go` - trigger handler
- `apps/twitch-mock/internal/handlers/conduits.go` - conduit CRUD endpoints
- `apps/eventsub/internal/manager/manager.go` - manager, uses TrimSuffix for apiBaseUrl
- `apps/eventsub/internal/manager/on_start.go` - createConduit (ShardCount:3 hardcoded!), twitchUpdateConduitShard
- `apps/eventsub/internal/manager/websocket.go` - startWebSocket, reconnect loop at lines 95-105

### Fixed bugs

1. `scope` in `/oauth2/token` changed from `""` (string) to `[]string{}` - auth.go
2. Double `/helix` in eventsub manager fixed with `strings.TrimSuffix(opts.Config.TwitchMockApiUrl, "/helix")` - manager.go

### Active bug: event duplication

Root cause: `Broadcast()` in websocket/eventsub.go sends to ALL clients in `s.clients` map.
The eventsub service creates a conduit with ShardCount:3 (hardcoded for prod) but only 1 shard
gets a WS session. If reconnect loop creates a second WS connection, old client is still in map.
Two clients → two deliveries.

TWO FIX OPTIONS:
A) In the mock WS server: track which sessionIDs are registered as conduit shards (via UpdateConduitShards),
and in Broadcast() only send to those sessions, not all clients.
B) In eventsub app: when TwitchMockEnabled=true, use ShardCount:1 instead of 3.

RECOMMENDED: Do BOTH for correctness.

- Option B prevents creating unused shard slots
- Option A ensures even with multiple WS clients, only the ones tied to active shards get events

### Constraints

- DO NOT mutate helix.AuthBaseURL package global - use RoundTripper
- DO NOT add twitch-mock to docker-compose.stack.yml (production)
- DO NOT use frontend framework for Admin UI - only embedded plain HTML in Go binary
- DO NOT hardcode ports in application code - only from config/env
- DO NOT touch docker-compose.stack.yml
- DO NOT change behavior when TWITCH_MOCK_ENABLED=false

## [2026-04-12] Session: admin UI event ordering

### Change

- Moved `channel.chat.message` to the top of the admin event grid in `apps/twitch-mock/internal/admin/templates/index.html`.

### Note

- Kept all other event sections in their original order and left handler logic untouched.
