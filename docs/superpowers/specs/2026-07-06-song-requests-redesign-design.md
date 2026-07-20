# Song Requests Redesign

**Date:** 2026-07-06
**Status:** Approved
**NOTE:** Do NOT commit any files during implementation. All changes are uncommitted.

---

## Overview

Complete rewrite of the YouTube song request system. The current implementation uses raw WebSocket (melody) in `apps/websockets` for real-time communication. The new system uses GQL subscriptions (backed by NATS wsRouter) for client events, GQL mutations for actions, and bus-core (NATS) for inter-service communication.

**Reference:** https://trula-music.ru/ (media donation service for streamers)

---

## Decisions Made

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Sync model | B: Master position, server stores state | Reconnect-safe, no client-side position computation |
| Communication | GQL WS (subscriptions + mutations) | Unified protocol, no raw WebSocket |
| Channel API key | A: Both work (channel key + user key fallback) | Backward compatibility |
| Melody YouTube namespace | A: Delete fully, migrate to GQL | Clean architecture, no dual support |
| OBS overlay | Visible YouTube player + title + author + progress bar | Video clip must be shown |
| Popup widget | Controls + queue, ~400x600, for moderators too | Compact external widget |
| Moderator auth | Trust API key | API key = access, streamer controls distribution |
| Widget layer | New `web/layers/widgets/` Nuxt layer | Separate from dashboard and overlays |
| Inter-service comms | bus-core (not wsRouter) | wsRouter only for pubsub to clients |

---

## Architecture

### Data Flow

```
Parser/EventSub/Events
        |
        v
   bus-core (NATS)
        |
        v
    api-gql (bridge: bus-core -> wsRouter)
        |
        v
   wsRouter (NATS)
        |
        v
   GQL Subscriptions -> Clients
```

### Sync Model

Server stores "master state" in Redis per channel:

```
songrequests:playback:{channelId} -> {
  videoId: string,
  title: string,
  position: number,      // server-computed, ready for clients
  isPlaying: boolean,
  volume: number,         // 0-100
  updatedAt: timestamp
}
```

**Client behavior:** On subscription event, client calls `player.seekTo(position)`, `player.setVolume(volume)`, then `playVideo()` or `pauseVideo()` based on `isPlaying`. Client does NOT compute position — server sends it ready.

**Playback ticker:** Goroutine in api-gql updates `position` in Redis every ~1s for all playing channels and publishes to wsRouter. This keeps clients in sync without them computing position.

**Position computation (server-side):**
- On play: `startedAt = now, position = 0`
- On tick: `position = (now - startedAt)`
- On pause: `position = position, isPlaying = false`
- On resume: `startedAt = now - position*1000`

---

## GraphQL Schema Changes

### New Subscriptions

```graphql
subscription SongRequestPlaybackState($channelId: UUID!) {
  songRequestPlaybackState(channelId: $channelId) {
    videoId
    title
    position
    isPlaying
    volume
    updatedAt
  }
}

subscription SongRequestQueueUpdated($channelId: UUID!) {
  songRequestQueueUpdated(channelId: $channelId) {
    id
    title
    songLink
    durationSeconds
    orderedByName
    orderedByDisplayName
    queuePosition
    createdAt
  }
}
```

### New Mutations

```graphql
mutation SongRequestPlay($channelId: UUID!, $videoId: String!) { ... }
mutation SongRequestPause($channelId: UUID!) { ... }
mutation SongRequestSkip($channelId: UUID!) { ... }
mutation SongRequestSetVolume($channelId: UUID!, $volume: Int!) { ... }
mutation SongRequestReorder($channelId: UUID!, $videoIds: [String!]!) { ... }
mutation SongRequestDeleteFromQueue($channelId: UUID!, $videoId: String!) { ... }
mutation SongRequestClearQueue($channelId: UUID!) { ... }
```

### New Queries

```graphql
query SongRequestWidgetData($channelId: UUID!) {
  songRequestPlaybackState(channelId: $channelId) { ... }
  songRequestQueue(channelId: $channelId) { ... }
}

query ChannelByApiKey($apiKey: String!) {
  channelByApiKey(apiKey: $apiKey) {
    id
    twitchUserId
    kickUserId
  }
}
```

### Updated Settings

Add `hideOnPause: Boolean!` to existing `ChannelSongRequestsSettings` type (default: `true`).

---

## Database Changes

### Migration 1: Channel API Key

```sql
ALTER TABLE channels ADD COLUMN api_key TEXT DEFAULT uuidv7();
CREATE UNIQUE INDEX channels_api_key_idx ON channels(api_key) WHERE api_key IS NOT NULL;
```

Existing channels get UUID7 key automatically via default.

### Migration 2: Hide On Pause Setting

```sql
ALTER TABLE channels_song_requests_settings ADD COLUMN hide_on_pause BOOL DEFAULT true;
```

---

## Bus-Core Changes

### New message types (`libs/bus-core/api/song_requests.go`)

