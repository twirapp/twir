# Task 8 report: secure GraphQL provider contract

Implemented the shared OAuth/import GraphQL contract for Nightbot, StreamElements, and Streamlabs.

## Delivered

- Added `IntegrationOAuthCodeInput`, `ImportFailureReason`, `ImportFailure`, and `ImportReport`.
- Added StreamElements data, auth-link, callback, logout, command-import, and timer-import fields with selected-dashboard permission directives.
- Changed Nightbot and Streamlabs callbacks to accept `{ code, state }`.
- Created provider/user/dashboard-bound OAuth attempts for all three authorization links and consumed them before every provider code exchange.
- Returned a sanitized `BAD_REQUEST` for expired, replayed, or mismatched OAuth attempts.
- Removed the raw `streamelementsExchangeDataByCode` contract and its unused DTO mappers.
- Preserved Nightbot importer failure reasons in the shared report shape.
- Replaced Streamlabs' variadic state API with an exact state argument in the service and provider library.
- Made the MCP Streamlabs toggle return an explicit dashboard-OAuth requirement; the old zero-state call already failed at runtime and is no longer representable.

## Verification

- `bun cli build gql`: pass.
- `cd apps/api-gql && go test ./...`: 323 passed across 118 packages.
- `cd apps/api-gql && go test -race ./internal/delivery/gql/resolvers ./internal/services/nightbot_integration ./internal/services/streamelements ./internal/services/streamlabs_integration`: 75 passed.
- `cd apps/api-gql && go vet ./...`: pass.
- `cd libs/integrations/streamlabs && go test -race ./...`: 15 passed.
- `cd libs/integrations/streamlabs && go vet ./...`: pass.
- `git diff --check`: pass.

## Follow-up boundary

The dashboard still references the previous Nightbot/Streamlabs mutation variables and result types. Task 12 must migrate those operations and regenerate the frontend GraphQL types before the whole web build is expected to pass.
