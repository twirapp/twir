# Song Request Modes (YouTube / Spotify)

Issue: https://github.com/twirapp/twir/issues/1011

Channels can choose how song requests are played via `channels_song_requests_settings.mode`:

- `YOUTUBE` (default) — the existing flow: YouTube search + dashboard player widget. Unchanged.
- `SPOTIFY` — requests are added to the streamer's **own Spotify Connect queue** via the Spotify Web API.

## Architecture

| Piece | Location |
| --- | --- |
| Mode enum | `libs/entities/song_request_mode` |
| Settings entity (`mode` field) | `libs/entities/song_requests_settings` |
| Spotify request entity/statuses | `libs/entities/spotify_song_request` |
| Persistence | `libs/repositories/spotify_song_requests` (table `spotify_song_requests`) |
| Spotify Web API client | `libs/integrations/spotify` (search, devices, queue, skip, currently-playing) |
| Domain service | `apps/api-gql/internal/services/spotify_song_requests` |
| Playback reconciler (5s tick / 60s idle) | same package, `reconciler.go` |
| Bus endpoints | `spotify.search`, `spotify.songRequest.create`, `spotify.songRequest.cancel` (`libs/bus-core/spotify`) |
| Chat commands (`!sr`, `sr wrong`, `voteskip`) | `apps/parser/internal/commands/songrequest/youtube` (mode dispatch) |
| Donation requests | `apps/events/internal/song_request` (mode dispatch) |
| GraphQL | `apps/api-gql/internal/delivery/gql/schema/song-requests.graphql` |
| Dashboard UI | `web/layers/dashboard/components/songRequests` (`settings.vue`, `spotify-queue.vue`) |

Data layer rule: everything new goes through `libs/entities` + `libs/repositories` (pgx); legacy
`libs/gomodels` is only read where existing YouTube code still requires it.

## Spotify mode behavior

- `!sr <query|link>` searches Spotify, adds the first match to the streamer's queue, and records a
  `spotify_song_requests` row (`queued`). Limits (`maxRequests`, `userMaxRequests`, min/max length)
  are enforced server-side.
- `sr wrong` marks the requester's latest request `cancelled_pending_skip`. Spotify does **not**
  allow removing arbitrary queue items, so the reconciler performs a **deferred skip**: when the
  cancelled track becomes the currently playing one, it calls `skip next` and marks the request
  `skipped_by_twir`.
- `voteskip` is not supported in Spotify mode (the command answers with a notice).
- The reconciler keeps request statuses in sync with the real player: `queued` → `playing` →
  `played`, `removed_or_reconciled` when a track disappears from the queue for >15s, `unknown` when
  Spotify is disconnected or scopes are missing.
- Device selection: the first active, non-restricted device is cached in KV for 5 minutes
  (`spotify:songrequests:device:<channelID>`). Dashboard "refresh device" re-resolves it.

## Requirements and limitations

- The streamer must connect the Spotify integration with the **new scope**
  `user-modify-playback-state` (plus existing `user-read-playback-state`). Tokens without the new
  scope surface as `hasPlaybackScope: false` in `spotifyCapabilities`; the dashboard shows a
  reconnect prompt.
- **Spotify Premium is required** for queue/playback control (Spotify API restriction).
- An **active device** must exist (open Spotify on any device) before requests can be queued.
- The dashboard player widget and YouTube-only settings (`playerNoCookieMode`, YouTube deny lists,
  min views, categories) do not apply in Spotify mode.

## GraphQL surface

- `songRequests.mode` + `songRequests.spotifyCapabilities { connected hasPlaybackScope canUseSpotify activeDevice selectedDevice }`
- `spotifySongRequestsQueue { requests currentDevice }`
- `spotifySongRequestsSearch(query, limit)` (limit capped at 20, default 5)
- Mutations: `spotifySongRequestSelectDevice`, `spotifySongRequestRefreshDevice`,
  `spotifySongRequestSkip`, `spotifySongRequestCancel`

## Migrations

- `20260805161157_song_request_mode.sql` — adds `mode` (`varchar(20)`, default `YOUTUBE`) with backfill.
- `20260805162057_spotify_song_requests.sql` — new `spotify_song_requests` table.

Both are additive and rollback-safe.

## Rollout checklist

1. Deploy migrations (additive; default `YOUTUBE` keeps existing behavior for all channels).
2. Deploy api-gql, parser, events together — the new bus endpoints must be served by api-gql before
   parser/events call them in Spotify mode (YouTube mode does not use the new endpoints, so a mixed
   deploy is safe for existing channels).
3. Verify the Spotify OAuth app exposes `user-modify-playback-state` in its scope list; no client
   config change is needed, scope is requested at (re)connect time.
4. Tell streamers who want the mode: connect/reconnect Spotify in dashboard integrations, open
   Spotify on a device, then switch the mode in Song Requests settings.
5. Watch api-gql logs for `spotify song request` / reconciler errors (rate limits, missing devices).
6. Follow-up: expose runtime metrics (`spotify_request_created`, `spotify_request_status_changed`,
   `spotify_device_missing`, `spotify_rate_limited`, `spotify_deferred_skip_executed`) once api-gql
   has a `/metrics` endpoint — it currently does not serve Prometheus.
