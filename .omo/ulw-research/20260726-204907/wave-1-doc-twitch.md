# Wave 1: Twitch Contracts

## Findings
- Twitch chat semantics are richer and non-portable: reply/shared-chat metadata, badges, cheer, fragment variants, granular moderation scopes.
- IDs are opaque strings; token principal equality matters on many endpoints.
- Badges are per-message state, not durable cross-platform roles.
- Stream/category/offline semantics are provider-specific live-state contracts.
- Rate limits include token buckets and endpoint-specific cooldowns.

## EXPAND
- Audit normalized model for Twitch semantics presented as universal.
- Verify shared-chat/reply fields are not misrouted to Kick/VK.
- Compare scope/capability model granularity.
