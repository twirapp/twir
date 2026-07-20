# Default Command Collision Migration Design

## Scope

Update scheduler default-command provisioning only. The parser lookup order remains unchanged.

The rollout protects the new Dota default names and applies to every default command provisioned by the scheduler.

## Candidate Detection

A default is already provisioned only when an existing command has all of these properties:

- `Default` is true.
- `DefaultName` exactly matches the required default identifier.
- `Name` exactly matches the current required default name.

This makes stale or renamed defaults candidates for replacement while leaving correctly provisioned defaults untouched on later scheduler runs.

## Collision Migration

For each candidate channel, load the existing command `ID`, `Name`, `Aliases`, `Default`, and `DefaultName` inside the provisioning transaction.

Before adding a default, archive every existing command whose `Name` equals the new default name under `strings.EqualFold`:

- Rename it to `<old-name>-<UnixMillis>`.
- Add `-2`, `-3`, and so on when that archive name is already occupied.
- Set `Default` to false and clear `DefaultName`.
- Preserve aliases and all other command data.

The timestamp distinguishes migrated commands for channel users. Archiving applies to both custom and stale default commands. It never deletes a command.

After archival, `hasCommandConflict(defaultName, existing)` checks the remaining command names and aliases case-insensitively. An alias conflict leaves the alias unchanged, skips the default, and emits a concise warning with only the channel ID and default command identifier.

## Persistence And Races

Perform candidate loading, archival updates, and creation in one GORM transaction. Bulk-create the resulting defaults with:

```go
clause.OnConflict{
	Columns: []clause.Column{{Name: "channelId"}, {Name: "name"}},
	DoNothing: true,
}
```

The conflict clause prevents a concurrent exact-name custom-command insert from failing the full provisioning batch. It does not change parser resolution or delete or rename aliases.

## Tests

- Table-test `hasCommandConflict` for no conflict, exact and case-insensitive name conflicts, exact and case-insensitive alias conflicts, and an unrelated alias.
- Table-test timestamp archive-name generation, including an occupied suffix.
- Use GORM PostgreSQL DryRun mode to assert that the bulk insert includes `ON CONFLICT ("channelId", "name") DO NOTHING`.
- Run scheduler service tests, scheduler build, and `git diff --check`.
