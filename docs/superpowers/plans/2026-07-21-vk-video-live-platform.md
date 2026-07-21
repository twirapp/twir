# VK Video Live Platform Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Replace channel-level Twitch/Kick coupling with normalized platform bindings, preserve both platforms, and add VK Video Live behind a feature flag using webhook ingestion and capability-driven adapters.

**Architecture:** `channels` remains the Twir aggregate; a new `channel_platforms` relation owns every provider binding. A registry supplies provider adapters and immutable capabilities. Webhooks normalize into the existing generic NATS contracts; outgoing actions resolve enabled bindings and delegate to adapters.

**Tech Stack:** Go, pgx v5, PostgreSQL 18/goose, NATS, gqlgen, Nuxt 3/Vue 3, Bun.

---

## File Map

- `libs/entities/platform/platform.go`: add `vk_video_live`, capabilities, and typed unsupported-capability errors.
- `libs/repositories/channel_platforms/{channel_platforms.go,model/model.go,pgx/pgx.go}`: normalized binding repository.
- `libs/migrations/postgres/20260721120000_channel_platforms.sql`: schema, Twitch/Kick backfill, constraints and indexes.
- `libs/repositories/channels/{channels.go,model/model.go,pgx/pgx.go}` and `libs/services/channels/channels.go`: remove platform columns and expose binding-based lookup.
- `apps/api-gql/internal/platform/{registry.go,vkvideo/provider.go}`: platform adapter registry and VK OAuth/API client.
- `apps/eventsub/internal/{platforms, vkvideo}`: transport contract, VK webhook verification and event normalizers.
- `apps/bots/internal/platforms` and `apps/bots/internal/services/channel/message.go`: capability-aware outgoing chat and moderation routing.
- `apps/api-gql/internal/delivery/{http/routes/auth,gql}`: binding OAuth, GraphQL types/mappers/resolvers.
- `web/layers/dashboard/features/channel-platforms`: binding cards and capability-gated controls.

### Task 1: Establish platform primitives and tests

**Files:**
- Modify: `libs/entities/platform/platform.go`
- Create: `libs/entities/platform/platform_test.go`

- [ ] Write tests for `PlatformVKVideoLive`, `Platform.All()`, validity, and `Capabilities.Supports`.
- [ ] Run `go test ./libs/entities/platform` and confirm the VK/capability assertions fail because the symbols do not exist.
- [ ] Add `PlatformVKVideoLive = "vk_video_live"`; replace Twitch/Kick predicates with only generally useful predicates; add `Capability` constants (`chat.read`, `chat.write`, `chat.reply`, `moderation.delete`, `streams.read`, `events.follow`, `events.raid`, `events.reward`) and `Capabilities.Supports`.
- [ ] Add `ErrUnsupportedCapability{Platform, Capability}` implementing `error`.
- [ ] Re-run `go test ./libs/entities/platform` and `gofmt` the package.

### Task 2: Create normalized binding storage

**Files:**
- Create: `libs/migrations/postgres/20260721120000_channel_platforms.sql`
- Create: `libs/repositories/channel_platforms/channel_platforms.go`
- Create: `libs/repositories/channel_platforms/model/model.go`
- Create: `libs/repositories/channel_platforms/pgx/pgx.go`
- Create: `libs/repositories/channel_platforms/pgx/pgx_test.go`

- [ ] Write repository tests for `Create`, `GetByChannelAndPlatform`, `GetByPlatformChannelID`, `ListByChannelID`, `Update`, and `Delete`; include duplicate `(channel_id, platform)` and a duplicate provider account linked to another Twir channel.
- [ ] Confirm these tests fail because the repository/package does not exist.
- [ ] Create `channel_platforms` with `id uuid primary key default uuidv7()`, `channel_id uuid not null references channels(id)`, `platform text not null`, `user_id uuid not null references users(id)`, `platform_channel_id text not null`, `enabled boolean not null default true`, `bot_user_id uuid null references users(id)`, `bot_config jsonb not null default '{}'::jsonb`, timestamps, `unique(channel_id, platform)`, and `unique(platform, platform_channel_id)`.
- [ ] Add B-tree indexes for `channel_id`, `user_id`, and `(platform, platform_channel_id)`. Keep all FK columns indexed.
- [ ] Implement pgx queries with explicit columns and map to a repository model containing `Nil`/`IsNil`.
- [ ] Run the repository tests against the migration database, then run `go test ./libs/repositories/channel_platforms/...`.

