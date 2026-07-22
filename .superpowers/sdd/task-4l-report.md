# Task 4 Slice 4l Report

## Scope

- Migrated only the REST dashboard ownership checks in:
  - `apps/api-gql/internal/delivery/http/middlewares/has_access_to_selected_dashboard.go`
  - `apps/api-gql/internal/delivery/http/middlewares/has_channel_roles_dashboard_permission.go`
- Did not modify GraphQL Task 11 directives, OAuth Task 9 routes, schemas,
  migrations, Docker, or databases.

## Binding Ownership Decision

- `Middlewares` now receives the existing `ChannelService` through Fx.
- The selected dashboard ID is resolved with `GetChannelByID`, and ownership is
  determined by membership in the returned normalized `Channel.Bindings`.
- Every binding is checked by internal `UserID`; no positional binding behavior
  or legacy Twitch/Kick owner field is used. This grants owner access for
  Twitch, Kick, and VK Video Live bindings.
- Authenticated-user lookup, selected-dashboard session lookup, Huma status
  codes/messages, bot-admin bypass, role queries, and role/stat authorization
  behavior are unchanged.

## TDD Evidence

### RED

After adding `dashboard_access_test.go` before production changes:

```sh
go test ./internal/delivery/http/middlewares -run 'Test(MiddlewaresIsSelectedDashboardOwnerUsesNormalizedBindings|HasChannelRolesDashboardAccess)' -count=1
```

Failed as expected because the generic channel getter, normalized owner helper,
and shared role-access helper did not exist:

```text
unknown field channelGetter in struct literal of type Middlewares
middlewares.isSelectedDashboardOwner undefined
undefined: hasChannelRolesDashboardAccess
```

### GREEN

After the minimal normalized lookup and role-policy extraction:

```sh
go test ./internal/delivery/http/middlewares -run 'Test(MiddlewaresIsSelectedDashboardOwnerUsesNormalizedBindings|HasChannelRolesDashboardAccess)' -count=1
```

Passed. Coverage includes Twitch, Kick, and VK Video Live owners with bindings
in reverse platform order, an unlinked nonowner, direct and stat-derived role
access, requested permission access, dashboard-wide access permission, and a
denied role permission.

## Changed Files

- `apps/api-gql/internal/delivery/http/middlewares/middlewares.go`
- `apps/api-gql/internal/delivery/http/middlewares/dashboard_access.go`
- `apps/api-gql/internal/delivery/http/middlewares/dashboard_access_test.go`
- `apps/api-gql/internal/delivery/http/middlewares/has_access_to_selected_dashboard.go`
- `apps/api-gql/internal/delivery/http/middlewares/has_channel_roles_dashboard_permission.go`

## Verification

- Passed: `gofmt -w` on all changed Go files.
- Passed: focused GREEN command after the production implementation.
- Passed: `go test ./internal/delivery/http/middlewares -count=1`.
- Passed: `go vet ./internal/delivery/http/middlewares`.
- Passed: static audit found no `IsOwner` or legacy `model.Channels` use in the
  two scoped REST middleware files.
- Passed: `git diff --check`.
- Ran: `go test ./internal/delivery/http/... -count=1`.

## Concerns

- The broader HTTP package check is blocked outside this slice. Unchanged OAuth,
  overlays, scheduled-VIP, and commands callers still reference removed legacy
  channel APIs/fields or absent generated GraphQL packages. The focused
  middleware package passes.
- No Docker, database, migration, OAuth, GraphQL, or provider integration work
  was performed.

## Commit

- `fix(api-gql): use bindings for dashboard ownership`
