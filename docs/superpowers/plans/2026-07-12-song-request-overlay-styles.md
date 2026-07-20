# Song Request Overlay Styles Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add six persisted song-request overlay styles, a dashboard settings modal with preview, and realtime overlay updates.

**Architecture:** A new pgx-backed overlay-settings domain owns visual configuration and publishes saved settings through `WsRouter`. Shared Vue renderer components are used by both dashboard previews and the public overlay, while the public overlay alone owns the YouTube iframe lifecycle.

**Tech Stack:** Go, pgx v5, gqlgen, GraphQL subscriptions, WsRouter, PostgreSQL/Goose, Nuxt 3, Vue 3, TypeScript, Tailwind CSS, shadcn-vue, Bun.

---

## File Map

**Database and domain**

- Create via CLI: `libs/migrations/postgres/*_song_request_overlay_settings.sql`
- Create: `libs/entities/song_request_overlay_settings/entity.go`
- Create: `libs/repositories/song_request_overlay_settings/song_request_overlay_settings.go`
- Create: `libs/repositories/song_request_overlay_settings/errors.go`
- Create: `libs/repositories/song_request_overlay_settings/pgx/pgx.go`

**Backend service and GraphQL**

- Create: `apps/api-gql/internal/services/song_request_overlay_settings/service.go`
- Create: `apps/api-gql/internal/delivery/gql/mappers/song-request-overlay-settings.go`
- Create or regenerate: `apps/api-gql/internal/delivery/gql/resolvers/song-request-overlay-settings.resolver.go`
- Modify: `apps/api-gql/internal/delivery/gql/schema/song-requests.graphql`
- Modify: `apps/api-gql/internal/delivery/gql/resolvers/resolver.go`
- Modify: `apps/api-gql/cmd/main.go`
- Modify: `apps/api-gql/internal/delivery/gql/resolvers/song-requests.resolver.go`
- Modify: `libs/gomodels/channels_song_requests_settings.go`

**Shared renderer**

- Create: `web/app/components/song-request-overlay/types.ts`
- Create: `web/app/components/song-request-overlay/SongRequestOverlayRenderer.vue`
- Create: `web/app/components/song-request-overlay/CinemaStyle.vue`
- Create: `web/app/components/song-request-overlay/CompactStyle.vue`
- Create: `web/app/components/song-request-overlay/TickerStyle.vue`
- Create: `web/app/components/song-request-overlay/StudioStyle.vue`
- Create: `web/app/components/song-request-overlay/PortraitStyle.vue`
- Create: `web/app/components/song-request-overlay/PillStyle.vue`

**Dashboard and public overlay**

- Create: `web/layers/dashboard/api/song-request-overlay-settings.ts`
- Create: `web/layers/dashboard/components/songRequests/overlay-settings.vue`
- Modify: `web/layers/dashboard/pages/dashboard/song-requests.vue`
- Modify: `web/layers/dashboard/components/songRequests/settings.vue`
- Modify: `web/layers/dashboard/api/song-requests.ts`
- Modify: `web/layers/overlays/pages/o/[apiKey]/song-requests.client.vue`
- Modify: `web/layers/dashboard/locales/en.json`
- Modify: `web/layers/dashboard/locales/ru.json`

### Task 1: Create and Populate the Overlay Settings Table

- [ ] Run the required CLI command:

```bash
bun cli m create --name song_request_overlay_settings --db postgres --type sql
```

- [ ] Replace the generated template with a Goose Up migration that creates `song_request_overlay_style`, creates `channels_song_requests_overlay_settings`, validates colors and speed, copies `hide_on_pause`, and drops the old column.

```sql
CREATE TYPE song_request_overlay_style AS ENUM (
  'CINEMA', 'COMPACT', 'TICKER', 'STUDIO', 'PORTRAIT', 'PILL'
);

CREATE TABLE channels_song_requests_overlay_settings (
  id UUID PRIMARY KEY DEFAULT uuidv7(),
  channel_id TEXT NOT NULL UNIQUE REFERENCES channels(id) ON DELETE CASCADE,
  style song_request_overlay_style NOT NULL DEFAULT 'CINEMA',
  accent_color VARCHAR(9) NOT NULL DEFAULT '#8B5CF6',
  ticker_background_color VARCHAR(9) NOT NULL DEFAULT '#111827E6',
  ticker_text_color VARCHAR(9) NOT NULL DEFAULT '#FFFFFF',
  ticker_speed INT NOT NULL DEFAULT 35 CHECK (ticker_speed BETWEEN 10 AND 100),
  hide_on_pause BOOLEAN NOT NULL DEFAULT true,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  CHECK (accent_color ~ '^#[0-9A-Fa-f]{6}([0-9A-Fa-f]{2})?$'),
  CHECK (ticker_background_color ~ '^#[0-9A-Fa-f]{6}([0-9A-Fa-f]{2})?$'),
  CHECK (ticker_text_color ~ '^#[0-9A-Fa-f]{6}([0-9A-Fa-f]{2})?$')
);

INSERT INTO channels_song_requests_overlay_settings (channel_id, hide_on_pause)
SELECT channel_id, hide_on_pause FROM channels_song_requests_settings
ON CONFLICT (channel_id) DO NOTHING;

ALTER TABLE channels_song_requests_settings DROP COLUMN hide_on_pause;
```

