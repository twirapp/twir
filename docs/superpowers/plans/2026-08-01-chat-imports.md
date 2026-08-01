# Chat Provider Imports and Donations Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Deliver secure Nightbot and StreamElements command/timer imports while making one StreamElements or Streamlabs OAuth connection power realtime donation integrations.

**Architecture:** Provider-aware Go libraries perform OAuth exchange, authenticated API calls, refresh, persistence, and one retry for GraphQL/import consumers. Focused TypeScript provider clients perform the equivalent lifecycle for `apps/integrations`, which owns Socket.IO connections; both runtimes serialize rotating-token refresh with the same Redis lock key and persist to the same channel integration record. A shared Go importer applies Twir business rules, while the Nuxt dashboard presents connected provider state and accurately keeps Streamlabs import unavailable.

**Tech Stack:** Go 1.24, gqlgen, pgx/PostgreSQL, go-redis, Fx, Bun, Bun Redis/SQL, Socket.IO client, Nuxt 4/Vue 3, TypeScript, Urql, shadcn-vue/Reka UI, Vitest.

## Global Constraints

- Imports are additive: never update or delete an existing Twir command or timer.
- StreamElements imports all supported custom commands and timers; there is no per-item selection UI.
- A StreamElements connection enables imports and realtime tips; request exactly `channel:read`, `bot:read`, and `tips:read`, with no write scope.
- A Streamlabs connection enables realtime donations, but no import action; explain that its public API does not expose Cloudbot commands or timers.
- Do not add StreamElements or Streamlabs refresh handling to `apps/tokens`.
- Provider libraries own refresh, rotated-token persistence, in-memory replacement, and at most one authenticated retry.
- Go and Bun refreshers use lock key `twir:integration-token-refresh:<provider>:<channelID>`, a random ownership value, a 30-second lease, and compare-and-delete release.
- After acquiring the refresh lock, reread persisted tokens and skip provider refresh if another process already changed the access token.
- OAuth state is random, atomically single-use, session-bound, dashboard-bound, provider-bound, user-bound, and expires after 15 minutes.
- Backend import reports expose stable reason codes; UI owns translations.
- Frontend uses Nuxt auto-imports, `<Icon />`, semantic Tailwind tokens, existing shadcn-vue components, and `DialogOrSheet`.
- Add dashboard copy to `de`, `en`, `es`, `ja`, `pt`, `ru`, `sk`, and `uk` with identical key shapes.
- Follow red-green-refactor for every production behavior change; generated GraphQL files and migrations are generated-artifact exceptions.

---

### Task 1: Register StreamElements as a persistent integration — complete

**Commit:** `be0930d4c feat(import): register StreamElements`

**Files:**
- Created: `libs/migrations/postgres/20260801000000_streamelements_integration.sql`
- Created: `libs/repositories/integrations/model/model_test.go`
- Modified: `libs/repositories/integrations/model/model.go`
- Modified: `libs/migrations/seeds/integrations.go`

**Interfaces produced:** `integrationsmodel.ServiceStreamElements == "STREAMELEMENTS"` and a seeded `integrations` row.

- [x] Add the enum value, seed, PostgreSQL migration, and repository test.
- [x] Run `cd libs/repositories && go test ./integrations/...`.
- [x] Pass implementation and quality review.

### Task 2: Bind OAuth callbacks to atomically single-use state — complete

**Commits:** `cbc7c48c8 fix(auth): bind integration OAuth state`, `1bcbd24a0 fix(auth): atomically claim OAuth state`

**Files:**
- Created: `apps/api-gql/internal/auth/oauth_attempt_redis_store.go`
- Modified: `apps/api-gql/internal/auth/sessions_user.go`
- Modified: `apps/api-gql/internal/auth/sessions_user_test.go`

**Interfaces produced:**
- `CreateIntegrationOAuthAttempt(ctx, service, channelID, initiatorUserID) (string, error)`.
- `ConsumeIntegrationOAuthAttempt(ctx, state, service, channelID, initiatorUserID, now) error`.

- [x] Bind state to provider, channel, initiator, expiry, and redirect target.
- [x] Atomically claim state through Redis before validation and release only mismatched attempts.
- [x] Cover success, mismatch, expiry, replay, and concurrent replay.
- [x] Pass implementation, fix-round, and quality review.

### Task 3: Build the provider-independent additive importer

**Files:**
- Create: `apps/api-gql/internal/services/importer/types.go`
- Create: `apps/api-gql/internal/services/importer/service.go`
- Create: `apps/api-gql/internal/services/importer/service_test.go`
- Modify: `apps/api-gql/internal/services/commands/create.go`
- Modify: `apps/api-gql/internal/services/timers/create.go`
- Modify: `apps/api-gql/cmd/main.go`

