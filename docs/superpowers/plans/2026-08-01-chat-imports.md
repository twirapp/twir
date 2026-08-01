# Chat Provider Imports Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Deliver secure, additive Nightbot and StreamElements command/timer imports and an honest Streamlabs unavailable state in the dashboard.

**Architecture:** Provider clients fetch and normalize external records into a shared backend importer, which applies Twir roles, limits, duplicate handling, auditing, caching, and timer lifecycle behavior through existing services. OAuth credentials persist in channel integrations and every provider callback consumes a single-use session-bound state. The Nuxt dashboard composes one shared import-card UI with provider-specific composables.

**Tech Stack:** Go 1.24, gqlgen, pgx/PostgreSQL, Fx, Bun, Nuxt 4/Vue 3, TypeScript, Urql, shadcn-vue/Reka UI, Vitest.

## Global Constraints

- Imports are additive: never update or delete an existing Twir command or timer.
- StreamElements imports all supported custom commands and timers; there is no per-item selection UI.
- Streamlabs exposes no import action and explains that its public API does not expose Cloudbot commands or timers.
- OAuth state is random, single-use, session-bound, dashboard-bound, provider-bound, user-bound, and expires after 15 minutes.
- Backend import reports expose stable reason codes; UI owns translations.
- Frontend uses Nuxt auto-imports, `<Icon />`, semantic Tailwind tokens, existing shadcn-vue components, and `DialogOrSheet`.
- All dashboard copy is translated in `de`, `en`, `es`, `ja`, `pt`, `ru`, `sk`, and `uk` locales.
- Follow red-green-refactor for every production behavior change; generated files and the SQL migration are generated-artifact exceptions.

---

### Task 1: Register StreamElements as a persistent channel integration

**Files:**
- Create: `libs/repositories/integrations/model/model_test.go`
- Create: `libs/migrations/postgres/20260801000000_streamelements_integration.sql`
- Modify: `libs/repositories/integrations/model/model.go`
- Modify: `libs/migrations/seeds/integrations.go`

**Interfaces:**
- Produces: `integrationsmodel.ServiceStreamElements` with database value `STREAMELEMENTS`.
- Produces: an `integrations` row to which `channels_integrations.integrationId` can refer.

- [ ] **Step 1: Write the failing enum test**

```go
func TestServiceStreamElementsValue(t *testing.T) {
	t.Parallel()
	if got := ServiceStreamElements; got != Service("STREAMELEMENTS") {
		t.Fatalf("ServiceStreamElements = %q", got)
	}
}
```

- [ ] **Step 2: Run the test and verify RED**

Run: `cd libs/repositories && go test ./integrations/model -run TestServiceStreamElementsValue`

Expected: compile failure because `ServiceStreamElements` does not exist.

- [ ] **Step 3: Add the service constant**

Add `ServiceStreamElements Service = "STREAMELEMENTS"` beside `ServiceNightbot`.

- [ ] **Step 4: Run the test and verify GREEN**

Run: `cd libs/repositories && go test ./integrations/model -run TestServiceStreamElementsValue`

Expected: PASS.

- [ ] **Step 5: Generate and fill the migration**

Run from the repository root:

```bash
bun cli m create --name streamelements_integration --db postgres --type sql
```

Rename the generated file to `20260801000000_streamelements_integration.sql` if the CLI timestamp differs, then implement:

```sql
-- +goose Up
-- +goose StatementBegin
ALTER TYPE integrations_service_enum ADD VALUE IF NOT EXISTS 'STREAMELEMENTS';
INSERT INTO integrations (service) VALUES ('STREAMELEMENTS') ON CONFLICT DO NOTHING;
-- +goose StatementEnd

-- +goose Down
-- PostgreSQL enum values cannot be removed safely while rows may reference them.
SELECT 1;
```

Also add `STREAMELEMENTS` to `CreateIntegrations` for clean seed runs.

- [ ] **Step 6: Verify formatting and migration ordering**

Run: `gofmt -w libs/repositories/integrations/model/model.go libs/repositories/integrations/model/model_test.go libs/migrations/seeds/integrations.go`

Run: `cd libs/repositories && go test ./integrations/...`

Expected: PASS.

- [ ] **Step 7: Commit**

```bash
git add libs/repositories/integrations/model/model.go libs/repositories/integrations/model/model_test.go libs/migrations/seeds/integrations.go libs/migrations/postgres/20260801000000_streamelements_integration.sql
git commit -m "feat(import): register StreamElements"
```

