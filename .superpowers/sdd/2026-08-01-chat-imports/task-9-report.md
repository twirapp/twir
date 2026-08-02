# Task 9 report: Bun token lifecycle primitives

## Implemented

- Added a shared OAuth refresh lock compatible with the Go lock key and Redis protocol:
  - opaque per-acquisition owner;
  - `SET NX PX 30000` acquisition with bounded retries;
  - owner-checked Lua renewal and release;
  - 10-second renewal, 5-second renewal deadline, and 25-second lease watchdog;
  - callback abort plus failure on lease loss;
  - bounded cleanup without persistent background timers.
- Added injected StreamElements and Streamlabs token stores:
  - StreamElements uses the enabled generic integration joined by `integrations.service = 'STREAMELEMENTS'`;
  - Streamlabs uses its enabled dedicated table;
  - both persist access and refresh tokens in one statement and fail when no enabled row matches.
- Added donation event claiming with a 24-hour Redis `NX`/`EX` key and pass-through for missing IDs.
- Exposed production provider token stores from the existing Bun SQL boundary without adding eager Redis/config dependencies to unit tests.

## TDD evidence

- RED: the three new test modules initially failed because their implementation modules did not exist.
- GREEN: `cd apps/integrations && bun test src/libs/oauth-lock.test.ts src/libs/provider-token-store.test.ts src/libs/donation-dedupe.test.ts`
  - 10 passed, 0 failed, 34 assertions.
- Typecheck: `cd apps/integrations && bun run prebuild`
  - passed (`tsc --noEmit`).
- Lint: scoped `oxlint` over the seven changed TypeScript files
  - 0 warnings, 0 errors.
- Whitespace: `git diff --check`
  - passed.

## Risks / follow-up

- Provider-specific refresh HTTP behavior and socket lifecycle intentionally remain for Tasks 10 and 11.
- Existing legacy Streamlabs update helpers remain until the Task 11 migration to the new token store.