**Interfaces:**
- Produces `FailureReason`, `Failure`, `Report`, `Command`, `Timer`, `RoleRequirement`.
- Produces `ImportCommands(ctx, channelID, actorID string, input []Command) (Report, error)`.
- Produces `ImportTimers(ctx, channelID, actorID string, input []Timer) (Report, error)`.
- Consumes existing command/timer create services so plan limits, audit, cache, and timer lifecycle stay centralized.

- [ ] **Step 1: Write failing importer classification tests**

Use fake command/timer creators and role lookup. Cover success, duplicate, plan limit, partial success, unsupported normalized records, and immediate return for infrastructure errors. Assert exact inputs and stable reasons:

```go
const (
	FailureDuplicate            FailureReason = "DUPLICATE"
	FailurePlanLimit            FailureReason = "PLAN_LIMIT"
	FailureUnsupportedRole      FailureReason = "UNSUPPORTED_ROLE"
	FailureUnsupportedResponse  FailureReason = "UNSUPPORTED_RESPONSE_TYPE"
	FailureIncompatibleInterval FailureReason = "INCOMPATIBLE_INTERVALS"
	FailureInvalidRecord        FailureReason = "INVALID_RECORD"
)
```

- [ ] **Step 2: Verify RED**

Run: `cd apps/api-gql && go test ./internal/services/importer`

Expected: package/symbol compile failure.

- [ ] **Step 3: Add machine-readable details to expected create errors**

Set `AppError.Details["reason"]` to `PLAN_LIMIT` for command/timer plan exhaustion and `DUPLICATE` for name/alias conflicts. Preserve existing codes and user-visible messages.

- [ ] **Step 4: Implement normalized types and item-by-item import**

Use these exact normalized fields:

```go
type Command struct {
	Name, Response string
	Enabled, Visible, IsReply bool
	Aliases []string
	Cooldown int
	Role RoleRequirement
	OnlineOnly, OfflineOnly bool
}

type Timer struct {
	Name, Message string
	Enabled, OnlineEnabled, OfflineEnabled bool
	TimeInterval, MessageInterval int
}
```

Resolve role UUIDs once per command batch. Append expected item failures, retain successful siblings, and return operation-level errors for repository/cache/lifecycle failures.

- [ ] **Step 5: Wire `importer.New` into Fx and verify GREEN**

Run: `cd apps/api-gql && go test ./internal/services/importer ./internal/services/commands ./internal/services/timers`

Expected: PASS.

- [ ] **Step 6: Commit**

```bash
git add apps/api-gql/internal/services/importer apps/api-gql/internal/services/commands/create.go apps/api-gql/internal/services/timers/create.go apps/api-gql/cmd/main.go
git commit -m "refactor(import): share provider importer"
```

### Task 4: Move Nightbot imports onto the shared importer

**Files:**
- Create: `apps/api-gql/internal/services/nightbot_integration/mapper.go`
- Create: `apps/api-gql/internal/services/nightbot_integration/mapper_test.go`
- Create: `apps/api-gql/internal/services/nightbot_integration/client_test.go`
- Modify: `apps/api-gql/internal/services/nightbot_integration/nightbot_integration.go`
- Modify: `apps/tokens/internal/bus_listener/bus_listener.go`
- Modify: `apps/tokens/internal/bus_listener/channel_integrations_test.go`

**Interfaces:**
- Produces `NormalizeCommands(...) ([]importer.Command, []importer.Failure)`.
- Produces `NormalizeTimers(...) ([]importer.Timer, []importer.Failure)`.
- Changes `GetAuthLink(ctx, state)` to encode the supplied state.
- Nightbot refresh remains in its existing path; the `apps/tokens` prohibition applies only to StreamElements and Streamlabs.

- [ ] **Step 1: Write failing mapper tests**

Assert prefix removal, lower-casing, aliases, cooldown, reply behavior, owner/moderator/subscriber/VIP mapping, explicit admin/regular failures, cron-to-minutes conversion, and invalid cron failures.

- [ ] **Step 2: Verify RED**

Run: `cd apps/api-gql && go test ./internal/services/nightbot_integration -run Normalize`

- [ ] **Step 3: Normalize provider DTOs and delegate persistence**

Fetch each provider collection once, convert it without DB writes, invoke the shared importer, and prepend normalization failures to the returned report.

- [ ] **Step 4: Add state and refresh redirect regressions**

Assert the authorization URL contains exactly the provided state. In the tokens test assert:

```go
require.Equal(t,
	"https://twir.test/dashboard/integrations/callbacks/nightbot",
	r.Form.Get("redirect_uri"),
)
```