```go
const (
    SongRequestAddToQueueSubject      = "api.songRequests.addToQueue"
    SongRequestRemoveFromQueueSubject = "api.songRequests.removeFromQueue"
    SongRequestPlaybackStateSubject   = "api.songRequests.playbackState"
)

type SongRequestAddToQueue struct {
    ChannelID   string
    SongRequest SongRequestData
}

type SongRequestRemoveFromQueue struct {
    ChannelID string
    VideoID   string
}

type SongRequestPlaybackState struct {
    ChannelID string
    VideoID   string
    Title     string
    Position  float64
    IsPlaying bool
    Volume    int
    UpdatedAt int64
}

type SongRequestData struct {
    ID                   string
    Title                string
    VideoID              string
    SongLink             string
    DurationSeconds      int
    OrderedByName        string
    OrderedByDisplayName string
    QueuePosition        int
    CreatedAt            string
}
```

### Registration

Add to `apiBus` struct in `bus-services.go`:
```go
SongRequestAddToQueue      Queue[api.SongRequestAddToQueue, struct{}]
SongRequestRemoveFromQueue Queue[api.SongRequestRemoveFromQueue, struct{}]
SongRequestPlaybackState   Queue[api.SongRequestPlaybackState, struct{}]
```

Instantiate in `NewNatsBus` in `bus.go`.

---

## Backend Changes

### Auth Middleware — Backward Compatibility

Current flow: `api-key header` -> `users.apiKey` -> user -> channel.

New flow:
1. Check `channels.api_key` -> if found -> channel
2. Fallback: `users.apiKey` -> user -> channel

Files:
- `apps/api-gql/internal/auth/api_key.go`
- `apps/api-gql/internal/delivery/http/middlewares/is_authenticated.go`
- `apps/api-gql/internal/delivery/gql/directives/is_authenticated.go`

### Playback State Service (new)

`apps/api-gql/internal/services/song_requests/playback_state.go`:
- Stores master state in Redis (`songrequests:playback:{channelId}`)
- `SetPlaying(channelId, videoId, title, position)`
- `SetPaused(channelId)` — saves position, sets `isPlaying: false`
- `SetVolume(channelId, volume)`
- `GetState(channelId)` — reads from Redis
- `ClearState(channelId)` — deletes on skip/stop

### Bridge Service (new)

`apps/api-gql/internal/services/song_requests/bridge.go`:
- In `fx.Hook.OnStart`: subscribes to bus-core topics via `SubscribeGroup("api", callback)`
- Callback: reads data, publishes to wsRouter with channel-scoped key
- `OnStop`: unsubscribus

### Playback Ticker (new)

Goroutine in api-gql (or within playback_state service):
- Every ~1s: updates `position` in Redis for all playing channels
- Publishes updated state to wsRouter -> clients receive fresh position

### GQL Resolvers

New resolvers in `apps/api-gql/internal/delivery/gql/resolvers/song-requests.resolver.go`:
- Subscriptions: subscribe to wsRouter keys, unmarshal, map to GQL models
- Mutations: update Redis state, publish to wsRouter, return result
- Queries: read from Redis + DB

### Parser Changes

`apps/parser/internal/commands/songrequest/youtube/sr.go`:
- Replace gRPC melody call with `twirBus.Api.SongRequestAddToQueue.Publish(...)`
- Include full `SongRequestData` in message

`apps/parser/internal/commands/songrequest/youtube/skip.go`:
- Replace gRPC melody call with `twirBus.Api.SongRequestRemoveFromQueue.Publish(...)`

`apps/parser/internal/commands/songrequest/youtube/wrong.go`:
- Same pattern as skip

### EventSub Changes

`apps/eventsub/internal/handler/redemption.go`:
- Replace gRPC melody call with bus-core publish

### Events Changes

`apps/events/internal/song_request/song_request.go`:
- Replace gRPC melody call with bus-core publish

---

## Frontend Changes

### Dashboard (`web/layers/dashboard/`)

**New composable:** `web/layers/dashboard/composables/useSongRequestGql.ts`
- Wraps GQL subscriptions (playback state + queue) and mutations (play, pause, skip, volume, reorder, delete, clear)
- Returns reactive state + action functions

**Refactored components:**
- `songRequests/player.vue` — uses `useSongRequestGql` instead of `useYoutubeSocket`. YouTube IFrame API applies state from subscription.
- `songRequests/queue.vue` — data from subscription instead of WebSocket events.
- `songRequests/settings.vue` — stays as is (already uses GQL).
- `layout/sidebar/sidebar-mini-player.vue` — connects to `useSongRequestGql`

**`useGlobalYoutubePlayer.ts`** — changes:
- Was: manages queue, decides when to switch tracks
- Now: receives state from subscription and applies. Does not decide — server is authority.

**Delete:** `songRequests/hook.ts` (old raw WebSocket composable)

### OBS Overlay (`web/layers/overlays/`)

**New page:** `web/layers/overlays/pages/o/[apiKey]/song-requests.vue`