### Task 2: Add integration-specific single-use OAuth attempts

**Files:**
- Modify: `apps/api-gql/internal/auth/sessions_user.go`
- Modify: `apps/api-gql/internal/auth/sessions_user_test.go`

**Interfaces:**
- Produces: `CreateIntegrationOAuthAttempt(ctx context.Context, service integrationsmodel.Service, channelID, initiatorUserID uuid.UUID) (string, error)`.
- Produces: `ConsumeIntegrationOAuthAttempt(ctx context.Context, state string, service integrationsmodel.Service, channelID, initiatorUserID uuid.UUID, now time.Time) error`.
- Consumes: existing session-backed `OAuthAttempt` map and 15-minute lifetime.

- [ ] **Step 1: Write failing session tests**

Add table tests proving:

```go
state, err := auth.CreateIntegrationOAuthAttempt(ctx, integrationsmodel.ServiceNightbot, channelID, userID)
require.NoError(t, err)
require.NotEmpty(t, state)
require.NoError(t, auth.ConsumeIntegrationOAuthAttempt(ctx, state, integrationsmodel.ServiceNightbot, channelID, userID, now))
require.ErrorIs(t, auth.ConsumeIntegrationOAuthAttempt(ctx, state, integrationsmodel.ServiceNightbot, channelID, userID, now), ErrOAuthAttemptNotFound)
```

Add separate assertions for wrong service, wrong channel, wrong initiator, and `now` after `ExpiresAt`.

- [ ] **Step 2: Run and verify RED**

Run: `cd apps/api-gql && go test ./internal/auth -run IntegrationOAuthAttempt`

Expected: compile failure because the two methods do not exist.

- [ ] **Step 3: Extend the stored attempt and implement create/consume**

Add an optional `IntegrationService *integrationsmodel.Service` to `OAuthAttempt`. Creation uses
`uuid.NewString()` and `time.Now().Add(15*time.Minute)`. Consumption reads without mutating first,
validates all bindings and expiry, and deletes/commits the attempt before returning success. Define
sentinels for mismatch and expiry so tests do not compare strings.

- [ ] **Step 4: Run and verify GREEN**

Run: `cd apps/api-gql && go test ./internal/auth -run 'OAuthAttempt|IntegrationOAuthAttempt'`

Expected: PASS, including existing platform OAuth-attempt coverage.

- [ ] **Step 5: Commit**

```bash
git add apps/api-gql/internal/auth/sessions_user.go apps/api-gql/internal/auth/sessions_user_test.go
git commit -m "fix(auth): bind integration OAuth state"
```

### Task 3: Build the normalized importer

**Files:**
- Create: `apps/api-gql/internal/services/importer/types.go`
- Create: `apps/api-gql/internal/services/importer/service.go`
- Create: `apps/api-gql/internal/services/importer/service_test.go`
- Modify: `apps/api-gql/internal/services/commands/create.go`
- Modify: `apps/api-gql/internal/services/timers/create.go`
- Modify: `apps/api-gql/cmd/main.go`

**Interfaces:**
- Produces: `type FailureReason string` constants `DUPLICATE`, `PLAN_LIMIT`, `UNSUPPORTED_ROLE`, `UNSUPPORTED_RESPONSE_TYPE`, `INCOMPATIBLE_INTERVALS`, `INVALID_RECORD`.
- Produces: `type Failure struct { Name string; Reason FailureReason }`.
- Produces: `type Report struct { ImportedCount int; Failures []Failure }` and `FailedCount() int`.
- Produces: `ImportCommands(ctx context.Context, channelID, actorID string, input []Command) (Report, error)`.
- Produces: `ImportTimers(ctx context.Context, channelID, actorID string, input []Timer) (Report, error)`.
- Consumes: existing `commands.Service.Create`, `timers.Service.Create`, and channel role repository.

The normalized structs are exact:

