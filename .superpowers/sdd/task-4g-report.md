# Task 4 Slice 4g Report

## Scope

- Migrated only the unowned emotes-cacher provider consumers identified for Task 4g.
- Changed startup cache filling, BTTV startup subscription discovery, and 7TV startup/profile binding selection.
- Did not modify provider API client contracts, EventSub, Events, Parser, Scheduler, Bots, OAuth, GraphQL, NATS, Docker, migrations, or database state.

## Mapping And Query Decisions

- `emotes_store/fill.go` starts from one `channel_platforms cp` row, explicitly limits `cp.platform` to Twitch and Kick, joins `channels c` only to retain `c."isEnabled" = true`, and reads the provider ID from `cp.platform_channel_id`.
- Fill intentionally processes one row per selected binding: a dual-bound channel fetches its Twitch and Kick 7TV caches independently without a cross-product. BTTV and FFZ remain Twitch-only; no VK cache path was added.
- BTTV startup selects only the Twitch binding with `cp.platform = twitch` and `cp.enabled = true`. This replaces the legacy Twitch bot-enabled condition and deliberately does not add a global channel-enabled filter.
- 7TV startup selects a Twir channel once with `EXISTS` over enabled Twitch or Kick bindings. This preserves the former any-enabled-provider eligibility while preventing dual-bound channel duplication.
- 7TV per-channel profile loading reads explicit Twitch and Kick binding rows without an `enabled` predicate, preserving the prior behavior of loading every connected supported provider after startup eligibility. It switches on each binding's platform and uses that binding's `platform_channel_id`; no binding slice position or users-repository lookup is used.
- Existing binding indexes support the `channel_id` and `(platform, platform_channel_id)` query predicates. No schema or index change was required.

## TDD Evidence

### RED

- Baseline: `go test -count=1 ./apps/emotes-cacher/...` failed because `seventv/add_channels.go` still referenced removed `Channel.TwitchUserID` and `Channel.KickUserID` fields.
- After adding query tests first:
  - `go test -count=1 ./apps/emotes-cacher/internal/emotes_store -run '^TestBuildStartupChannelsQueryUsesExplicitNormalizedBindings$'` failed because `buildStartupChannelsQuery` was absent.
  - `go test -count=1 ./apps/emotes-cacher/internal/services/bttv -run '^TestBuildEnabledTwitchChannelsQueryUsesTwitchBindingEligibility$'` failed because `buildEnabledTwitchChannelsQuery` was absent.
  - `go test -count=1 ./apps/emotes-cacher/internal/services/seventv -run '^(TestBuildStartupChannelsQueryUsesEnabledTwitchOrKickBindingWithoutDuplicateChannels|TestBuildChannelBindingsQueryUsesExplicitTwitchAndKickBindingsRegardlessOfBindingEnablement)$'` failed from the removed legacy model fields and absent normalized query helpers.
- The compile-time RED state is expected for Go query-builder tests: packages compile before dry-run SQL assertions can execute.

### GREEN

- Focused emote-store, BTTV, and 7TV query tests passed after the minimal migration.
- `go test -count=1 ./apps/emotes-cacher/...` passed.
- `gofmt -d` for all changed Go files produced no output.
- `git diff --check` passed.

## Behavior Coverage

- Query tests assert normalized binding-sourced provider IDs, explicit Twitch/Kick filters, no legacy channel columns, and a single `channel_platforms` source.
- Fill tests retain global channel enablement and reject an added binding-enabled filter.
- BTTV tests retain binding enablement and reject an added global channel-enabled filter.
- 7TV startup tests require `EXISTS` over enabled supported bindings, so a dual-bound channel is selected once.
- 7TV profile-binding tests require both supported platform bindings without a binding or global enablement predicate, preventing accidental behavior changes and VK support.

## Changed Files

- `apps/emotes-cacher/go.mod`
- `apps/emotes-cacher/internal/emotes_store/fill.go`
- `apps/emotes-cacher/internal/emotes_store/fill_test.go`
- `apps/emotes-cacher/internal/services/bttv/service.go`
- `apps/emotes-cacher/internal/services/bttv/service_test.go`
- `apps/emotes-cacher/internal/services/seventv/service.go`
- `apps/emotes-cacher/internal/services/seventv/add_channels.go`
- `apps/emotes-cacher/internal/services/seventv/bindings_test.go`
- `.superpowers/sdd/task-4g-report.md`

## Limitations

- No Docker, database, migration, EXPLAIN, or live Twitch/Kick/7TV/BTTV provider call was used.
- Query coverage uses GORM PostgreSQL DryRun to validate SQL shape and identifier provenance; it does not execute against local database rows.
- The direct PostgreSQL GORM driver is a test dependency and is declared directly in `apps/emotes-cacher/go.mod`, matching the existing Scheduler query-test convention.

## Commit

- `96c1ecd2a` `fix(emotes-cacher): use channel platform bindings`

## Important Test Follow-up

### Scope

- Addressed only the Task 4g review finding that DryRun tests asserted SQL placeholders but not their bound values.
- Production code and dependencies are unchanged.

### RED

- Deliberately incorrect ordered `statement.Vars` expectations failed under the existing query builders:
  - Fill returned `{true, twitch, kick}` rather than `{false, kick, twitch}`.
  - BTTV returned `{twitch, true}` rather than `{kick, false}`.
  - 7TV startup returned `{twitch, kick, true}` rather than `{kick, twitch, false}`.
  - 7TV profile binding returned `{channel-id, twitch, kick}` rather than `{wrong-channel-id, kick, twitch}`.
- Commands:
  - `go test -count=1 ./apps/emotes-cacher/internal/emotes_store -run '^TestBuildStartupChannelsQueryUsesExplicitNormalizedBindings$'`
  - `go test -count=1 ./apps/emotes-cacher/internal/services/bttv -run '^TestBuildEnabledTwitchChannelsQueryUsesTwitchBindingEligibility$'`
  - `go test -count=1 ./apps/emotes-cacher/internal/services/seventv -run '^(TestBuildStartupChannelsQueryUsesEnabledTwitchOrKickBindingWithoutDuplicateChannels|TestBuildChannelBindingsQueryUsesExplicitTwitchAndKickBindingsRegardlessOfBindingEnablement)$'`

### GREEN

- Each query test now compares its complete ordered `statement.Vars` with exact supported platform and eligibility values using `reflect.DeepEqual`.
- The focused emote-store, BTTV, and 7TV test commands passed.
- `go test -count=1 ./apps/emotes-cacher/...` passed.
- `gofmt -d` and `git diff --check` produced no output.

### Changed Files

- `apps/emotes-cacher/internal/emotes_store/fill_test.go`
- `apps/emotes-cacher/internal/services/bttv/service_test.go`
- `apps/emotes-cacher/internal/services/seventv/bindings_test.go`
- `.superpowers/sdd/task-4g-report.md`

### Commit

- `5ab09d783` `test(emotes-cacher): assert query vars`
