# Chat Provider Imports Design

## Goal

Finish the dashboard import experience for Nightbot, StreamElements, and Streamlabs. Nightbot and
StreamElements must support authenticated, server-side import of commands and timers with the same
user experience. Streamlabs must clearly explain that import is unavailable because its public API
does not expose Cloudbot commands or timers.

## Scope

- Finish and harden Nightbot OAuth, command import, timer import, reporting, and UI.
- Add persistent StreamElements OAuth and server-side command and timer import.
- Replace the Streamlabs placeholder with an informative, localized unavailable state.
- Add automated backend and frontend coverage for the new behavior.
- Keep imports additive. Existing Twir commands and timers are never overwritten or deleted.

The work does not add file-based Streamlabs import, selection of individual external items, default
bot-command import, or automatic variable-syntax translation between providers and Twir.

## Architecture

### Normalized importer

A provider-independent importer in `apps/api-gql/internal/services/importer` owns Twir business
rules. Provider services fetch and normalize external records; they do not write commands or timers
directly. The importer exposes separate command and timer operations and returns a stable report:

- imported count;
- failed count;
- failed item names with a localized-at-the-UI reason code;
- no silent overwrites.

The importer uses the existing command and timer services so plan limits, audit records, repository
constraints, and entity creation remain centralized. It resolves channel roles once per command
batch, invalidates the command cache after a command batch, and publishes timer lifecycle events via
the existing timer service. Each item is independent: one conflict does not roll back successfully
created siblings. Infrastructure failures that make the whole batch unreliable return a GraphQL
error instead of being reported as an item conflict.

Nightbot is refactored to normalize its API responses and delegate persistence to this importer.
StreamElements uses the same importer, so both providers have identical duplicate, limit, and report
semantics.

### Normalized command

The normalized command carries name, response, enabled state, aliases, global cooldown, role
requirement, reply behavior, online/offline response restrictions, and visibility. Provider-specific
values that Twir cannot represent produce a failed-item result instead of silently widening access.

StreamElements permission mapping is conservative:

| StreamElements level | Twir roles allowed |
| --- | --- |
| 100 Everyone | unrestricted |
| 250 Subscriber | broadcaster, moderator, subscriber |
| 300 Regular | unsupported; item fails |
| 400 VIP | broadcaster, moderator, VIP |
| 500 Moderator | broadcaster, moderator |
| 1000 or higher | broadcaster only |

Nightbot keeps its current supported role mapping. Nightbot `admin` and `regular`, and unknown
provider levels, remain explicit failed items.

StreamElements `reply` becomes a Twir reply. `say` and `action` become normal responses because Twir
has no distinct action response. Whisper/custom response types that cannot be represented fail the
item. Hidden commands are imported as not visible. Online/offline flags are retained on the response.

### Normalized timer

The normalized timer carries name, message, enabled state, online/offline enablement, time interval,
and message interval. Nightbot cron intervals continue to be parsed into minutes. StreamElements
timers use the enabled online or offline interval. If both modes are enabled with different
intervals, the timer fails explicitly because Twir has one shared interval and choosing one would
silently change behavior.

### StreamElements integration persistence

`STREAMELEMENTS` is added to `integrations_service_enum` and seeded in `integrations`. Per-channel
access token, refresh token, username, and avatar use the existing `channels_integrations` repository
and model, matching Nightbot. The StreamElements client gains an injectable HTTP client and token
refresh operation. The tokens service refreshes StreamElements tokens and persists rotations before
imports.

Only custom channel commands and timers are fetched. Default StreamElements commands are excluded.
The minimal OAuth scopes are `channel:read` and `bot:read` unless the provider rejects that set in
contract tests against documented behavior.

### GraphQL contract

StreamElements receives the same lifecycle and import surface as Nightbot:

- integration data;
- authorization URL;
- post authorization code with OAuth state;
- logout;
- import commands;
- import timers.

Import mutations require `MANAGE_COMMANDS` or `MANAGE_TIMERS`. OAuth lifecycle mutations require
`MANAGE_INTEGRATIONS`; integration data requires `VIEW_INTEGRATIONS`. All operations resolve the
selected dashboard on the server. The old query that exchanges a code and returns raw external DTOs
is removed after the frontend is migrated because it does not perform an import.

