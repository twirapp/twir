# Multiplatform Robustness Audit

Status: wave 3 verification in progress against `0c9ed8bdc`.

## Executive Summary

Twir has working Twitch, Kick, and VK Video Live paths, but multiplatform support is not governed by one end-to-end contract. The highest-risk failures arise where older Twitch assumptions survive behind newer platform bindings: workflow filtering, parser channel resolution, greetings, provider profile routing, and outbound drop semantics.

The committed VK parser fix in `0c9ed8bdc` is sound at the ingress boundary, but runtime exposed two downstream failures: provider message IDs are written to a UUID persistence field, and outbound VK stream discovery depends on an endpoint that returned 404. These are independent from the rejected claim that provider channel/user IDs reach UUID consumers; all live chat ingress paths normalize channel and user identities before publication.

## Priority Matrix

| Priority | Finding | Evidence | User impact | Minimal fix direction |
|---|---|---|---|---|
| P0 | Unban request invokes redemption workflow | execution + source | wrong automation runs; intended workflow skipped | pass `EventTypeChannelUnbanRequestCreate`; add listener/workflow test |
| P0 | Greetings always target Twitch, mark processed first, discard errors | source; wave-3 execution pending | Kick/VK greeting misroute or permanent loss | platform dispatch; persist processed only after successful outcome |
| P0 | Bare `/ban` or `/timeout` response can panic | source; wave-3 execution pending | bot process crash from user-configurable response | parse arity before indexing; regression tests |
| P1 | VK provider message ID enters UUID persistence key | live runtime + source | ClickHouse message persistence fails | generate/internalize persistence UUID; retain provider ID separately |
| P1 | VK active-stream discovery returns 404 | live runtime + source/docs | VK bot cannot reply | replace discovery with documented VK Video Live flow or explicit unsupported/offline result |
| P1 | Parser response slice race/order loss | source; wave-3 execution pending | lost/reordered command responses, data race | indexed result slots or serial construction |
| P1 | Role-less stat requirements bypassed | source/UI; wave-3 execution pending | restricted commands become public | evaluate role and stat predicates independently |
| P1 | Twitch ban suppresses eligible Kick/VK events | execution + source | cross-platform automations silently skipped | scope ban predicate to Twitch-dependent operations/events |
| P1 | Parser channel lookup hardcodes Twitch | source; wave-3 execution pending | Kick/VK mentioned-user variables fail | use internal `DBChannelID` |
| P1 | Kick UTF-8 split tests fail | committed tests + source | malformed/oversized outbound messages | split by UTF-8-safe byte budget |
| P1 | Twitch SkipRateLimits is dead; rejection reports success | source; wave-3 execution pending | requested bypass ignored; delivery loss hidden | honor flag and return typed dropped outcome |
| P2 | 7TV profile cache unkeyed | execution + source | wrong user/platform profile in reachable command path | key by platform + user ID and lock shared access |
| P2 | Capability declarations under-report ingress/features | source/tests/UI | API/dashboard support view diverges from runtime | expand capability vocabulary and conformance tests |
| P2 | VK command/profile/chat UI support fragmented | schema + source | unavailable selectors, wrong Twitch links/branding | typed platform fields and explicit VK branches |
| P2 | Kick 429 converted to successful send | source/docs | intentional drops invisible upstream | typed sent/dropped/rate-limited outcome + metric |
| P2 | Watched stats can double-count during start-first rollout | source/config | temporary stat inflation on deploy | distributed tick lock or idempotent time windows |

## Contract Findings

### Identity

- `ChannelID`, `ChannelBindingID`, and `UserID` are internal UUIDs at bus consumer boundaries.
- `PlatformChannelID`, `ChatterUserId`, and `MessageID` are provider-native identifiers.
- `generic.ChatMessage.ID` is used as the persistence key, but VK currently supplies its numeric provider message ID. This is the live ClickHouse failure.
- `events.BaseInfo.ChannelPlatformID` is stringly typed and has split semantics: Twitch generally supplies provider ID while Kick moderation supplies internal channel UUID and API code compensates.

### Capabilities

- Kick/VK chat ingress exists but `chat.read` is absent from the canonical capability model.
- Twitch stream-edit and moderation-ban behavior cannot be represented by current capability names.
- VK timer exclusion and missing-bot setup no-op are explicit, test-covered unsupported states.
- Robust OSS comparators centralize typed capabilities, return explicit unsupported errors, and run adapter conformance/golden tests.

### Error Semantics

- Several outbound paths log failures but return success to callers.
- Kick 429 and Twitch limiter rejection are message-drop outcomes represented as `nil` errors.
- A boolean/error-only adapter contract cannot distinguish sent, deliberately dropped, unsupported, retryable, and permanent failure outcomes.

## Rejected Or Narrowed Claims

- Rejected: raw Twitch provider channel/user IDs reach UUID consumers. The handler normalizes both before bus publication.
- Narrowed: `GetGbUserStats` has an unkeyed cache defect, but current production callers use one target per fresh cacher.
- Narrowed: VK exclusion is not universally accidental; timers explicitly exclude VK and tests require it.
- Rejected: missing VK `BotUserID` success is accidental; a focused test defines it as intentional setup behavior.
- Narrowed: Kick 429 handling is not silent or accidental; it explicitly logs an intentional drop, but upstream still sees success.

## Regression Architecture

1. Define one field-role contract for internal UUID, provider user/channel ID, provider message ID, and persistence event ID.
2. Add ingress golden fixtures for Twitch/Kick/VK that assert all identity roles before publication.
3. Add platform adapter conformance tests for capabilities, unsupported behavior, and outbound outcomes.
4. Add parser tests for response ordering/races, stat-only permissions, platform channel resolution, and cache keys.
5. Add event listener/workflow tests for every event-type mapping and platform-scoped gate.
6. Add UI conformance tests deriving selectors/profile links/icons from the same platform model.
7. Preserve explicit unsupported cells, especially VK timers, rather than adding implicit fallbacks.

## Recommended Fix Order

1. Correct deterministic wrong behavior and crashes: unban mapping, greetings, moderation parsing.
2. Repair live VK persistence and outbound discovery.
3. Fix parser concurrency, permissions, and platform lookup.
4. Make outbound outcomes explicit and repair Kick/Twitch rate/splitting semantics.
5. Align capability declarations and UI platform typing.
6. Add the cross-platform conformance matrix before expanding provider support further.

## Verification Snapshot

- `apps/eventsub`: VK package, 30 tests passed under `-race`.
- `apps/events`: workflow package, 1 focused test passed under `-race`.
- `apps/parser`: 40 tests passed across 73 internal packages under `-race`; contested branches still lacked native tests.
- `apps/bots`: 18 tests passed and 2 committed Kick split tests failed under `-race`.