### Task 3: Backfill and remove legacy channel columns

**Files:**
- Modify: migration from Task 2
- Modify: `libs/repositories/channels/model/model.go`
- Modify: `libs/repositories/channels/channels.go`
- Modify: `libs/repositories/channels/pgx/pgx.go`
- Modify: `libs/services/channels/channels.go`
- Create: `libs/repositories/channels/pgx/multiplatform_migration_test.go`

- [ ] Write migration tests with Twitch-only, Kick-only, and dual-bound channels; assert each backfilled binding retains user ID, provider ID, enabled flag, and Kick bot configuration.
- [ ] Run the tests and confirm they fail before backfill exists.
- [ ] Backfill Twitch/Kick rows with `INSERT ... SELECT ... ON CONFLICT DO NOTHING`; validate no source row is silently omitted before dropping legacy columns.
- [ ] Change `Channel` to contain only Twir fields plus `Bindings []channel_platformsmodel.ChannelPlatform`; remove `TwitchUserID`, `KickUserID`, provider IDs and bot columns.
- [ ] Replace repository lookups with `GetByBindingUserID`, `GetByPlatformChannelID`, and `GetBySlug` joins against `channel_platforms`/`users`; update `ChannelService` to call those generic methods without a `switch`.
- [ ] Only after all compile-time callers are migrated, add a follow-up migration dropping legacy columns and obsolete indexes.
- [ ] Run `go test ./libs/repositories/channels/... ./libs/services/channels/...`.

### Task 4: Migrate existing provider consumers

**Files:**
- Modify: callers returned by `rg 'GetByTwitch|GetByKick|TwitchUserID|KickUserID' apps libs`
- Modify: `libs/cache/channel/channel.go`
- Modify: `apps/eventsub/internal/kick/{handlers.go,subscription_manager.go}`
- Modify: `apps/events/internal/{listener,workflows}`
- Modify: `apps/parser/internal/{cacher,variables}`
- Modify: `apps/scheduler/internal/timers/streams.go`
- Modify: corresponding `*_test.go` fakes.

- [ ] Add regression tests for generic channel lookup by platform user ID and platform channel ID for Twitch and Kick.
- [ ] Verify the tests fail after legacy repository methods are removed.
- [ ] Replace every provider-specific lookup with `GetChannelByBindingUserID` or `GetChannelByPlatformChannelID`; obtain provider metadata from the selected binding.
- [ ] Convert fake repositories to implement the new generic interface; do not retain compatibility methods.
- [ ] Run Go tests for each changed app and `go test ./apps/eventsub/... ./apps/events/... ./apps/parser/... ./apps/scheduler/...`.

### Task 5: Introduce provider registry and capability-aware actions

**Files:**
- Create: `libs/platforms/registry.go`
- Create: `libs/platforms/registry_test.go`
- Create: `apps/bots/internal/platforms/chat.go`
- Create: `apps/bots/internal/platforms/chat_test.go`
- Modify: `apps/bots/internal/services/channel/message.go`
- Modify: `apps/bots/internal/kick/chat_client.go`
- Modify: `apps/bots/internal/twitchactions/sendmessage.go`

- [ ] Write tests that dispatch an enabled Twitch/Kick binding through the registry, skip bindings lacking `chat.write`, and preserve one platform's failure without blocking another.
- [ ] Confirm the tests fail because the registry and adapter interface are missing.
- [ ] Define a narrow `ChatAdapter` accepting a binding, message, reply ID, and options; register Twitch and Kick implementations with their capability sets.
- [ ] Replace `switch p` in `SendMessage` with binding selection plus registry dispatch. An empty platform request selects all enabled `chat.write` bindings.
- [ ] Return `ErrUnsupportedCapability` for explicitly selected but unsupported operations and aggregate per-binding errors for independent failures.
- [ ] Run `go test ./apps/bots/... ./libs/platforms/...`.

### Task 6: Decouple canonical NATS messages from repository models

**Files:**
- Modify: `libs/bus-core/generic/chat-message.go`
- Modify: `apps/eventsub/internal/{handler,kick/handlers.go,mappers/chat_message.go}`
- Modify: `apps/bots/internal/messagehandler/**/*.go`
- Modify: tests for chat handling and parser publication.