- [ ] **Step 5: Add `redirect_uri` to Nightbot refresh and verify GREEN**

Run: `cd apps/api-gql && go test ./internal/services/nightbot_integration`

Run: `cd apps/tokens && go test ./internal/bus_listener -run Nightbot`

Expected: PASS.

- [ ] **Step 6: Commit**

```bash
git add apps/api-gql/internal/services/nightbot_integration apps/tokens/internal/bus_listener/bus_listener.go apps/tokens/internal/bus_listener/channel_integrations_test.go
git commit -m "fix(import): harden Nightbot imports"
```

### Task 5: Add reusable Go OAuth refresh libraries

**Files:**
- Create: `libs/integrations/oauthlock/redis.go`
- Create: `libs/integrations/oauthlock/redis_test.go`
- Create: `libs/integrations/streamlabs/streamlabs.go`
- Create: `libs/integrations/streamlabs/responses.go`
- Create: `libs/integrations/streamlabs/streamlabs_test.go`
- Create: `libs/integrations/streamelements/streamelements_test.go`
- Modify: `libs/integrations/streamelements/streamelements.go`
- Modify: `libs/integrations/streamelements/responses.go`
- Modify: `libs/integrations/go.mod`
- Modify: `libs/integrations/go.sum`

**Interfaces:**
- Produces `oauthlock.Locker.WithLock(ctx context.Context, key string, fn func(context.Context) error) error`.
- Produces provider `TokenStore` interfaces with `GetTokens` and `UpdateTokens` by channel ID.
- Produces `streamelements.NewAuthorized(...)` and `streamlabs.NewAuthorized(...)` repository-aware clients.
- Produces static OAuth clients for auth-link/code-exchange operations that do not require persisted channel credentials.

- [ ] **Step 1: Write failing Redis lock tests**

Use a fake Redis command interface. Assert `SET key owner NX PX 30000`, retry-until-context behavior, callback exclusion, and compare/delete release:

```lua
if redis.call("GET", KEYS[1]) == ARGV[1] then
	return redis.call("DEL", KEYS[1])
end
return 0
```

- [ ] **Step 2: Implement the lock and verify GREEN**

The full key passed by clients is `twir:integration-token-refresh:<provider>:<channelID>`. Generate a UUID owner per acquisition. Never delete a lock owned by another process.

Run: `cd libs/integrations && go test ./oauthlock`

- [ ] **Step 3: Write failing provider contract tests**

With `httptest.Server`, fake stores, and fake locker, cover for each provider:

- authorization query and form-encoded code exchange;
- bearer-authenticated profile/API calls;
- sanitized non-2xx errors and malformed JSON;
- first `401`, locked token reread, refresh, atomic access/refresh persistence, in-memory replacement, and one successful retry;
- another process changed access token while waiting, so provider refresh is skipped and the request retries with the reread token;
- second `401` returns `ErrUnauthorized` without looping.

StreamElements auth scopes must equal `channel:read bot:read tips:read`. Refresh posts to `/oauth2/token` with `grant_type=refresh_token`, client credentials, and refresh token. Streamlabs refresh posts form data to `/api/v2.0/token` with `grant_type=refresh_token`, client credentials, refresh token, and `redirect_uri`.

- [ ] **Step 4: Implement common authenticated request flow inside each provider package**

Use a 15-second `http.Client` timeout. The retry algorithm is exact:

```go
response, err := request(currentAccessToken)
if !errors.Is(err, ErrUnauthorized) { return response, err }
err = locker.WithLock(ctx, lockKey(provider, channelID), func(ctx context.Context) error {
	fresh, err := store.GetTokens(ctx, channelID)
	if err != nil { return err }
	if fresh.AccessToken != currentAccessToken { current = fresh; return nil }
	rotated, err := refresh(ctx, fresh.RefreshToken)
	if err != nil { return err }
	if rotated.RefreshToken == "" { rotated.RefreshToken = fresh.RefreshToken }
	if err := store.UpdateTokens(ctx, channelID, rotated); err != nil { return err }
	current = rotated
	return nil
})
if err != nil { return zero, err }
return request(current.AccessToken)
```

- [ ] **Step 5: Verify provider packages**

Run: `cd libs/integrations && go test ./oauthlock ./streamelements ./streamlabs`

Expected: PASS.

- [ ] **Step 6: Commit**

```bash
git add libs/integrations/oauthlock libs/integrations/streamelements libs/integrations/streamlabs libs/integrations/go.mod libs/integrations/go.sum
git commit -m "feat(integrations): own provider token refresh"
```

### Task 6: Implement persistent StreamElements lifecycle and imports

