- Goose SQL migrations in this repo use `-- +goose Up/Down` with `-- +goose StatementBegin/End` wrappers.
- GraphQL schema files for api-gql are generated from `apps/api-gql/internal/delivery/gql/schema/`.
- Platform enum foundation is now shared across Go, Postgres, and GraphQL with matching `twitch`/`kick` values.
- `users.id` can be swapped to internal UUID in-place by staging a temp `internal_id`, converting all direct user FKs plus `channels.id`, then renaming old `users.id` to `twitch_id` for backward compatibility.
- `user_platform_accounts.platform` should use the existing Postgres `platform` enum and the Go entity should use `libs/entities/platform.Platform`.
- Deploying the user UUID migration requires a Redis `FLUSHDB` maintenance step because stored sessions still contain legacy Twitch IDs.
- `kick_bots` is a dedicated bot-token table with no platform column; migration uses Goose SQL wrappers and the entity follows the repo's nil-pattern with `isNil` + `IsNil()` + `Nil` singleton.

- Added platforms[] columns for commands/timers/keywords using the shared platform enum; migration runner completed successfully.
- pgx scanning/writing worked with the new platforms field in repository code, and the repo-level build checks passed in the affected Go modules.
