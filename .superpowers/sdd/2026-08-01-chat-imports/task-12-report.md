# Task 12 report

## Delivered

- Added shared, responsive import provider cards, import settings, and partial-result reporting.
- Migrated Nightbot to the unified integrations query and the `{ code, state }` OAuth callback.
- Added StreamElements OAuth, command/timer import UI, donation status, callback route, and mutations.
- Added Streamlabs donation connection state with an explicit Cloudbot import-unavailable explanation.
- Added a state-aware Streamlabs callback at its configured `/dashboard/integrations/streamlabs` route.
- Removed placeholder import routes and corrected contradictory public Streamlabs import copy.
- Added locale parity for `imports` across de/en/es/ja/pt/ru/sk/uk.

## Verification

- `cd web && bun run graphql-codegen` — pass.
- `cd web && bun run nuxt-prepare` — pass.
- `cd web && bun run test -- layers/dashboard/features/import` — 18 tests pass in 6 files.
- Scoped oxlint over all touched TypeScript/Vue files — 0 warnings, 0 errors.
- `git diff --check` — pass.
- Full `vue-tsc --noEmit -p .nuxt/tsconfig.json` remains globally red because the generated
  config includes unrelated monorepo packages and sibling worktrees with existing errors. Filtering
  its output to this worktree's touched import/API/page paths produces no diagnostics.

## Notes

- Generated GraphQL output lives under ignored `web/app/gql`; it is rebuilt by codegen/prepare.
- The Nuxt test runner emits existing sitemap and SSR-disabled OG-image warnings.

## Review fixes

- OAuth-link fields protected by `MANAGE_INTEGRATIONS` are conditionally selected from the unified
  query using the viewer's reactive permission. Integration data remains available to import roles
  while login/logout stays disabled. Spotify remains unconditional because its schema has no role
  restriction.
- Provider cards now receive an explicit connection state. Streamlabs uses the persisted `enabled`
  flag, so a disabled row with retained profile data renders disconnected/reconnect state and does
  not claim donation listening is active.
- Added permission-query and retained-profile regressions. The scoped suite now passes 21 tests in
  7 files; codegen, Nuxt prepare, scoped lint, filtered type diagnostics, and diff checks pass.