## OAuth Security

Nightbot, StreamElements, and Streamlabs authorization URLs include a cryptographically random,
single-use `state`. The state is stored in the authenticated session using the existing OAuth attempt
mechanism and records provider, selected dashboard, initiator, creation time, and redirect target.
The callback mutation accepts both `code` and `state`, verifies provider/dashboard/user/lifetime,
deletes the attempt before exchanging the code, and rejects missing, expired, replayed, or mismatched
states.

OAuth attempts expire after 15 minutes. Callback failures remain visible in the popup and never send
a success refresh notification. Successful callbacks notify the opener and close the popup.

Nightbot refresh requests include the registered `redirect_uri`. Nightbot logout clears local data
and attempts provider token revocation; failure to revoke is logged but does not keep a stale local
integration enabled. StreamElements follows the same local logout behavior and revokes remotely only
if its public API supports revocation.

## Frontend

The import page renders three equal provider cards in the dashboard layout.

Nightbot and StreamElements share an import-card presentation while keeping provider API state in
separate composables. Each card shows its icon, localized description, connection state, connected
account, login/logout action, and settings action. The settings action opens the existing responsive
`DialogOrSheet` and contains separate command and timer sections.

Each import section has:

- a localized title and description;
- an import button disabled by connection state, permission, or active request;
- a spinner while importing;
- a success/partial-success summary;
- a scrollable list of failed item names and translated reason codes;
- a toast for transport or GraphQL errors.

The UI never treats a GraphQL response containing an error as success. OAuth popup callbacks display
an actionable localized error and a close button if connection fails.

The Streamlabs card uses the same card structure and brand icon but has no login or import action. It
shows a localized unavailable badge and explanation that the public Streamlabs API does not expose
Cloudbot commands or timers. It does not imply that adding Streamlabs as a donation integration will
enable import.

All copy is added to every dashboard locale. Vue and Nuxt composables use auto-imports. UI uses
existing shadcn-vue components, semantic Tailwind tokens, Nuxt `<Icon />`, and `DialogOrSheet`; no new
component library or custom CSS is introduced.

## Error Semantics

Expected per-item failures are returned in the report with stable reason codes:

- duplicate name or alias;
- plan limit reached;
- unsupported role level;
- unsupported response type;
- incompatible online/offline timer intervals;
- invalid provider record.

Authentication failure, expired token, provider outage, malformed provider response, repository
failure, cache failure that risks stale behavior, and missing plan/roles are operation failures and
surface as GraphQL errors. Provider response bodies containing credentials or user data are not
included verbatim in user-facing errors or logs.

## Testing

Backend tests use fakes and `httptest.Server`:

- normalized command role and response mappings for both providers;
- timer interval conversions and incompatible interval rejection;
- additive import, duplicate conflicts, plan limits, partial success, cache invalidation, and timer
  lifecycle publication;
- StreamElements authorization, code exchange, profile, commands, timers, refresh, malformed JSON,
  and non-2xx responses;
- Nightbot refresh includes `redirect_uri`;
- OAuth state success, mismatch, expiry, replay, and cross-provider rejection;
- GraphQL permission and selected-dashboard boundaries.

Frontend Vitest coverage verifies:

- disconnected and connected Nightbot/StreamElements cards;
- command and timer loading, success, partial success, and error states;
- permission-disabled actions;
- callback success and failure behavior;
- Streamlabs unavailable content and absence of import/login actions.

Final verification runs GraphQL generation, Go tests for affected modules, scoped dashboard Vitest,
web typecheck, lint for touched code, and relevant backend/web builds.

## Success Criteria

- A permitted dashboard user can connect Nightbot or StreamElements and independently import all
  supported custom commands and timers.
- Existing Twir data is preserved and conflicts are reported without hiding successful imports.
- OAuth callbacks cannot be replayed or attached to another provider, user, or dashboard.
- Nightbot continues importing after access-token expiry through a valid refresh request.
- Streamlabs clearly explains why no import action exists.
- The three cards are localized, responsive, visually consistent, and covered by automated tests.
