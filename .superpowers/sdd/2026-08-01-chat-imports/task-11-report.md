# Task 11 report: hardened Streamlabs runtime

## Implemented

- Added a typed Streamlabs provider client for `/api/v2.0/socket/token`. A valid access token is
  used directly; only the first `401` enters the shared per-channel OAuth lock, rereads persisted
  credentials, refreshes if still necessary, persists both tokens, and retries exactly once.
- Streamlabs refresh now uses `application/x-www-form-urlencoded`, includes `redirect_uri`, keeps
  the previous refresh token when rotation omits one, and updates in-memory credentials only after
  successful persistence.
- Provider requests have a 15-second deadline and 1 MiB response bound. Public errors are typed
  and sanitized; response bodies, credentials, socket tokens, and client secrets are never logged.
- Replaced the legacy Streamlabs socket wrapper with an injected Socket.IO 2.x-compatible
  connection using WebSocket-only transport, disabled built-in reconnection, and a forced new
  manager while preserving the provider's query-token protocol.
- Donation payloads are validated message-by-message. Stable message IDs (or event-ID/index
  fallbacks) pass through the shared Redis dedupe before the existing donation pipeline; malformed,
  duplicate, non-donation, and independently failing messages do not suppress valid siblings.
- Added generation-guarded jittered exponential reconnects capped at 30 seconds. Replacement,
  removal, and shutdown cancel stale timers and prevent late donation publication.
- Added a serialized per-channel Streamlabs lifecycle store with startup load, authoritative
  Add-by-row-ID reread, Remove-by-channel-ID, one live socket per channel, and full shutdown drain.
- Removed the obsolete `updateStreamlabsIntegration` helper and the legacy paths that refreshed on
  every Add, logged provider bodies, or disabled integrations after transient failures.

## TDD and verification

- RED confirmed before each layer: the new client module/exports and lifecycle store were absent.
- Targeted suite: 17 passed, 0 failed, 57 assertions.
- Full `apps/integrations` suite: 50 passed, 0 failed, 148 assertions.
- `bun run prebuild`: passed (`tsc --noEmit`).
- `bun run build`: passed and compiled `.out/twir-integrations`.
- Scoped oxlint: 0 warnings, 0 errors across all eight touched runtime/test files.
- `git diff --check`: passed.

## Self-review

- Socket.IO usage matches the installed 2.3.1 API (`io.connect`, `forceNew`, explicit event
  listeners, and `close`) and does not use the v4-only `auth` option.
- A second socket-token `401` cannot trigger another refresh; network socket reconnects never
  refresh OAuth credentials.
- Token-store predicates continue to require an enabled Streamlabs row, so logout cannot be undone
  by an in-flight refresh or a stale Add lifecycle message.
- No Streamlabs runtime path writes `enabled=false`; auth/network failures remain recoverable.
