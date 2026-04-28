# Kick Stream Persistence And Parser Vars Design

**Goal:** Add Kick stream persistence using official webhook events first, extend existing stream storage to support multi-platform nullable fields, and unlock Kick parser variable parity for live stream fields and history-based category time.

## Current State

- Parser already has limited Kick-backed vars through `apps/parser/internal/variables/shared/kick.go`.
- Current Kick-backed vars use official REST reads for:
  - `stream.title`
  - `stream.category`
  - `subscribers.count`
- Kick parser coverage is missing parity for:
  - `stream.viewers`
  - `stream.uptime`
  - `stream.category.time`
  - richer stream/channel metadata vars
  - gifted kicks leaderboard vars
- `apps/eventsub/internal/kick/subscription_manager.go` already subscribes to official `livestream.status.updated`.
- `gokick` already knows official `livestream.metadata.updated`, but current Kick subscription manager does not subscribe to it yet.
- Twitch stream state currently uses `channels_streams`, `streams` repository, and `apps/scheduler/internal/timers/streams.go` polling.

## Constraints

- Official Kick API/webhooks are primary source of truth.
- Unofficial/private endpoints may be used only when official API cannot provide a needed field.
- Existing `channels_streams` path should be reused instead of creating a separate Kick-only stream table.
- Cross-platform stream fields must become nullable where not all platforms can provide them consistently.

## Decision

### Source Of Truth

- Primary event source for Kick stream lifecycle:
  - `livestream.status.updated`
  - `livestream.metadata.updated`
- Primary read fallback for current state:
  - official Kick REST channel/livestream endpoints already exposed through `gokick`
- Recovery/reconcile:
  - use scheduler polling only to repair missed webhook state or bootstrap current state after downtime
- Unofficial/private endpoints:
  - allowed only for fields still unavailable from official API after implementation

### Persistence Strategy

- Reuse existing `channels_streams` model and `streams` repository.
- Expand the persistence contract from Twitch-only assumptions to multi-platform support.
- Add platform-aware semantics to the stream row so Kick current stream data can live in the same storage path.
- Make non-universal fields nullable in DB and repository model to support partial platform coverage.

### Nullable Field Policy

Fields that can be missing for some platform/event source combinations should be nullable in persistence and repository update flow. Expected nullable candidates:

- `gameId`
- `gameName`
- `communityIds`
- `type`
- `viewerCount`
- `startedAt`
- `language`
- `thumbnailUrl`
- `tagIds`
- `tags`
- `isMature`

Kick-specific current metadata should populate these fields when available. Missing fields should remain null instead of forcing synthetic defaults.

### Kick Event Handling

#### `livestream.status.updated`

- On `is_live=true`:
  - resolve internal channel/user IDs
  - upsert current stream row in `channels_streams`
  - populate at minimum:
    - platform/channel identity
    - title
    - startedAt
    - viewerCount if present from fallback fetch/reconcile
  - publish existing Kick online bus event
- On `is_live=false`:
  - remove or close current stream row using same conceptual contract as Twitch current-stream storage
  - publish existing Kick offline bus event

#### `livestream.metadata.updated`

- Subscribe and handle officially.
- Update current persisted stream row with latest:
  - title
  - category id/name
  - language
  - mature flag
  - thumbnail if available
  - tags/custom tags if available from read fallback

### Scheduler Responsibility

`apps/scheduler/internal/timers/streams.go` remains reconciliation path, not primary state engine for Kick.

Scheduler additions:

- identify Kick channels needing reconciliation
- fetch current official Kick live state
- repair missing online rows if webhook was missed
- repair stale online rows if webhook missed offline transition
- optionally enrich current row with latest viewer count / thumbnail / tags when safe

Polling is not used as the normal source of timestamps when webhook events are available.

## Parser Variable Scope

### Direct Current-State Vars

Implement or extend Kick support for:

- `stream.viewers`
- `stream.uptime`
- `stream.language`
- `stream.tags`
- `stream.slug`
- `stream.description`
- `stream.thumbnail`

Preferred resolution order:

1. current persisted `channels_streams` row when field exists there
2. official Kick live/channel REST fallback
3. unofficial/private fallback only if official source still lacks the field

### History-Based Vars

- `stream.category.time` should be computed from persisted stream/history data, not from current live endpoints.
- Kick must share the same conceptual history path as Twitch, rather than adding parser-only custom logic.

### Gift/Kicks Vars

Add official `kicks` leaderboard-backed vars using `public/v1/kicks/leaderboard` when token scope allows:

- lifetime gifted kicks leaderboard
- monthly gifted kicks leaderboard
- weekly gifted kicks leaderboard

These do not require unofficial APIs because official support exists.

## Unofficial API Use Policy

Known unofficial/private endpoints can help with:

- richer live metadata
- extra channel presentation fields
- some leaderboard/subscriber fields

But they should only be adopted after official coverage is exhausted. In this design, unofficial/private use is acceptable for:

- fields like `description`, `slug`, `banner`, or tag variants if official read path proves insufficient in practice
- not as primary source for stream online/offline lifecycle

## Files Likely To Change

### EventSub

- `apps/eventsub/internal/kick/subscription_manager.go`
- `apps/eventsub/internal/kick/handlers.go`

### Scheduler

- `apps/scheduler/internal/timers/streams.go`
- possibly scheduler support code for Kick reconciliation

### Repositories / Models / Migrations

- `libs/repositories/streams/model/stream.go`
- `libs/repositories/streams/streams.go`
- `libs/repositories/streams/datasource/postgres/postgres.go`
- migration for `channels_streams` nullable/multi-platform adjustments
- legacy `libs/gomodels/channels_streams.go` only if still required by current scheduler/event paths

### Parser Vars

- `apps/parser/internal/variables/shared/kick.go`
- `apps/parser/internal/variables/stream/*.go`
- new metadata variable files if needed
- `apps/parser/internal/variables/top` only if kicks leaderboard formatting belongs there

## Error Handling

- Webhook handlers must remain idempotent.
- Metadata events arriving before current live row exists should not crash; they may trigger best-effort enrichment or be ignored until reconcile creates the row.
- Scheduler reconcile should be safe to rerun repeatedly.
- Missing fields from Kick APIs should yield empty variable results, not hard errors, unless the whole source fetch failed unexpectedly.

## Testing

- EventSub tests:
  - `livestream.metadata.updated` subscription and handler
  - online/offline persistence updates
  - idempotent repeated delivery
- Repository tests:
  - nullable field read/write behavior
  - multi-platform update semantics
- Parser tests:
  - Kick `stream.viewers`
  - Kick `stream.uptime`
  - Kick metadata vars
  - Kick `stream.category.time` using persisted data
- Scheduler tests:
  - reconcile creates missing Kick current stream row
  - reconcile removes stale Kick current stream row

## Risks

- Existing `channels_streams` code assumes Twitch-only non-null values.
- Nullability migration may impact current Gorm and pgx reads if not updated consistently.
- Kick metadata webhook payload may not include every field needed for parser vars, so reconcile/fallback enrichment is necessary.
- `stream.category.time` depends on consistent persisted lifecycle boundaries; missed webhook events must be repairable.

## Recommended Implementation Order

1. Add Kick metadata event subscription and handler wiring.
2. Migrate `channels_streams` and repository model to nullable/multi-platform-safe semantics.
3. Persist Kick current stream state through EventSub.
4. Add scheduler reconciliation for Kick.
5. Add parser vars backed by persisted/current Kick stream data.
6. Add official kicks leaderboard vars.
7. Add unofficial/private fallbacks only for fields still blocked after official path lands.