```go
type RoleRequirement string

const (
	RoleEveryone    RoleRequirement = "EVERYONE"
	RoleSubscriber  RoleRequirement = "SUBSCRIBER"
	RoleVIP         RoleRequirement = "VIP"
	RoleModerator   RoleRequirement = "MODERATOR"
	RoleBroadcaster RoleRequirement = "BROADCASTER"
)

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

- [ ] **Step 1: Write failing importer tests with fakes**

Cover one behavior per test: successful command mapping, role UUID expansion, duplicate AppError,
plan-limit AppError, unexpected infrastructure error, partial success, successful timer mapping, and
invalid normalized record. Assert exact `commands.CreateInput`/`timers.CreateInput` received by fakes.

- [ ] **Step 2: Run and verify RED**

Run: `cd apps/api-gql && go test ./internal/services/importer`

Expected: package or symbols missing.

- [ ] **Step 3: Add stable reason metadata to existing create services**

Return AppErrors with `Details["reason"]` set to `PLAN_LIMIT` for command/timer plan exhaustion and
`DUPLICATE` for conflicts. Keep existing user-visible messages and error codes.

- [ ] **Step 4: Implement normalized conversion and classification**

Resolve broadcaster/moderator/subscriber/VIP roles once. Convert `RoleRequirement` to the same
inclusive role UUID sets already used by Nightbot. For each item call the existing service; append a
stable failure for expected AppErrors and return immediately for internal/unknown errors. Do not wrap
the whole batch in one transaction because partial success is part of the contract.

- [ ] **Step 5: Wire the importer into Fx**

Add `importer.New` to the service providers in `apps/api-gql/cmd/main.go`.

- [ ] **Step 6: Run and verify GREEN**

Run: `cd apps/api-gql && go test ./internal/services/importer ./internal/services/commands ./internal/services/timers`

Expected: PASS.

- [ ] **Step 7: Commit**

```bash
git add apps/api-gql/internal/services/importer apps/api-gql/internal/services/commands/create.go apps/api-gql/internal/services/timers/create.go apps/api-gql/cmd/main.go
git commit -m "refactor(import): share provider importer"
```

### Task 4: Refactor and harden Nightbot against the shared importer

**Files:**
- Create: `apps/api-gql/internal/services/nightbot_integration/mapper.go`
- Create: `apps/api-gql/internal/services/nightbot_integration/mapper_test.go`
- Create: `apps/api-gql/internal/services/nightbot_integration/client_test.go`
- Modify: `apps/api-gql/internal/services/nightbot_integration/nightbot_integration.go`
- Modify: `apps/tokens/internal/bus_listener/bus_listener.go`
- Modify: `apps/tokens/internal/bus_listener/channel_integrations_test.go`

**Interfaces:**
- Consumes: `importer.Service` and its normalized `Command`, `Timer`, and `Report`.
- Produces: `NormalizeCommands(nightbotCustomCommandsResponse) ([]importer.Command, []importer.Failure)`.
- Produces: `NormalizeTimers(nightbotTimersResponse) ([]importer.Timer, []importer.Failure)`.
- Preserves: existing `ImportCommands` and `ImportTimers` service entry points while returning the shared report.
- Changes: `GetAuthLink(ctx, state)` appends the supplied OAuth state without changing the registered redirect URI.

- [ ] **Step 1: Write failing mapping tests**

Assert command prefix removal, lower-casing, alias association, cooldown, owner/moderator/subscriber/VIP
roles, unsupported admin/regular failures, cron-minute parsing, enabled timer state, and invalid cron
failure.

- [ ] **Step 2: Run and verify RED**

Run: `cd apps/api-gql && go test ./internal/services/nightbot_integration -run Normalize`

Expected: compile failure because mapper functions do not exist.

- [ ] **Step 3: Implement provider-only normalization**

Move provider mapping out of the persistence loops. Fetch commands/timers exactly once, normalize,
call the shared importer, then prepend normalization failures to its report.

Add an auth-link test that passes `nightbot-state` and asserts the encoded URL contains exactly that
`state`; update `GetAuthLink` to accept and encode it.

- [ ] **Step 4: Write the failing refresh regression test**

Extend `TestRequestChannelIntegrationToken_NightbotRefreshesViaTokensService` to require:

```go
if got := r.Form.Get("redirect_uri"); got != "https://twir.test/dashboard/integrations/callbacks/nightbot" {
	t.Fatalf("redirect_uri = %q", got)
}
```

- [ ] **Step 5: Run the refresh test and verify RED**

Run: `cd apps/tokens && go test ./internal/bus_listener -run NightbotRefreshesViaTokensService`

Expected: FAIL because the request omits `redirect_uri`.

- [ ] **Step 6: Add Nightbot redirect URI to refresh**

Require `integration.RedirectURL` and send it in the refresh form. Keep rotated refresh-token
persistence behavior.

- [ ] **Step 7: Run and verify GREEN**

Run: `cd apps/api-gql && go test ./internal/services/nightbot_integration`

Run: `cd apps/tokens && go test ./internal/bus_listener -run Nightbot`

Expected: PASS.

- [ ] **Step 8: Commit**

```bash
git add apps/api-gql/internal/services/nightbot_integration apps/tokens/internal/bus_listener/bus_listener.go apps/tokens/internal/bus_listener/channel_integrations_test.go
git commit -m "fix(import): harden Nightbot imports"
```

### Task 5: Implement persistent StreamElements OAuth and imports

**Files:**
- Create: `libs/integrations/streamelements/streamelements_test.go`
- Create: `apps/api-gql/internal/services/streamelements/mapper.go`
- Create: `apps/api-gql/internal/services/streamelements/mapper_test.go`
- Create: `apps/api-gql/internal/services/streamelements/service_test.go`
- Modify: `libs/integrations/streamelements/streamelements.go`
- Modify: `libs/integrations/streamelements/responses.go`
- Modify: `apps/api-gql/internal/services/streamelements/dto.go`
- Modify: `apps/api-gql/internal/services/streamelements/streamelements.go`
- Modify: `apps/tokens/internal/bus_listener/bus_listener.go`
- Modify: `apps/tokens/internal/bus_listener/channel_integrations_test.go`

**Interfaces:**
- Produces client methods `SetAccessToken(string)`, `RefreshToken(ctx, refreshToken, redirectURL string) (*TokenResponse, error)`, and existing profile/commands/timers methods through an injected `*http.Client` and base URL.
- Produces service methods `GetData`, `GetAuthLink(state)`, `PostCode`, `Logout`, `ImportCommands`, and `ImportTimers` mirroring Nightbot.
- Consumes: `channels_integrations.Repository`, `integrations.Repository`, token bus, and shared importer.

- [ ] **Step 1: Write failing HTTP client contract tests**

Use `httptest.Server` to assert authorization query parameters, code exchange form, bearer header,
profile/commands/timers decoding, refresh form, non-2xx sanitization, and malformed JSON handling.
Authorization scopes must contain each of `channel:read` and `bot:read` exactly once.

- [ ] **Step 2: Run and verify RED**

Run: `cd libs/integrations && go test ./streamelements`

Expected: failures for injection and refresh APIs that do not exist.

- [ ] **Step 3: Refactor the client and verify GREEN**

Add constructor options without breaking the production `New(clientID, clientSecret)` call. Centralize
authenticated JSON requests and redact response bodies from returned errors.

Run: `cd libs/integrations && go test ./streamelements`

Expected: PASS.

- [ ] **Step 4: Write failing provider mapping tests**

Assert access-level mapping `100/250/300/400/500/1000+`, visibility, aliases, cooldown, reply/say/action,
unsupported response types, online/offline command flags, single-mode timers, equal dual-mode timer
intervals, and incompatible dual-mode intervals.

- [ ] **Step 5: Implement provider mapping and verify GREEN**

Run: `cd apps/api-gql && go test ./internal/services/streamelements -run Normalize`

Expected: PASS after `mapper.go` converts to shared importer records and failures.

- [ ] **Step 6: Write failing lifecycle service tests**

With repository, importer, and token-bus fakes, assert PostCode create/update, profile persistence,
GetData nil behavior, Logout token clearing, refreshed access-token use, and command/timer import report
composition.

- [ ] **Step 7: Implement the lifecycle service and verify GREEN**

Replace the old ephemeral `ExchangeDataByCode` flow with persistent channel integration behavior.
Return a configuration error instead of `(nil, nil)` when credentials are missing.

Run: `cd apps/api-gql && go test ./internal/services/streamelements`

Expected: PASS.

- [ ] **Step 8: Write and implement StreamElements token refresh tests**

Add a `ServiceStreamElements` case to the token listener. Test client credentials, refresh token,
redirect URI, rotated refresh-token persistence, access-token persistence, and non-2xx behavior.

Run: `cd apps/tokens && go test ./internal/bus_listener -run StreamElements`

Expected: RED before the switch case, GREEN after implementation.

- [ ] **Step 9: Commit**

```bash
git add libs/integrations/streamelements apps/api-gql/internal/services/streamelements apps/tokens/internal/bus_listener/bus_listener.go apps/tokens/internal/bus_listener/channel_integrations_test.go
git commit -m "feat(import): add StreamElements backend"
```

### Task 6: Expose the secure symmetric GraphQL contract

**Files:**
- Modify: `apps/api-gql/internal/delivery/gql/schema/import.graphql`
- Modify: `apps/api-gql/internal/delivery/gql/schema/integrations/integrations-streamlabs.graphql`
- Modify: `apps/api-gql/internal/delivery/gql/resolvers/import.resolver.go`
- Modify: `apps/api-gql/internal/delivery/gql/resolvers/integrations-streamlabs.resolver.go`
- Modify: `apps/api-gql/internal/delivery/gql/resolvers/resolver.go`
- Modify: `apps/api-gql/internal/services/streamlabs_integration/service.go`
- Create: `apps/api-gql/internal/services/streamlabs_integration/service_test.go`
- Create: `apps/api-gql/internal/delivery/gql/resolvers/import_test.go`
- Regenerate: `apps/api-gql/internal/delivery/gql/gqlmodel/models_gen.go`
- Regenerate: `apps/api-gql/internal/delivery/gql/generated.go`

**Interfaces:**
- Produces: `ImportFailureReason`, `ImportFailure`, and symmetric Nightbot/StreamElements report outputs.
- Produces: StreamElements auth/data/logout/import fields with dashboard permission directives.
- Changes: Nightbot and Streamlabs post-code inputs require `{ code, state }`.
- Changes: `streamlabs_integration.Service.GetAuthLink(ctx, state)` encodes state in the provider URL.

- [ ] **Step 1: Write failing resolver tests**

Test that auth-link resolvers create state for the selected dashboard/user, post-code resolvers consume
the matching state, replay fails, StreamElements import maps stable failures, and selected-dashboard
errors are handled through `gqlerrors.HandleError`. Add a Streamlabs service test asserting that a
supplied state is encoded exactly once in its authorization URL.

- [ ] **Step 2: Run and verify RED**

Run: `cd apps/api-gql && go test ./internal/delivery/gql/resolvers -run Import`

Expected: compile failure for the new contract.

- [ ] **Step 3: Update schema and regenerate gqlgen**

Run from the repository root: `bun cli build gql`.

Re-read generated resolver signatures before implementing them. Do not hand-edit generated model or
execution files.

- [ ] **Step 4: Implement resolvers and service wiring**

Resolve selected dashboard and authenticated user for every OAuth start/finish. Convert shared
reports to GraphQL DTOs. Remove `streamelementsExchangeDataByCode` and its raw external DTO types.

- [ ] **Step 5: Run and verify GREEN**

Run: `cd apps/api-gql && go test ./internal/delivery/gql/resolvers ./internal/services/nightbot_integration ./internal/services/streamelements`

Expected: PASS.

- [ ] **Step 6: Commit**

```bash
git add apps/api-gql/internal/delivery/gql
git commit -m "feat(import): expose provider import API"
```

### Task 7: Build the shared localized Nuxt import experience

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
- Modify: `web/layers/dashboard/pages/dashboard/integrations/[name].vue`
- Modify: `web/layers/dashboard/pages/dashboard/integrations/callbacks/[name].vue`
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
- `ImportProviderCard` props: provider title/icon/description, loading, account, auth link, availability, and callbacks; slots the shared settings UI.
- `ImportSettings` accepts command/timer permission, loading, report, and import callbacks.
- `ImportResult` renders GraphQL `ImportFailure[]` using `imports.failureReasons.<enum>` translations.
- Both provider composables expose auth/data/postCode/logout/importCommands/importTimers and a refresh broadcaster.

- [ ] **Step 1: Inspect installed shadcn-vue APIs before authoring components**

Run from `web/`: `bun run shadcn-vue docs card button badge alert`.

Read the returned component documentation and the existing `DialogOrSheet` API. Reuse installed
components; do not add or overwrite registry files.

- [ ] **Step 2: Write failing shared-component tests**

Test disconnected/connected actions, unavailable state without login/settings buttons, permission
disabled buttons, loading spinner, partial result counts, translated failure reasons, and emitted
import actions.

- [ ] **Step 3: Run and verify RED**

Run: `cd web && bun run test -- layers/dashboard/features/import/components`

Expected: component modules missing.

- [ ] **Step 4: Implement shared components and verify GREEN**

Use full Card composition, Badge, Alert, Button, toast via `vue-sonner`, and `DialogOrSheet`. Use gaps,
semantic colors, Nuxt icons, and auto-imported Vue primitives.

Run: `cd web && bun run test -- layers/dashboard/features/import/components`

Expected: PASS.

- [ ] **Step 5: Write failing Nightbot and StreamElements feature tests**

Mock each composable module. Assert OAuth link use, callback `{code,state}`, no close/broadcast on
mutation error, close/broadcast on success, separate command/timer requests, and report rendering.

- [ ] **Step 6: Implement provider components and callbacks**

Refactor Nightbot onto shared UI. Implement StreamElements with the same behavior. Remove stale generic
callback branches and the debug `console.log('here')`. Keep Streamlabs donation integration routing
separate from the import page.

- [ ] **Step 7: Write and implement the Streamlabs unavailable test**

Assert the brand icon, unavailable badge, explanatory copy, and absence of login/import/settings
actions. Implement `streamlabs-import.vue` as presentation only.

- [ ] **Step 8: Add all locale keys**

Use one `imports` namespace with provider descriptions, actions, states, counts, OAuth errors, and
every stable backend failure reason. Ensure every one of the eight locale files has the same key
shape.

- [ ] **Step 9: Regenerate frontend GraphQL types**

Run: `cd web && bun run graphql-codegen`.

Re-read generated types and remove all references to the old raw StreamElements exchange query.

- [ ] **Step 10: Run scoped frontend tests**

Run: `cd web && bun run test -- layers/dashboard/features/import`

Expected: PASS.

- [ ] **Step 11: Commit**

```bash
git add web/layers/dashboard/features/import web/layers/dashboard/pages/dashboard/import.vue web/layers/dashboard/pages/dashboard/integrations web/layers/dashboard/pages/import web/layers/dashboard/api/integrations/integrations.ts web/layers/dashboard/locales web/gql
git commit -m "feat(web): finish chat provider imports"
```

### Task 8: Verify the complete feature and documentation

**Files:**
- Modify if implementation details changed: `docs/superpowers/specs/2026-08-01-chat-imports-design.md`
- Modify if commands or exact paths changed: `docs/superpowers/plans/2026-08-01-chat-imports.md`

**Interfaces:**
- Validates every success criterion from the approved design.

- [ ] **Step 1: Run backend test suites**

```bash
cd libs/repositories && go test ./integrations/...
cd ../../../libs/integrations && go test ./streamelements
cd ../../apps/tokens && go test ./internal/bus_listener
cd ../api-gql && go test ./internal/auth ./internal/services/importer ./internal/services/nightbot_integration ./internal/services/streamelements ./internal/delivery/gql/resolvers
```

Expected: all commands exit 0 with no failures.

- [ ] **Step 2: Run frontend generation, tests, and typecheck**

```bash
cd web
bun run nuxt-prepare
bun run graphql-codegen
bun run test -- layers/dashboard/features/import
bunx --bun vue-tsc --noEmit -p .nuxt/tsconfig.json
```

Expected: all commands exit 0 with no failures.

- [ ] **Step 3: Run formatting and lint checks**

Run `gofmt -w` on every touched Go file, then:

```bash
bun lint
git diff --check
```

Expected: no lint or whitespace errors attributable to the change.

- [ ] **Step 4: Run relevant builds**

```bash
bun cli b app api-gql
cd web && bun run build
```

Expected: both builds exit 0.

- [ ] **Step 5: Perform browser QA through Caddy**

Open `${SITE_BASE_URL}/dashboard/import`, falling back to `http://localhost:3005/dashboard/import`.
Verify desktop and mobile layouts, Nightbot/StreamElements disconnected and connected states, both
import reports, permission-disabled controls, OAuth error popup, and Streamlabs unavailable copy.

- [ ] **Step 6: Review the final diff against the design**

Check additive behavior, permission directives, state consumption before token exchange, no secret
response bodies in errors, locale parity, no `lucide-vue-next`, no manual Vue primitive imports, and
no placeholder pages.

- [ ] **Step 7: Commit verification-only documentation adjustments if present**

```bash
git add docs/superpowers/specs/2026-08-01-chat-imports-design.md docs/superpowers/plans/2026-08-01-chat-imports.md
git commit -m "docs: align chat import implementation"
```
