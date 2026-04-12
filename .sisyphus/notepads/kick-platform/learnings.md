- Goose SQL migrations in this repo use `-- +goose Up/Down` with `-- +goose StatementBegin/End` wrappers.
- GraphQL schema files for api-gql are generated from `apps/api-gql/internal/delivery/gql/schema/`.
- Platform enum foundation is now shared across Go, Postgres, and GraphQL with matching `twitch`/`kick` values.
- `users.id` can be swapped to internal UUID in-place by staging a temp `internal_id`, converting all direct user FKs plus `channels.id`, then renaming old `users.id` to `twitch_id` for backward compatibility.
- `user_platform_accounts.platform` should use the existing Postgres `platform` enum and the Go entity should use `libs/entities/platform.Platform`.
- Deploying the user UUID migration requires a Redis `FLUSHDB` maintenance step because stored sessions still contain legacy Twitch IDs.
- `kick_bots` is a dedicated bot-token table with no platform column; migration uses Goose SQL wrappers and the entity follows the repo's nil-pattern with `isNil` + `IsNil()` + `Nil` singleton.

- Added platforms[] columns for commands/timers/keywords using the shared platform enum; migration runner completed successfully.
- pgx scanning/writing worked with the new platforms field in repository code, and the repo-level build checks passed in the affected Go modules.
- Channels multi-platform migration can safely reuse `channels.id` as the legacy-to-new UUID mapping source after renaming it to `user_id`, then re-point all child FKs by updating UUID values before re-adding the original FK definitions.
- QA for the new `(user_id, platform)` uniqueness must copy required non-null channel fields like `botId` from an existing Twitch row; bare `INSERT (user_id, platform)` will fail on existing table constraints unrelated to the new schema.
- Tokens repository callers now need to parse legacy string user IDs into `uuid.UUID` at the boundary before hitting the pgx repository API.

## [T5] PlatformProvider interface + Twitch implementation
- `platform.PlatformProvider` interface lives in `apps/api-gql/internal/platform/provider.go`; Twitch impl in `internal/platform/twitch/provider.go`.
- Twitch provider creates a plain `helix.Client` (no bus/app-token) using only ClientID, ClientSecret, RedirectURI — sufficient for OAuth exchange flows.
- `GetAuthURL` ignores `codeChallenge` (Twitch doesn't use PKCE); scopes are embedded in the provider.
- `ExchangeCode` ignores `codeVerifier` param — pass `_` to satisfy interface.
- Registered in `cmd/main.go` via `fx.Annotate(twitchplatform.New, fx.As(new(platform.PlatformProvider)))`.
- `TwitchMockEnabled` → `apiBaseURL` override pattern mirrors libs/twitch/twitch.go.

- Kick provider uses exported PKCE helpers so the authorize route can store a session-only verifier and return a Kick OAuth URL with S256 challenge params.
- Kick config now supports KICK_CLIENT_ID, KICK_CLIENT_SECRET, and optional KICK_REDIRECT_URL; default callback falls back to {SITE_BASE_URL}/login/kick.
- Kick /auth/kick/code needs session code verifier retrieval via Auth.Get plus explicit Auth.Commit after multiple Put calls to avoid per-key helper proliferation.
- Kick bus-core topics live in a dedicated `libs/bus-core/kick` package with raw event structs and subject constants; `go build ./libs/bus-core/...` passes after adding them.