- [ ] Write a serialization test proving `generic.ChatMessage` contains binding/channel IDs and role flags but no repository `Channel`, `Stream`, `User`, or stats model.
- [ ] Confirm the test fails due to `EnrichedData` serializing repository models.
- [ ] Split transport fields from local enrichment: retain canonical IDs, content, fragments, badges, roles, platform, and external event ID; load channel/user/stream state inside each consumer before it needs it.
- [ ] Keep NATS subjects unchanged so downstream transport topology does not change.
- [ ] Run `go test ./libs/bus-core/... ./apps/eventsub/... ./apps/bots/internal/messagehandler/...`.

### Task 7: Generalize event transport subscription lifecycle

**Files:**
- Create: `apps/eventsub/internal/platforms/transport.go`
- Create: `apps/eventsub/internal/platforms/transport_test.go`
- Modify: `apps/eventsub/internal/webhook/manager.go`
- Modify: `apps/eventsub/internal/kick/subscription_manager.go`
- Modify: `libs/bus-core/eventsub/eventsub.go`

- [ ] Write tests that a transport receives only bindings for its platform, subscribes enabled bindings, and continues after one binding fails.
- [ ] Confirm they fail because the manager queries Kick-specific channel fields.
- [ ] Define `EventTransport` with `Platform`, `Capabilities`, `Subscribe`, `Unsubscribe`, and callback-base URL setup; make the manager query generic bindings and route by registry.
- [ ] Move Kick's subscription state key to include binding ID, event type, and transport kind, preventing future webhook/WS collisions.
- [ ] Run `go test ./apps/eventsub/...`.

### Task 8: Add VK configuration and identity adapter

**Files:**
- Modify: `libs/config/config.go` and sample env documentation
- Create: `apps/api-gql/internal/platform/vkvideo/provider.go`
- Create: `apps/api-gql/internal/platform/vkvideo/provider_test.go`
- Modify: `apps/api-gql/internal/platform/registry.go`
- Modify: `apps/api-gql/internal/platform/platform.go`

- [ ] Obtain the authenticated VK developer-cabinet contract before writing fixtures: authorization URL parameters, token endpoint, profile endpoint, refresh semantics, scopes, and error response schemas.
- [ ] Record sanitized HTTP fixtures for authorization-code exchange, refresh, and profile lookup; write tests against an `http.RoundTripper` fixture transport.
- [ ] Confirm each test fails before the VK provider exists.
- [ ] Add `VKVideoClientID`, `VKVideoClientSecret`, callback URL, webhook secret, API base URL, and feature flag config. Mark secrets as required only when the feature flag is enabled.
- [ ] Implement the same `PlatformProvider` interface used by Kick with PKCE/state validation, typed provider errors, context-bound HTTP requests, and no credentials in logs.
- [ ] Run `go test ./apps/api-gql/internal/platform/...`.

### Task 9: Make OAuth bind any platform without provider switches

**Files:**
- Modify: `apps/api-gql/internal/delivery/http/routes/auth/oauth-platform.go`
- Modify: `apps/api-gql/internal/delivery/http/routes/auth/post-platform-code.go`
- Modify: `apps/api-gql/internal/delivery/http/routes/auth/oauth-platform_test.go`
- Modify: `apps/api-gql/internal/auth/sessions_user.go`

- [ ] Write tests for creating a VK-only Twir channel, linking VK to an existing Twitch/Kick channel, rejecting an account linked to another channel, and preserving existing bindings.
- [ ] Confirm they fail because the auth flow creates/updates Twitch/Kick fields.
- [ ] Make channel creation independent of platform; create or update a `channel_platforms` row inside the auth transaction, then store user tokens as today.
- [ ] Keep Kick's special bot selection in a binding-specific configuration path, not in the shared OAuth function.
- [ ] Publish generic EventSub subscription requests after a binding transaction commits.
- [ ] Run `go test ./apps/api-gql/internal/delivery/http/routes/auth/...`.

### Task 10: Add VK webhook transport and canonical normalizers

**Files:**
- Create: `apps/eventsub/internal/vkvideo/{transport.go,webhook.go,normalizer.go,webhook_test.go,normalizer_test.go}`
- Modify: `apps/eventsub/internal/http/server.go`
- Modify: `apps/eventsub/internal/platforms/transport.go`

