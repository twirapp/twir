# Task 10 report: StreamElements realtime tips

## Implemented

- Added a refresh-aware StreamElements Bun client using the Task 9 token store and distributed
  lock. It rereads credentials under the lock, skips a redundant provider refresh after another
  process rotates tokens, sends a form-encoded refresh request, preserves an omitted rotated
  refresh token, persists both credentials before changing memory, and exposes only a sanitized
  failure.
- Bounded refresh requests to 15 seconds and response bodies to 1 MiB.
- Added a Socket.IO 2.x-compatible StreamElements connection:
  - connects to `https://realtime.streamelements.com` with WebSocket transport only and the
    built-in reconnect disabled;
  - emits `authenticate` with `{ method: 'oauth2', token }` after `connect`;
  - accepts only tip events and normalizes them into the existing `onDonation` pipeline;
  - applies Task 9's stable event ID claim before publishing while passing events without an ID;
  - refreshes at most once after `unauthorized`, then rebuilds the socket;
  - uses jittered exponential reconnect delays capped at 30 seconds for network failures without
    refreshing credentials;
  - uses generation checks and cancellable timers so stale sockets cannot reconnect after removal.
- Added a per-channel serialized lifecycle store. Concurrent Add operations cannot leave parallel
  sockets, Remove consumes the channel ID, and `closeAll` waits for pending lifecycle operations.
- Added enabled StreamElements integration lookup through the generic integration table, startup
  loading, Add-by-row-ID and Remove-by-channel-ID bus handling, and SIGTERM/SIGINT cleanup.
- Added the existing Go-side production credential names (`STREAM_ELEMENTS_CLIENT_ID` and
  `STREAM_ELEMENTS_CLIENT_SECRET`) to the TypeScript config schema so both runtimes share one
  deployment configuration.

## TDD evidence

- RED: provider and socket tests initially failed because `streamelements-client.ts` and
  `streamElements.ts` did not exist; store tests likewise failed before `store/streamelements.ts`.
- Targeted GREEN:
  `cd apps/integrations && bun test src/libs/streamelements-client.test.ts src/services/streamElements.test.ts src/store/streamelements.test.ts`
  passed.
- Full suite: `cd apps/integrations && bun test`
  - 33 passed, 0 failed, 91 assertions after review regressions.
- Typecheck: `cd apps/integrations && bun run prebuild`
  - passed (`tsc --noEmit`).
- Build: `cd apps/integrations && bun run build`
  - passed and compiled `.out/twir-integrations`.
- Whitespace: `git diff --check`
  - passed.

## Self-review

- Inspected installed `socket.io-client@2.3.1` and its v1-compatible type package; no v4 `auth`
  option is used.
- No provider response body, access token, refresh token, client secret, or event payload is logged.
- Refresh persistence occurs before in-memory replacement, including the persistence-error path.
- Ordinary reconnects never invoke refresh and cannot exceed the 30-second delay cap.
- Startup and bus lifecycle operations await StreamElements store changes.

## Follow-up

- Task 11 will move the legacy Streamlabs connection onto the same lifecycle and refresh/dedupe
  primitives and can broaden shutdown cleanup to all provider stores.

## Review fixes

- Bus callbacks now detach through an explicit rejection handler because `Queue.subscribe` does
  not await returned promises.
- Add-by-row-ID performs an initial channel discovery and then rereads the authoritative enabled
  row inside that channel's serialized lifecycle queue. A Remove/logout that wins the race cannot
  be undone by a delayed stale Add, and shutdown also rejects Add operations still discovering a
  channel.
- The one-refresh guard resets only after the provider emits `authenticated`, allowing a future
  token expiry to start a new refresh cycle without looping on a failed authentication attempt.
- Tip delivery captures the socket generation and rechecks it after the awaited Redis claim, so a
  removed or replaced connection cannot publish a late donation.