- [ ] Implement Down in reverse order: restore the old column, copy values back, drop the settings table, then drop the enum.

### Task 2: Add Entity and pgx Repository

- [ ] Add `song_request_overlay_settings.Settings` with constants/default constructor for all six styles and persisted fields.
- [ ] Add repository `GetByChannelID(ctx, channelID)` and `Upsert(ctx, input)` contracts.
- [ ] Add `ErrNotFound` and ensure missing rows map to it.
- [ ] Implement pgx SQL using `$1` placeholders, `updated_at = now()`, and `RETURNING` all fields.
- [ ] Map repository rows to entities before returning them from the service.

### Task 3: Add Service, DI, Validation, and Realtime Publication

- [ ] Implement service defaults without inserting absent rows.
- [ ] Validate colors with `^#[0-9A-Fa-f]{6}([0-9A-Fa-f]{2})?$` and speed 10 through 100.
- [ ] Upsert through the repository and publish the resulting entity on `songrequests:overlay-settings:<channelId>`.
- [ ] Register the repository with `fx.Annotate(... fx.As(new(Repository)))` in `cmd/main.go`.
- [ ] Register the service constructor and add it to GraphQL resolver dependencies.

### Task 4: Add GraphQL Query, Mutation, and Subscription

- [ ] Add the enum, settings type, update input, authenticated query, permission-protected mutation, and public `@noRateLimit` subscription to `song-requests.graphql`.
- [ ] Map the entity to GraphQL DTO in a focused mapper.
- [ ] Query and mutation resolve the selected dashboard channel.
- [ ] Subscription resolves `apiKey` to channel, subscribes before loading initial settings, sends initial/default settings, then forwards updates.
- [ ] Remove `hideOnPause` from `SongRequestsSettings`, `SongRequestsSettingsOpts`, resolver mapping, and legacy gomodel.
- [ ] Run:

```bash
bun cli build gql
gofmt -w apps/api-gql libs/entities/song_request_overlay_settings libs/repositories/song_request_overlay_settings
```

### Task 5: Build Shared Style Components

- [ ] Define shared typed props and defaults in `types.ts`, including thumbnail URL derivation.
- [ ] Implement Cinema as full viewport with media slot and bottom status bar.
- [ ] Implement Compact as a bottom-left 480 px vertical card with 16:9 media.
- [ ] Implement Ticker with static text when it fits and ResizeObserver-driven marquee duration when it overflows.
- [ ] Implement Studio with square thumbnail, requester line, decorative playing waveform, times, and progress.
- [ ] Implement Portrait with thumbnail, title, requester, times, progress, and accent glow.
- [ ] Implement Pill with mini-thumbnail, `title • @requester`, and right-side progress.
- [ ] Implement renderer enum dispatch with unknown-style fallback to Cinema and media slot forwarding.

### Task 6: Add Dashboard Modal, Preview, and Localization

- [ ] Add dashboard query and update mutation composable with cache invalidation.
- [ ] Build modal using a native form, style selection cards, shared renderer preview, ColorPicker controls, numeric Ticker speed, hide-on-pause switch, Cancel, and Save.
- [ ] Keep preview local until Save; show success/error toast and keep modal open on error.
- [ ] Add settings icon beside the Overlay link controls and wire modal state.
- [ ] Remove hide-on-pause from the general settings modal/query/input.
- [ ] Add English and Russian `songRequests.links.*`, modal, style name/description, color, speed, preview, and feedback translations.
- [ ] Run frontend GraphQL generation:

```bash
bunx --bun graphql-codegen
```

### Task 7: Integrate Realtime Settings Into the Public Overlay

- [ ] Add the settings subscription by route `apiKey` and initialize renderer defaults before the first response.
- [ ] Keep one YouTube player instance regardless of style changes.
- [ ] Pass iframe media to Cinema/Compact; keep it full-sized and opacity-hidden as audio host for Ticker/Studio/Portrait/Pill.
- [ ] Feed title, requester, position, duration, video ID, playing state, and settings into the shared renderer.
- [ ] Apply hide-on-pause through opacity/pointer-events rather than `display:none`.
- [ ] Ensure settings updates alter renderer props without seeking, reloading, or recreating the iframe.

### Task 8: Verify End-to-End Behavior

- [ ] Run backend generation/build:

```bash
bun cli build gql
rtk go build ./...
```

- [ ] Run frontend generation/typecheck and document unrelated pre-existing failures:

```bash
bunx --bun graphql-codegen
bunx --bun nuxt typecheck
```

- [ ] In Playwright, verify all six styles against the public overlay URL.
- [ ] Open dashboard and overlay simultaneously, Save each style, and verify immediate subscription-driven updates.
- [ ] Verify accent updates all styles and Ticker colors/speed update in realtime.
- [ ] Verify long Ticker text scrolls slowly and short text remains static.
- [ ] Verify pause/resume with hide enabled and disabled without iframe detachment or position reset.
- [ ] Verify reconnect returns persisted settings and paused background remains transparent.

## Execution Constraints

- Do not create commits.
- Do not add automated test files.
- Do not modify or revert unrelated worktree changes.
- Use `apply_patch` for manual edits and the project CLI for migration creation and GraphQL generation.
