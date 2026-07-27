# PR #1004 Refactor Incompleteness Audit

Scope: `origin/main...0c9ed8bdc`. This report covers functional regressions introduced by the PR and consumers the PR design explicitly required to migrate but left behind. Unrelated pre-existing defects are excluded. A supplementary security gate is reported separately because it found a PR-introduced merge blocker.

## Supplementary Merge Blocker

### P0: A dashboard collaborator can turn a platform binding into permanent ownership

- `channelPlatformConnect` and `channelPlatformDisconnect` require only `MANAGE_BOT_SETTINGS`, so a delegated collaborator can start targeted OAuth and mutate bindings.
- `dashboard_access.IsOwner` treats the user behind any normalized binding as an unconditional channel owner.
- A collaborator can connect their own provider account, become an owner, disconnect the original owner's binding, and retain owner access after the delegated role is revoked. The last-binding guard preserves the attacker's binding rather than preventing the takeover.
- Identity-changing connect/disconnect must require an existing owner or bot administrator; add a regression test covering collaborator link, original-owner disconnect, logout/login, and role revocation.

## Confirmed Branch Regressions

### P0: PR CI is blocked by new frontend lint errors

- Full PR build succeeds, then `bun lint` fails.
- The current GitHub run reports 26 errors. At least seven are in PR-owned channel-platform binding specs and VK bot setup/callback files: imports after executable declarations and unsorted imported members.
- Other lint errors are unrelated baseline files and should be tracked separately, but the PR-specific errors alone keep the check red.

### P0: VK chat persistence rejects provider message IDs

- VK normalizer assigns numeric provider message ID to both `generic.ChatMessage.ID` and `MessageID`.
- Bots forwards `msg.ID` unchanged into `chat_messages.id`.
- Postgres and ClickHouse schemas both define `chat_messages.id` as UUID.
- Live runtime failed with `id invalid UUID length: 9`.
- Required correction: separate internal persistence UUID from external provider message ID; retain provider ID for replies/deletion/deduplication.

### P0: VK outbound stream discovery uses the wrong endpoint

- PR calls `GET /v1/channels/active`; runtime receives 404.
- The implemented route does not match the provider's catalog/channel lookup contract; the exact production route must be pinned from the authoritative VK application contract rather than the repository's self-mocking tests.
- `POST /v1/chat/message/send` request shape is otherwise aligned.
- Bots construct the client with empty options, so it ignores the configured API base URL and silently selects the development host.

### P1: Kick messages no longer count third-party emotes

- Before the PR, the Kick webhook handler fetched channel/global emotes and attached the resulting user-stat enrichment before publishing chat messages.
- The PR moved enrichment into the shared bots handler, but `chatMessageCountEmotes` immediately returns for every non-Twitch message.
- Kick chat still carries native Kick emote fragments, but 7TV/BTTV/FFZ usage and the associated user emote counters are no longer enriched or counted.

## Required Consumers Left Unmigrated

### P0: Required webhook-first VK ingress is absent

- The design and Task 10 require a verified `/webhook/vk-video-live` transport with challenge handling, signature validation, bounded processing, idempotency, and a future-compatible transport boundary.
- The HTTP server registers only `/webhook/kick`; there is no VK webhook handler or route.
- VK ingress is instead wired exclusively through a hard-coded development Centrifugo WebSocket endpoint, reversing the planned webhook-first/future-WebSocket sequence and leaving no production endpoint configuration.

### P1: Five Kick event publishers were not migrated to canonical channel identity

- Task 4 explicitly includes Kick handlers and the events listener; Task 6 requires canonical payloads to carry both Twir channel identity and external provider IDs.
- Subscribe, resubscribe, subgift, redemption, and ban handlers put the internal channel UUID in `BaseInfo.ChannelPlatformID` and leave `ChannelDBID` empty.
- The listener interprets a missing `ChannelDBID` by looking up `ChannelPlatformID` as a Kick provider ID, logs the lookup error, and acknowledges the event without executing workflows or chat alerts.
- The malformed publisher behavior predates the branch, but leaving it intact is a missed migration requirement rather than an unrelated issue; Follow demonstrates the correct dual-ID shape.

### P1: Twitch ban state remains globally coupled to every event platform

- `events_flow` skips every matching event when the Twitch binding's `IsTwitchBanned` flag is true, including Kick/VK events and platform-neutral operations.
- The aggregate gate existed on `origin/main`; this PR migrated its storage location without removing the cross-platform behavior.
- Classify this as a missed platform-neutral migration requirement, not a newly introduced predicate.

### P1: Greetings bypass capability-aware dispatch

- Task 5 requires outgoing actions to select enabled bindings and use adapters.
- `handleGreetings` still explicitly selects the Twitch binding and directly calls `twitchActions.SendMessage`.
- Once VK stream state exists, Kick/VK-origin greetings would still route to Twitch or disappear when no active Twitch binding exists. Today VK greetings return even earlier because VK has no stream state.
- Greeting is marked processed before parsing/send and send errors are discarded, making a missed cross-platform greeting permanent.

### P1: Parser legacy cacher still resolves every channel as Twitch

- Task 4 requires parser consumers to use generic binding lookup.
- PR changed `getDbChannel` to `GetChannelByPlatformChannelID(PlatformTwitch, parseCtxChannel.ID)` rather than using the already available internal `DBChannelID` or active platform binding.
- Kick/VK mentioned-user stats and legacy integration variables therefore cannot resolve their channel correctly.