**Files:**
- Create: `apps/api-gql/internal/services/streamelements/mapper.go`
- Create: `apps/api-gql/internal/services/streamelements/mapper_test.go`
- Create: `apps/api-gql/internal/services/streamelements/service_test.go`
- Modify: `apps/api-gql/internal/services/streamelements/dto.go`
- Modify: `apps/api-gql/internal/services/streamelements/streamelements.go`
- Modify: `libs/repositories/channels_integrations/repository.go`
- Modify: `libs/repositories/channels_integrations/datasource/postgres/pgx.go`
- Modify: `libs/bus-core/integrations/integrations.go`
- Modify: `libs/bus-core/src/integrations/integrations.ts`

**Interfaces:**
- Produces service methods `GetData`, `GetAuthLink(state)`, `PostCode`, `Logout`, `ImportCommands`, and `ImportTimers`.
- Produces repository token-store adapter methods without routing through `apps/tokens`.
- Produces `integrations.StreamElements` and `IntegrationService.STREAMELEMENTS` bus values.

- [ ] **Step 1: Write failing StreamElements mapper tests**

Assert access levels `100/250/300/400/500/1000+`, aliases, cooldown, visibility, reply/say/action handling, unsupported response types, online/offline command flags, one-mode timers, equal dual-mode timers, and incompatible intervals.

- [ ] **Step 2: Implement normalization and verify GREEN**

Run: `cd apps/api-gql && go test ./internal/services/streamelements -run Normalize`

- [ ] **Step 3: Write failing lifecycle tests**

Use fake integration repositories, provider client factory, importer, and bus. Assert:

```go
stateURL, _ := service.GetAuthLink(ctx, "oauth-state")
require.Equal(t, "oauth-state", mustQuery(stateURL, "state"))
```

Also cover create/update after code exchange, username/avatar persistence, enabled state, Add bus publication, GetData nil behavior, Logout disable/Remove publication, and command/timer report composition using the repository-aware client.

- [ ] **Step 4: Add explicit token clearing semantics to the generic repository**

Extend `UpdateInput` with booleans `ClearAccessToken` and `ClearRefreshToken`; reject a request that both sets and clears one field. Generate SQL `SET "accessToken" = NULL` / `SET "refreshToken" = NULL` only when the clear flag is true. Test set, clear, and unrelated-field preservation.

- [ ] **Step 5: Implement persistent service lifecycle**

Replace `ExchangeDataByCode`. Create/update the `channels_integrations` row, store profile data, build `NewAuthorized` with repository and Redis locker, publish Add/Remove, and return configuration errors instead of `(nil, nil)`.

- [ ] **Step 6: Verify GREEN**

Run: `cd libs/repositories && go test ./channels_integrations/...`

Run: `cd apps/api-gql && go test ./internal/services/streamelements`

Run: `cd libs/bus-core && go test ./integrations/...`

Expected: PASS.

- [ ] **Step 7: Commit**

```bash
git add apps/api-gql/internal/services/streamelements libs/repositories/channels_integrations libs/bus-core
git commit -m "feat(import): persist StreamElements connection"
```

### Task 7: Move Streamlabs OAuth HTTP behavior into its Go provider library

**Files:**
- Create: `apps/api-gql/internal/services/streamlabs_integration/service_test.go`
- Modify: `apps/api-gql/internal/services/streamlabs_integration/service.go`
- Modify: `libs/repositories/streamlabs_integration/repository.go`
- Modify: `libs/repositories/streamlabs_integration/datasource/postgres/pgx.go`

**Interfaces:**
- Service consumes `libs/integrations/streamlabs` for auth URL, code exchange, profile, and refresh-aware authenticated requests.
- `GetAuthLink(ctx, state)` encodes state exactly once.
- Existing `channels_integrations_streamlabs` records and Add/Remove bus behavior remain compatible.

- [ ] **Step 1: Write failing service tests**

Assert state encoding, configuration failure, code exchange and profile mapping through a fake client factory, create/update behavior, Add publication, GetIntegrationData, Logout, and Remove publication. Assert provider response bodies are absent from returned errors.

- [ ] **Step 2: Add repository token-store adapter behavior**

Expose read/update methods needed by the provider client while preserving the existing repository API. A refresh update writes access and refresh tokens together in one SQL statement.

- [ ] **Step 3: Replace inline HTTP with the library**

Delete `streamlabsTokensResponse`, `streamlabsProfileResponse`, and `getProfileData` from the service. Keep GraphQL-facing mapping and bus lifecycle in the service.

- [ ] **Step 4: Verify GREEN**

Run: `cd libs/repositories && go test ./streamlabs_integration/...`

