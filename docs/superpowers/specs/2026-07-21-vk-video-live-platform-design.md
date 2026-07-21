# VK Video Live Platform Design

## Goal

Make Twir channels platform-neutral and add VK Video Live as a third platform alongside Twitch and Kick. A Twir channel may have one binding for each platform; none is primary. VK uses webhooks initially while preserving a transport boundary that can accept WebSockets later.

## Scope

- Fully replace Twitch/Kick-specific channel bindings with a normalized platform-binding model.
- Migrate existing Twitch and Kick functionality to the new model before adding VK.
- Keep platform-specific API code inside adapters and pass only canonical messages through NATS.
- Support every feature the provider documents and expose unavailable features through capabilities rather than platform checks.
- Add OAuth/API application configuration, VK binding management, webhook ingestion, outgoing chat actions, GraphQL, and dashboard UI.

## Data Model

`channels` remains the platform-neutral Twir-channel aggregate: ownership, commands, timers, roles, API key, and shared settings belong here. It does not retain provider user IDs, provider bot flags, or provider-specific moderation fields.

Create `channel_platforms` with:

- `channel_id` referencing `channels`.
- `platform`: `twitch`, `kick`, or `vk_video_live`.
- `user_id` referencing the linked platform user.
- `platform_channel_id`, login, display name, and avatar snapshot.
- Binding-specific status, enabled state, bot identity/configuration, and token association.
- A unique constraint on `(channel_id, platform)`.
- Lookup indexes on `(platform, platform_channel_id)` and foreign-key columns.

Existing Twitch and Kick fields move to bindings through an additive migration. Once every caller uses bindings, remove the legacy columns and platform-specific repository APIs. The migration must retain a channel's exact enabled state, bot configuration, and lookup identity.

## Platform Contracts

Platform code implements narrow, provider-specific adapters:

- `Identity`: authorization, token renewal, and resolving the authenticated account.
- `Chat`: sending messages and replies.
- `Moderation`: supported moderation operations.
- `Streams`: stream state and metadata, when available.
- `EventTransport`: subscription lifecycle and incoming provider events.

A registry resolves adapters and their immutable capability sets by `platform.Platform`. Consumers operate on a capability, not on a provider name. Missing support returns a typed `unsupported capability` error; the dashboard hides or disables the corresponding control.

VK implements `EventTransport` with verified webhooks. A future WebSocket transport plugs into the same adapter and feeds the same normalizers, subscription state, idempotency, and canonical events. Webhooks remain available as fallback.

## Event and Action Flows

Incoming events follow one path:

`Webhook or WebSocket -> platform adapter -> normalizer -> canonical NATS event -> consumers`

Normalizers publish universal `generic.ChatMessage` and `events.*` payloads containing Twir channel ID, platform, external IDs, actor data, content, and roles. Repository models and other local database state are not serialized into NATS. Consumers enrich state locally or through a dedicated enrichment stage.

Duplicate deliveries are identified by `(platform, external_event_id)` and must be safe for webhook retries and future WebSocket overlap.

Outgoing actions follow one path:

`command, timer, or action -> select channel bindings -> capability check -> adapter -> provider API`

An empty platform filter selects every enabled binding that supports the capability. A specified filter selects only its enabled, supported bindings. Failure on one platform is logged and isolated from the other platform actions.

## API and Dashboard

GraphQL returns the list of a channel's platform bindings and their capabilities. It exposes connect, disconnect, status, and binding-specific bot authorization flows without representing a primary platform.

The dashboard renders independent Twitch, Kick, and VK Video Live binding cards. Each supports connect/disconnect, status, bot authorization state, and the provider's available capabilities. Shared features such as commands and timers retain platform selection but use the capability data to prevent invalid choices.

## Webhooks and Operations

VK validates every webhook signature and challenge before accepting an event. It completes provider verification synchronously; processing is bounded by a timeout and made idempotent. Subscription state is reconciled periodically so deleted or failed subscriptions are recreated.

Provider credentials and application keys are configuration secrets. VK stays behind a feature flag and is activated only after a real test channel passes provider integration checks. VK errors cannot degrade Twitch or Kick handling.

## Delivery Sequence

1. Add normalized bindings and migrate current Twitch/Kick data.
2. Convert repositories, channel lookup, chat actions, subscription managers, GraphQL, and dashboard to bindings and capabilities.
3. Remove legacy channel columns and APIs after all consumers are converted.
4. Add the VK OAuth/API client, adapter registry entry, webhook transport, normalizers, and outgoing actions.
5. Add VK dashboard binding UI and provider-backed integration tests.
6. Enable VK for a test channel under a feature flag, then release broadly.

## Verification

- Migration tests cover Twitch-only, Kick-only, and dual-platform channels.
- Repository and service tests cover uniqueness, lookup, disconnect, and binding deletion.
- Adapter contract tests cover capabilities and canonical NATS normalization.
- Webhook tests cover signatures, challenges, retries, duplicate delivery, malformed payloads, and reconciliation.
- Twitch and Kick regression tests cover chat, commands, events, and moderation after migration.
- GraphQL/dashboard tests cover binding management and capability-gated controls.
- VK integration tests run against a real test channel before enabling the feature flag.

## Constraints and Open Dependencies

The public VK Video Live page confirms an HTTP API, API-key application registration, and webhook event categories, but detailed API contracts are access-gated. Exact OAuth parameters, chat endpoints, webhook signing, event schemas, rate limits, and any WebSocket support must be confirmed from the developer cabinet before their adapter methods are implemented.
