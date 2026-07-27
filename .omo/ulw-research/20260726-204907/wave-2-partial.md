# Wave 2 Partial Convergence

## Message Identity

- `PlatformChannelID` is consistently provider-native and is the safest routing key across services.
- `ChannelID` and `UserID` are consumed as internal UUIDs by bots. The raw Twitch mapper temporarily writes provider IDs, but `processChannelChatMessage` replaces both with internal UUIDs before either bus publish.
- `ID` is the persistence key. `MessageID` is the provider deletion/reply key and is not persisted in `chat_messages`.
- Kick, VK, Twitch, redemption, and song-request paths all normalize internal IDs before publishing. The raw Twitch mapper is only an intermediate representation.

## Capability Model

- Kick and VK implement chat ingress while the capability registry cannot declare `chat.read`.
- Twitch implements stream editing and moderation bans while the registry cannot name those capabilities.
- Kick outbound chat deletion is not implemented despite a documented provider endpoint.
- VK timer exclusion is explicit and test-covered; it is an intentional unsupported cell, not an accidental omission.

## UI And API

- The backend enum includes VK, but command/event/timer/keyword platform fields remain `String[]` and the shared selector exposes only Twitch and Kick.
- `resolveProfile` maps Kick explicitly and every other platform to Twitch, producing wrong VK profile links.
- Dashboard profile state and frontend-chat branding only model Twitch/Kick; VK falls through to Twitch or blank state.
- Channel-platform binding and bot-status surfaces do support VK, proving support is fragmented rather than globally absent.

## Persistence Writers

- Runtime direct writers outside canonical repositories exist in scheduler online/watched timers, admin bans, community stat reset, and EventSub VIP handling.
- Schema uniqueness largely prevents duplicates, but blind updates can silently miss absent rows. The scheduler defaults to one replica, yet `start-first` rollouts can overlap old/new instances and double-increment watched counters during deployment.
- Migration-only legacy owner/channel-platform writers are guarded by constraints and integration tests; they are lower operational risk than runtime writers.

## Provider Contracts

- Kick reply shape is aligned, but HTTP 429 is logged and returned as success, hiding message loss from callers.
- Twitch moderation and replies align with first-party APIs; search is first-page-only.
- A separate legacy VK integration still uses old host/query-token conventions. This is distinct from VK Video Live and must not be conflated with the active-stream 404.

## Open Verification

- Execute parser cache contamination proofs.
- Execute or statically isolate event miswire and Twitch-ban suppression.
- Complete SHA-pinned OSS comparison.
- Apply hostile skeptic review to all high-risk claims.