Layout:
```
┌─────────────────────────────────────┐
│ ┌─────────────────────────────────┐ │
│ │                                 │ │
│ │      YouTube Video Player       │ │
│ │         (16:9 aspect)           │ │
│ │                                 │ │
│ └─────────────────────────────────┘ │
│  Track Title                        │
│  by Author                          │
│  ████████░░░░░░░░ 1:23 / 3:45      │
└─────────────────────────────────────┘
```

Logic:
1. Get `apiKey` from route params
2. `channelByApiKey` query -> `channelId`
3. Subscribe to `songRequestPlaybackState(channelId)`
4. Apply state: seekTo, play/pause, update UI
5. If `isPlaying` -> show overlay with title, progress bar
6. If `!isPlaying` -> check `hideOnPause` setting. If true -> hide. Otherwise -> show paused state
7. If state = null (no current track) -> hide completely

YouTube IFrame API: visible player (16:9), shows video clip. Progress bar from subscription position (not computed client-side).

### Widget Layer (`web/layers/widgets/`) — NEW

**New Nuxt layer:** `web/layers/widgets/`

**Route:** `web/layers/widgets/pages/w/[channelApiKey]/song-requests.vue`

Layout (compact, ~400x600):
```
┌────────────────────────────────────────┐
│ ┌────────────────────────────────────┐ │
│ │                                    │ │
│ │     YouTube Video Player (16:9)    │ │
│ │                                    │ │
│ └────────────────────────────────────┘ │
│  Track Title — Author                  │
│  ████████░░░░░░ 1:23 / 3:45           │
│  [⏮] [⏯] [⏭]  🔊━━━━░░ 70%          │
├────────────────────────────────────────┤
│ Queue (3)                              │
│ 1. Song Name — Author    3:45 [✕]     │
│ 2. Song Name — Author    2:10 [✕]     │
│ 3. Song Name — Author    4:20 [✕]     │
├────────────────────────────────────────┤
│ [Clear All]                            │
└────────────────────────────────────────┘
```

Logic:
1. Get `channelApiKey` from route params
2. `channelByApiKey` query -> `channelId`
3. `SongRequestWidgetData(channelId)` query — initial load (state + queue)
4. Subscriptions: `songRequestPlaybackState` + `songRequestQueueUpdated`
5. YouTube IFrame API — visible player (16:9), shows video clip
6. Mutations: play, pause, skip, volume, reorder, delete, clear

Differences from dashboard:
- Compact layout (no sidebar, no header)
- No settings (controls only)
- Auth via channel API key (not user session)
- Accessible to moderators (anyone with the key)

**Layer config:** `web/layers/widgets/nuxt.config.ts` — minimal, route prefix `/w/`

---

## Files to Delete

- `apps/websockets/internal/namespaces/youtube/` (entire directory)
- `apps/websockets/internal/grpc_impl/youtube.go`
- `web/layers/dashboard/components/songRequests/hook.ts` (old raw WebSocket composable)

---

## Implementation Order

1. Channel API key migration + auth middleware backward compat
2. Bus-core message types + registration
3. Playback state service (Redis)
4. GQL schema + resolvers (subscriptions, mutations, queries)
5. Bridge service (bus-core -> wsRouter)
6. Playback ticker (server-side position updates)
7. Parser/eventsub/events -> bus-core publish (replace gRPC)
8. Dashboard refactoring (GQL subscriptions instead of raw WS)
9. OBS overlay (`web/layers/overlays/`)
10. Widget layer (`web/layers/widgets/`)
11. Delete melody YouTube namespace
12. i18n keys

---

## Key File Paths

| Area | Path |
|------|------|
| Migrations | `libs/migrations/postgres/` |
| Auth | `apps/api-gql/internal/auth/api_key.go` |
| Bus-core messages | `libs/bus-core/api/song_requests.go` |
| Bus-core registration | `libs/bus-core/bus-services.go`, `libs/bus-core/bus.go` |
| Playback service | `apps/api-gql/internal/services/song_requests/playback_state.go` |
| Bridge service | `apps/api-gql/internal/services/song_requests/bridge.go` |
| GQL schema | `apps/api-gql/internal/delivery/gql/schema/song-requests.graphql` |
| GQL resolvers | `apps/api-gql/internal/delivery/gql/resolvers/song-requests.resolver.go` |
| Parser | `apps/parser/internal/commands/songrequest/youtube/sr.go` |
| EventSub | `apps/eventsub/internal/handler/redemption.go` |
| Events | `apps/events/internal/song_request/song_request.go` |
| Dashboard composables | `web/layers/dashboard/composables/useSongRequestGql.ts` |
| Dashboard components | `web/layers/dashboard/components/songRequests/` |
| OBS overlay | `web/layers/overlays/pages/o/[apiKey]/song-requests.vue` |
| Widget layer | `web/layers/widgets/` |
| Delete | `apps/websockets/internal/namespaces/youtube/` |
