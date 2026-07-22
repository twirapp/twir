# Task 8 Report: VK Configuration and Identity Adapter

## Status

DONE_WITH_CONCERNS

## Implementation Summary

- Added feature-gated `VK_VIDEO_*` configuration and sample environment entries. VK-specific credentials are validated only when `VK_VIDEO_ENABLED=true`; legacy `VK_CLIENT_ID`, `VK_CLIENT_SECRET`, and `VK_APP_ACCESS_TOKEN` remain separate.
- Added `libs/integrations/vk/id.go`, a narrow shared VK ID OAuth client used by API-GQL and the tokens service. It builds the confirmed PKCE authorization URL, posts form-encoded code/refresh/profile requests, requires `device_id` before token requests, sends only `service_token` as the confidential-app request secret, maps the documented profile fields, and exposes typed `ProviderError` values.
- Added `PlatformTokens.DeviceID`, typed exchange and refresh inputs, VK Video Live provider/registry wiring, and adapted Twitch/Kick to the typed inputs without changing their OAuth request behavior.
- Created additive nullable `tokens."deviceID"` storage. API-GQL encrypts a supplied device ID before create/update persistence; the repository treats it as opaque encrypted data; tokens decrypts it only for an expired VK refresh.
- Added the explicit VK refresh branch that rejects a missing persisted device ID before client construction/HTTP, preserves the old refresh token when VK omits a replacement, and persists rotated access/refresh values.
- Did not add callback/session capture, generic OAuth binding routes, VK linking, Video-specific scope/API calls, webhooks, GraphQL, or dashboard work.

## Files Changed

- Configuration and docs: `.env.example`, `libs/config/config.go`, `libs/config/config_test.go`.
- Shared integration: `libs/integrations/vk/id.go`, `libs/integrations/vk/id_test.go`.
- Token storage: `libs/migrations/postgres/20260722184627_add_tokens_device_id.sql`, `libs/repositories/tokens/{repository.go,repository_test.go,model/model.go,datasources/postgres/pgx.go}`.
- API-GQL: `apps/api-gql/internal/platform/{provider.go,registry.go,registry_test.go,kick/provider.go,twitch/provider.go,vkvideo/provider.go,vkvideo/provider_test.go}`, `apps/api-gql/cmd/main.go`, and the existing token persistence helper/caller in `internal/delivery/http/routes/auth`.
- Tokens: `apps/tokens/go.mod`, `apps/tokens/internal/bus_listener/{bus_listener.go,bus_listener_test.go}`.

## TDD Evidence

### RED

1. `go test ./libs/config -run 'TestNewWithEnvPath_(ValidatesVKVideoConfigurationOnlyWhenEnabled|LoadsVKVideoConfiguration)$' -count=1`
   - Failed as expected: enabled VK configuration accepted missing credentials and all seven `VKVideo*` fields were missing.
2. `go test ./libs/integrations/vk -run '^TestIDClient' -count=1`
   - Failed as expected before implementation: missing `IDClient`, typed inputs, `ErrDeviceIDRequired`, and `ProviderError`.
3. `go test ./apps/api-gql/internal/platform/... -run 'Test(Provider|Registry)' -count=1`
   - Failed as expected before implementation: missing `NewRegistry`, `ExchangeCodeInput`, `RefreshTokenInput`, and `vkvideo.Provider`.
4. `go test ./apps/api-gql/internal/platform -run '^Test(NewFeatureGatedRegistry|Registry)' -count=1`
   - Failed as expected before feature gating: missing `NewFeatureGatedRegistry`.
5. `go test ./libs/repositories/tokens/... ./apps/tokens/internal/bus_listener ./apps/api-gql/internal/delivery/http/routes/auth -run '(TestTokenInputsAndModelSupportNullableEncryptedDeviceID|TestUpsertPlatformUserToken_EncryptsDeviceID|TestRequestUserToken_VK)' -count=1`
   - Failed as expected before implementation: `DeviceID` fields and the VK token refresher did not exist.

### GREEN

1. `go test ./libs/config -count=1`
   - Passed.
2. `go test ./libs/integrations/vk/... ./libs/repositories/tokens/... ./apps/tokens/... ./apps/api-gql/internal/platform/...`
   - Passed, including all new shared-client, registry/provider, and VK refresh tests.
3. `go test ./libs/migrations/...`
   - Passed; no migration package tests exist.
4. `git diff --check`
   - Passed with no whitespace errors.

## Required Suite

Command:

```sh
go test ./libs/integrations/vk/... ./libs/repositories/tokens/... ./apps/tokens/... ./apps/api-gql/internal/platform/...
```

Result: passed. The packages with tests completed successfully; remaining listed packages report `[no test files]`.

## Migration

Created with the required CLI command:

```sh
bun cli m create --name add_tokens_device_id --db postgres --type sql
```

Generated path: `libs/migrations/postgres/20260722184627_add_tokens_device_id.sql`.

The up migration adds nullable `tokens."deviceID" TEXT`; it has no default and no index. The down migration drops that column.

## Database Integration Status

Not run by explicit approval: Docker/Postgres integration is unavailable. The migration was generated via the project CLI, reviewed for additive nullable/no-index behavior, and the repository/token behavior is covered with unit/fake-repository tests. No database migration execution or live PostgreSQL persistence test was attempted.

## Commits

- `a6b90fd25 feat(vk): add VK ID token adapter`

## Self-Review

- The new VK client does not call or import the legacy VK API client.
- Token/profile forms are context-bound, form-encoded, and do not log request credentials or user tokens.
- `service_token` is the only confidential-app secret sent by the shared client; `client_secret` is not sent.
- Profile conversion intentionally leaves `PlatformUser.Login` empty and composes only the documented first/last names into `DisplayName`.
- Existing Twitch/Kick code paths compile in the focused platform suite and consume the new typed inputs unchanged in behavior.
- Existing token rows remain compatible because `DeviceID` is nullable and updates leave it untouched when no new device ID is supplied.

## Concerns

- `go test ./apps/api-gql/internal/delivery/http/routes/auth` remains blocked by the known deferred Task 9 migration: `oauth-platform.go` still references removed `Channel` legacy fields and `GetChannelByConnectedUser`. This was present at the supplied base and is outside Task 8. The new device-ID encryption tests in that package therefore cannot execute until Task 9 repairs the existing package.
- `go test ./apps/api-gql/cmd` is also unavailable in this worktree because generated GraphQL `gqlmodel` and `graph` packages are absent. Registry feature gating is instead directly unit-tested in `internal/platform`.
- No Docker/Postgres integration verification was run, as noted above.