- [ ] Obtain signed example payloads and verification/challenge behavior from VK documentation; store only redacted fixtures.
- [ ] Write tests for signature rejection, verification response, malformed JSON, duplicate `(vk_video_live,event_id)`, chat-message normalization, and each documented supported event.
- [ ] Confirm the endpoint tests fail before routing is registered.
- [ ] Register `/webhook/vk-video-live`; verify requests before enqueueing/normalizing; use the common deduplication store.
- [ ] Map only documented events to `generic.ChatMessage` and `events.*`; publish no provider payloads to NATS.
- [ ] Register the transport only when the VK feature flag is enabled; leave its WebSocket implementation absent rather than adding a speculative protocol.
- [ ] Run `go test ./apps/eventsub/internal/vkvideo/... ./apps/eventsub/...`.

### Task 11: Expose bindings and capabilities through GraphQL

**Files:**
- Modify: `apps/api-gql/internal/delivery/gql/schema/{platform.graphql,channels.graphql}`
- Create: `apps/api-gql/internal/delivery/gql/mappers/channel_platform.go`
- Modify: `apps/api-gql/internal/delivery/gql/mappers/channels.go`
- Modify: relevant resolvers and resolver tests.

- [ ] Write resolver tests for a channel returning all bindings, provider profile data, enabled state, and capability strings; ensure VK appears only when enabled by config.
- [ ] Confirm they fail before the schema type exists.
- [ ] Add `ChannelPlatformBinding` and `PlatformCapability` schema types; add authenticated binding connect/disconnect/status mutations that route to the generic service.
- [ ] Run `bun cli build gql`, re-read generated contracts, implement only generated resolver interfaces, then run `go test ./apps/api-gql/internal/delivery/gql/...`.

### Task 12: Build dashboard binding management

**Files:**
- Create: `web/layers/dashboard/features/channel-platforms/{api.ts,composables/use-channel-platforms.ts,ui/platform-binding-card.vue,ui/platform-bindings.vue}`
- Modify: `web/layers/dashboard/pages/dashboard/bot-settings.vue`
- Modify: generated GraphQL client artifacts through the project's generation command.
- Test: feature/component tests in the dashboard's existing test framework.

- [ ] Load `nuxt`, `reka-ui`, and `shadcn-vue` skills before implementation because this changes Nuxt UI and component primitives.
- [ ] Write component tests for connected/disconnected bindings, capability-disabled actions, and OAuth redirect initiation.
- [ ] Confirm the tests fail before the feature is mounted.
- [ ] Query the new bindings once, render one independent card per platform with connect/disconnect/status controls, and use Lucide icons plus existing shadcn components.
- [ ] Use the capability list to hide unsupported controls; do not branch templates on provider names except display metadata in a centralized platform presentation map.
- [ ] Run the dashboard's focused tests and `bun --cwd web run build` once dependencies are installed correctly.

### Task 13: Complete migration cleanup and release controls

**Files:**
- Create: `libs/migrations/postgres/20260721120001_drop_legacy_channel_platform_columns.sql`
- Modify: deployment/config documentation
- Modify: tests and fixtures still referring to legacy fields.

- [ ] Add a schema assertion test that no production query references `twitch_user_id`, `kick_user_id`, `twitch_bot_enabled`, `kick_bot_enabled`, `kick_bot_id`, `botId`, or `isTwitchBanned` after the cleanup migration.
- [ ] Confirm it fails while legacy callers remain.
- [ ] Remove columns, obsolete indexes, and legacy repository methods only after Tasks 3 through 12 pass.
- [ ] Document feature-flag rollout: disabled, internal VK test channel, monitored enablement, then general availability; define rollback as disabling VK without rolling back shared binding migrations.
- [ ] Run focused Go tests, GraphQL generation/build, dashboard checks, `bun lint`, and finally `bun cli build` after repairing the baseline dependency installation.

## Plan Review

- Spec coverage: Tasks 1-4 normalize data and consumers; Tasks 5-7 isolate actions, NATS, and transports; Tasks 8-10 implement VK identity and webhooks; Tasks 11-12 expose dashboard/API; Task 13 completes safe rollout and cleanup.
- Known external dependency: Tasks 8 and 10 require the access-gated VK provider contract. They must not invent OAuth endpoints, signing rules, or WebSocket details.
- Baseline: `bun cli build` currently fails in this worktree before implementation because existing symlinked dependencies cannot resolve `@twir/*`, Nuxt, and overlay tooling. This must be repaired or reproduced outside this feature before accepting full-build verification.