Run: `cd apps/api-gql && go test ./internal/services/streamlabs_integration`

Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add apps/api-gql/internal/services/streamlabs_integration libs/repositories/streamlabs_integration
git commit -m "refactor(streamlabs): use provider OAuth client"
```

### Task 8: Expose the secure GraphQL import and connection contract

**Files:**
- Modify: `apps/api-gql/internal/delivery/gql/schema/import.graphql`
- Modify: `apps/api-gql/internal/delivery/gql/schema/integrations/integrations-streamlabs.graphql`
- Modify: `apps/api-gql/internal/delivery/gql/resolvers/import.resolver.go`
- Modify: `apps/api-gql/internal/delivery/gql/resolvers/integrations-streamlabs.resolver.go`
- Modify: `apps/api-gql/internal/delivery/gql/resolvers/resolver.go`
- Create: `apps/api-gql/internal/delivery/gql/resolvers/import_test.go`
- Regenerate: `apps/api-gql/internal/delivery/gql/gqlmodel/models_gen.go`
- Regenerate: `apps/api-gql/internal/delivery/gql/generated.go`

**Interfaces:**
- Produces shared `ImportFailureReason`, `ImportFailure`, and `ImportReport` GraphQL types.
- Produces StreamElements data/auth/post-code/logout/import operations with dashboard permissions.
- Changes Nightbot and Streamlabs callback inputs to `{ code, state }`.
- Removes the raw `streamelementsExchangeDataByCode` DTO surface.

- [ ] **Step 1: Write failing resolver tests**

Cover state creation for the selected dashboard/user, state consumption before code exchange, replay rejection, wrong-provider rejection, permission boundaries, stable report conversion, and `gqlerrors.HandleError` for operation failures.

- [ ] **Step 2: Define the symmetric schema**

Use this shared result shape for both providers and resource types:

```graphql
enum ImportFailureReason {
  DUPLICATE
  PLAN_LIMIT
  UNSUPPORTED_ROLE
  UNSUPPORTED_RESPONSE_TYPE
  INCOMPATIBLE_INTERVALS
  INVALID_RECORD
}

type ImportFailure { name: String!, reason: ImportFailureReason! }
type ImportReport { importedCount: Int!, failedCount: Int!, failures: [ImportFailure!]! }
input IntegrationOAuthCodeInput { code: String!, state: String! }
```

- [ ] **Step 3: Regenerate gqlgen**

Run from root: `bun cli build gql`.

Re-read generated resolver signatures before implementation; never hand-edit generated execution/model files.

- [ ] **Step 4: Implement resolvers and consume state before provider calls**

Resolve authenticated user and selected dashboard for every start/finish. Create state in auth-link resolvers. Consume matching state before `PostCode`. Map service reports without provider-specific failure lists.

- [ ] **Step 5: Verify GREEN**

Run: `cd apps/api-gql && go test ./internal/delivery/gql/resolvers ./internal/services/nightbot_integration ./internal/services/streamelements ./internal/services/streamlabs_integration`

Expected: PASS.

- [ ] **Step 6: Commit**

```bash
git add apps/api-gql/internal/delivery/gql
git commit -m "feat(import): expose secure provider API"
```

### Task 9: Add Bun token locking, persistence, and donation deduplication

**Files:**
- Create: `apps/integrations/src/libs/oauth-lock.ts`
- Create: `apps/integrations/src/libs/oauth-lock.test.ts`
- Create: `apps/integrations/src/libs/provider-token-store.ts`
- Create: `apps/integrations/src/libs/provider-token-store.test.ts`
- Create: `apps/integrations/src/libs/donation-dedupe.ts`
- Create: `apps/integrations/src/libs/donation-dedupe.test.ts`
- Modify: `apps/integrations/src/libs/db.ts`

**Interfaces:**
- Produces `withOAuthRefreshLock(provider, channelID, callback)` with the same key/lease/release contract as Go.
- Produces StreamElements and Streamlabs `getTokens`/`updateTokens` stores over their existing tables.
- Produces `claimDonation(provider, eventID): Promise<boolean>` using a 24-hour Redis NX/EX key.

- [ ] **Step 1: Write failing lock and dedupe tests**

Inject fake Redis clients. Assert raw commands:

```ts
await redis.send('SET', [key, owner, 'NX', 'PX', '30000'])
await redis.send('EVAL', [releaseScript, '1', key, owner])
await redis.send('SET', [`twir:donation:${provider}:${eventId}`, '1', 'NX', 'EX', '86400'])
```

`claimDonation` returns `true` only when SET returns `OK`; an empty event ID returns `true` without touching Redis.

- [ ] **Step 2: Write failing token-store tests**

Inject a query function. Assert StreamElements reads `channels_integrations` joined to `integrations.service = 'STREAMELEMENTS'`; Streamlabs reads `channels_integrations_streamlabs`; each update changes access and refresh token atomically and scopes by channel ID.

- [ ] **Step 3: Implement focused libraries**

Use Bun `RedisClient.send` for lock/dedupe. Keep SQL in the DB/token-store boundary. Never log access tokens, refresh tokens, lock owners, or provider bodies.

- [ ] **Step 4: Verify GREEN**

Run: `cd apps/integrations && bun test src/libs/oauth-lock.test.ts src/libs/provider-token-store.test.ts src/libs/donation-dedupe.test.ts`

Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add apps/integrations/src/libs
git commit -m "feat(integrations): add token lifecycle primitives"
```

