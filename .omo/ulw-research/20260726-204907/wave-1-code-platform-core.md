# Wave 1: Platform Core

## Findings
- Platform truth is duplicated across `capabilitiesByPlatform`, `IsValid`, GraphQL schema enumeration, and `All()` in `libs/entities/platform/platform.go`; no consistency test binds them together.
- `ShouldExecute([], current)` is allow-all and is used by commands, keywords, events, and timers.
- Channel binding lookup is first-match in memory; DB uniqueness protects persisted rows but not malformed assembled entities.
- VK enablement is gated only by client ID/secret, while operational prerequisites fail later.
- Registry errors conflate unavailable providers and unsupported capabilities.

## EXPAND
- Parser variable platform branches and Twitch fallback semantics.
- Concrete VK/Kick provider capability declarations vs implementations.
- OAuth/provider switches duplicating registry logic.
