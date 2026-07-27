# Wave 1: Identity and Stats

## Findings
- Canonical identity is internal UUID; provider IDs belong in `(platform, platform_id)` and `platform_channel_id`.
- `UnsureUser.UserID` is ambiguous: UUID-first internal lookup, then provider-scoped lookup.
- Legacy channel owner columns duplicate normalized bindings and rely on repository synchronization.
- `userswithstats` can return user + nil stats, making parts of UserCreator's ErrNotFound flow ineffective.
- Stats have multiple pgx/GORM/raw-SQL writers; DB constraints carry most safety.

## EXPAND
- Inventory every direct stats/binding writer.
- Verify UUID-shaped provider-ID collision behavior.
- Verify nil-stats creation/update path against real repositories.