### Task 10: Listen for StreamElements tips in `apps/integrations`

**Files:**
- Create: `apps/integrations/src/libs/streamelements-client.ts`
- Create: `apps/integrations/src/libs/streamelements-client.test.ts`
- Create: `apps/integrations/src/services/streamElements.ts`
- Create: `apps/integrations/src/services/streamElements.test.ts`
- Create: `apps/integrations/src/store/streamelements.ts`
- Create: `apps/integrations/src/store/streamelements.test.ts`
- Modify: `apps/integrations/src/libs/db.ts`
- Modify: `apps/integrations/src/index.ts`

**Interfaces:**
- Produces `getStreamElementsIntegrations()` returning enabled generic channel integrations.
- Produces a refresh-aware TypeScript client with `refresh()` and current credentials.
- Produces one `StreamElementsConnection` per channel and lifecycle `addIntegration`/`removeIntegration`.

- [ ] **Step 1: Write failing OAuth refresh client tests**

With injected fetch, token store, and lock, assert form-encoded refresh, token reread after lock, skip-refresh when another process changed access token, rotated-token fallback when `refresh_token` is omitted, atomic persistence, and sanitized failures.

- [ ] **Step 2: Implement the TypeScript provider client**

Use the exact refresh body:

```ts
new URLSearchParams({
	grant_type: 'refresh_token',
	client_id: config.STREAMELEMENTS_CLIENT_ID,
	client_secret: config.STREAMELEMENTS_CLIENT_SECRET,
	refresh_token: current.refreshToken,
})
```

- [ ] **Step 3: Write failing socket tests**

Inject a Socket.IO factory. Assert URL `https://realtime.streamelements.com`, WebSocket-only transport, connect-time `authenticate` payload `{ method: 'oauth2', token }`, `event.type === 'tip'` normalization, ignored non-tip events, stable-ID dedupe, and pass-through for missing IDs.

- [ ] **Step 4: Implement socket lifecycle and bounded reconnect**

On `unauthorized`, refresh once and rebuild the socket. On network disconnect, reconnect with jittered delays capped at 30 seconds without refreshing. A generation counter or abort signal must prevent stale reconnect callbacks after removal and guarantee one socket per channel.

- [ ] **Step 5: Wire startup and bus lifecycle**

Load all enabled StreamElements records at startup. Handle `IntegrationService.STREAMELEMENTS` Add by integration ID and Remove by channel ID. Await `addIntegration` calls and close all live sockets on SIGTERM/SIGINT.

- [ ] **Step 6: Verify GREEN**

Run: `cd apps/integrations && bun test src/libs/streamelements-client.test.ts src/services/streamElements.test.ts src/store/streamelements.test.ts`

Run: `cd apps/integrations && bun run prebuild`

Expected: PASS.

- [ ] **Step 7: Commit**

```bash
git add apps/integrations/src
git commit -m "feat(integrations): listen for StreamElements tips"
```

### Task 11: Refactor Streamlabs socket lifecycle onto provider libraries

**Files:**
- Create: `apps/integrations/src/libs/streamlabs-client.ts`
- Create: `apps/integrations/src/libs/streamlabs-client.test.ts`
- Create: `apps/integrations/src/services/streamLabs.test.ts`
- Create: `apps/integrations/src/store/streamlabs.test.ts`
- Modify: `apps/integrations/src/services/streamLabs.ts`
- Modify: `apps/integrations/src/store/streamlabs.ts`
- Modify: `apps/integrations/src/index.ts`

**Interfaces:**
- Provider client obtains `/api/v2.0/socket/token`, refreshes only after `401`, persists rotations, and retries socket-token acquisition once.
- Store keeps exactly one connection per channel and never refreshes unconditionally during add.

- [ ] **Step 1: Write failing client tests**

Assert bearer socket-token request; successful first request performs no refresh; first `401` uses lock/reread/refresh/persist/retry; provider refresh uses form encoding and `redirect_uri`; second `401` stops. Sanitize all provider-body errors.

