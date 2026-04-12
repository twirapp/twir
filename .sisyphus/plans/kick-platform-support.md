# Kick.com Platform Support — Multi-Platform Bot

## TL;DR

> **Quick Summary**: Transform Twir from a Twitch-only bot into a platform-agnostic streaming bot by adding Kick.com as a second platform, with a clean architecture that allows adding more platforms in the future.
>
> **Deliverables**:
>
> - New `user_platform_accounts` table + `users.id` → internal UUID migration (big-bang with maintenance window)
> - Kick OAuth 2.1 + PKCE auth flow; Twitch auth refactored to PlatformProvider interface
> - Kick EventSub webhooks endpoint in `apps/eventsub`; broadcaster's own token used for event subscriptions
> - Generic `ChatMessage` struct in `libs/bus-core` replacing Twitch-coupled struct (parallel queues strategy)
> - `platforms text[]` column on commands, timers, keywords; parser filters execution by platform
> - Built-in `$(platform)` variable
> - `KickProfile` + `linkedAccounts` added to GraphQL schema
> - 7TV Kick profile lookup (`GetProfileByKickId`)
> - `kick_bots` table for the single Kick bot account
> - Frontend: Kick login button, linked accounts UI, platform selector on command/timer/keyword forms
>
> **Estimated Effort**: XL
> **Parallel Execution**: YES — 9 phases, heavily parallelized within each
> **Critical Path**: Task 1 (DB migration) → Task 8 (Auth refactor) → Task 15 (Kick EventSub webhooks) → Task 22 (Generic ChatMessage bus) → Task 29 (Parser platform filtering) → Task 35 (GraphQL schema) → Task 43 (Frontend auth)

---

## Context

### Original Request

> "Давай накидаем поддержки kick.com в боте. Мы должны принимать события через тот же @apps/eventsub/ приложение по вебхукам, мы должны уметь в наших переменных и командах в @apps/parser/ определять с какой платформы, мы должны сделать вход и по кику и по твичу, мы должны иметь возможность привязать к своему профилю и кик и твич, команды, таймеры, и всё такое должны быть со свитчерами платформ, типо чтобы они работали только на кике, только на твиче, или и там и там. Мы должны уметь определять 7tv сэты с кика, то есть искать профиль по платформе. Это очень большая таска. Думай максимально хорошо как сделать бот не привязанным к платформе. По хорошему я хочу в будущем добавлять и другие платформы."

### Interview Summary

**Key Decisions**:

