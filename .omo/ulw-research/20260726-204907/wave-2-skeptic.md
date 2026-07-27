# Wave 2 Skeptic Review

## Accepted High-Risk Claims

- Capability declarations under-report Kick/VK chat ingress.
- Multi-response parser construction races on a shared result slice and cannot preserve deterministic order.
- Commands with only statistical requirements bypass those requirements when `RolesIDS` is empty.
- Parser `getDbChannel` ignores the available internal `DBChannelID` and resolves provider IDs as Twitch, breaking legacy cacher paths on Kick/VK.
- Unban-request-created events execute redemption workflows.
- Twitch ban state suppresses platform-neutral operations for eligible Kick/VK events.
- Greetings are marked processed before send, always select the Twitch binding, and discard send/shoutout errors.
- Twitch `SkipRateLimits` is propagated but never honored; limiter rejection returns success.
- Bare `/ban` or `/timeout` custom responses can panic via unchecked token indexing.
- Existing Kick Unicode split tests fail.
- VK profile URL/current-profile mapping falls through to Twitch.

## Partial Or Conditional

- Event `BaseInfo.ChannelPlatformID` is semantically inconsistent: Twitch publishers use provider IDs while Kick moderation uses an internal channel UUID. API code explicitly compensates for the Kick exception.
- 7TV cache contamination is reachable through the profile command, not through every ordinary variable path.
- Outbound failures are often caller-invisible but usually logged; they are not uniformly silent.
- VK selector exclusion is a defect for commands/keywords/chat filtering, but timers are deliberately unsupported and tested.
- Kick 429 dropping is intentional overload shedding, though callers cannot distinguish it from delivery.

## Rejected

- Provider IDs reaching UUID-only chat consumer fields: all live ingress paths normalize before publication.
- Current `GetGbUserStats` cross-user leakage: method defect exists, but current callers do not supply two targets to one cacher.
- VK missing `BotUserID` no-op as accidental failure: behavior is explicit and tested.
- Duplicate same-platform binding dispatch hazard as normally reachable: database uniqueness prevents the state.