- [ ] **Step 2: Implement the TypeScript Streamlabs client**

Move every token request out of `store/streamlabs.ts`. Return a typed `{ socketToken: string }`; keep credentials private and replace them only after persistence succeeds.

- [ ] **Step 3: Write failing socket/store tests**

Assert donation normalization, event ID passed to dedupe, multiple messages handled independently, ignored event types, replacement closes the old connection, removal cancels reconnect, and auth failure does not disable the DB integration.

- [ ] **Step 4: Implement lifecycle and reconnect**

Use the same capped reconnect controller as StreamElements. Do not write `enabled=false` for transient provider failures. Log only provider/channel identifiers and typed error categories.

- [ ] **Step 5: Verify GREEN**

Run: `cd apps/integrations && bun test src/libs/streamlabs-client.test.ts src/services/streamLabs.test.ts src/store/streamlabs.test.ts`

Run: `cd apps/integrations && bun run prebuild`

Expected: PASS.

- [ ] **Step 6: Commit**

```bash
git add apps/integrations/src/libs/streamlabs-client.ts apps/integrations/src/libs/streamlabs-client.test.ts apps/integrations/src/services/streamLabs.ts apps/integrations/src/services/streamLabs.test.ts apps/integrations/src/store/streamlabs.ts apps/integrations/src/store/streamlabs.test.ts apps/integrations/src/index.ts
git commit -m "refactor(integrations): harden Streamlabs socket"
```

### Task 12: Build the localized Nuxt import experience

**Files:**
- Create: `web/layers/dashboard/features/import/components/import-provider-card.vue`
- Create: `web/layers/dashboard/features/import/components/import-settings.vue`
- Create: `web/layers/dashboard/features/import/components/import-result.vue`
- Create: `web/layers/dashboard/features/import/components/import-provider-card.spec.ts`
- Create: `web/layers/dashboard/features/import/components/import-settings.spec.ts`
- Create: `web/layers/dashboard/features/import/streamelements/composables/use-streamelements-integration.ts`
- Create: `web/layers/dashboard/features/import/streamelements/streamelements-import.vue`
- Create: `web/layers/dashboard/features/import/streamelements/streamelements-callback.vue`
- Create: `web/layers/dashboard/features/import/streamelements/streamelements-import.spec.ts`
- Create: `web/layers/dashboard/features/import/streamlabs/streamlabs-import.vue`
- Create: `web/layers/dashboard/features/import/streamlabs/streamlabs-import.spec.ts`
- Create: `web/layers/dashboard/pages/dashboard/integrations/callbacks/streamelements.vue`
- Modify: `web/layers/dashboard/features/import/nightbot/composables/use-nightbot-integration.ts`
- Modify: `web/layers/dashboard/features/import/nightbot/nightbot-import.vue`
- Modify: `web/layers/dashboard/features/import/nightbot/nightbot-callback.vue`
- Create: `web/layers/dashboard/features/import/nightbot/nightbot-import.spec.ts`
- Modify: `web/layers/dashboard/pages/dashboard/import.vue`
- Delete: `web/layers/dashboard/pages/import/streamelements.vue`
- Delete: `web/layers/dashboard/pages/import/streamlabs.vue`
- Modify: `web/layers/dashboard/api/integrations/integrations-page.ts`
- Modify: `web/layers/dashboard/api/integrations/integrations.ts`
- Modify: `web/layers/dashboard/locales/de.json`
- Modify: `web/layers/dashboard/locales/en.json`
- Modify: `web/layers/dashboard/locales/es.json`
- Modify: `web/layers/dashboard/locales/ja.json`
- Modify: `web/layers/dashboard/locales/pt.json`
- Modify: `web/layers/dashboard/locales/ru.json`
- Modify: `web/layers/dashboard/locales/sk.json`
- Modify: `web/layers/dashboard/locales/uk.json`
- Regenerate: `web/gql/gql.ts`
- Regenerate: `web/gql/graphql.ts`

**Interfaces:**
- `ImportProviderCard` renders icon, description, connection/account state, auth/logout, availability, and settings slot.
- `ImportSettings` emits independent command/timer imports and renders `ImportReport`.
- StreamElements connection state enables both donation explanatory copy and import controls.
- Streamlabs connection state enables donation explanatory copy but never an import action.

- [ ] **Step 1: Load required frontend skills and inspect installed component APIs**

Read the `nuxt`, `shadcn-vue`, and `reka-ui` skills. Inspect existing `DialogOrSheet`, Card, Button, Badge, Alert, and integration-page query usage. Do not add or overwrite registry components.

- [ ] **Step 2: Write failing shared component tests**

