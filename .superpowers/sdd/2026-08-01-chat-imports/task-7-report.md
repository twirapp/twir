# Task 7 Report: Streamlabs Go provider-library integration

## Status

Complete. The API service delegates authorization URL generation, code exchange, and profile
fetching to `libs/integrations/streamlabs`. The existing repository now implements the provider
`TokenStore`, and logout serializes with provider refresh before deleting all persisted
credentials.

## Implementation

- Added a provider-client factory seam to the Streamlabs API service.
- Required exactly one nonblank OAuth state and delegated its encoding to the provider once.
- Removed service-local OAuth HTTP requests and token/profile response DTOs.
- Preserved create/update behavior and Streamlabs Add/Remove bus payloads.
- Added construction of refresh-aware authorized clients with the repository as `TokenStore` and
  the shared Redis OAuth locker.
- Exported `streamlabs.RefreshLockKey(channelID)` and made the provider use it internally.
- Serialized logout with that refresh key and deleted the integration row inside the lock, leaving
  no recoverable access token, refresh token, or profile data.
- Added enabled-only token reads and a single conditional SQL statement that updates access and
  refresh tokens together.
- Rejected disabled/missing rows and incomplete token pairs without writing.

The Postgres best-practices guidance was applied by keeping refresh persistence to one short,
conditional statement rather than holding a database transaction across provider work.

## TDD evidence

### RED: service

Command:

```text
cd apps/api-gql && go test ./internal/services/streamlabs_integration
```

Observed before production changes:

```text
Go test: 0 passed, 1 failed in 1 packages
streamlabs_integration [build failed]
service_test.go: too many arguments in call to service.GetAuthLink
service_test.go: service.authorizedClient undefined
service_test.go: unknown field clientFactory in struct literal of type Service
service_test.go: unknown field locker in struct literal of type Service
service_test.go: unknown field events in struct literal of type Service
```

### RED: repository

Command:

```text
cd libs/repositories && go test ./streamlabs_integration/...
```

Observed before production changes:

```text
Go test: 0 passed, 1 failed in 3 packages
postgres [build failed]
pgx_test.go: unknown field tokenStoreDB in struct literal of type Pgx
pgx_test.go: repo.GetTokens undefined
pgx_test.go: repo.UpdateTokens undefined
```

### GREEN: required verification

Command:

```text
cd libs/repositories && go test ./streamlabs_integration/...
```

Exact output:

```text
Go test: 5 passed in 3 packages
```

Command:

```text
cd apps/api-gql && go test ./internal/services/streamlabs_integration ./cmd
```

Exact output:

```text
Go test: 7 passed in 2 packages
```

### Additional verification

```text
cd libs/repositories && go test -race ./streamlabs_integration/...
Go test: 5 passed in 3 packages

cd libs/repositories && go vet ./streamlabs_integration/...
Go vet: No issues found

cd apps/api-gql && go test -race ./internal/services/streamlabs_integration ./cmd
Go test: 7 passed in 2 packages

cd apps/api-gql && go vet ./internal/services/streamlabs_integration ./cmd
Go vet: No issues found

cd libs/integrations && go test ./streamlabs
Go test: 15 passed in 1 packages

cd libs/integrations && go test -race ./streamlabs
Go test: 15 passed in 1 packages

cd libs/integrations && go vet ./streamlabs
Go vet: No issues found
```

`git diff --check` also completed without output.

## Files

- `apps/api-gql/internal/services/streamlabs_integration/service.go`
- `apps/api-gql/internal/services/streamlabs_integration/service_test.go`
- `libs/integrations/streamlabs/streamlabs.go`
- `libs/repositories/streamlabs_integration/repository.go`
- `libs/repositories/streamlabs_integration/datasource/postgres/pgx.go`
- `libs/repositories/streamlabs_integration/datasource/postgres/pgx_test.go`
- `.superpowers/sdd/2026-08-01-chat-imports/task-7-report.md`

## Self-review

- Confirmed the API service contains no direct Streamlabs HTTP request code.
- Confirmed provider errors report operation/status only and do not expose response bodies or
  credentials.
- Confirmed both token columns are updated in one statement with `enabled = TRUE` in the predicate.
- Confirmed logout and provider refresh use the identical provider-generated lock key.
- Confirmed deletion occurs while the refresh lock is held and Remove publication occurs only
  after the locked delete succeeds.
- Confirmed `apps/tokens` is untouched.

## Concerns / handoff

- Task 8 must update the GraphQL resolver to supply provider-bound OAuth state. The service keeps a
  variadic compatibility signature so the current generated resolver compiles, but a zero-state
  call intentionally returns `streamlabs OAuth state is required`.
- Standalone `GOWORK=off` verification cannot resolve the unpublished Task 5 Streamlabs provider
  package. All required monorepo workspace verification is green; module publication/version
  updates remain part of the repository's normal release process.