### P1: Timers hardcode Twitch/Kick despite VK `chat.write`

- Design says outgoing actions use capability checks and shared features such as timers retain platform selection through capability data.
- `getTimerSendTargets` iterates a literal Twitch/Kick list and tests require VK to be skipped.
- VK has an enabled chat adapter and declares `chat.write`, so timer delivery should flow through the same adapter contract or the capability declaration must explicitly model why it cannot.

### P1: Dispatcher does not enforce reply/announce capabilities

- `Dispatch` always requires only `chat.write`, regardless of non-empty `replyID` or announce options.
- VK declares no `chat.reply`; its adapter discards `replyID` and sends a plain message, so reply-mode commands silently change behavior.
- Kick/VK adapters also ignore announce options instead of returning an explicit unsupported result or applying a documented fallback.
- Existing dispatcher tests verify option forwarding to Twitch, but not capability selection based on requested operation.

### P1: Shared feature selector cannot select VK

- Backend `Platform` includes VK and VK ingress publishes commands.
- Shared `platform-selector.vue` exposes only Twitch/Kick and is used by commands, keywords, timers, events, and chat-history filters.
- This prevents configuring VK for features that the new transport/adapter can execute and bypasses the design requirement for capability-gated shared controls.

### P2: Capability table under-reports implemented ingress

- Kick subscribes to incoming chat and VK publishes incoming chat/commands.
- Both omit `chat.read`; VK declares only `chat.write`.
- GraphQL/dashboard expose this table, so runtime support and displayed support diverge.

### P2: Parser deduplication is not platform-scoped

- Task 6 introduced parser-side deduplication for canonical chat messages, but the Redis key is only `parser:dedup:<messageID>`.
- The design requires duplicate identity to be `(platform, external_event_id)`; equal provider IDs received from different platforms within the 60-second TTL suppress one another.
- VK transport's own 24-hour key is platform-scoped, so this defect is specifically in the new shared downstream deduplication layer.

### P2: VK profile/current-user presentation falls through to Twitch

- `resolveProfile` and current-profile mapping were not migrated by the PR.
- Any non-Kick profile becomes a Twitch URL/profile, so VK-linked users can receive wrong links or blank/wrong dashboard identity.
- This is a missed consumer of the newly expanded `Platform` enum, not a newly introduced helper bug.

### P2: Frontend chat/overlay platform type remains Twitch/Kick-only

- `libs/frontend-chat` still defines only `twitch | kick`; overlay consumers cast API platform into that narrower type.
- VK messages therefore cannot render correct platform icon/badge semantics and may fall through to Twitch behavior.

## Deferred Or Lower-Confidence Gaps

- VK never receives persisted stream state. Scheduler polling covers Twitch/Kick and VK ingress publishes chat only, so online-only commands, stream variables/counters, first-stream-user logic, and greetings treat VK as offline. The plan defers access-gated Video/stream contracts, so this is an observable limitation requiring an explicit product decision rather than a proven Task 10 regression.
- Task 13 cleanup is incomplete: no follow-up migration drops legacy channel platform columns/APIs. This is planned migration debt, not an independent runtime blocker while canonical reads use bindings.

## Excluded As Unrelated To This PR

- Existing Kick Unicode splitting failures.
- Existing bare `/ban` or `/timeout` panic.
- Existing unban-request workflow type miswire.
- Existing role-less command stat permission behavior.
- Method-level `GetGbUserStats` cache key defect that current callers do not trigger.
- Existing donation publisher uses the legacy `base_info.channel_id` shape; neither that producer nor the Go `BaseInfo` contract changed in this PR, and the VK plan does not include donation migration.

## Rejected Audit Candidates

- VK subscriptions are restored on process start: `webhook.Manager.subscribeAllPlatforms` iterates all registered transports, including VK. The Twitch/Kick-only `reinitBoundChannels` belongs to a separate Twitch EventSub cleanup operation.
- VK OAuth does not require the `device_id` flow described by an early draft; the implemented server OAuth contract uses `auth.live.vkvideo.ru`.
- Kick subscription-manager binding lookup is correct; the broken identity is limited to five event payload producers.

## Verification

- Full GitHub build completes; PR check fails at lint.
- Local `bun lint` reproduces the same 26 errors; at least seven belong to files added or changed for channel-platform/VK UI in this PR.
- New platform entity packages: 33 race-enabled tests passed.
- Platform registry: 2 race-enabled tests passed.
- Timer manager: 14 race-enabled tests passed, including the current VK-skip expectation.
- API-GQL channel-platform/auth packages: 60 race-enabled tests passed.
- Dashboard channel-platform feature: 9 tests passed.
- VK EventSub package: 30 race-enabled tests passed.
- Parser internal suite: 40 tests passed across 73 packages; no VK stream-state end-to-end coverage exists.
- Live ClickHouse probe rejects numeric VK ID `123456789` with `CANNOT_PARSE_UUID`.
- Temporary producer-consumer probes reproduced all five malformed Kick event identities and the pre-existing donation JSON mismatch; all temporary artifacts were removed.
- Five independent review lanes all returned FAIL for merge readiness. The supplementary security lane identified the binding-ownership takeover above.