Cover disconnected/connected actions, loading, permission-disabled controls, partial results, translated failure reasons, and an unavailable import area that may still show provider login/logout.

- [ ] **Step 3: Implement shared components and verify GREEN**

Run: `cd web && bun run test -- layers/dashboard/features/import/components`

- [ ] **Step 4: Write failing provider feature tests**

For Nightbot and StreamElements assert auth URL usage, callback `{ code, state }`, no broadcast/close on mutation error, broadcast/close on success, command/timer separation, and report rendering. For Streamlabs assert connected/disconnected donation state, login/logout availability, unavailable badge, explanatory copy, and no import/settings button.

- [ ] **Step 5: Implement provider components and unified data access**

Add StreamElements and Streamlabs fields to `integrations-page.ts`; do not introduce per-card initial queries. Keep mutations separate and invalidate `integrationsPageCacheKey`. Use `DialogOrSheet`, Nuxt `<Icon />`, semantic tokens, and no custom CSS.

- [ ] **Step 6: Add locale parity and regenerate GraphQL**

Add one `imports` namespace containing provider descriptions, donation-enabled text, unavailable text, actions, counts, callback errors, and all failure reasons to all eight locales.

Run: `cd web && bun run graphql-codegen`.

- [ ] **Step 7: Run scoped tests and typecheck**

Run: `cd web && bun run test -- layers/dashboard/features/import`

Run: `cd web && bun run nuxt-prepare && bunx --bun vue-tsc --noEmit -p .nuxt/tsconfig.json`

Expected: PASS.

- [ ] **Step 8: Commit**

```bash
git add web/layers/dashboard/features/import web/layers/dashboard/pages/dashboard/import.vue web/layers/dashboard/pages/dashboard/integrations/callbacks/streamelements.vue web/layers/dashboard/pages/import web/layers/dashboard/api/integrations web/layers/dashboard/locales web/gql
git commit -m "feat(web): connect imports and donations"
```

### Task 13: Verify the complete feature

**Files:**
- Modify only if implementation changed documented behavior: `docs/superpowers/specs/2026-08-01-chat-imports-design.md`
- Modify only if paths/commands changed: `docs/superpowers/plans/2026-08-01-chat-imports.md`

**Interfaces:** Validates every success criterion from the approved design.

- [ ] **Step 1: Run Go tests**

```bash
cd libs/repositories && go test ./integrations/... ./channels_integrations/... ./streamlabs_integration/...
cd ../integrations && go test ./oauthlock ./streamelements ./streamlabs
cd ../../apps/tokens && go test ./internal/bus_listener
cd ../api-gql && go test ./internal/auth ./internal/services/importer ./internal/services/nightbot_integration ./internal/services/streamelements ./internal/services/streamlabs_integration ./internal/delivery/gql/resolvers
```

Expected: all commands exit 0.

- [ ] **Step 2: Run integrations runtime verification**

```bash
cd apps/integrations
bun test
bun run prebuild
bun run build
```

Expected: all commands exit 0; socket tests leave no open handles.

- [ ] **Step 3: Run frontend verification**

```bash
cd web
bun run nuxt-prepare
bun run graphql-codegen
bun run test -- layers/dashboard/features/import
bunx --bun vue-tsc --noEmit -p .nuxt/tsconfig.json
```

Expected: all commands exit 0.

- [ ] **Step 4: Run formatting, lint, and builds**

Run `gofmt -w` on touched Go files, then from the repository root:

```bash
bun lint
git diff --check
bun cli b app api-gql
bun cli b app integrations
```

Expected: no failures attributable to the change.

- [ ] **Step 5: Perform browser and runtime QA**

Open `${SITE_BASE_URL}/dashboard/import`, falling back to `http://localhost:3005/dashboard/import`. Verify desktop/mobile connected and disconnected states, both import reports, permission-disabled controls, callback errors, StreamElements donation-enabled copy, and Streamlabs donation-enabled/import-unavailable copy. With provider test events, verify one stored/published donation per stable provider event ID and no socket survives logout.

- [ ] **Step 6: Review the final diff against security and scope**

Check state consumption before exchange, exact scopes, no StreamElements/Streamlabs switch in `apps/tokens`, one refresh retry, compare-and-delete locks, rotated-token persistence, no credentials/provider bodies in logs, one socket per channel, locale parity, no `lucide-vue-next`, no handwritten generated GraphQL code, and no placeholder import pages.

- [ ] **Step 7: Commit documentation alignment only if required**

```bash
git add docs/superpowers/specs/2026-08-01-chat-imports-design.md docs/superpowers/plans/2026-08-01-chat-imports.md
git commit -m "docs: align import integration behavior"
```