- **Auth model**: Universal login screen with platform selection (Twitch / Kick). First linked platform is primary; others are "linked accounts"
- **Platform switchers**: `platforms text[] DEFAULT '{}'` on commands/timers/keywords. Empty = all platforms; `['twitch']` = Twitch only
- **Kick EventSub**: Official Kick EventSub Webhooks via HTTP. **Each broadcaster must authorize Twir** (broadcaster's own token is used for event subscriptions)
- **ID strategy**: New `user_platform_accounts` table (internal UUID FK → users, platform, platform_user_id, platform_login, tokens). `channels` gets a `platform` column
- **users.id migration**: Big-bang approach with maintenance window (not zero-downtime)
- **Token storage**: Platform-specific tokens in `user_platform_accounts`; existing `tokens` table repurposed for bot tokens only
- **Channels schema**: One row per platform per user (Twitch channel + Kick channel = 2 rows per user)
- **Kick bot account**: Single Kick bot account stored in new `kick_bots` table; sends messages via `POST /public/v1/chat` using bot's own token
- **Kick EventSub subscription trigger**: Automatic on Kick OAuth link — broadcaster's token used to subscribe to their events

**Research Findings**:

- `users.id` currently = Twitch ID (TEXT string), referenced as FK in 40+ tables across 15+ direct tables
- `channels.id` currently = Twitch ID (same value as `users.id`), referenced in 45+ FK relationships
- Auth handler at `apps/api-gql/internal/delivery/http/routes/auth/post-code.go` uses GORM + `libs/gomodels` — must remain untouched for Twitch; new Kick handler uses pgx only
- Session stores `model.Users{ID: "twitch_id"}` + `helix.User` in Redis via SCS — all sessions must be invalidated at migration time
- `tokens` repo: `GetByUserID(userID string)` — currently queries by Twitch ID; must be updated to use internal UUID post-migration
- EventSub app uses Twitch WebSocket conduits — it does NOT currently expose an HTTP endpoint for webhooks
- 7TV `GetProfileByTwitchId` is Twitch-only; Kick profile endpoint exists: `https://7tv.io/v3/users/kick/{channelID}`
- `libs/bus-core/twitch/chat-message.go`: `TwitchChatMessage` struct — platform-coupled; parallel queue strategy required

### Metis Review

**Identified Gaps** (addressed):

- Kick App Token subscription model → resolved: each broadcaster must authorize; broadcaster's token used for subscriptions
- Kick bot sending model → resolved: single `kick_bots` table; bot account sends
- Token architecture → resolved: `user_platform_accounts` stores platform tokens
- `channels` schema → resolved: one row per platform per user
- UUID migration strategy → resolved: big-bang with maintenance window
- Kick OAuth 2.1 PKCE requirement → incorporated into auth task
- Session invalidation on UUID migration → explicitly planned in migration task
- Webhook idempotency (Kick-Event-Message-Id) → required from day one
- Parallel NATS queue strategy → explicitly planned before removing old queue
- `kick_bots` table for bot account → separate table confirmed
- SSE Kick updates for 7TV → OUT of scope Phase 1

---

## Work Objectives

### Core Objective

Make Twir a platform-agnostic bot by adding Kick.com as the second supported platform, laying clean architectural foundations (internal UUID identity, generic bus messages, PlatformProvider interface) so that future platforms can be added with minimal friction.

### Concrete Deliverables

- `user_platform_accounts` table (platform identity + tokens) with migration from `users.id` Twitch ID → internal UUID
- `kick_bots` table for single Kick bot account
- `platforms text[]` on `channels_commands`, `channels_timers`, `channels_keywords`
- `platform text` on `channels` table (with separate row per platform)
- `PlatformProvider` interface + Twitch + Kick implementations in `apps/api-gql`
- `POST /auth/kick/code` HTTP endpoint with OAuth 2.1 + PKCE
- `POST /webhook/kick` HTTP endpoint in `apps/eventsub` with HMAC verification + idempotency
- Kick EventSub subscription manager (subscribe/unsubscribe using broadcaster's token)
- Generic `ChatMessage` struct in `libs/bus-core` with parallel queue migration
- `ParseContext.Platform` field + `$(platform)` built-in variable
- Platform filtering in command/timer/keyword execution
- `KickProfile` GraphQL type + `AuthenticatedUser.kickProfile` + `AuthenticatedUser.linkedAccounts`
- `GetProfileByKickId` in `libs/integrations/seventv`
- Frontend: Kick login button, linked accounts settings section, platform selectors on forms

### Definition of Done

- [ ] `bun cli build` passes with zero errors
- [ ] `bun lint` passes
- [ ] All existing tests pass (no regressions)
- [ ] Kick user can log in and their `user_platform_accounts` row is created
- [ ] Twitch user login still works end-to-end (no regression)
- [ ] Kick EventSub webhook `chat.message.sent` triggers a `ChatMessage{Platform:"kick"}` on NATS
- [ ] Command with `platforms:["twitch"]` does NOT execute for Kick messages
- [ ] `$(platform)` variable returns correct platform string

### Must Have

- Internal UUID primary key for `users` (replacing Twitch ID as PK)
- `user_platform_accounts` table storing platform identity + tokens per platform
- Kick OAuth 2.1 + PKCE auth flow
- Kick EventSub HTTP webhooks with HMAC-SHA256 signature verification
- Webhook idempotency via `Kick-Event-Message-Id`
- Generic `ChatMessage` bus type (parallel queue, no old queue removed until all consumers migrated)
- Parser platform awareness (`ParseContext.Platform`)
- Platform filter on commands/timers/keywords
- `$(platform)` built-in variable
- 7TV `GetProfileByKickId`
- `KickProfile` in GraphQL
- Frontend Kick login + linked accounts UI

### Must NOT Have (Guardrails)

- **No GORM in new code** — all new Go code uses pgx via `libs/repositories` only
- **No `libs/gomodels` in new code** — use `libs/entities` for new domain objects
- **No removal of `twitch.TwitchChatMessage` NATS queue** until ALL consumers (bots, timers, parser, events) confirmed migrated to generic `ChatMessage`
- **No Kick equivalents of Twitch-only features** in this scope: redemptions, predictions, polls, hype trains, raids — these stay Twitch-only
- **No additional built-in parser variables** beyond `$(platform)` in this scope
- **No 7TV SSE Kick live emote updates** — only `GetProfileByKickId` REST call is in scope
- **No advanced account management UI** — only: Kick login button, linked accounts list (connect/disconnect), platform selector on forms
- **No AI slop patterns**: no excessive comments, no over-abstraction, no generic variable names (data/result/item), no empty error catches

---

## Verification Strategy

> **ZERO HUMAN INTERVENTION** - ALL verification is agent-executed. No exceptions.
> Acceptance criteria requiring "user manually tests/confirms" are FORBIDDEN.

### Test Decision

- **Infrastructure exists**: YES (Go tests, bun test)
- **Automated tests**: Tests-after (no full TDD, but each task includes integration/unit tests where applicable)
- **Framework**: Go `testing` package + bun test for frontend
- **Agent-Executed QA**: ALWAYS mandatory for every task

### QA Policy

Every task MUST include agent-executed QA scenarios.
Evidence saved to `.sisyphus/evidence/task-{N}-{scenario-slug}.{ext}`.

- **Backend HTTP**: `curl` — send requests, assert status + response JSON fields
- **DB state**: `psql` or `postgres_query` — `SELECT` queries validating row counts / field values
- **NATS bus**: Go test publishing + subscribing to verify message flow
- **Frontend/UI**: Playwright — navigate, interact, assert DOM, screenshot

---

## Execution Strategy

### Parallel Execution Waves

```
Wave 1 — DB Foundation (independent migrations, run in parallel):
├── Task 1: users.id → internal UUID migration + user_platform_accounts table [deep]
├── Task 2: kick_bots table + bots table platform column [quick]
├── Task 3: platforms text[] on commands/timers/keywords [quick]
└── Task 4: platform column on channels + one-row-per-platform schema prep [quick]

Wave 2 — Core Interfaces + Entities (after Wave 1):
├── Task 5: PlatformProvider interface + Twitch implementation [unspecified-high]
├── Task 6: Platform entities in libs/entities (PlatformAccount, linked account types) [quick]
├── Task 7: tokens repository: update GetByUserID to use internal UUID [quick]
└── Task 8: user_platform_accounts repository (pgx) [unspecified-high]

Wave 3 — Auth Refactor (after Wave 2):
├── Task 9: Kick PlatformProvider (OAuth 2.1 + PKCE) [deep]
├── Task 10: Generic auth handler POST /auth/:platform/code [deep]
├── Task 11: Session refactor: internal UUID + current_platform (+ session invalidation strategy) [unspecified-high]
└── Task 12: kick_bots repository (pgx) [quick]

Wave 4 — EventSub / Kick Webhooks (after Task 1, independent from Wave 3):
├── Task 13: HTTP server setup in apps/eventsub [unspecified-high]
├── Task 14: Kick webhook HMAC-SHA256 verification middleware [unspecified-high]
├── Task 15: Kick EventSub subscription manager (subscribe/unsubscribe using broadcaster token) [deep]
└── Task 16: Kick bus topics in libs/bus-core/eventsub (KickSubscribeToAllEvents, KickUnsubscribe) [quick]

Wave 5 — Generic Bus + Kick Event Handlers (after Wave 4, Task 6):
├── Task 17: Generic ChatMessage struct in libs/bus-core (with Platform field) [unspecified-high]
├── Task 18: Kick webhook event handlers (chat.message.sent, stream.online, stream.offline, channel.follow) [deep]
├── Task 19: BusListener: subscribe to Kick bus topics + dual-publish Twitch to both queues [unspecified-high]
└── Task 20: Auto-resubscription health-check job for Kick webhooks [unspecified-high]

Wave 6 — Parser + Execution (after Wave 5, Task 3):
├── Task 21: ParseContext.Platform field + sender platform [quick]
├── Task 22: Parser: dual-subscribe generic + twitch queues with Redis MessageID dedup [unspecified-high]
├── Task 23: Built-in variable $(platform) [quick]
└── Task 24: Command/timer/keyword platform filtering in execution [unspecified-high]

Wave 7 — GraphQL + 7TV (after Wave 3, Task 4, independent):
├── Task 25: KickProfile GraphQL type + schema update [unspecified-high]
├── Task 26: AuthenticatedUser.kickProfile + linkedAccounts resolvers [deep]
├── Task 27: GetProfileByKickId in libs/integrations/seventv [quick]
└── Task 28: emotes-cacher: platform-aware 7TV profile lookup [unspecified-high]

Wave 8 — Frontend (after Wave 7, Task 11):
├── Task 29: Kick login button + /auth/kick/callback route (Vue dashboard) [visual-engineering]
├── Task 30: Linked accounts settings section (connect/disconnect) [visual-engineering]
├── Task 31: Platform selector on command/timer/keyword forms [visual-engineering]
└── Task 32: Profile header shows current platform profile (not always Twitch) [visual-engineering]

Wave 9 — Kick Bot Message Sending (after Task 12, Task 18):
├── Task 33: Kick chat client (POST /public/v1/chat with kick_bots token) [deep]
└── Task 34: Bot service: route send-message by platform [unspecified-high]

Wave FINAL — Review (after ALL implementation tasks):
├── Task F1: Plan compliance audit [oracle]
├── Task F2: Code quality review (build + lint + tests) [unspecified-high]
├── Task F3: Real manual QA (all flows) [unspecified-high]
└── Task F4: Scope fidelity check [deep]
```

### Dependency Matrix

- **Task 1**: none → blocks 5,6,7,8,10,11,13,16
- **Task 2**: none → blocks 12,33
- **Task 3**: none → blocks 24
- **Task 4**: none → blocks 10,15,26
- **Task 5**: 1 → blocks 9,10
- **Task 6**: 1 → blocks 8,17,18
- **Task 7**: 1 → blocks 10
- **Task 8**: 1,6 → blocks 10,11,15
- **Task 9**: 5 → blocks 10
- **Task 10**: 4,5,7,8,9 → blocks 11,29
- **Task 11**: 8,10 → blocks 29,30
- **Task 12**: 2 → blocks 33
- **Task 13**: 1 → blocks 14,15,18
- **Task 14**: 13 → blocks 18
- **Task 15**: 4,8,13 → blocks 18,20
- **Task 16**: 1 → blocks 19
- **Task 17**: 6 → blocks 18,19,22
- **Task 18**: 13,14,15,16,17 → blocks 19,34
- **Task 19**: 16,17,18 → blocks 22
- **Task 20**: 15 → none (independent health-check)
- **Task 21**: none → blocks 22,23,24
- **Task 22**: 17,19,21 → blocks 24
- **Task 23**: 21 → blocks 24 (optional)
- **Task 24**: 3,21,22 → none
- **Task 25**: 1 → blocks 26,29
- **Task 26**: 4,11,25 → blocks 30,32
- **Task 27**: none → blocks 28
- **Task 28**: 27 → none
- **Task 29**: 10,11,25 → blocks 30
- **Task 30**: 11,26,29 → none
- **Task 31**: 25 → none
- **Task 32**: 26 → none
- **Task 33**: 12,18 → blocks 34
- **Task 34**: 18,33 → none

### Agent Dispatch Summary

- **Wave 1 (4 tasks)**: T1→`deep`, T2→`quick`, T3→`quick`, T4→`quick`
- **Wave 2 (4 tasks)**: T5→`unspecified-high`, T6→`quick`, T7→`quick`, T8→`unspecified-high`
- **Wave 3 (4 tasks)**: T9→`deep`, T10→`deep`, T11→`unspecified-high`, T12→`quick`
- **Wave 4 (4 tasks)**: T13→`unspecified-high`, T14→`unspecified-high`, T15→`deep`, T16→`quick`
- **Wave 5 (4 tasks)**: T17→`unspecified-high`, T18→`deep`, T19→`unspecified-high`, T20→`unspecified-high`
- **Wave 6 (4 tasks)**: T21→`quick`, T22→`unspecified-high`, T23→`quick`, T24→`unspecified-high`
- **Wave 7 (4 tasks)**: T25→`unspecified-high`, T26→`deep`, T27→`quick`, T28→`unspecified-high`
- **Wave 8 (4 tasks)**: T29→`visual-engineering`, T30→`visual-engineering`, T31→`visual-engineering`, T32→`visual-engineering`
- **Wave 9 (2 tasks)**: T33→`deep`, T34→`unspecified-high`
- **FINAL (4 tasks)**: F1→`oracle`, F2→`unspecified-high`, F3→`unspecified-high`, F4→`deep`

---

## TODOs

- [ ] 1. DB Migration: users.id → Internal UUID + user_platform_accounts table

  **What to do**:
  - Create migration: `bun cli migrations create --name add_user_platform_accounts_and_uuid_migration --db postgres --type sql`
  - In the UP migration (single big-bang transaction):
    1. Add column `internal_id UUID DEFAULT gen_random_uuid()` to `users` table
    2. Backfill `users.internal_id` for all existing rows (each gets a new UUID)
    3. Create `user_platform_accounts` table:
       ```sql
       CREATE TABLE user_platform_accounts (
         id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
         user_id UUID NOT NULL REFERENCES users(internal_id) ON DELETE CASCADE,
         platform TEXT NOT NULL,  -- 'twitch' | 'kick'
         platform_user_id TEXT NOT NULL,
         platform_login TEXT NOT NULL,
         platform_display_name TEXT NOT NULL DEFAULT '',
         platform_avatar TEXT NOT NULL DEFAULT '',  -- profile picture URL; populated at OAuth time from PlatformUser.Avatar
         access_token TEXT NOT NULL,   -- encrypted
         refresh_token TEXT NOT NULL,  -- encrypted
         scopes TEXT[] NOT NULL DEFAULT '{}',
         expires_in INT NOT NULL DEFAULT 0,
         obtainment_timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW(),
         UNIQUE(platform, platform_user_id)
       );
       CREATE INDEX idx_user_platform_accounts_user_id ON user_platform_accounts(user_id);
       CREATE INDEX idx_user_platform_accounts_platform ON user_platform_accounts(platform, platform_user_id);
       ```
    4. Populate `user_platform_accounts` from existing `tokens` + `users` data:
       - For each user in `users` that has a token, insert a row with `platform='twitch'`, `platform_user_id=users.id` (old text ID), and tokens from the `tokens` table
    5. Update ALL FK references from `users.id` (text) to `users.internal_id` (UUID):
       - For each FK table (`channels`, `users_stats`, `users_files`, `users_viewed_notifications`, `channels_permits`, `channels_dashboard_access`, `notifications`, `users_online`, `channels_requested_songs`, `channels_modules_settings`, `channels_messages`, `channels_emotes_usages`, `dudes_user_settings`, `channels_badges`, `channels_eventsub_subscriptions`, `tts_user_settings`, `channels_scheduled_vips`, `shortened_urls`, `channels_giveaways`, `pastebins`, `toxic_messages`, `channels_chat_wall`, `channels_short_links_custom_domains`):
         - Add temp UUID column, populate by JOIN on old text ID, DROP old text FK column, RENAME temp column
    6. Drop `tokens` table FK on `users.id`, update `tokens` to reference `users.internal_id`
    7. Drop old `users.id` TEXT primary key, rename `users.internal_id` → `users.id`, set as PK
    8. Update `channels.id` similarly (it equals old users.id for Twitch channels)
    9. Add Redis session invalidation note: all existing sessions MUST be manually flushed at deploy time (include `FLUSHDB` command in deploy runbook)
  - Create entity `libs/entities/user_platform_account/entity.go` with `UserPlatformAccount` struct
  - Update `libs/gomodels/users.go` internal_id awareness (if needed for compatibility — do NOT use GORM in new code)
  - Verify pre/post row counts match

  **Must NOT do**:
  - Do NOT use GORM in migration code
  - Do NOT leave any table still referencing `users.id` as a TEXT value after migration
  - Do NOT remove the `tokens` table — just repurpose it for bot tokens only (remove the user token rows that have been migrated)

  **Recommended Agent Profile**:
  - **Category**: `deep`
    - Reason: Complex multi-step SQL migration with FK graph restructuring and data migration; requires careful ordering and verification at each step
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES (with T2, T3, T4 — all are independent migrations)
  - **Parallel Group**: Wave 1
  - **Blocks**: Tasks 5, 6, 7, 8, 10, 11, 13, 16, 25
  - **Blocked By**: None (can start immediately)

  **References**:

  **Pattern References**:
  - `libs/migrations/postgres/20260105022225_ulid_to_uuidv7.sql` — example of a UUID migration in this codebase
  - `libs/migrations/postgres/20260105035953_migrate_uuid_v4_to_v7.sql` — another UUID migration pattern
  - `libs/gomodels/users.go` — current `Users` struct (id = Twitch ID string)
  - `libs/gomodels/channels.go` — current `Channels` struct (id = Twitch ID string)
  - `libs/repositories/tokens/repository.go` — `GetByUserID(userID string)` — must be updated after migration

  **API/Type References**:
  - `libs/entities/` — existing entity patterns to follow for `UserPlatformAccount`

  **Acceptance Criteria**:

  **QA Scenarios (MANDATORY)**:

  ```
  Scenario: Migration runs without error and row counts match
    Tool: Bash (psql)
    Preconditions: Database has existing users + tokens rows
    Steps:
      1. Record pre-migration: SELECT COUNT(*) FROM users; SELECT COUNT(*) FROM tokens;
      2. Run migration: bun cli migrations run
      3. Post-migration: SELECT COUNT(*) FROM users; SELECT COUNT(*) FROM user_platform_accounts;
      4. Assert: users count unchanged; user_platform_accounts count = pre-migration users-with-tokens count
      5. Assert: SELECT COUNT(*) FROM users WHERE id::text ~ '^[0-9]+$' = 0 (no Twitch IDs as PK)
    Expected Result: All counts match; no Twitch IDs remain as PKs
    Failure Indicators: Row count mismatch; migration error; FK constraint violation
    Evidence: .sisyphus/evidence/task-1-migration-counts.txt

  Scenario: All FK constraints valid after migration
    Tool: Bash (psql)
    Preconditions: Migration complete
    Steps:
      1. SELECT COUNT(*) FROM information_schema.referential_constraints WHERE constraint_schema='public'
      2. Run: SELECT conname, conrelid::regclass FROM pg_constraint WHERE contype='f' AND NOT convalidated;
    Expected Result: All FK constraints validated; zero unvalidated constraints
    Failure Indicators: Any constraint listed as not validated
    Evidence: .sisyphus/evidence/task-1-fk-validation.txt
  ```

  **Evidence to Capture**:
  - [ ] task-1-migration-counts.txt
  - [ ] task-1-fk-validation.txt

  **Commit**: YES
  - Message: `feat(db): migrate users.id to internal UUID + add user_platform_accounts`
  - Files: `libs/migrations/postgres/{timestamp}_add_user_platform_accounts_and_uuid_migration.sql`, `libs/entities/user_platform_account/entity.go`

- [ ] 2. DB Migration: kick_bots table

  **What to do**:
  - Create migration: `bun cli migrations create --name add_kick_bots_table --db postgres --type sql`
  - Create `kick_bots` table:
    ```sql
    CREATE TABLE kick_bots (
      id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
      type TEXT NOT NULL DEFAULT 'DEFAULT',  -- 'DEFAULT' | 'CUSTOM'
      access_token TEXT NOT NULL,   -- encrypted
      refresh_token TEXT NOT NULL,  -- encrypted
      scopes TEXT[] NOT NULL DEFAULT '{}',
      expires_in INT NOT NULL DEFAULT 0,
      obtainment_timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW(),
      kick_user_id TEXT NOT NULL,
      kick_user_login TEXT NOT NULL,
      created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );
    ```
  - Create entity `libs/entities/kick_bot/entity.go`

  **Must NOT do**:
  - Do NOT reuse the existing `bots` table — separate table for Kick bots

  **Recommended Agent Profile**:
  - **Category**: `quick`
    - Reason: Straightforward new table creation, no FK complexity
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES (Wave 1, with T1, T3, T4)
  - **Blocks**: Tasks 12, 33
  - **Blocked By**: None

  **References**:
  - `libs/gomodels/bots.go` — existing bots model for reference
  - `libs/migrations/postgres/20250802072510_twitch_conduits.sql` — recent migration example

  **Acceptance Criteria**:

  **QA Scenarios**:

  ```
  Scenario: kick_bots table created with correct schema
    Tool: Bash (psql)
    Preconditions: Migration applied
    Steps:
      1. SELECT column_name, data_type FROM information_schema.columns WHERE table_name='kick_bots';
    Expected Result: All 9 columns present with correct types
    Failure Indicators: Missing columns or wrong types
    Evidence: .sisyphus/evidence/task-2-kick-bots-schema.txt
  ```

  **Commit**: YES
  - Message: `feat(db): add kick_bots table`
  - Files: `libs/migrations/postgres/{timestamp}_add_kick_bots_table.sql`, `libs/entities/kick_bot/entity.go`

- [ ] 3. DB Migration: platforms[] on commands/timers/keywords

  **What to do**:
  - Create migration: `bun cli migrations create --name add_platforms_to_commands_timers_keywords --db postgres --type sql`
  - Add column to three tables:
    ```sql
    ALTER TABLE channels_commands ADD COLUMN IF NOT EXISTS platforms TEXT[] NOT NULL DEFAULT '{}';
    ALTER TABLE channels_timers ADD COLUMN IF NOT EXISTS platforms TEXT[] NOT NULL DEFAULT '{}';
    ALTER TABLE channels_keywords ADD COLUMN IF NOT EXISTS platforms TEXT[] NOT NULL DEFAULT '{}';
    ```
  - Update entities in `libs/entities/` for commands, timers, keywords to include `Platforms []string` field
  - Update repositories in `libs/repositories/` for commands/timers/keywords to read/write the new column

  **Must NOT do**:
  - Do NOT change existing command/timer/keyword behavior — empty `platforms` = all platforms (existing rows have `{}` = runs everywhere)

  **Recommended Agent Profile**:
  - **Category**: `quick`
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES (Wave 1, with T1, T2, T4)
  - **Blocks**: Task 24
  - **Blocked By**: None

  **References**:
  - `libs/gomodels/channels_commands.go` — current ChannelsCommands struct
  - `libs/gomodels/channels_timers.go` — current ChannelsTimers struct
  - `libs/repositories/channels/` — pattern for repository updates

  **Acceptance Criteria**:

  **QA Scenarios**:

  ```
  Scenario: platforms column exists with default value
    Tool: Bash (psql)
    Preconditions: Migration applied
    Steps:
      1. SELECT platforms FROM channels_commands LIMIT 1;
      2. INSERT INTO channels_commands(..., platforms) VALUES (..., '{twitch}');
      3. SELECT platforms FROM channels_commands ORDER BY created_at DESC LIMIT 1;
    Expected Result: Default is '{}'; explicit value '{"twitch"}' stored correctly
    Evidence: .sisyphus/evidence/task-3-platforms-column.txt
  ```

  **Commit**: YES
  - Message: `feat(db): add platforms[] to commands/timers/keywords`
  - Files: migration SQL, updated entity + repository files

- [ ] 4. DB Migration: channels multi-platform schema (new PK + platform column)

  **What to do**:

  The current `channels` table has `id TEXT PRIMARY KEY` = Twitch user ID (1:1 with `users`). To support one row per platform per user (Twitch + Kick = 2 rows per user), the schema must change.

  **IMPORTANT: This task depends on T1 completing first** — channels.id must already be migrated to internal UUID before this task runs.
  - Create migration: `bun cli migrations create --name channels_multi_platform --db postgres --type sql`
  - Migration steps (single transaction):
    1. `ALTER TABLE channels RENAME COLUMN id TO channel_internal_id;` — keep existing UUIDs
    2. `ALTER TABLE channels ADD COLUMN id UUID PRIMARY KEY DEFAULT gen_random_uuid();` — new surrogate PK
    3. `ALTER TABLE channels ADD COLUMN platform TEXT NOT NULL DEFAULT 'twitch';`
    4. `ALTER TABLE channels ADD COLUMN user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE;`
    5. Backfill `user_id` from the old `channel_internal_id` (which equals `users.id` after T1 migration): `UPDATE channels SET user_id = channel_internal_id;`
    6. Add unique constraint: `ALTER TABLE channels ADD CONSTRAINT channels_user_platform_unique UNIQUE (user_id, platform);`
    7. Update all FK columns in other tables that reference `channels.id` (old text column) to now reference the new surrogate UUID `channels.id`. These tables include: `channels_commands`, `channels_timers`, `channels_keywords`, `channels_roles`, `channels_alerts`, `channels_permits`, `channels_badges`, `channels_eventsub_subscriptions`, `channels_modules_settings`, `channels_chat_wall`, `channels_scheduled_vips`, `channels_giveaways`, etc. — use the same temp-column-swap pattern as T1.
    8. Drop `channel_internal_id` column after all FK migrations are confirmed.
  - Update `libs/repositories/channels/model/model.go` entity to add `Platform string`, `UserID string` fields.
  - Update channels repository pgx queries: `GetByID` still works; add `GetByUserIDAndPlatform(userID, platform)`.

  **Must NOT do**:
  - Do NOT run this before T1 completes (channels.id must already be UUID).
  - Do NOT add a `platform` column without also adding the `user_id` FK + unique constraint — without the constraint, two rows per user cannot be correctly enforced.
  - Do NOT change the `user_id` approach: channels are associated to users via `user_id UUID REFERENCES users(id)`.
  - Do NOT use GORM in migration or repository code.

  **Recommended Agent Profile**:
  - **Category**: `deep`
    - Reason: Second-stage FK migration with composite unique constraint; must execute after T1; complex ordering.
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: NO — must run after T1 completes (channels.id depends on UUID migration)
  - **Parallel Group**: Wave 1b (sequential after T1, but T2 and T3 can be parallel with T1)
  - **Blocks**: T10 (auth handler needs channels lookup), T15 (EventSub subscription needs channel), T26 (resolver needs channels)
  - **Blocked By**: T1

  **References**:
  - `libs/repositories/channels/model/model.go` — current Channel entity to update
  - `libs/gomodels/channels.go` — existing GORM model (reference only; do NOT use GORM in new code)
  - `libs/migrations/postgres/20260105022225_ulid_to_uuidv7.sql` — FK migration pattern (temp column swap)
  - `apps/api-gql/internal/delivery/gql/resolvers/` — callers of channels repo that need updating

  **Acceptance Criteria**:

  **QA Scenarios**:

  ```
  Scenario: channels table supports two rows for same user (Twitch + Kick)
    Tool: Bash (psql)
    Preconditions: T1 + T4 migrations applied; test user UUID = 'aaaaaaaa-0000-0000-0000-000000000001'
    Steps:
      1. INSERT INTO channels (user_id, platform) VALUES ('aaaaaaaa-0000-0000-0000-000000000001', 'twitch');
      2. INSERT INTO channels (user_id, platform) VALUES ('aaaaaaaa-0000-0000-0000-000000000001', 'kick');
      3. SELECT COUNT(*) FROM channels WHERE user_id = 'aaaaaaaa-0000-0000-0000-000000000001';
      4. Assert COUNT = 2
    Expected Result: Both rows inserted without unique constraint violation
    Evidence: .sisyphus/evidence/task-4-channels-multi-platform.txt

  Scenario: Duplicate platform for same user is rejected
    Tool: Bash (psql)
    Steps:
      1. INSERT INTO channels (user_id, platform) VALUES ('aaaaaaaa-0000-0000-0000-000000000001', 'twitch');
      2. INSERT INTO channels (user_id, platform) VALUES ('aaaaaaaa-0000-0000-0000-000000000001', 'twitch') — duplicate
      3. Assert psql returns: "duplicate key value violates unique constraint channels_user_platform_unique"
    Expected Result: Second insert rejected with unique constraint error
    Evidence: .sisyphus/evidence/task-4-channels-duplicate-rejected.txt
  ```

  **Evidence to Capture**:
  - [ ] task-4-channels-multi-platform.txt
  - [ ] task-4-channels-duplicate-rejected.txt

  **Commit**: YES
  - Message: `feat(db): channels multi-platform schema — new surrogate PK + user_id FK + platform column`
  - Files: migration SQL, `libs/repositories/channels/model/model.go`

- [ ] 5. PlatformProvider Interface + Twitch Implementation

  **What to do**:
  - Create `apps/api-gql/internal/platform/provider.go` with the `PlatformProvider` interface:
    ```go
    type PlatformTokens struct {
        AccessToken  string
        RefreshToken string
        ExpiresIn    int
        Scopes       []string
    }
    type PlatformUser struct {
        ID          string  // platform-specific user ID
        Login       string
        DisplayName string
        Avatar      string
    }
    type PlatformProvider interface {
        Name() string  // "twitch" | "kick"
        GetAuthURL(state, codeChallenge string) string
        ExchangeCode(ctx context.Context, code, codeVerifier string) (*PlatformTokens, error)
        RefreshToken(ctx context.Context, refreshToken string) (*PlatformTokens, error)
        GetUser(ctx context.Context, accessToken string) (*PlatformUser, error)
    }
    ```
  - Create `apps/api-gql/internal/platform/twitch/provider.go` implementing `PlatformProvider` for Twitch:
    - `Name()` → `"twitch"`
    - `GetAuthURL(state, _)` → Twitch OAuth URL (no PKCE for Twitch)
    - `ExchangeCode(ctx, code, _)` → wraps existing `helix.RequestUserAccessToken`
    - `RefreshToken(ctx, refreshToken)` → wraps existing helix refresh logic
    - `GetUser(ctx, accessToken)` → wraps `twitchClient.GetUsers`
  - Register provider in FX DI container

  **Must NOT do**:
  - Do NOT modify the existing `post-code.go` Twitch auth handler yet (done in Task 10)
  - Do NOT use GORM

  **Recommended Agent Profile**:
  - **Category**: `unspecified-high`
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES (Wave 2, with T6, T7, T8)
  - **Blocks**: Tasks 9, 10
  - **Blocked By**: Task 1

  **References**:
  - `apps/api-gql/internal/delivery/http/routes/auth/post-code.go` — existing Twitch auth logic to wrap
  - `libs/twitch/` — existing helix client helpers

  **Acceptance Criteria**:

  **QA Scenarios**:

  ```
  Scenario: Twitch PlatformProvider compiles and Name() returns "twitch"
    Tool: Bash
    Preconditions: Task 1 complete
    Steps:
      1. cd apps/api-gql && go build ./...
    Expected Result: Zero compilation errors; interface satisfied by Twitch implementation
    Evidence: .sisyphus/evidence/task-5-build.txt
  ```

  **Commit**: YES
  - Message: `feat(auth): PlatformProvider interface + Twitch implementation`
  - Files: `apps/api-gql/internal/platform/provider.go`, `apps/api-gql/internal/platform/twitch/provider.go`

- [ ] 6. Platform Entities in libs/entities

  **What to do**:
  - Create `libs/entities/user_platform_account/entity.go`:
    ```go
    type UserPlatformAccount struct {
        ID                  uuid.UUID
        UserID              uuid.UUID
        Platform            string  // "twitch" | "kick"
        PlatformUserID      string
        PlatformLogin       string
        AccessToken         string  // encrypted
        RefreshToken        string  // encrypted
        Scopes              []string
        ExpiresIn           int
        ObtainmentTimestamp time.Time
        isNil               bool
    }
    func (u UserPlatformAccount) IsNil() bool { return u.isNil }
    var Nil = UserPlatformAccount{isNil: true}
    ```
  - Create `libs/entities/kick_bot/entity.go` (if not done in Task 2)

  **Must NOT do**:
  - Do NOT add business logic to entities — data structures only

  **Recommended Agent Profile**:
  - **Category**: `quick`
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES (Wave 2, with T5, T7, T8)
  - **Blocks**: Tasks 8, 17, 18
  - **Blocked By**: Task 1

  **References**:
  - `libs/repositories/channels/model/model.go` — nil pattern to follow
  - `libs/repositories/users/model/model.go` — nil pattern example

  **Acceptance Criteria**:

  **QA Scenarios**:

  ```
  Scenario: Entities compile cleanly
    Tool: Bash
    Steps:
      1. cd libs/entities && go build ./...
    Expected Result: Zero errors
    Evidence: .sisyphus/evidence/task-6-build.txt
  ```

  **Commit**: YES
  - Message: `feat(entities): add platform account entities`
  - Files: `libs/entities/user_platform_account/entity.go`, `libs/entities/kick_bot/entity.go`

- [ ] 7. tokens Repository: Update GetByUserID to Use Internal UUID

  **What to do**:
  - After Task 1 migrates `users.id` to UUID, the `tokens` table now references `users.id` (UUID) instead of the old Twitch ID string
  - Update `libs/repositories/tokens/repository.go`: change `GetByUserID(ctx, userID string)` signature to `GetByUserID(ctx, userID uuid.UUID)`
  - Update the pgx implementation in `libs/repositories/tokens/datasources/pgx.go`
  - Find all callers with `lsp_find_references` on `GetByUserID` and update them to pass `uuid.UUID` instead of string
  - Update `CreateInput.UserID` from `string` to `uuid.UUID`
  - Update `UpdateTokenInput` if needed

  **Must NOT do**:
  - Do NOT use GORM
  - Do NOT break the existing bot token lookup path (`GetByBotID` stays as string if bots still have string IDs)

  **Recommended Agent Profile**:
  - **Category**: `quick`
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES (Wave 2, with T5, T6, T8)
  - **Blocks**: Task 10
  - **Blocked By**: Task 1

  **References**:
  - `libs/repositories/tokens/repository.go` — interface to update
  - `libs/repositories/tokens/datasources/` — pgx implementation
  - `apps/api-gql/internal/delivery/http/routes/auth/post-code.go` — main caller

  **Acceptance Criteria**:

  **QA Scenarios**:

  ```
  Scenario: GetByUserID accepts UUID and returns correct token
    Tool: Bash
    Steps:
      1. go build ./libs/repositories/tokens/...
      2. go vet ./libs/repositories/tokens/...
    Expected Result: Zero errors
    Evidence: .sisyphus/evidence/task-7-tokens-build.txt
  ```

  **Commit**: YES
  - Message: `fix(tokens): update GetByUserID to use internal UUID`
  - Files: `libs/repositories/tokens/repository.go`, `libs/repositories/tokens/datasources/pgx.go`

- [ ] 8. user_platform_accounts Repository (pgx)

  **What to do**:
  - Create `libs/repositories/user_platform_accounts/repository.go` interface:
    ```go
    type Repository interface {
        GetByUserIDAndPlatform(ctx context.Context, userID uuid.UUID, platform string) (entity.UserPlatformAccount, error)
        GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]entity.UserPlatformAccount, error)
        GetByPlatformUserID(ctx context.Context, platform, platformUserID string) (entity.UserPlatformAccount, error)
        Upsert(ctx context.Context, input UpsertInput) (entity.UserPlatformAccount, error)
        Delete(ctx context.Context, id uuid.UUID) error
    }
    type UpsertInput struct {
        UserID              uuid.UUID
        Platform            string
        PlatformUserID      string
        PlatformLogin       string
        AccessToken         string
        RefreshToken        string
        Scopes              []string
        ExpiresIn           int
        ObtainmentTimestamp time.Time
    }
    ```
  - Create `libs/repositories/user_platform_accounts/pgx/pgx.go` with pgx v5 implementation
  - Use `ON CONFLICT (platform, platform_user_id) DO UPDATE` for Upsert
  - Register in FX DI

  **Must NOT do**:
  - Do NOT use GORM

  **Recommended Agent Profile**:
  - **Category**: `unspecified-high`
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES (Wave 2, with T5, T6, T7)
  - **Blocks**: Tasks 10, 11, 15
  - **Blocked By**: Tasks 1, 6

  **References**:
  - `libs/repositories/channels/pgx/pgx.go` — pgx repository pattern to follow
  - `libs/entities/user_platform_account/entity.go` — entity to return (Task 6)

  **Acceptance Criteria**:

  **QA Scenarios**:

  ```
  Scenario: Upsert creates new platform account row
    Tool: Bash (psql + Go test)
    Steps:
      1. Call Upsert with platform='kick', platform_user_id='12345'
      2. SELECT * FROM user_platform_accounts WHERE platform='kick' AND platform_user_id='12345';
    Expected Result: One row returned with all fields matching UpsertInput
    Evidence: .sisyphus/evidence/task-8-upsert.txt

  Scenario: Upsert updates existing row (idempotent)
    Tool: Bash (psql)
    Steps:
      1. Call Upsert twice with same platform+platform_user_id but different access_token
      2. SELECT COUNT(*) FROM user_platform_accounts WHERE platform='kick' AND platform_user_id='12345';
      3. SELECT access_token FROM user_platform_accounts WHERE platform='kick' AND platform_user_id='12345';
    Expected Result: COUNT=1; access_token = second call's value
    Evidence: .sisyphus/evidence/task-8-upsert-idempotent.txt
  ```

  **Commit**: YES
  - Message: `feat(repo): user_platform_accounts repository`
  - Files: `libs/repositories/user_platform_accounts/repository.go`, `libs/repositories/user_platform_accounts/pgx/pgx.go`

- [ ] 9. Kick PlatformProvider (OAuth 2.1 + PKCE)

  **What to do**:
  - Create `apps/api-gql/internal/platform/kick/provider.go` implementing `PlatformProvider` for Kick:
    - `Name()` → `"kick"`
    - `GetAuthURL(state, codeChallenge string)` → `https://id.kick.com/oauth/authorize?client_id=...&redirect_uri=...&response_type=code&scope=user:read+events:subscribe+chat:write&state={state}&code_challenge={codeChallenge}&code_challenge_method=S256`
    - `ExchangeCode(ctx, code, codeVerifier string)` → `POST https://id.kick.com/oauth/token` with `grant_type=authorization_code&code_verifier={codeVerifier}`
    - `RefreshToken(ctx, refreshToken string)` → `POST https://id.kick.com/oauth/token` with `grant_type=refresh_token`
    - `GetUser(ctx, accessToken string)` → `GET https://api.kick.com/public/v1/users` with Bearer token
  - PKCE code verifier generation: implement `generateCodeVerifier()` (random 43-128 char URL-safe string) and `generateCodeChallenge(verifier)` (SHA256 base64url)
  - Add `GET /api/auth/kick/authorize` endpoint that:
    1. Generates `codeVerifier` and `codeChallenge`
    2. Stores `codeVerifier` in SCS session under key `"kick_code_verifier"`
    3. Returns `{"authorize_url": "https://id.kick.com/oauth/authorize?..."}` with `codeChallenge` in the URL
  - Required OAuth scopes: `user:read events:subscribe chat:write channel:read`

  **Must NOT do**:
  - Do NOT store `codeVerifier` in a cookie or URL param — only in the SCS Redis session
  - Do NOT skip PKCE — Kick OAuth 2.1 requires it

  **Recommended Agent Profile**:
  - **Category**: `deep`
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES (Wave 3, with T10, T11, T12)
  - **Blocks**: Task 10
  - **Blocked By**: Task 5

  **References**:
  - `apps/api-gql/internal/platform/provider.go` — interface to implement (Task 5)
  - `apps/api-gql/internal/auth/auth.go` — SCS session manager usage
  - Kick OAuth docs: `https://docs.kick.com/getting-started/authenticating-users`

  **Acceptance Criteria**:

  **QA Scenarios**:

  ```
  Scenario: GET /api/auth/kick/authorize returns valid Kick OAuth URL with code_challenge
    Tool: Bash (curl)
    Preconditions: Kick client_id configured in .env
    Steps:
      1. curl -c cookies.txt -X GET http://localhost:3009/api/auth/kick/authorize -H 'Content-Type: application/json'
      2. Assert HTTP response = 200
      3. Assert response JSON contains key "authorize_url" with value containing "id.kick.com/oauth/authorize"
      4. Assert authorize_url contains "code_challenge=" and "code_challenge_method=S256"
      5. Assert session cookie set (Set-Cookie header present)
    Expected Result: 200 with Kick OAuth URL containing PKCE params
    Evidence: .sisyphus/evidence/task-9-kick-auth-url.txt

  Scenario: Code verifier stored in session after GET /api/auth/kick/authorize
    Tool: Bash (Redis CLI)
    Steps:
      1. curl -c cookies.txt -X GET http://localhost:3009/api/auth/kick/authorize
      2. Extract session ID from Set-Cookie header value
      3. redis-cli GET {session_key} | grep "kick_code_verifier"
      4. Assert kick_code_verifier key exists and is a non-empty string (43+ chars)
    Expected Result: Session contains non-empty code_verifier string
    Failure Indicators: key absent, or empty string value
    Evidence: .sisyphus/evidence/task-9-session-verifier.txt
  ```

  **Commit**: YES
  - Message: `feat(auth): Kick OAuth 2.1 + PKCE provider`
  - Files: `apps/api-gql/internal/platform/kick/provider.go`

- [ ] 10. Generic Auth Handler POST /auth/:platform/code

  **What to do**:
  - Create `apps/api-gql/internal/delivery/http/routes/auth/post-platform-code.go` (new file, do NOT modify existing `post-code.go`)
  - Handler signature: `POST /api/auth/:platform/code` (or `/api/auth/kick/code`)
  - Handler logic (platform-agnostic):
    1. Look up `PlatformProvider` by `platform` parameter
    2. Retrieve `codeVerifier` from SCS session (for Kick; empty string for Twitch)
    3. Call `provider.ExchangeCode(ctx, code, codeVerifier)` → `PlatformTokens`
    4. Call `provider.GetUser(ctx, accessToken)` → `PlatformUser`
    5. Encrypt access + refresh tokens
    6. Look up `user_platform_accounts` by `(platform, platformUser.ID)`
    7. If no existing account: look up or create `users` row (with new UUID PK)
    8. Upsert `user_platform_accounts` row with new tokens **AND** `platform_avatar = platformUser.Avatar`, `platform_display_name = platformUser.DisplayName`, `platform_login = platformUser.Login` (keep profile info fresh on every login)
    9. If new user: create `channels` row for this platform (with `platform` column set)
    10. Publish `scheduler.CreateDefaultRoles` + `CreateDefaultCommands` for new users
    11. Subscribe to platform events (Kick: `KickSubscribeToAllEvents` bus topic)
    12. Store `internal_user_id UUID` + `current_platform` in SCS session
    13. Store platform user data in session (gob-registered type)
    14. Return `redirect_to`
  - Keep old `post-code.go` for Twitch at `/api/auth/code` — do NOT delete it until all clients are migrated to new endpoint
  - Register route with dependency injection

  **Must NOT do**:
  - Do NOT use GORM in new handler — pgx via repositories only
  - Do NOT delete `post-code.go` (backward compat)
  - Do NOT import `libs/gomodels` in new handler

  **Recommended Agent Profile**:
  - **Category**: `deep`
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES (Wave 3, with T9, T11, T12)
  - **Blocks**: Task 11, Task 29
  - **Blocked By**: Tasks 4, 5, 7, 8, 9

  **References**:
  - `apps/api-gql/internal/delivery/http/routes/auth/post-code.go` — existing Twitch handler (reference only, do not modify)
  - `libs/repositories/user_platform_accounts/repository.go` — Upsert (Task 8)
  - `libs/repositories/channels/` — channel creation pattern
  - `libs/bus-core/eventsub/` — EventSub bus topics

  **Acceptance Criteria**:

  **QA Scenarios**:

  ```
  Scenario: POST /api/auth/kick/code with valid code creates user_platform_accounts row
    Tool: Bash (curl + psql)
    Preconditions: Valid Kick OAuth code obtained from test app
    Steps:
      1. curl -X POST http://localhost:3009/api/auth/kick/code -d '{"code":"VALID_CODE","state":"aHR0cDovL2xvY2FsaG9zdA==","code_verifier":"VERIFIER"}'
      2. SELECT * FROM user_platform_accounts WHERE platform='kick';
    Expected Result: 200 with redirect_to; one new row in user_platform_accounts with platform='kick'
    Failure Indicators: 500 error; no new row created
    Evidence: .sisyphus/evidence/task-10-kick-login.txt

  Scenario: POST /api/auth/kick/code is idempotent (second login updates token)
    Tool: Bash (psql)
    Steps:
      1. Login twice with same Kick account
      2. SELECT COUNT(*) FROM user_platform_accounts WHERE platform='kick' AND platform_user_id='{id}';
    Expected Result: COUNT=1 (no duplicates)
    Evidence: .sisyphus/evidence/task-10-kick-login-idempotent.txt
  ```

  **Commit**: YES
  - Message: `feat(auth): generic /auth/:platform/code handler`
  - Files: `apps/api-gql/internal/delivery/http/routes/auth/post-platform-code.go`

- [ ] 11. Session Refactor: Internal UUID + current_platform

  **What to do**:
  - Update `apps/api-gql/internal/auth/sessions_user.go`:
    - Add session keys: `internalUserIdKey = "internalUserId"`, `currentPlatformKey = "currentPlatform"`
    - Add `SetSessionInternalUserID(ctx, id uuid.UUID)` and `GetInternalUserID(ctx) (uuid.UUID, error)`
    - Add `SetSessionCurrentPlatform(ctx, platform string)` and `GetCurrentPlatform(ctx) (string, error)`
    - Update `GetAuthenticatedUserModel(ctx)` to first try `internalUserIdKey` (UUID); fall back to old `dbUserKey` for backward compat during transition
    - Add `gob.Register(uuid.UUID{})` in `NewSessions`
  - Add `gob.Register` for Kick user data struct (define `KickSessionUser{ ID string; Login string; Avatar string }`)
  - Add `SetSessionKickUser` + `GetSessionKickUser` methods
  - **Session invalidation strategy**: document in migration runbook: after deploying Task 1 (UUID migration) + Task 11, run `redis-cli FLUSHDB` on the sessions Redis DB to force all users to re-login. This is mandatory.

  **Must NOT do**:
  - Do NOT remove existing `dbUserKey`/`twitchUserKey` session paths — keep for backward compat temporarily
  - Do NOT break `GetSelectedDashboard`

  **Recommended Agent Profile**:
  - **Category**: `unspecified-high`
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES (Wave 3, with T9, T10, T12)
  - **Blocks**: Tasks 29, 30
  - **Blocked By**: Tasks 8, 10

  **References**:
  - `apps/api-gql/internal/auth/sessions_user.go` — file to modify
  - `apps/api-gql/internal/auth/auth.go` — gob.Register location

  **Acceptance Criteria**:

  **QA Scenarios**:

  ```
  Scenario: New session stores UUID and current platform
    Tool: Bash (Go test)
    Steps:
      1. go test ./apps/api-gql/internal/auth/... -run TestSessionInternalUserID
    Expected Result: PASS
    Evidence: .sisyphus/evidence/task-11-session-test.txt

  Scenario: Existing Twitch session fallback works
    Tool: Bash (curl)
    Steps:
      1. Login via old /api/auth/code (Twitch)
      2. Call GET /api/users/me (authenticated endpoint)
    Expected Result: 200 with user data (fallback to old session path works)
    Evidence: .sisyphus/evidence/task-11-twitch-fallback.txt
  ```

  **Commit**: YES
  - Message: `feat(session): store internal UUID + current platform; add Kick session support`
  - Files: `apps/api-gql/internal/auth/sessions_user.go`, `apps/api-gql/internal/auth/auth.go`

- [ ] 12. kick_bots Repository (pgx)

  **What to do**:
  - Create `libs/repositories/kick_bots/repository.go`:
    ```go
    type Repository interface {
        GetDefault(ctx context.Context) (entity.KickBot, error)
        GetByID(ctx context.Context, id uuid.UUID) (entity.KickBot, error)
        Create(ctx context.Context, input CreateInput) (entity.KickBot, error)
        UpdateToken(ctx context.Context, id uuid.UUID, input UpdateTokenInput) (entity.KickBot, error)
    }
    ```
  - Create `libs/repositories/kick_bots/pgx/pgx.go` implementation using pgx v5
  - `GetDefault` fetches the row where `type = 'DEFAULT'`

  **Must NOT do**:
  - Do NOT use GORM

  **Recommended Agent Profile**:
  - **Category**: `quick`
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES (Wave 3, with T9, T10, T11)
  - **Blocks**: Task 33
  - **Blocked By**: Task 2

  **References**:
  - `libs/repositories/channels/pgx/pgx.go` — pgx repository pattern
  - `libs/entities/kick_bot/entity.go` — entity to return (Task 6)

  **Acceptance Criteria**:

  **QA Scenarios**:

  ```
  Scenario: GetDefault returns the DEFAULT kick bot
    Tool: Bash (psql + Go test)
    Steps:
      1. INSERT INTO kick_bots(type, ...) VALUES ('DEFAULT', ...)
      2. Call GetDefault() → assert returned entity has type='DEFAULT'
    Expected Result: Entity returned with all fields populated
    Evidence: .sisyphus/evidence/task-12-kick-bot-repo.txt
  ```

  **Commit**: YES
  - Message: `feat(repo): kick_bots repository`
  - Files: `libs/repositories/kick_bots/repository.go`, `libs/repositories/kick_bots/pgx/pgx.go`

---

- [ ] 13. HTTP Server Setup in apps/eventsub

  **What to do**:
  - Add an HTTP server to `apps/eventsub` (alongside or wrapping the existing Twitch WebSocket conduit logic). The server must listen on a configurable port (e.g., `EVENTSUB_HTTP_PORT`, default `3030`).
  - Register two route groups: `/webhook/twitch/...` (existing Twitch HTTP callbacks if any) and `/webhook/kick` (new).
  - Use `net/http` or the existing fiber/gin/echo framework — check `apps/eventsub/main.go` to match the existing choice.
  - Wire the HTTP server lifecycle into the app's `fx` or manual lifecycle (start/stop gracefully).
  - No business logic here — just the skeleton: router, port binding, graceful shutdown.

  **Must NOT do**:
  - Do NOT break the existing Twitch WebSocket conduit — it must continue running in parallel.
  - Do NOT introduce a second HTTP framework if one already exists in the app.
  - Do NOT add route handlers here — that is T14's job.

  **Recommended Agent Profile**:
  - **Category**: `quick`
  - **Skills**: []
    - Reason: Pure scaffolding of an HTTP listener in an existing Go app; no domain complexity.

  **Parallelization**:
  - **Can Run In Parallel**: YES (Wave 4, with T14, T15, T16)
  - **Parallel Group**: Wave 4 (with Tasks 14, 15, 16) — all four are independent; T14/T15 need T13 to merge before full integration, but each can be written independently
  - **Blocks**: T14, T15 (need the server to exist to register routes)
  - **Blocked By**: T8 (app structure must be stable), T12 (kick_bots repo must exist to wire later)

  **References**:
  - `apps/eventsub/main.go` — understand existing startup/lifecycle and framework in use
  - `apps/api-gql/main.go` — example of HTTP server wiring pattern in this repo
  - `apps/eventsub/internal/manager/manager.go` — how manager is initialized; must keep working
  - Framework docs (if fiber): `https://docs.gofiber.io/` (use context7 MCP if needed)

  **Acceptance Criteria**:

  **QA Scenarios**:

  ```
  Scenario: HTTP server starts and responds to health check
    Tool: Bash (curl)
    Preconditions: apps/eventsub compiled and running locally
    Steps:
      1. Run: curl -s -o /dev/null -w "%{http_code}" http://localhost:3030/health
      2. Assert response code is 200 or 404 (any non-5xx proves server is up)
    Expected Result: HTTP server responds; process does not crash
    Failure Indicators: connection refused, or existing Twitch conduit logs show errors
    Evidence: .sisyphus/evidence/task-13-http-server-up.txt

  Scenario: Existing Twitch WebSocket conduit still connects
    Tool: Bash (go test / log inspection)
    Preconditions: apps/eventsub running; Twitch env vars set (or mocked)
    Steps:
      1. Tail logs after startup: look for "conduit" or "websocket" connection log line
      2. Assert no "panic" or "fatal" lines in first 10 seconds
    Expected Result: Twitch conduit connects; HTTP server also up simultaneously
    Evidence: .sisyphus/evidence/task-13-twitch-conduit-still-runs.txt
  ```

  **Evidence to Capture**:
  - [ ] task-13-http-server-up.txt
  - [ ] task-13-twitch-conduit-still-runs.txt

  **Commit**: YES
  - Message: `feat(eventsub): add HTTP server skeleton to apps/eventsub`
  - Files: `apps/eventsub/main.go` or new `apps/eventsub/internal/http/server.go`

---

- [ ] 14. Kick Webhook HMAC-SHA256 Verification Middleware

  **What to do**:
  - Implement an HTTP middleware (or handler-level check) for the `/webhook/kick` route that verifies the `Kick-Event-Signature` header.
  - Kick signs webhooks with **RSA-SHA256**: the signature is base64-encoded. The public key is retrieved from `https://api.kick.com/public/v1/public-key` and should be cached (TTL 1 hour in Redis) to avoid hammering Kick's API.
  - Verification steps:
    1. Read the raw request body (do NOT drain it before the handler can read it — use `io.TeeReader` or buffer).
    2. Base64-decode `Kick-Event-Signature` header value.
    3. Fetch/cache Kick's RSA public key.
    4. `rsa.VerifyPKCS1v15(publicKey, crypto/sha256, sha256(body), decodedSig)` — if error → return HTTP 403.
  - Also extract and pass through these Kick headers to downstream handlers: `Kick-Event-Message-Id`, `Kick-Event-Type`, `Kick-Event-Version`, `Kick-Event-Subscription-Id`.
  - Store `Kick-Event-Message-Id` in middleware context for idempotency checks downstream.

  **Must NOT do**:
  - Do NOT hardcode the public key — it must be fetched from Kick's API and cached.
  - Do NOT use HMAC — Kick uses RSA-SHA256, not HMAC-SHA256 (the task name says HMAC for short but the actual algorithm is RSA-SHA256).
  - Do NOT swallow body bytes before the route handler reads them.

  **Recommended Agent Profile**:
  - **Category**: `unspecified-high`
  - **Skills**: []
    - Reason: Cryptographic verification + HTTP middleware + Redis caching — needs careful, correct implementation.

  **Parallelization**:
  - **Can Run In Parallel**: YES (Wave 4, alongside T13, T15, T16)
  - **Parallel Group**: Wave 4
  - **Blocks**: T18 (Kick event handlers need verified route)
  - **Blocked By**: T13 (HTTP server must exist)

  **References**:
  - Kick webhook signature docs: `https://docs.kick.com/events/webhooks` — RSA-SHA256 verification
  - Kick public key endpoint: `https://api.kick.com/public/v1/public-key`
  - `apps/eventsub/internal/handler/handler.go` — existing Twitch verification pattern (for structural reference)
  - `libs/cache/` — existing Redis cache pattern for TTL-based caching
  - Go crypto: `crypto/rsa`, `crypto/sha256`, `encoding/base64`

  **Acceptance Criteria**:

  **QA Scenarios**:

  ```
  Scenario: Request with invalid signature is rejected with 403
    Tool: Bash (curl)
    Preconditions: apps/eventsub HTTP server running (T13 done)
    Steps:
      1. curl -s -o /dev/null -w "%{http_code}" -X POST http://localhost:3030/webhook/kick \
           -H 'Kick-Event-Signature: aW52YWxpZA==' \
           -H 'Kick-Event-Type: chat.message.sent' \
           -H 'Content-Type: application/json' \
           -d '{"test":true}'
      2. Assert HTTP status = 403
    Expected Result: 403 Forbidden
    Failure Indicators: 200 (signature check skipped), 500 (panic), 404 (route missing)
    Evidence: .sisyphus/evidence/task-14-invalid-sig-403.txt

  Scenario: Request with missing signature header is rejected
    Tool: Bash (curl)
    Steps:
      1. curl -X POST http://localhost:3030/webhook/kick -H 'Content-Type: application/json' -d '{}'
      2. Assert HTTP 400 or 403
    Expected Result: non-2xx response
    Evidence: .sisyphus/evidence/task-14-missing-sig-rejected.txt
  ```

  **Evidence to Capture**:
  - [ ] task-14-invalid-sig-403.txt
  - [ ] task-14-missing-sig-rejected.txt

  **Commit**: YES
  - Message: `feat(eventsub): Kick RSA-SHA256 webhook signature verification middleware`
  - Files: `apps/eventsub/internal/kick/middleware.go`

---

- [ ] 15. Kick EventSub Subscription Manager

  **What to do**:
  - Create `apps/eventsub/internal/kick/subscription_manager.go` that manages Kick EventSub subscriptions for each broadcaster.
  - **Subscribe**: `POST https://api.kick.com/public/v1/events/subscriptions` with broadcaster's OAuth token (from `user_platform_accounts`). Subscribe to these event types on connect: `chat.message.sent`, `channel.follow`, `stream.online`, `stream.offline`.
  - **Unsubscribe**: `DELETE https://api.kick.com/public/v1/events/subscriptions` with `subscription_id`.
  - **List**: `GET https://api.kick.com/public/v1/events/subscriptions` to check existing subscriptions.
  - Store `subscription_id` → `channel_id` mapping in Redis (TTL-less, or TTL ~25 hours).
  - Trigger subscriptions: called from T10 (Kick OAuth callback) after tokens are saved.
  - Expose method `SubscribeAll(ctx, kickChannelID, broadcasterToken string) error` for use by auth flow.

  **Must NOT do**:
  - Do NOT use app-level Kick token — use the **broadcaster's own token** from `user_platform_accounts`.
  - Do NOT subscribe to event types not listed above (no scope creep).
  - Do NOT block the OAuth callback — subscription should be async (goroutine with context).

  **Recommended Agent Profile**:
  - **Category**: `unspecified-high`
  - **Skills**: []
    - Reason: External API integration with token management + Redis state + async subscription.

  **Parallelization**:
  - **Can Run In Parallel**: YES (Wave 4, alongside T13, T14, T16)
  - **Parallel Group**: Wave 4
  - **Blocks**: T18 (event handlers need subscriptions to exist), T20 (health check needs manager)
  - **Blocked By**: T8 (user_platform_accounts repo), T9 (Kick OAuth tokens available)

  **References**:
  - Kick EventSub subscribe API: `https://docs.kick.com/events/webhooks` — POST /public/v1/events/subscriptions, required fields: `broadcaster_user_id`, `type`, `method: "webhook"`, `callback_url`
  - `apps/eventsub/internal/manager/subscribe.go` — Twitch subscription pattern (structural reference)
  - `libs/repositories/user_platform_accounts/` (from T8) — to fetch broadcaster's Kick token
  - `libs/cache/` — Redis storage pattern for subscription ID mapping

  **Acceptance Criteria**:

  **QA Scenarios**:

  ```
  Scenario: SubscribeAll returns no error with valid mock broadcaster token
    Tool: Bash (go test)
    Preconditions: Kick API mocked (httptest.NewServer); broadcaster token = "valid-token"; broadcaster_user_id = "12345"
    Steps:
      1. go test ./apps/eventsub/internal/kick/... -run TestSubscribeAll -v
      2. Assert: mock server received POST /public/v1/events/subscriptions for each of: chat.message.sent, channel.follow, stream.online, stream.offline (4 calls total)
      3. Assert: subscription IDs stored in Redis (check mock Redis)
    Expected Result: 4 subscriptions created, no errors
    Evidence: .sisyphus/evidence/task-15-subscribe-all-test.txt

  Scenario: SubscribeAll with expired token returns error (not panic)
    Tool: Bash (go test)
    Steps:
      1. Mock Kick API returns 401
      2. Assert SubscribeAll returns non-nil error with message containing "401" or "unauthorized"
      3. Assert no goroutine leak
    Evidence: .sisyphus/evidence/task-15-subscribe-expired-token-error.txt
  ```

  **Evidence to Capture**:
  - [ ] task-15-subscribe-all-test.txt
  - [ ] task-15-subscribe-expired-token-error.txt

  **Commit**: YES
  - Message: `feat(eventsub): Kick EventSub subscription manager`
  - Files: `apps/eventsub/internal/kick/subscription_manager.go`, `apps/eventsub/internal/kick/subscription_manager_test.go`

---

- [ ] 16. Kick EventSub Bus Topics in libs/bus-core

  **What to do**:
  - Add new NATS subjects for Kick-specific events in `libs/bus-core`. Mirror the structure of existing Twitch topics.
  - Create `libs/bus-core/kick/` package with:
    - `chat-message.go`: subject `kick.chat-message`, struct `KickChatMessage` (minimal: ChannelID, SenderID, SenderLogin, MessageID, Text, Badges, Color)
    - `follow.go`: subject `kick.channel-follow`, struct `KickChannelFollow`
    - `stream-online.go` / `stream-offline.go`: subjects `kick.stream-online` / `kick.stream-offline`
  - Keep as thin typed wrappers — no business logic in bus-core.
  - These topics are the **Kick-specific raw event bus**. The generic `ChatMessage` topic (T17) is separate.

  **Must NOT do**:
  - Do NOT modify existing `libs/bus-core/twitch/` package — Kick topics live in a new `kick/` sub-package.
  - Do NOT add business logic or transformations here.
  - Do NOT remove or rename existing Twitch NATS subjects.

  **Recommended Agent Profile**:
  - **Category**: `quick`
  - **Skills**: []
    - Reason: Pure struct definitions + NATS subject constants; trivial scaffolding.

  **Parallelization**:
  - **Can Run In Parallel**: YES (Wave 4, alongside T13, T14, T15)
  - **Parallel Group**: Wave 4
  - **Blocks**: T18 (Kick handlers publish to these topics), T19 (BusListener subscribes)
  - **Blocked By**: None (fully independent)

  **References**:
  - `libs/bus-core/twitch/chat-message.go` — exact pattern to mirror for Kick
  - `libs/bus-core/twitch/` directory — naming and subject string conventions
  - `libs/bus-core/` package structure — how topics are exported

  **Acceptance Criteria**:

  **QA Scenarios**:

  ```
  Scenario: Kick bus-core compiles without errors
    Tool: Bash
    Steps:
      1. cd libs/bus-core && go build ./kick/...
      2. Assert exit code = 0
    Expected Result: Compiles clean
    Evidence: .sisyphus/evidence/task-16-bus-core-build.txt

  Scenario: NATS subject constants are correct strings
    Tool: Bash (go test)
    Steps:
      1. go test ./libs/bus-core/kick/... -run TestSubjectConstants -v
      2. Assert kick.KickChatMessageSubject == "kick.chat-message"
      3. Assert kick.KickChannelFollowSubject == "kick.channel-follow"
    Expected Result: All subject strings match expected values
    Evidence: .sisyphus/evidence/task-16-subject-constants-test.txt
  ```

  **Evidence to Capture**:
  - [ ] task-16-bus-core-build.txt
  - [ ] task-16-subject-constants-test.txt

  **Commit**: YES
  - Message: `feat(bus): Kick EventSub bus topics in libs/bus-core`
  - Files: `libs/bus-core/kick/chat-message.go`, `libs/bus-core/kick/follow.go`, `libs/bus-core/kick/stream-online.go`, `libs/bus-core/kick/stream-offline.go`

---

- [ ] 17. Generic ChatMessage Struct + Bus Fields in libs/bus-core

  **What to do**:
  - Create `libs/bus-core/generic/chat-message.go` with a `ChatMessage` struct:
    ```go
    type ChatMessage struct {
      Platform          string  // "twitch" | "kick"
      ChannelID         string  // channels.id surrogate UUID (T4 PK — used for all channel-scoped repo lookups in parser)
      UserID            string  // users.id internal UUID (T1 PK)
      PlatformChannelID string  // platform-native broadcaster ID (Twitch broadcaster_user_id or Kick broadcaster_user_id)
      SenderID          string
      SenderLogin       string
      SenderDisplayName string
      MessageID         string
      Text              string
      Badges            []Badge
      Color             string
      Emotes            []Emote
    }
    ```
  - `ChannelID` MUST be `channels.id` (the new surrogate UUID from T4), **not** `users.id`. This is the key used by ALL channel-scoped repository lookups in the parser (commands, timers, keywords — all keyed by `channel_id`). `UserID` carries `users.id` for user-scoped operations.
  - Add two new Bus fields to `libs/bus-core/bus.go`:
    1. `ChatMessagesGeneric` (field) — type `Queue[generic.ChatMessage, struct{}]`, subject `"chat.messages.generic"` — replaces `ChatMessages` for platform-agnostic consumers
    2. In the `Parser` sub-struct: add `ProcessGenericMessage` — type `Queue[generic.ChatMessage, struct{}]`, subject `"parser.process_generic_message"` — the new fire-and-forget subject for parser command execution from any platform
  - Do NOT remove or modify the existing `bus.ChatMessages` field (`"chat.messages"`, type `twitch.TwitchChatMessage`) — Twitch consumers still use it.
  - Do NOT remove `bus.Parser.ProcessMessageAsCommand` (`"parser.process_message_as_command"`) — existing Twitch parser subscription stays.

  **Must NOT do**:
  - Do NOT change the type of `bus.ChatMessages` from `twitch.TwitchChatMessage` — that would break existing code.
  - Do NOT remove `libs/bus-core/twitch/TwitchChatMessage` — it must remain until ALL consumers migrated.
  - Do NOT add platform-specific fields to the generic struct — use the `Platform` string discriminator.
  - Do NOT add business logic here.

  **Recommended Agent Profile**:
  - **Category**: `quick`
  - **Skills**: []
    - Reason: Struct definition + two Bus field additions; mechanical.

  **Parallelization**:
  - **Can Run In Parallel**: YES — can be done alongside T13-T16 (Wave 4) or as first item of Wave 5
  - **Parallel Group**: Wave 5 start (T18, T19, T22 all depend on this)
  - **Blocks**: T18, T19, T22, T24
  - **Blocked By**: None (fully independent)

  **References**:
  - `libs/bus-core/bus.go` — MUST READ: existing Bus struct to understand where to add new fields (ChatMessages, Parser sub-struct pattern). This is the authoritative wiring file.
  - `libs/bus-core/twitch/chat-message.go` — existing struct to draw field names from
  - `libs/bus-core/parser/commands.go` — existing Parser sub-struct fields for reference
  - `libs/bus-core/twitch/` — naming and package structure to mirror

  **Acceptance Criteria**:
  - [ ] `libs/bus-core/generic/chat-message.go` exists with all fields above
  - [ ] `bus.ChatMessagesGeneric` field exists in `libs/bus-core/bus.go`
  - [ ] `bus.Parser.ProcessGenericMessage` field exists in `libs/bus-core/bus.go`
  - [ ] `go build ./libs/bus-core/...` exits 0

  **QA Scenarios**:

  ```
  Scenario: Generic chat-message package compiles
    Tool: Bash
    Steps:
      1. go build ./libs/bus-core/generic/...
      2. Assert exit code = 0
    Expected Result: Clean build
    Evidence: .sisyphus/evidence/task-17-generic-bus-build.txt

  Scenario: New Bus fields exist and compile in bus.go
    Tool: Bash
    Steps:
      1. grep -n "ChatMessagesGeneric\|ProcessGenericMessage" libs/bus-core/bus.go
      2. Assert both names appear
      3. go build ./libs/bus-core/...
      4. Assert exit code = 0
    Expected Result: Both fields compile; no regression in bus package
    Evidence: .sisyphus/evidence/task-17-bus-fields.txt
  ```

  **Evidence to Capture**:
  - [ ] task-17-generic-bus-build.txt
  - [ ] task-17-bus-fields.txt

  **Commit**: YES
  - Message: `feat(bus): generic ChatMessage struct + bus fields for platform-agnostic processing`
  - Files: `libs/bus-core/generic/chat-message.go`, `libs/bus-core/bus.go`

---

- [ ] 18. Kick Webhook Event Handlers

  **What to do**:
  - Create `apps/eventsub/internal/kick/handlers.go` with HTTP handlers for each Kick event type. Route is `/webhook/kick` (single endpoint); event type is discriminated by `Kick-Event-Type` header (set by T14 middleware in request context).
  - Handle these event types:
    - `chat.message.sent` → parse body → build `generic.ChatMessage{Platform: "kick", ...}` → publish to BOTH:
      - `bus.ChatMessagesGeneric.Publish(ctx, genericMsg)` — broadcasts the message on the generic chat topic
      - `bus.Parser.ProcessGenericMessage.Publish(ctx, genericMsg)` — triggers parser command execution (fire-and-forget, mirrors how Twitch uses `bus.Parser.ProcessMessageAsCommand`)
    - `channel.follow` → parse body → publish to `bus.Events.ChannelFollow.Publish(ctx, ...)` (using existing events bus topics where they exist, or new Kick-specific topics from T16 where they don't)
    - `stream.online` → publish to appropriate events bus field
    - `stream.offline` → publish to appropriate events bus field
  - **Idempotency**: before processing, check `Kick-Event-Message-Id` (from middleware context) against Redis cache (TTL 10 min). If already seen → return 200 immediately (Kick requires 200 ack).
  - Always return HTTP 200 to Kick (even on internal errors — log the error but ack).
  - Resolve internal IDs from Kick `broadcaster_user_id` using a two-step lookup:
    1. Fetch `user_platform_accounts` row by `(platform='kick', platform_user_id=broadcaster_user_id)` → get `user_id` (internal UUID = `users.id`)
    2. Call `channels.GetByUserIDAndPlatform(userID, "kick")` → get `channels.id` (surrogate UUID from T4)
    3. Set `generic.ChatMessage.ChannelID = channels.id` and `generic.ChatMessage.UserID = user_id`

  **Must NOT do**:
  - Do NOT skip idempotency check — duplicate events will cause duplicate chat responses.
  - Do NOT return non-200 to Kick on business logic errors — Kick will retry → infinite loop.
  - Do NOT block on NATS publish — use async with deadline context.
  - Do NOT import `libs/gomodels`.

  **Recommended Agent Profile**:
  - **Category**: `unspecified-high`
  - **Skills**: []
    - Reason: Multi-concern handler: HTTP parsing + NATS publishing + Redis idempotency + UUID lookup.

  **Parallelization**:
  - **Can Run In Parallel**: NO — depends on T14 (middleware), T16 (Kick bus topics), T17 (generic ChatMessage)
  - **Parallel Group**: Wave 5 (after T14, T16, T17 complete)
  - **Blocks**: T19 (BusListener subscribes to topics published here)
  - **Blocked By**: T14, T16, T17

  **References**:
  - `apps/eventsub/internal/handler/chat_message.go` — Twitch event handler pattern
  - `apps/eventsub/internal/handler/handler.go` — how handlers are registered and dispatched
  - `libs/bus-core/kick/` (T16) — Kick-specific topics to publish to
  - `libs/bus-core/generic/chat-message.go` (T17) — generic struct to build; note `ChannelID` = `channels.id` surrogate UUID
  - `libs/repositories/user_platform_accounts/` (T8) — step 1: look up `user_id` from `(platform, platform_user_id)`
  - `libs/repositories/channels/` (T4) — step 2: call `GetByUserIDAndPlatform("kick")` to get `channels.id`
  - `libs/cache/` — Redis idempotency check pattern
  - Kick webhook payload schemas: `https://docs.kick.com/events/event-types`

  **Acceptance Criteria**:

  **QA Scenarios**:

  ```
  Scenario: chat.message.sent event is published to NATS
    Tool: Bash (go test)
    Preconditions: NATS test server running; Redis mock; user_platform_accounts mock returns UUID "uuid-abc"
    Steps:
      1. go test ./apps/eventsub/internal/kick/... -run TestChatMessageHandler -v
      2. POST /webhook/kick with headers Kick-Event-Type: chat.message.sent, Kick-Event-Message-Id: msg-001
         body: {"broadcaster_user_id": "12345", "sender": {"user_id": "67890", "username": "user1"}, "content": "hello"}
      3. Assert: NATS subject "generic.chat-message" received exactly 1 message
      4. Assert: published ChatMessage.Platform == "kick", Text == "hello", ChannelID == "uuid-abc"
      5. Assert: HTTP response = 200
    Expected Result: Message published, 200 returned
    Evidence: .sisyphus/evidence/task-18-chat-message-handler.txt

  Scenario: Duplicate event (same Kick-Event-Message-Id) is acknowledged but NOT re-published
    Tool: Bash (go test)
    Steps:
      1. Send same event twice with Kick-Event-Message-Id: msg-001
      2. Assert NATS receives exactly 1 message (second is deduplicated)
      3. Assert both HTTP responses = 200
    Expected Result: Idempotency works; NATS gets 1 message only
    Evidence: .sisyphus/evidence/task-18-idempotency.txt
  ```

  **Evidence to Capture**:
  - [ ] task-18-chat-message-handler.txt
  - [ ] task-18-idempotency.txt

  **Commit**: YES
  - Message: `feat(eventsub): Kick webhook event handlers with idempotency`
  - Files: `apps/eventsub/internal/kick/handlers.go`, `apps/eventsub/internal/kick/handlers_test.go`

---

- [ ] 19. BusListener Kick Support + Twitch Dual-Publish

  **What to do**:
  - In `apps/eventsub/internal/bus-listener/bus-listener.go` (existing file), add support for Kick channel subscriptions triggered via existing EventSub bus topics. Specifically:
    - Kick subscriptions should be triggered by `bus.EventSub.SubscribeToAllEvents` (existing subject `"eventsub.subscribeAll"`) — extend the existing `subscribeToAllEvents` handler to also call `KickSubscriptionManager.SubscribeAll(...)` if the request includes a Kick channel. Alternatively, add a new dedicated handler if the request DTO needs a platform discriminator.
    - Add `bus.EventSub.Unsubscribe` handler extension for Kick if needed (subject `"eventsub.unsubscribe"`).
    - Reuse the existing eventsub bus field pattern in `libs/bus-core/eventsub/eventsub.go` — do NOT invent new NATS subjects for this control plane; use `"eventsub.subscribeAll"` and `"eventsub.unsubscribe"`.
  - Update the **Twitch `chat.message` handler** in `apps/eventsub/internal/handler/chat_message.go` to dual-publish:
    - Keep existing: `c.twirBus.ChatMessages.Publish(ctx, twitchMsg)` (subject `"chat.messages"`) — unchanged
    - Keep existing: `c.twirBus.Parser.ProcessMessageAsCommand.Publish(ctx, twitchMsg)` — unchanged
    - ADD: `c.twirBus.ChatMessagesGeneric.Publish(ctx, genericMsg)` — new bus field from T17 (subject `"chat.messages.generic"`)
    - ADD: `c.twirBus.Parser.ProcessGenericMessage.Publish(ctx, genericMsg)` — new bus field from T17 (subject `"parser.process_generic_message"`)
    - The `genericMsg` is built by mapping `twitch.TwitchChatMessage` → `generic.ChatMessage{Platform: "twitch", ...}`.
  - This dual-publish is temporary scaffolding — parser will eventually drop its `ProcessMessageAsCommand` subscriber and use only `ProcessGenericMessage`.
  - The `ChannelID` field in `genericMsg` for Twitch messages must be `channels.id` (the T4 surrogate UUID), NOT the Twitch broadcaster ID. Look up `channels.id` via `channels.GetByUserIDAndPlatform(userID, "twitch")` where `userID` = internal `users.id` from session/event.

  **Must NOT do**:
  - Do NOT remove publishing to `bus.ChatMessages` — it must continue (existing consumers still use it).
  - Do NOT remove `bus.Parser.ProcessMessageAsCommand.Publish` — it must continue.
  - Do NOT invent new NATS subject strings for the EventSub control plane — use the existing `"eventsub.subscribeAll"` and `"eventsub.unsubscribe"` subjects from `libs/bus-core/eventsub/eventsub.go`.
  - Do NOT rewrite the bus-listener — extend it minimally.
  - Do NOT change the signature of existing BusListener struct methods.

  **Recommended Agent Profile**:
  - **Category**: `unspecified-high`
  - **Skills**: []
    - Reason: Modifying a critical existing component (bus-listener) requires care to avoid regressions.

  **Parallelization**:
  - **Can Run In Parallel**: NO — depends on T15 (Kick subscription manager), T17 (generic ChatMessage)
  - **Parallel Group**: Wave 5 (after T15, T17)
  - **Blocks**: T22 (parser can start consuming generic.chat-message)
  - **Blocked By**: T15, T17

  **References**:
  - `apps/eventsub/internal/bus-listener/bus-listener.go` — MUST READ; existing NATS subscription patterns
  - `apps/eventsub/internal/handler/chat_message.go` — Twitch chat message shape (for dual-publish mapping)
  - `libs/bus-core/generic/chat-message.go` (T17) — target struct for dual-publish
  - `libs/bus-core/twitch/chat-message.go` — source struct for mapping

  **Acceptance Criteria**:

  **QA Scenarios**:

  ```
  Scenario: Twitch chat message is published to BOTH twitch.chat-message AND generic.chat-message
    Tool: Bash (go test)
    Preconditions: NATS test server; mock Twitch chat message event
    Steps:
      1. go test ./apps/eventsub/internal/bus-listener/... -run TestDualPublish -v
      2. Trigger a Twitch chat message event
      3. Assert NATS topic "twitch.chat-message" received 1 message
      4. Assert NATS topic "generic.chat-message" received 1 message
      5. Assert generic message Platform == "twitch"
    Expected Result: Both topics receive the message
    Evidence: .sisyphus/evidence/task-19-dual-publish.txt

  Scenario: kick.subscribe-channel NATS message triggers Kick subscription
    Tool: Bash (go test)
    Steps:
      1. Publish to "kick.subscribe-channel" with {channel_id: "uuid-abc", kick_broadcaster_id: "12345"}
      2. Assert KickSubscriptionManager.SubscribeAll was called with "12345"
    Expected Result: Subscription triggered
    Evidence: .sisyphus/evidence/task-19-kick-subscribe-trigger.txt
  ```

  **Evidence to Capture**:
  - [ ] task-19-dual-publish.txt
  - [ ] task-19-kick-subscribe-trigger.txt

  **Commit**: YES
  - Message: `feat(eventsub): BusListener Kick support + Twitch dual-publish to generic queue`
  - Files: `apps/eventsub/internal/bus-listener/bus-listener.go`

---

- [ ] 20. Kick Webhook Auto-Resubscription Health-Check Job

  **What to do**:
  - Create `apps/eventsub/internal/kick/resubscribe_job.go` — a background goroutine (or ticker) that runs every **23 hours** (Kick auto-unsubscribes after 24 hours of failure/expiry).
  - For each Kick channel stored in `user_platform_accounts` (platform = "kick"), call `GET /public/v1/events/subscriptions` using broadcaster's token.
  - If any of the 4 required subscriptions (chat.message.sent, channel.follow, stream.online, stream.offline) is missing → call `SubscribeAll(...)` to re-register.
  - Log: `"re-subscribed kick eventsub for channel %s"` on re-registration.
  - Graceful shutdown: stop ticker on context cancellation.

  **Must NOT do**:
  - Do NOT re-subscribe all channels every tick unconditionally — only re-subscribe if subscription is actually missing (avoid unnecessary API calls).
  - Do NOT hardcode the tick interval — read from config (default 23h).
  - Do NOT use `context.TODO()`.

  **Recommended Agent Profile**:
  - **Category**: `unspecified-high`
  - **Skills**: []
    - Reason: Background job with external API calls + conditional logic + lifecycle management.

  **Parallelization**:
  - **Can Run In Parallel**: NO — depends on T15 (KickSubscriptionManager) and T8 (user_platform_accounts repo)
  - **Parallel Group**: Wave 5 (can be done in parallel with T18, T19 if T15 and T8 are done)
  - **Blocks**: Nothing (terminal in its wave)
  - **Blocked By**: T15, T8

  **References**:
  - `apps/eventsub/internal/kick/subscription_manager.go` (T15) — SubscribeAll + ListSubscriptions
  - `libs/repositories/user_platform_accounts/` (T8) — list all Kick channels
  - `libs/config/` — pattern for reading config values in this app
  - Example ticker pattern: `apps/scheduler/` — how tickers are managed with context

  **Acceptance Criteria**:

  **QA Scenarios**:

  ```
  Scenario: Missing subscription triggers re-registration
    Tool: Bash (go test)
    Preconditions: Mock Kick API returns only 2 of 4 subscriptions for channel "12345"; mock SubscribeAll records calls
    Steps:
      1. go test ./apps/eventsub/internal/kick/... -run TestResubscribeJob -v
      2. Trigger one tick of the job
      3. Assert SubscribeAll was called for channel "12345"
      4. Assert log contains "re-subscribed kick eventsub"
    Expected Result: Re-subscription triggered for the channel with missing subscriptions
    Evidence: .sisyphus/evidence/task-20-resubscribe-triggered.txt

  Scenario: All subscriptions present → no re-subscription
    Tool: Bash (go test)
    Steps:
      1. Mock Kick API returns all 4 subscriptions for all channels
      2. Trigger one tick
      3. Assert SubscribeAll was NOT called
    Expected Result: No unnecessary re-subscription
    Evidence: .sisyphus/evidence/task-20-no-unnecessary-resubscribe.txt
  ```

  **Evidence to Capture**:
  - [ ] task-20-resubscribe-triggered.txt
  - [ ] task-20-no-unnecessary-resubscribe.txt

  **Commit**: YES
  - Message: `feat(eventsub): Kick webhook auto-resubscription health-check job`
  - Files: `apps/eventsub/internal/kick/resubscribe_job.go`, `apps/eventsub/internal/kick/resubscribe_job_test.go`

---

- [ ] 21. ParseContext.Platform Field + Sender Platform

  **What to do**:
  - In `apps/parser/internal/types/parse_context.go`, add a `Platform string` field to `ParseContext` (and any embedded sender/channel structs) so the parser knows which platform originated the message.
  - Platform value: `"twitch"` or `"kick"` (matching the `generic.ChatMessage.Platform` string).
  - Update all constructors/builders that create `ParseContext` to propagate the platform field.
  - The parser still receives messages via NATS — when it receives from `generic.chat-message`, it sets `Platform` from the message. When it receives from `twitch.chat-message` (old queue, still active), it sets `Platform = "twitch"`.

  **Must NOT do**:
  - Do NOT remove the old `twitch.chat-message` subscription from parser yet (T22 handles migration).
  - Do NOT change the `ParseContext` interface in a breaking way — add the field, keep existing fields.
  - Do NOT import `libs/gomodels`.

  **Recommended Agent Profile**:
  - **Category**: `quick`
  - **Skills**: []
    - Reason: Simple struct field addition + propagation; mechanical change.

  **Parallelization**:
  - **Can Run In Parallel**: YES — can start as soon as T17 (generic ChatMessage) is done
  - **Parallel Group**: Wave 6 (with T22, T23, T24)
  - **Blocks**: T22, T23, T24 (all need Platform in ParseContext)
  - **Blocked By**: T17

  **References**:
  - `apps/parser/internal/types/parse_context.go` — MUST READ; existing ParseContext fields
  - `apps/parser/internal/` — all files that construct ParseContext (grep for `ParseContext{`)
  - `libs/bus-core/generic/chat-message.go` (T17) — Platform field source

  **Acceptance Criteria**:

  **QA Scenarios**:

  ```
  Scenario: ParseContext built from generic ChatMessage has correct Platform
    Tool: Bash (go test)
    Steps:
      1. go test ./apps/parser/... -run TestParseContextPlatform -v
      2. Build ParseContext from ChatMessage{Platform: "kick", ...}
      3. Assert ctx.Platform == "kick"
    Expected Result: Platform field correctly propagated
    Evidence: .sisyphus/evidence/task-21-parse-context-platform.txt

  Scenario: Parser compiles without errors after adding Platform field
    Tool: Bash
    Steps:
      1. cd apps/parser && go build ./...
      2. Assert exit code = 0
    Expected Result: Clean build
    Evidence: .sisyphus/evidence/task-21-parser-build.txt
  ```

  **Evidence to Capture**:
  - [ ] task-21-parse-context-platform.txt
  - [ ] task-21-parser-build.txt

  **Commit**: YES
  - Message: `feat(parser): add Platform field to ParseContext`
  - Files: `apps/parser/internal/types/parse_context.go`, affected constructor files

---

- [ ] 22. Parser: Dual-Subscribe to twitch.chat-message AND generic.chat-message with MessageID Deduplication

  **What to do**:
  - In `apps/parser/internal/commands-bus/commands-bus.go`, add a **second** NATS subscription to `generic.chat-message` (in addition to the existing `twitch.chat-message` subscription — keep both active from day one).
  - The `generic.chat-message` handler reads `Platform`, `MessageID`, `ChannelID`, `UserID`, `Text` etc. from `libs/bus-core/generic/chat-message.go` (T17), constructs `ParseContext` (T21) with `Platform` set correctly, and feeds into the same execution pipeline.
  - The `twitch.chat-message` handler remains untouched — but now that T19 dual-publishes Twitch messages to BOTH `twitch.chat-message` AND `generic.chat-message`, Twitch messages will arrive on both subscriptions. **Prevent double-processing** using Redis deduplication:
    - Before processing any message, compute the Redis key `parser:dedup:{MessageID}` and call `SET NX EX 60`.
    - If the key already exists → **skip silently** (already processed by the other subscription).
    - If the key was just set (NX succeeded) → **process normally**.
    - This dedup key TTL of 60 seconds is sufficient because duplicate delivery always arrives within milliseconds.
  - Kick messages arrive **only** on `generic.chat-message` (T16/T18 never publish to `twitch.chat-message`), so they are always processed exactly once — no dedup key collision possible.
  - **No feature flags.** Kick commands/timers/keywords work in production from the moment this task is deployed (after T18 and T19 are also deployed).
  - Add a `// TODO(Phase-2): remove twitch.chat-message subscription once all consumers migrated off it` comment on the old Twitch handler so it is discoverable but not forgotten.

  **Must NOT do**:
  - Do NOT introduce a feature flag (`PARSER_USE_GENERIC_QUEUE` or any variant) — this was the old flawed approach. Remove any trace if it exists.
  - Do NOT remove the `twitch.chat-message` subscription in this task — it stays for backward compatibility.
  - Do NOT change the command execution pipeline itself — only the input source is changing.
  - Do NOT use an in-memory map for dedup — it won't work across multiple parser replicas. Use Redis only.
  - Do NOT set the dedup TTL to less than 10 seconds — duplicates arrive within milliseconds but we want headroom.

  **Recommended Agent Profile**:
  - **Category**: `unspecified-high`
  - **Skills**: []
    - Reason: Critical migration path; Redis dedup logic; touching the main execution pipeline entrypoint; must not introduce double-processing in production.
  - **Skills Evaluated but Omitted**:
    - `deep`: Overkill for a well-scoped subscription addition.

  **Parallelization**:
  - **Can Run In Parallel**: NO — depends on T19 (Twitch dual-publish must be in place) and T21 (Platform field in ParseContext)
  - **Parallel Group**: Wave 6 (after T19, T21)
  - **Blocks**: T23, T24 ($(platform) variable and platform-filtering both need parser to process generic messages)
  - **Blocked By**: T19, T21

  **References**:

  **Pattern References** (existing code to follow):
  - `apps/parser/internal/commands-bus/commands-bus.go` — existing `twitch.chat-message` NATS subscriber; replicate its structure for the second subscriber
  - `apps/parser/internal/types/parse_context.go` (T21) — `ParseContext` with Platform field

  **API/Type References** (contracts to implement against):
  - `libs/bus-core/generic/chat-message.go` (T17) — `GenericChatMessage{MessageID, Platform, ChannelID, UserID, Text, ...}`
  - `libs/bus-core/twitch/chat-message.go` — existing Twitch message type for comparison
  - `libs/cache/` or existing Redis client in `apps/parser` — Redis `SET NX EX` usage pattern

  **External References**:
  - Redis `SET NX EX` semantics — key expires automatically; NX means "only set if not exists"; returns OK if set, nil if already existed

  **WHY Each Reference Matters**:
  - `commands-bus.go` — this is the exact file being modified; the new subscription must match the existing handler structure
  - `GenericChatMessage` — the new message schema the handler must unmarshal; `MessageID` is the dedup key
  - Redis NX pattern — prevents double-processing across parser replicas when Twitch publishes to both topics

  **Acceptance Criteria**:
  - [ ] `apps/parser/internal/commands-bus/commands-bus.go` subscribes to both `twitch.chat-message` AND `generic.chat-message`
  - [ ] Redis dedup key `parser:dedup:{MessageID}` is set with `NX EX 60` before each message is processed
  - [ ] No `PARSER_USE_GENERIC_QUEUE` env var or feature flag exists anywhere in the parser codebase
  - [ ] `go build ./apps/parser/...` succeeds with zero errors

  **QA Scenarios**:

  ```
  Scenario: Kick chat message on generic.chat-message triggers command execution exactly once
    Tool: Bash (go test)
    Preconditions: NATS test server; Redis test instance; T17 GenericChatMessage type available
    Steps:
      1. go test ./apps/parser/internal/commands-bus/... -run TestGenericQueueKickMessage -v
      2. Publish to "generic.chat-message": {MessageID: "kick-msg-001", Platform: "kick", Text: "!hello", ChannelID: "uuid-abc"}
      3. Assert command execution pipeline called once with ParseContext{Platform: "kick", MessageID: "kick-msg-001"}
      4. Assert Redis key "parser:dedup:kick-msg-001" exists with TTL between 1 and 60 seconds
    Expected Result: Kick message processed exactly once; dedup key stored in Redis
    Failure Indicators: Pipeline not called, or called more than once, or Redis key absent
    Evidence: .sisyphus/evidence/task-22-kick-message-parsed.txt

  Scenario: Twitch message dual-published to both topics is processed exactly once (dedup works)
    Tool: Bash (go test)
    Preconditions: NATS test server; Redis test instance (fresh, no pre-existing dedup keys)
    Steps:
      1. go test ./apps/parser/internal/commands-bus/... -run TestDedupTwitchDualPublish -v
      2. Publish same message to BOTH "twitch.chat-message" AND "generic.chat-message":
         MessageID="twitch-msg-002", Platform="twitch", Text="!points", ChannelID="uuid-xyz"
      3. Wait 100ms for both subscribers to receive
      4. Assert command execution pipeline called exactly once (not twice)
      5. Assert Redis key "parser:dedup:twitch-msg-002" exists
    Expected Result: Dedup prevents double-processing; pipeline invoked once despite dual delivery
    Failure Indicators: Pipeline called twice, or Redis key absent
    Evidence: .sisyphus/evidence/task-22-dedup-twitch-dual.txt

  Scenario: Absence of PARSER_USE_GENERIC_QUEUE env var — parser starts successfully
    Tool: Bash
    Steps:
      1. grep -r "PARSER_USE_GENERIC_QUEUE" apps/parser/ libs/config/
      2. Assert zero matches
    Expected Result: No feature flag of this name exists anywhere
    Failure Indicators: Any match found
    Evidence: .sisyphus/evidence/task-22-no-feature-flag.txt
  ```

  **Evidence to Capture**:
  - [ ] task-22-kick-message-parsed.txt
  - [ ] task-22-dedup-twitch-dual.txt
  - [ ] task-22-no-feature-flag.txt

  **Commit**: YES
  - Message: `feat(parser): subscribe to generic.chat-message with Redis MessageID deduplication`
  - Files: `apps/parser/internal/commands-bus/commands-bus.go`
  - Pre-commit: `go build ./apps/parser/...`

---

- [ ] 23. Built-in Variable $(platform)

  **What to do**:
  - Add a new built-in variable `$(platform)` to the parser's variable registry.
  - When evaluated, returns the platform name from `ParseContext.Platform` (e.g., `"twitch"` or `"kick"`).
  - Follow the existing pattern for built-in variables (grep for `$(sender)` or similar to find the registration pattern).
  - Variable name: `platform` (so it resolves as `$(platform)` in command responses).

  **Must NOT do**:
  - Do NOT hardcode `"twitch"` — must read from ParseContext.Platform dynamically.
  - Do NOT add platform-specific logic to the variable itself — it just returns the string.
  - Do NOT add new variables beyond `$(platform)` in this task.

  **Recommended Agent Profile**:
  - **Category**: `quick`
  - **Skills**: []
    - Reason: Single variable registration; pattern already exists; mechanical addition.

  **Parallelization**:
  - **Can Run In Parallel**: YES (Wave 6, alongside T24)
  - **Parallel Group**: Wave 6 (after T21)
  - **Blocks**: Nothing (terminal in its concern)
  - **Blocked By**: T21 (Platform field in ParseContext)

  **References**:
  - `apps/parser/internal/` — grep for existing built-in variable registration (e.g., `$(sender)`, `$(channel)`)
  - `apps/parser/internal/types/parse_context.go` (T21) — Platform field source
  - Variable registration pattern: look for a `variables` map or `RegisterVariable` function

  **Acceptance Criteria**:

  **QA Scenarios**:

  ```
  Scenario: $(platform) returns "twitch" when platform = twitch
    Tool: Bash (go test)
    Steps:
      1. go test ./apps/parser/... -run TestPlatformVariable -v
      2. Build ParseContext{Platform: "twitch"}
      3. Evaluate "Hello from $(platform)"
      4. Assert result == "Hello from twitch"
    Expected Result: Variable resolves correctly
    Evidence: .sisyphus/evidence/task-23-platform-variable-twitch.txt

  Scenario: $(platform) returns "kick" when platform = kick
    Tool: Bash (go test)
    Steps:
      1. ParseContext{Platform: "kick"}
      2. Evaluate "Hello from $(platform)"
      3. Assert result == "Hello from kick"
    Expected Result: Variable resolves correctly for Kick
    Evidence: .sisyphus/evidence/task-23-platform-variable-kick.txt
  ```

  **Evidence to Capture**:
  - [ ] task-23-platform-variable-twitch.txt
  - [ ] task-23-platform-variable-kick.txt

  **Commit**: YES (group with T21)
  - Message: `feat(parser): add $(platform) built-in variable`
  - Files: `apps/parser/internal/` (variable registration file)

---

- [ ] 24. Command/Timer/Keyword Platform Filtering in Execution

  **What to do**:
  - In the parser's command, timer, and keyword execution logic, apply platform filtering using the `platforms text[]` column added in T3.
  - Logic: if `command.Platforms` is empty (or `{}`) → execute on ALL platforms. If `command.Platforms` contains one or more values (e.g., `["twitch"]`) → only execute if `ParseContext.Platform` is in the list.
  - Apply the same logic to timers and keywords.
  - The `platforms` column is already on the DB model (from T3) — just need to use it in the filter.

  **Must NOT do**:
  - Do NOT change the command/timer/keyword DB model structure — T3 already handles that.
  - Do NOT filter by platform in SQL — filter in the execution layer (the list is already loaded).
  - Do NOT break existing behavior: commands without `platforms` must still execute on all platforms.

  **Recommended Agent Profile**:
  - **Category**: `quick`
  - **Skills**: []
    - Reason: Simple conditional check addition to existing execution flow; no new infrastructure.

  **Parallelization**:
  - **Can Run In Parallel**: YES (Wave 6, alongside T23)
  - **Parallel Group**: Wave 6 (after T22)
  - **Blocks**: Nothing (terminal)
  - **Blocked By**: T22 (Platform must be in ParseContext and flowing through), T3 (platforms column must exist)

  **References**:
  - `apps/parser/internal/` — grep for command execution trigger (e.g., `handleCommand`, `executeCommand`)
  - `apps/parser/internal/types/parse_context.go` — Platform field
  - `libs/repositories/channels_commands/` or `libs/gomodels/channels_commands.go` — existing command model with `Platforms` field (from T3)
  - Same for timers: `libs/gomodels/channels_timers.go` and keywords

  **Acceptance Criteria**:

  **QA Scenarios**:

  ```
  Scenario: Command with platforms=["twitch"] does NOT execute on Kick
    Tool: Bash (go test)
    Steps:
      1. go test ./apps/parser/... -run TestPlatformFilter -v
      2. Create command with Platforms=["twitch"]; ParseContext{Platform: "kick"}
      3. Trigger command execution
      4. Assert command response was NOT sent
    Expected Result: Filtered out; no response
    Evidence: .sisyphus/evidence/task-24-platform-filter-kick-blocked.txt

  Scenario: Command with platforms=[] executes on both Twitch and Kick
    Tool: Bash (go test)
    Steps:
      1. Create command with Platforms=[] (empty); test with Platform="kick" AND Platform="twitch"
      2. Assert command response sent in BOTH cases
    Expected Result: All-platforms execution works
    Evidence: .sisyphus/evidence/task-24-platform-filter-all-platforms.txt
  ```

  **Evidence to Capture**:
  - [ ] task-24-platform-filter-kick-blocked.txt
  - [ ] task-24-platform-filter-all-platforms.txt

  **Commit**: YES
  - Message: `feat(parser): platform filter on command/timer/keyword execution`
  - Files: `apps/parser/internal/` (command/timer/keyword execution files)

---

- [ ] 25. KickProfile GraphQL Type + Schema Update

  **What to do**:
  - Add a `KickProfile` GraphQL type to the schema in `apps/api-gql/internal/delivery/gql/schema/`:
    ```graphql
    type KickProfile {
    	id: String!
    	slug: String!
    	displayName: String!
    	profilePicture: String
    	isLive: Boolean!
    	followersCount: Int!
    }
    ```
  - Add to the `AuthenticatedUser` type (in `user.graphql` or `shared-users.graphql`):
    - `kickProfile: KickProfile` (nullable — only present if Kick is linked)
    - `linkedAccounts: [LinkedAccount!]!` where:
      ```graphql
      type LinkedAccount {
      	platform: String!
      	platformUserId: String!
      	platformLogin: String!
      	platformAvatar: String
      }
      ```
  - After schema changes, run `bun cli build gql` to regenerate Go resolvers. Refresh file reads after regeneration.
  - Do NOT implement the resolver logic here — just the schema. T26 handles resolvers.

  **Must NOT do**:
  - Do NOT implement resolver logic in this task.
  - Do NOT remove or rename existing Twitch-related fields (e.g., `twitchProfile`) from `AuthenticatedUser`.
  - Do NOT forget to run `bun cli build gql` after changing `.graphql` files.

  **Recommended Agent Profile**:
  - **Category**: `quick`
  - **Skills**: []
    - Reason: Schema-only change; well-defined output; no logic.

  **Parallelization**:
  - **Can Run In Parallel**: YES — independent of parser wave
  - **Parallel Group**: Wave 7 (with T26, T27, T28)
  - **Blocks**: T26 (resolvers need the schema to exist)
  - **Blocked By**: None (schema is self-contained)

  **References**:
  - `apps/api-gql/internal/delivery/gql/schema/user.graphql` — existing AuthenticatedUser type; add kickProfile + linkedAccounts here
  - `apps/api-gql/internal/delivery/gql/schema/twitch.graphql` — TwitchProfile pattern to mirror for KickProfile
  - `apps/api-gql/internal/delivery/gql/schema/shared-users.graphql` — check if LinkedAccount-like types exist

  **Acceptance Criteria**:

  **QA Scenarios**:

  ```
  Scenario: bun cli build gql succeeds after schema update
    Tool: Bash
    Steps:
      1. bun cli build gql
      2. Assert exit code = 0
      3. Assert generated file `apps/api-gql/internal/delivery/gql/generated.go` updated (check mtime)
    Expected Result: Clean codegen
    Evidence: .sisyphus/evidence/task-25-gql-codegen.txt

  Scenario: KickProfile type appears in GraphQL introspection
    Tool: Bash (curl GraphQL introspection)
    Preconditions: api-gql service running
    Steps:
      1. curl -X POST http://localhost:3009/query -H 'Content-Type: application/json' \
           -d '{"query":"{ __type(name: \"KickProfile\") { fields { name } } }"}'
      2. Assert response contains fields: id, slug, displayName, profilePicture, isLive, followersCount
    Expected Result: KickProfile type exists in schema
    Evidence: .sisyphus/evidence/task-25-kick-profile-introspection.txt
  ```

  **Evidence to Capture**:
  - [ ] task-25-gql-codegen.txt
  - [ ] task-25-kick-profile-introspection.txt

  **Commit**: YES
  - Message: `feat(gql): add KickProfile type + linkedAccounts to AuthenticatedUser schema`
  - Files: `apps/api-gql/internal/delivery/gql/schema/user.graphql` (or shared-users.graphql), generated files

---

- [ ] 26. AuthenticatedUser kickProfile + linkedAccounts Resolvers

  **What to do**:
  - Implement GraphQL resolvers for:
    - `AuthenticatedUser.kickProfile` → fetch Kick profile from Kick's `/public/v1/channels/{slug}` API using the broadcaster's Kick token from `user_platform_accounts`. Return `nil` if no Kick account linked.
    - `AuthenticatedUser.linkedAccounts` → query `user_platform_accounts` for all rows matching the current user's `users.id` (internal UUID). Map each row to `LinkedAccount{platform, platformUserId, platformLogin, platformAvatar}`.
  - Current user's UUID is available from the session (post T11 migration).
  - Create `apps/api-gql/internal/services/kick_profile/service.go` for the Kick API call.

  **Must NOT do**:
  - Do NOT call Kick API for `linkedAccounts` — that is a pure DB query.
  - Do NOT cache Kick profile indefinitely — use short TTL (5 min) or no cache (acceptable for profile).
  - Do NOT import `libs/gomodels`.
  - Do NOT use GORM.

  **Recommended Agent Profile**:
  - **Category**: `unspecified-high`
  - **Skills**: []
    - Reason: New service + resolver wiring + external API call + DB query combination.

  **Parallelization**:
  - **Can Run In Parallel**: NO — depends on T25 (schema must exist) and T8 (user_platform_accounts repo)
  - **Parallel Group**: Wave 7 (after T25)
  - **Blocks**: T29 (frontend needs this resolver to show profile)
  - **Blocked By**: T25, T8, T11

  **References**:
  - `apps/api-gql/internal/delivery/gql/resolvers/user.resolver.go` — existing resolver file; add new resolvers here
  - `apps/api-gql/internal/services/` — existing service patterns (e.g., twitch_profile service) to mirror
  - `libs/repositories/user_platform_accounts/` (T8) — for linkedAccounts query
  - Kick API: `GET /public/v1/channels/{slug}` — returns channel info
  - `apps/api-gql/internal/auth/sessions_user.go` — how to get current user UUID from session context

  **Acceptance Criteria**:

  **QA Scenarios**:

  ```
  Scenario: linkedAccounts resolver returns all linked platforms for user
    Tool: Bash (curl GraphQL)
    Preconditions: User has Twitch + Kick linked; api-gql running; user authenticated
    Steps:
      1. curl -X POST http://localhost:3009/query \
           -H 'Content-Type: application/json' -H 'Cookie: session=...' \
           -d '{"query":"{ authenticatedUser { linkedAccounts { platform platformLogin } } }"}'
      2. Assert response contains 2 entries: one with platform="twitch", one with platform="kick"
    Expected Result: Both accounts returned
    Evidence: .sisyphus/evidence/task-26-linked-accounts.txt

  Scenario: kickProfile resolver returns nil for user without Kick linked
    Tool: Bash (curl GraphQL)
    Preconditions: User has Twitch only
    Steps:
      1. Query authenticatedUser { kickProfile { id slug } }
      2. Assert kickProfile == null
    Expected Result: null (not an error)
    Evidence: .sisyphus/evidence/task-26-kick-profile-null.txt
  ```

  **Evidence to Capture**:
  - [ ] task-26-linked-accounts.txt
  - [ ] task-26-kick-profile-null.txt

  **Commit**: YES
  - Message: `feat(gql): kickProfile and linkedAccounts resolvers`
  - Files: `apps/api-gql/internal/delivery/gql/resolvers/user.resolver.go`, `apps/api-gql/internal/services/kick_profile/service.go`

---

- [ ] 27. 7TV GetProfileByKickId in libs/integrations/seventv

  **What to do**:
  - In `libs/integrations/seventv/client.go`, add a new method `GetProfileByKickId(ctx context.Context, kickChannelID string) (*SevenTVProfile, error)`.
  - Endpoint: `GET https://7tv.io/v3/users/kick/{kickChannelID}`
  - Parse the response into the existing `SevenTVProfile` struct (or create a minimal struct if it doesn't exist).
  - If 404 → return `nil, nil` (user not found is not an error).
  - If other non-200 → return `nil, fmt.Errorf(...)`.
  - The existing `GetProfileByTwitchId` method is the direct pattern to follow — it likely calls `https://7tv.io/v3/users/twitch/{twitchID}`.

  **Must NOT do**:
  - Do NOT change the signature of `GetProfileByTwitchId`.
  - Do NOT add caching here — caching is the caller's responsibility (emotes-cacher in T28).
  - Do NOT add the Kick channel ID to a Twitch-specific path.

  **Recommended Agent Profile**:
  - **Category**: `quick`
  - **Skills**: []
    - Reason: Single HTTP GET method addition; exact pattern exists; trivial.

  **Parallelization**:
  - **Can Run In Parallel**: YES (Wave 7, alongside T25, T26)
  - **Parallel Group**: Wave 7
  - **Blocks**: T28 (emotes-cacher calls this)
  - **Blocked By**: None (independent library addition)

  **References**:
  - `libs/integrations/seventv/client.go` — MUST READ; `GetProfileByTwitchId` is the exact pattern to follow
  - 7TV API: `https://7tv.io/v3/users/kick/{channelID}` — same response shape as the Twitch endpoint

  **Acceptance Criteria**:

  **QA Scenarios**:

  ```
  Scenario: GetProfileByKickId returns profile for known Kick channel
    Tool: Bash (go test)
    Preconditions: httptest.NewServer mocking 7TV response for kick channel "xqc"
    Steps:
      1. go test ./libs/integrations/seventv/... -run TestGetProfileByKickId -v
      2. Mock returns valid 7TV profile JSON for kick channel "xqc"
      3. Assert returned profile.ID != "" and profile.Username != ""
    Expected Result: Profile returned correctly
    Evidence: .sisyphus/evidence/task-27-7tv-kick-profile.txt

  Scenario: GetProfileByKickId returns nil, nil for unknown channel (7TV returns 404)
    Tool: Bash (go test)
    Steps:
      1. Mock returns HTTP 404
      2. Assert result == (nil, nil)
    Expected Result: Not-found treated as nil without error
    Evidence: .sisyphus/evidence/task-27-7tv-kick-not-found.txt
  ```

  **Evidence to Capture**:
  - [ ] task-27-7tv-kick-profile.txt
  - [ ] task-27-7tv-kick-not-found.txt

  **Commit**: YES (group with T28)
  - Message: `feat(7tv): add GetProfileByKickId method`
  - Files: `libs/integrations/seventv/client.go`

---

- [ ] 28. emotes-cacher Platform-Aware 7TV Profile Lookup

  **What to do**:
  - In `apps/emotes-cacher`, update the 7TV profile lookup logic to be platform-aware:
    - If channel platform = "twitch" → call `GetProfileByTwitchId(twitchID)` (existing behavior)
    - If channel platform = "kick" → call `GetProfileByKickId(kickChannelID)` (new T27 method)
  - The channel's platform is determined from `channels.platform` column (added in T4) or from `user_platform_accounts`.
  - The emote set caching and response building logic is unchanged — just the lookup source changes.

  **Must NOT do**:
  - Do NOT break Twitch emote caching — existing Twitch behavior must continue unchanged.
  - Do NOT call both Twitch and Kick 7TV endpoints for the same channel — only call the one matching the channel's platform.
  - Do NOT import `libs/gomodels`.

  **Recommended Agent Profile**:
  - **Category**: `quick`
  - **Skills**: []
    - Reason: Conditional branch addition to existing lookup; the hard part (GetProfileByKickId) is in T27.

  **Parallelization**:
  - **Can Run In Parallel**: NO — depends on T27
  - **Parallel Group**: Wave 7 (last item, after T27)
  - **Blocks**: Nothing (terminal)
  - **Blocked By**: T27, T4 (channels.platform column)

  **References**:
  - `apps/emotes-cacher/` — MUST READ; find the 7TV lookup call site (grep for `GetProfileByTwitchId`)
  - `libs/integrations/seventv/client.go` — both methods (Twitch + Kick)
  - `libs/repositories/channels/` — to read channel platform

  **Acceptance Criteria**:

  **QA Scenarios**:

  ```
  Scenario: Kick channel fetches 7TV profile via Kick endpoint
    Tool: Bash (go test)
    Preconditions: Mock channel with platform="kick", kick_channel_id="12345"; mock 7TV Kick endpoint returns valid profile
    Steps:
      1. go test ./apps/emotes-cacher/... -run TestPlatformAware7TV -v
      2. Assert GetProfileByKickId was called with "12345"
      3. Assert GetProfileByTwitchId was NOT called
    Expected Result: Correct endpoint used based on platform
    Evidence: .sisyphus/evidence/task-28-emotes-cacher-kick-7tv.txt

  Scenario: Twitch channel still uses Twitch 7TV endpoint
    Tool: Bash (go test)
    Steps:
      1. Channel with platform="twitch"
      2. Assert GetProfileByTwitchId called; GetProfileByKickId NOT called
    Expected Result: Backward compatibility preserved
    Evidence: .sisyphus/evidence/task-28-emotes-cacher-twitch-unchanged.txt
  ```

  **Evidence to Capture**:
  - [ ] task-28-emotes-cacher-kick-7tv.txt
  - [ ] task-28-emotes-cacher-twitch-unchanged.txt

  **Commit**: YES (group with T27)
  - Message: `feat(emotes): platform-aware 7TV profile lookup for Kick channels`
  - Files: `apps/emotes-cacher/` (7TV lookup file)

---

- [ ] 29. Frontend: Kick Login Button + /auth/kick/callback Route (Nuxt web app)

  **What to do**:

  The login surface is the **Nuxt web app** (`web/`), NOT `frontend/dashboard`. The real login/callback page is `web/layers/landing/pages/login.client.vue`. The existing Twitch callback flow calls `api.auth.authPostCode({code, state})`. The Kick flow follows the same pattern but with a different endpoint and adds PKCE.

  **PKCE ownership**: The `code_verifier` is generated **server-side in the Go auth handler** (T9/T10), not in the browser. The Kick authorize URL is generated by the backend (`GET /api/auth/kick/authorize`) which creates the PKCE pair, stores `code_verifier` in the user's server-side session, and returns the full redirect URL. The browser simply navigates to that URL.

  **Changes to `web/` (Nuxt app)**:
  1. **Add Kick login button** to `web/layers/landing/pages/` login UI (same file or page where the Twitch login button lives — read `web/` to find it). The button calls `GET /api/auth/kick/authorize` → redirects the browser to the returned URL.

  2. **Update `web/layers/landing/pages/login.client.vue`** to handle Kick callback:
     - The page already handles `code` + `state` from URL params.
     - Add detection of `kick_callback=true` query param (or use the `state` to discriminate).
     - When it's a Kick callback: call `api.auth.authKickPostCode({code, state})` (new API endpoint from T10).
     - On success: redirect to `redirect_to` URL.
     - Twitch flow is unchanged.

  3. **New API client method** in `web/` (or `libs/api/`) for the Kick code exchange: `POST /api/auth/kick/code`.

  4. **T9/T10 backend** must implement `GET /api/auth/kick/authorize` which:
     - Generates PKCE pair (`code_verifier`, `code_challenge`)
     - Stores `code_verifier` in server-side session
     - Returns `{authorize_url: "https://id.kick.com/oauth/authorize?..."}`

  **Must NOT do**:
  - Do NOT put the login button in `frontend/dashboard/src/` — the login page is in `web/`.
  - Do NOT generate `code_verifier` in the browser — all PKCE is server-side.
  - Do NOT store `code_verifier` in `localStorage` or `sessionStorage`.
  - Do NOT modify T9 Kick PKCE to be client-side — keep server-side ownership.
  - Do NOT remove or change the Twitch login button.

  **Recommended Agent Profile**:
  - **Category**: `visual-engineering`
  - **Skills**: []
    - Reason: Nuxt 3 component addition (NOT Vue dashboard) + API client method; follows existing Twitch pattern.

  **Parallelization**:
  - **Can Run In Parallel**: YES — Wave 8 (with T30, T31, T32); needs T9/T10 backend done first
  - **Parallel Group**: Wave 8
  - **Blocks**: Nothing (user-facing, terminal)
  - **Blocked By**: T10 (backend `/auth/kick/authorize` + `/auth/kick/code` handlers must exist)

  **References**:
  - `web/layers/landing/pages/login.client.vue` — **MUST READ**; existing Twitch callback handler; add Kick here
  - `web/layers/landing/pages/` — find the page with the Twitch login button (index.vue or login.vue)
  - `web/app/stores/user.ts` — auth store; `api.auth.authPostCode` call pattern
  - `web/AGENTS.md` — Nuxt 3 conventions for this app
  - Kick OAuth authorize URL shape: `https://id.kick.com/oauth/authorize?client_id=...&redirect_uri=...&response_type=code&scope=user:read+events:subscribe+chat:write+channel:read&code_challenge=...&code_challenge_method=S256&state=...`

  **Acceptance Criteria**:

  **QA Scenarios**:

  ```
  Scenario: Kick login button is visible on the web landing login page
    Tool: Playwright
    Preconditions: web app running (Nuxt dev server); backend running; user not logged in
    Steps:
      1. Navigate to http://localhost:3000/ (or wherever the landing page runs)
      2. Find the page with the Twitch login button
      3. Assert a "Login with Kick" button (or Kick logo button) is also visible
      4. Take screenshot
    Expected Result: Kick login button visible alongside Twitch button
    Evidence: .sisyphus/evidence/task-29-kick-login-button.png

  Scenario: Clicking Kick login initiates server-side PKCE authorize URL
    Tool: Playwright + network inspection
    Steps:
      1. Click the Kick login button
      2. Monitor network: assert GET /api/auth/kick/authorize was called
      3. Assert browser navigated to a URL containing "id.kick.com/oauth/authorize"
      4. Assert URL contains "code_challenge=" and "code_challenge_method=S256"
    Expected Result: Redirect to Kick OAuth with server-generated PKCE params
    Evidence: .sisyphus/evidence/task-29-kick-oauth-redirect.png

  Scenario: Kick callback page calls /api/auth/kick/code and redirects on success
    Tool: Bash (curl) — simulate callback
    Steps:
      1. curl -s -X POST http://localhost:3000/api/auth/kick/code \
           -H 'Content-Type: application/json' \
           -d '{"code":"test-code","state":"aHR0cDovL2xvY2FsaG9zdA=="}'
      2. Assert response: non-500 (will be 400 "invalid code" from Kick, not internal error)
    Expected Result: Endpoint exists and handles request (even if code is invalid)
    Evidence: .sisyphus/evidence/task-29-kick-code-endpoint.txt
  ```

  **Evidence to Capture**:
  - [ ] task-29-kick-login-button.png
  - [ ] task-29-kick-oauth-redirect.png
  - [ ] task-29-kick-code-endpoint.txt

  **Commit**: YES
  - Message: `feat(web): Kick login button + callback handler in Nuxt web app`
  - Files: `web/layers/landing/pages/login.client.vue`, login page component in `web/`, API client update

---

- [ ] 30. Frontend: Linked Accounts Settings Section

  **What to do**:
  - In `frontend/dashboard/src/pages/settings/` (or equivalent profile settings page), add a "Linked Accounts" section that shows all platforms connected to the current user's account.
  - Uses the `linkedAccounts` GraphQL resolver (from T26) via urql query.
  - Each linked account shows: platform icon (Twitch/Kick logo), `platformLogin`, a "Disconnect" button.
  - Show a "Connect" button for platforms NOT yet linked (e.g., if only Twitch is linked, show "Connect Kick" which initiates Kick OAuth).
  - "Disconnect" calls a new mutation `unlinkPlatformAccount(platform: String!): Boolean!` — add this to the GraphQL schema and implement the resolver (deletes row from `user_platform_accounts`, also unsubscribes Kick EventSub if platform=kick).
  - Minimum viable: show linked accounts list + connect/disconnect actions. No profile editing here.

  **Must NOT do**:
  - Do NOT allow disconnecting the primary account (the account used to log in) if it's the only one — show a tooltip "Cannot remove primary account".
  - Do NOT show account passwords or tokens.
  - Do NOT add this section to a separate page — add it to the existing profile/settings page.

  **Recommended Agent Profile**:
  - **Category**: `visual-engineering`
  - **Skills**: []
    - Reason: Vue 3 component + GraphQL query/mutation + UI design.

  **Parallelization**:
  - **Can Run In Parallel**: YES (Wave 8, alongside T29, T31, T32)
  - **Parallel Group**: Wave 8
  - **Blocks**: Nothing (terminal)
  - **Blocked By**: T26 (linkedAccounts resolver must exist), T25 (schema must include unlinkPlatformAccount mutation)

  **References**:
  - `frontend/dashboard/src/pages/settings/` — existing settings page; add section here
  - `frontend/dashboard/AGENTS.md` — urql pattern, vee-validate, defineModel
  - `apps/api-gql/internal/delivery/gql/schema/user.graphql` — add unlinkPlatformAccount mutation
  - `frontend/dashboard/src/api/integrations/integrations-page.ts` — unified query pattern for reference

  **Acceptance Criteria**:

  **QA Scenarios**:

  ```
  Scenario: Settings page shows linked accounts for authenticated user
    Tool: Playwright
    Preconditions: User logged in with Twitch + Kick linked; dashboard running
    Steps:
      1. Navigate to http://localhost:3005/dashboard/settings (or /profile)
      2. Assert section with heading "Linked Accounts" (or equivalent) is visible
      3. Assert two account entries: one for "twitch", one for "kick"
      4. Assert each entry shows platformLogin
      5. Take screenshot
    Expected Result: Both linked accounts displayed
    Evidence: .sisyphus/evidence/task-30-linked-accounts-ui.png

  Scenario: Disconnect Kick account removes it from the list
    Tool: Playwright
    Steps:
      1. Click "Disconnect" next to the Kick account entry
      2. Confirm dialog (if any)
      3. Assert Kick entry disappears from the list
      4. Assert Twitch entry remains
    Expected Result: Kick account unlinked from UI
    Evidence: .sisyphus/evidence/task-30-disconnect-kick.png
  ```

  **Evidence to Capture**:
  - [ ] task-30-linked-accounts-ui.png
  - [ ] task-30-disconnect-kick.png

  **Commit**: YES
  - Message: `feat(frontend): linked accounts settings section with connect/disconnect`
  - Files: `frontend/dashboard/src/pages/settings/` (settings component), `apps/api-gql/internal/delivery/gql/schema/user.graphql` (unlinkPlatformAccount mutation), resolver file

---

- [ ] 31. Frontend: Platform Selector on Command/Timer/Keyword Forms

  **What to do**:
  - Add a "Platforms" multi-select field to the command creation/edit form, the timer form, and the keyword form.
  - The field allows selecting: "Twitch", "Kick", or both (maps to `["twitch"]`, `["kick"]`, or `[]` = both).
  - UI: use a `<Checkbox>` group or `<ToggleGroup>` with platform icons + labels. Empty selection = all platforms (default). Make this clear in the UI with helper text.
  - Bind to the `platforms` field on the command/timer/keyword GraphQL mutation (existing mutation must be updated to accept `platforms: [String!]` — update schema + resolver if not already done in T3 backend tasks).
  - On save, send `platforms: ["twitch"]` or `platforms: []` etc.

  **Must NOT do**:
  - Do NOT use a plain text input for platform selection — use a proper multi-select UI.
  - Do NOT show this field if the user has only Twitch linked (can optionally hide or show as disabled).
  - Do NOT change the existing form validation schema — add the `platforms` field to the existing zod schema.

  **Recommended Agent Profile**:
  - **Category**: `visual-engineering`
  - **Skills**: []
    - Reason: Vue 3 form component addition; vee-validate + zod schema update.

  **Parallelization**:
  - **Can Run In Parallel**: YES (Wave 8, alongside T29, T30, T32)
  - **Parallel Group**: Wave 8
  - **Blocks**: Nothing (terminal)
  - **Blocked By**: T3 (platforms[] column in DB), GraphQL schema must include `platforms` on command mutation

  **References**:
  - `frontend/dashboard/src/` — MUST READ; find command edit form (grep for "command" + "form" or "useForm")
  - `frontend/dashboard/AGENTS.md` — vee-validate + zod pattern (AGENTS.md section 3)
  - `apps/api-gql/internal/delivery/gql/schema/commands.graphql` — command mutation input type; add `platforms` field
  - `frontend/dashboard/src/components/ui/` — existing checkbox/toggle components to use

  **Acceptance Criteria**:

  **QA Scenarios**:

  ```
  Scenario: Platform selector appears on command edit form
    Tool: Playwright
    Preconditions: Dashboard running; user logged in; navigate to a command's edit page
    Steps:
      1. Navigate to http://localhost:3005/dashboard/commands/[id]/edit (or modal)
      2. Assert a "Platforms" field with checkboxes/toggles for "Twitch" and "Kick" is visible
      3. Take screenshot
    Expected Result: Platform selector visible
    Evidence: .sisyphus/evidence/task-31-platform-selector-visible.png

  Scenario: Selecting only Twitch saves platforms=["twitch"] to backend
    Tool: Playwright + network inspection
    Steps:
      1. Uncheck Kick; ensure only Twitch is selected
      2. Save form
      3. Inspect outgoing GraphQL mutation body
      4. Assert mutation contains platforms: ["twitch"]
    Expected Result: Correct platforms array sent in mutation
    Evidence: .sisyphus/evidence/task-31-platforms-saved.png
  ```

  **Evidence to Capture**:
  - [ ] task-31-platform-selector-visible.png
  - [ ] task-31-platforms-saved.png

  **Commit**: YES
  - Message: `feat(frontend): platform selector on command/timer/keyword forms`
  - Files: `frontend/dashboard/src/` (command/timer/keyword form components)

---

- [ ] 32. Frontend: Platform-Aware Profile Header

  **What to do**:
  - Update the dashboard profile header/navigation to show the profile data for the **current active platform** (the platform the user logged in with / currently selected).
  - If logged in via Twitch → show Twitch avatar + display name (existing behavior, preserve it).
  - If logged in via Kick → show Kick avatar + display name from `kickProfile` resolver.
  - Add a small platform badge/icon next to the avatar showing which platform is active.
  - Use `authenticatedUser { twitchProfile { ... } kickProfile { ... } }` query (both fields) and display based on `currentPlatform` (available from session or a new `authenticatedUser.currentPlatform: String!` field — add this to the GraphQL schema if needed).

  **Must NOT do**:
  - Do NOT remove Twitch profile display — fallback to Twitch if Kick profile unavailable.
  - Do NOT make a separate API call for profile — reuse the existing `authenticatedUser` query.
  - Do NOT add a platform switcher here (that is a separate future task) — just display the current active platform.

  **Recommended Agent Profile**:
  - **Category**: `visual-engineering`
  - **Skills**: []
    - Reason: Vue 3 component update; conditional rendering based on platform.

  **Parallelization**:
  - **Can Run In Parallel**: YES (Wave 8, alongside T29, T30, T31)
  - **Parallel Group**: Wave 8
  - **Blocks**: Nothing (terminal)
  - **Blocked By**: T26 (kickProfile resolver must exist)

  **References**:
  - `frontend/dashboard/src/` — MUST READ; find the profile header / navigation component (grep for "avatar" or "display name")
  - `apps/api-gql/internal/delivery/gql/schema/user.graphql` — add `currentPlatform: String!` field to AuthenticatedUser if not already present
  - `frontend/dashboard/AGENTS.md` — Vue 3 + Tailwind CSS conventions

  **Acceptance Criteria**:

  **QA Scenarios**:

  ```
  Scenario: Twitch user sees Twitch avatar in profile header
    Tool: Playwright
    Preconditions: User logged in via Twitch; dashboard running
    Steps:
      1. Navigate to http://localhost:3005/dashboard
      2. Assert profile header shows Twitch avatar (img src contains twitch CDN domain)
      3. Assert platform badge shows Twitch icon
      4. Take screenshot
    Expected Result: Twitch profile displayed correctly
    Evidence: .sisyphus/evidence/task-32-twitch-profile-header.png

  Scenario: Kick user sees Kick avatar in profile header
    Tool: Playwright
    Preconditions: User logged in via Kick; dashboard running
    Steps:
      1. Navigate to http://localhost:3005/dashboard
      2. Assert profile header shows Kick avatar
      3. Assert platform badge shows Kick icon
    Expected Result: Kick profile displayed correctly
    Evidence: .sisyphus/evidence/task-32-kick-profile-header.png
  ```

  **Evidence to Capture**:
  - [ ] task-32-twitch-profile-header.png
  - [ ] task-32-kick-profile-header.png

  **Commit**: YES
  - Message: `feat(frontend): platform-aware profile header with platform badge`
  - Files: `frontend/dashboard/src/` (header/nav component)

---

- [ ] 33. Kick Chat Client (POST /public/v1/chat)

  **What to do**:
  - Create `apps/bots/internal/kick/chat_client.go` (or `libs/integrations/kick/chat_client.go` if bot is a shared lib — check `apps/bots` structure first).
  - The client wraps Kick's `POST https://api.kick.com/public/v1/chat` endpoint:
    - Auth: bot's OAuth token from `kick_bots` table (fetched via T12 repository)
    - Required fields: `broadcaster_user_id` (Kick channel ID), `content` (message text), optionally `reply_to_message_id`
    - On 401 (token expired) → attempt token refresh using Kick's refresh token endpoint, update `kick_bots` table with new token, retry once.
    - On rate limit (429) → log warning + drop message (do NOT retry in a tight loop).
    - Returns `error` only on unexpected failures; 429 returns `nil` (graceful drop).
  - Method signature: `SendMessage(ctx context.Context, broadcasterKickID string, text string) error`

  **Must NOT do**:
  - Do NOT use the broadcaster's token — use the **bot's token** from `kick_bots`.
  - Do NOT retry more than once on 401.
  - Do NOT panic on rate limit — log and return nil.
  - Do NOT import `libs/gomodels`.

  **Recommended Agent Profile**:
  - **Category**: `unspecified-high`
  - **Skills**: []
    - Reason: External API client + token refresh logic + error handling nuances.

  **Parallelization**:
  - **Can Run In Parallel**: YES — can be built independently once T12 (kick_bots repo) exists
  - **Parallel Group**: Wave 9 (with T34)
  - **Blocks**: T34
  - **Blocked By**: T12 (kick_bots repo for token fetching)

  **References**:
  - `apps/bots/` — MUST READ; find existing Twitch chat message sending code (pattern to mirror)
  - `libs/repositories/kick_bots/` (T12) — for bot token retrieval + update
  - Kick chat API: `https://docs.kick.com/messaging/send-chat-message` — POST /public/v1/chat payload
  - Kick refresh token: `POST https://id.kick.com/oauth/token` with `grant_type=refresh_token`

  **Acceptance Criteria**:

  **QA Scenarios**:

  ```
  Scenario: SendMessage sends correct HTTP request to Kick chat API
    Tool: Bash (go test)
    Preconditions: httptest.NewServer mocking Kick API; kick_bots mock returns token "bot-token-123"
    Steps:
      1. go test ./apps/bots/internal/kick/... -run TestSendMessage -v
      2. Call SendMessage(ctx, "12345", "Hello from Twir!")
      3. Assert mock server received POST /public/v1/chat
      4. Assert request body: {"broadcaster_user_id": "12345", "content": "Hello from Twir!"}
      5. Assert Authorization header: "Bearer bot-token-123"
    Expected Result: Correct request sent
    Evidence: .sisyphus/evidence/task-33-send-message-request.txt

  Scenario: 401 response triggers token refresh and retry
    Tool: Bash (go test)
    Steps:
      1. Mock: first call returns 401; refresh endpoint returns new token "new-token-456"; second call returns 200
      2. Call SendMessage
      3. Assert: refresh endpoint was called
      4. Assert: second POST uses "new-token-456"
      5. Assert: kick_bots repo updated with new token
      6. Assert: no error returned
    Expected Result: Token refresh flow works
    Evidence: .sisyphus/evidence/task-33-token-refresh.txt
  ```

  **Evidence to Capture**:
  - [ ] task-33-send-message-request.txt
  - [ ] task-33-token-refresh.txt

  **Commit**: YES
  - Message: `feat(kick): Kick chat client with token refresh`
  - Files: `apps/bots/internal/kick/chat_client.go`, `apps/bots/internal/kick/chat_client_test.go`

---

- [ ] 34. Bot Service Route Send-Message by Platform

  **What to do**:
  - In `apps/bots`, update the "send chat message" logic to route by platform:
    - If `platform == "twitch"` → use existing Twitch IRC/helix client (unchanged)
    - If `platform == "kick"` → use `KickChatClient.SendMessage(...)` (from T33)
  - The routing point is wherever the bot currently decides to send a message. This is triggered by the parser's command response flow (via NATS bus or direct call).
  - The send-message NATS topic/handler must accept a `Platform` field (add to the existing send-message bus struct or create a new `generic.SendChatMessage` struct in `libs/bus-core/generic/`).
  - If platform is unknown → log error + drop (do NOT panic).

  **Must NOT do**:
  - Do NOT change the Twitch message sending path — only add the routing branch.
  - Do NOT remove the existing Twitch bot code.
  - Do NOT use the broadcaster's token for sending — always use the bot's token (Twitch bot or Kick bot account).

  **Recommended Agent Profile**:
  - **Category**: `unspecified-high`
  - **Skills**: []
    - Reason: Modifying the critical message dispatch path; needs careful routing + regression safety.

  **Parallelization**:
  - **Can Run In Parallel**: NO — depends on T33 (Kick chat client)
  - **Parallel Group**: Wave 9 (last)
  - **Blocks**: Nothing (terminal)
  - **Blocked By**: T33, T17 (generic bus struct may be needed for platform field)

  **References**:
  - `apps/bots/` — MUST READ; find the NATS subscription for send-chat-message and the dispatch logic
  - `libs/bus-core/` — existing send-message bus topic (if any); or add `libs/bus-core/generic/send-chat-message.go`
  - `apps/bots/internal/kick/chat_client.go` (T33) — the Kick send method
  - `apps/parser/` — how parser triggers bot to send (the interface between parser and bots)

  **Acceptance Criteria**:

  **QA Scenarios**:

  ```
  Scenario: Send message with platform="kick" routes to KickChatClient
    Tool: Bash (go test)
    Preconditions: Mock KickChatClient.SendMessage; mock Twitch client
    Steps:
      1. go test ./apps/bots/... -run TestPlatformRoutingSendMessage -v
      2. Trigger send-message with {platform: "kick", channel_id: "uuid-abc", text: "hi"}
      3. Assert KickChatClient.SendMessage was called with text "hi"
      4. Assert Twitch client was NOT called
    Expected Result: Kick message sent via Kick client
    Evidence: .sisyphus/evidence/task-34-kick-message-routed.txt

  Scenario: Send message with platform="twitch" routes to Twitch client
    Tool: Bash (go test)
    Steps:
      1. Trigger send-message with {platform: "twitch", channel_id: "uuid-abc", text: "hello"}
      2. Assert Twitch client.SendMessage was called
      3. Assert KickChatClient was NOT called
    Expected Result: Backward compatibility preserved
    Evidence: .sisyphus/evidence/task-34-twitch-message-routed.txt
  ```

  **Evidence to Capture**:
  - [ ] task-34-kick-message-routed.txt
  - [ ] task-34-twitch-message-routed.txt

  **Commit**: YES
  - Message: `feat(bots): route send-message by platform (Twitch/Kick)`
  - Files: `apps/bots/` (dispatch file), `libs/bus-core/generic/send-chat-message.go` (if new)

---

## Final Verification Wave (MANDATORY — after ALL T1–T34 tasks complete)

> 4 review agents run in PARALLEL. ALL must APPROVE before work is considered done.
> Present consolidated results to user and get explicit "okay" before marking complete.
> **Do NOT auto-proceed. Rejection or user feedback → fix → re-run → present again → wait for okay.**

- [ ] F1. **Plan Compliance Audit** — `oracle`

  Read the full plan file `.sisyphus/plans/kick-platform-support.md` end-to-end.

  **Must Have** — verify each exists:
  - `user_platform_accounts` table exists: `psql -c "\d user_platform_accounts"` → shows columns platform, platform_user_id, platform_login
  - `kick_bots` table exists: `psql -c "\d kick_bots"` → shows token fields
  - `platforms text[]` column on channels_commands: `psql -c "\d channels_commands"` → shows `platforms` column
  - `channels.platform` column: `psql -c "\d channels"` → shows `platform` column
  - `apps/eventsub/internal/kick/` directory exists with: handlers.go, middleware.go, subscription_manager.go, resubscribe_job.go
  - `libs/bus-core/generic/chat-message.go` exists with `Platform string` field: grep for `Platform` in that file
  - `libs/bus-core/kick/` directory exists with chat-message.go, follow.go
  - `apps/parser/internal/types/parse_context.go` has `Platform` field: grep file for `Platform`
  - `$(platform)` variable registered: grep parser for `"platform"` variable registration
  - `libs/integrations/seventv/client.go` has `GetProfileByKickId` method: grep file
  - `apps/api-gql/internal/delivery/gql/schema/` has KickProfile type: grep for `KickProfile`
  - Frontend Kick login button: grep `web/layers/landing/pages/` for "kick" in login-related components (T29 is in the Nuxt `web/` app, NOT `frontend/dashboard/`)
  - `apps/bots/internal/kick/chat_client.go` exists
  - All evidence files under `.sisyphus/evidence/` exist (check directory listing)

  **Must NOT Have** — search for forbidden patterns:
  - No GORM in new files: `grep -r "gorm.DB\|gorm.Open\|gorm.Model" apps/eventsub/internal/kick/ libs/repositories/kick_bots/ libs/repositories/user_platform_accounts/` → must return 0 results
  - No `libs/gomodels` imports in new files: `grep -r "libs/gomodels" apps/eventsub/internal/kick/ apps/bots/internal/kick/ libs/repositories/kick_bots/` → must return 0 results
  - `twitch.TwitchChatMessage` NATS subject still present: `grep -r "twitch.chat-message\|TwitchChatMessage" libs/bus-core/twitch/` → must return results (it must still exist)
  - No `context.TODO()` in new production files: `grep -r "context.TODO()" apps/eventsub/internal/kick/ apps/bots/internal/kick/` → must return 0 results

  Output: `Must Have [N/N] | Must NOT Have [N/N] | Evidence files [N/N] | VERDICT: APPROVE/REJECT`

---

- [ ] F2. **Code Quality Review** — `unspecified-high`

  Run build + lint + tests and inspect all new/changed files.

  **Steps**:
  1. `bun cli build` → record exit code
  2. `bun lint` → record exit code and any warnings
  3. `go test ./apps/eventsub/... ./apps/parser/... ./apps/bots/... ./libs/... 2>&1 | tail -50` → record pass/fail counts
  4. Review every new `.go` file in: `apps/eventsub/internal/kick/`, `apps/bots/internal/kick/`, `libs/repositories/kick_bots/`, `libs/repositories/user_platform_accounts/`, `libs/bus-core/generic/`, `libs/bus-core/kick/`, `libs/integrations/seventv/`
  5. Check each file for:
     - `as any` or `@ts-ignore` (frontend files)
     - Empty `catch` blocks (`catch(_) {}` or `catch {}`)
     - `console.log` in Vue components (use logger instead)
     - Commented-out code blocks (> 3 consecutive lines of `//`)
     - Unused imports: run `goimports -l` on Go files; check Vue files manually
     - AI slop patterns: overly generic names like `data`, `result`, `item`, `temp`, `manager2`
     - Over-abstraction: interfaces with only one implementation that adds no value

  Output: `Build [PASS/FAIL] | Lint [PASS/FAIL] | Tests [N pass/N fail] | Files reviewed [N] | Issues found [N] | VERDICT: APPROVE/REJECT`

---

- [ ] F3. **Real Manual QA** — `unspecified-high` + `playwright` skill

  Execute every QA scenario from every task. Start from a clean state. Save ALL evidence to `.sisyphus/evidence/final-qa/`.

  **Auth flow**:
  - Kick login button visible at `/login` — Playwright screenshot
  - Clicking Kick login redirects to `id.kick.com/oauth/authorize` with `code_challenge` in URL
  - Twitch login still works unchanged — complete a Twitch login flow end-to-end

  **Webhook ingestion**:
  - POST to `/webhook/kick` with bad RSA signature → assert 403
  - POST to `/webhook/kick` with missing `Kick-Event-Signature` → assert 400/403
  - POST to `/webhook/kick` with duplicate `Kick-Event-Message-Id` → assert 200 but NATS receives 0 duplicate messages (check via test NATS subscriber)

  **Parser + Platform Filter**:
  - `go test ./apps/parser/... -run TestPlatformFilter` → PASS
  - `go test ./apps/parser/... -run TestPlatformVariable` ($(platform) variable) → PASS
  - `go test ./apps/parser/... -run TestGenericQueueKickMessage` → PASS

  **GraphQL resolvers**:
  - `linkedAccounts` query returns correct platform list for a user with both Twitch + Kick linked — curl + assert
  - `kickProfile` returns null for Twitch-only user — curl + assert
  - `KickProfile` type exists in introspection — curl + assert

  **7TV**:
  - `go test ./libs/integrations/seventv/... -run TestGetProfileByKickId` → PASS
  - `go test ./apps/emotes-cacher/... -run TestPlatformAware7TV` → PASS

  **Bot message routing**:
  - `go test ./apps/bots/... -run TestPlatformRoutingSendMessage` → PASS (Kick message routes to KickChatClient; Twitch routes to Twitch client)

  **Frontend**:
  - Platform selector visible on command form — Playwright screenshot
  - Linked accounts section visible on settings page — Playwright screenshot
  - Profile header shows correct platform badge — Playwright screenshot

  Output: `Auth [PASS/FAIL] | Webhooks [PASS/FAIL] | Parser [N/N tests] | GraphQL [PASS/FAIL] | 7TV [PASS/FAIL] | Bots [PASS/FAIL] | Frontend [N/N scenarios] | VERDICT: APPROVE/REJECT`

---

- [ ] F4. **Scope Fidelity Check** — `deep`

  For each of the 34 tasks, read "What to do" from the plan and compare against actual git diff.

  **Steps**:
  1. `git log --oneline main..HEAD` → list all commits made
  2. `git diff main..HEAD --name-only` → list all changed files
  3. For each task T1–T34:
     a. Read the task's "What to do" requirements
     b. Check the corresponding changed files contain the expected implementation
     c. Check the task's "Must NOT do" items are NOT violated (search changed files)
  4. Cross-task contamination check:
     - T1–T4 (migrations) — only `libs/migrations/postgres/` changed
     - T5–T8 (interfaces/repos) — only `libs/auth/`, `libs/entities/`, `libs/repositories/` changed
     - T13–T16 (eventsub HTTP/webhook) — only `apps/eventsub/internal/kick/` changed
     - T29 (Kick login) — only `web/layers/landing/pages/` and `web/app/stores/` changed (Nuxt web app, NOT frontend/dashboard)
     - T30–T32 (dashboard frontend) — only `frontend/dashboard/src/` changed
  5. Unaccounted files: list any files in git diff NOT covered by a task plan entry. Flag each one.

  Output: `Tasks compliant [N/34] | Must-NOT violations [N] | Contamination issues [N] | Unaccounted files [N] | VERDICT: APPROVE/REJECT`

---

## Commit Strategy

Group by phase. Each task = 1 atomic commit (or grouped with sibling if trivially small).

- T1: `feat(db): migrate users.id to internal UUID + add user_platform_accounts`
- T2: `feat(db): add kick_bots table`
- T3: `feat(db): add platforms[] to commands/timers/keywords`
- T4: `feat(db): add platform column to channels`
- T5: `feat(auth): PlatformProvider interface + Twitch implementation`
- T6: `feat(entities): add platform account entities`
- T7: `fix(tokens): update GetByUserID to use internal UUID`
- T8: `feat(repo): user_platform_accounts repository`
- T9: `feat(auth): Kick OAuth 2.1 + PKCE provider`
- T10: `feat(auth): generic /auth/:platform/code handler`
- T11: `feat(session): store internal UUID + current platform; invalidate old sessions`
- T12: `feat(repo): kick_bots repository`
- T13: `feat(eventsub): add HTTP server to apps/eventsub`
- T14: `feat(eventsub): Kick webhook HMAC-SHA256 verification middleware`
- T15: `feat(eventsub): Kick EventSub subscription manager`
- T16: `feat(bus): Kick EventSub bus topics`
- T17: `feat(bus): generic ChatMessage struct with Platform field`
- T18: `feat(eventsub): Kick webhook event handlers`
- T19: `feat(eventsub): BusListener Kick support + Twitch dual-publish`
- T20: `feat(eventsub): Kick webhook auto-resubscription health check`
- T21: `feat(parser): ParseContext.Platform field`
- T22: `feat(parser): dual-subscribe generic + twitch queues with Redis MessageID dedup`
- T23: `feat(parser): $(platform) built-in variable`
- T24: `feat(parser): platform filter on command/timer/keyword execution`
- T25: `feat(gql): KickProfile type + schema`
- T26: `feat(gql): kickProfile + linkedAccounts resolvers`
- T27: `feat(7tv): GetProfileByKickId`
- T28: `feat(emotes): platform-aware 7TV profile lookup`
- T29: `feat(frontend): Kick login button + auth callback`
- T30: `feat(frontend): linked accounts settings section`
- T31: `feat(frontend): platform selector on forms`
- T32: `feat(frontend): platform-aware profile header`
- T33: `feat(kick): Kick chat client + kick_bots token usage`
- T34: `feat(bots): route send-message by platform`

---

## Success Criteria

### Verification Commands

```bash
# Build
bun cli build  # Expected: exit 0

# Lint
bun lint  # Expected: exit 0

# DB: no Twitch IDs left as users.id
psql -c "SELECT COUNT(*) FROM users WHERE id::text ~ '^[0-9]+$'"
# Expected: count = 0

# Auth: Kick login returns redirect_to
curl -X POST http://localhost:3009/api/auth/kick/code \
  -H 'Content-Type: application/json' \
  -d '{"code":"test","state":"aHR0cDovL2xvY2FsaG9zdA==","code_verifier":"test"}'
# Expected: 200 {"data":{"redirect_to":"http://localhost"}} or 400 with specific OAuth error (not 500)

# Webhook: bad signature returns 403
curl -X POST http://localhost:3000/webhook/kick \
  -H 'Kick-Event-Signature: bad-sig' \
  -H 'Content-Type: application/json' \
  -d '{}'
# Expected: 403

# Platform filter: verified via integration test in apps/parser
go test ./apps/parser/... -run TestPlatformFilter
# Expected: PASS
```

### Final Checklist

- [ ] All "Must Have" deliverables present
- [ ] No GORM in any new files
- [ ] No `libs/gomodels` imports in any new files
- [ ] `twitch.TwitchChatMessage` NATS queue still present (not removed)
- [ ] All existing Twitch tests pass
- [ ] Kick login flow works end-to-end
- [ ] Kick webhook ingestion verified
- [ ] Platform filter verified on commands
- [ ] `$(platform)` variable works
