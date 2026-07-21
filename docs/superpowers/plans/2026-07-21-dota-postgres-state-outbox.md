# Dota PostgreSQL State Outbox Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Make Dota GSI state transitions and Twitch prediction actions durable, ordered, and safe across six replicas.

**Architecture:** PostgreSQL stores the authoritative per-channel match snapshot and prediction outbox in one transaction. GSI handlers use optimistic revision checks plus a locked state row. Workers claim ordered outbox actions with lock tokens and `FOR UPDATE SKIP LOCKED`; the existing Redis prediction record remains a second idempotency boundary for Twitch calls.

**Tech Stack:** Go 1.26, pgx v5, PostgreSQL 18, transaction-manager, Fx, Redis (prediction ownership/cache only), Twitch Helix.

---

## File Structure

- `libs/migrations/postgres/<timestamp>_dota_match_state_outbox.sql`: authoritative state and durable action tables.
- `libs/repositories/dota/model/match_lifecycle.go`: database-facing match state and outbox models.
- `libs/repositories/dota/repository.go`: transition, action claim, retry, and completion contracts.
- `libs/repositories/dota/pgx/pgx.go`: transactional pgx implementation.
- `apps/dota/internal/match/actions.go`: domain lifecycle action and payload mapping.
- `apps/dota/internal/match/state_machine.go`: side-effect-free GSI transitions over repository state.
- `apps/dota/internal/predictions/lifecycle_worker.go`: polling worker for claimed durable actions.
- `apps/dota/internal/predictions/predictions.go`: direct create/resolve/cancel methods, no NATS lifecycle subscriptions.
- `apps/dota/internal/predictions/redis_store.go`: retained create and terminal ownership claims.
- `apps/dota/app/app.go`: Fx wiring for repository-backed state and worker.

### Task 1: Add Authoritative Match State and Prediction Outbox Schema

**Files:**
- Create: `libs/migrations/postgres/<timestamp>_dota_match_state_outbox.sql`
- Modify: `libs/repositories/dota/model/model.go`
- Create: `libs/repositories/dota/model/match_lifecycle.go`
- Modify: `libs/repositories/dota/repository.go`
- Modify: `libs/repositories/dota/repository_test.go`

- [ ] **Step 1: Create the migration through the CLI**

Verify the destination and generate the migration:

```bash
ls libs/migrations/postgres
bun cli m create --name dota_match_state_outbox --db postgres --type sql
```

Use the generated filename and write this schema:

```sql
-- +goose Up
-- +goose StatementBegin
CREATE TABLE dota_channel_match_states (
    channel_id UUID PRIMARY KEY REFERENCES channels(id) ON DELETE CASCADE,
    revision BIGINT NOT NULL DEFAULT 0 CHECK (revision >= 0),
    provider_timestamp BIGINT NOT NULL DEFAULT 0,
    snapshot JSONB NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE dota_prediction_outbox (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    channel_id UUID NOT NULL REFERENCES channels(id) ON DELETE CASCADE,
    match_id BIGINT NOT NULL CHECK (match_id > 0),
    action TEXT NOT NULL CHECK (action IN ('create', 'resolve', 'cancel')),
    sequence BIGINT NOT NULL CHECK (sequence > 0),
    payload JSONB NOT NULL,
    attempts INT NOT NULL DEFAULT 0 CHECK (attempts >= 0),
    available_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    locked_at TIMESTAMPTZ,
    lock_token UUID,
    completed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (channel_id, match_id, action)
);

CREATE INDEX dota_prediction_outbox_claim_idx
    ON dota_prediction_outbox (available_at, sequence, created_at)
    WHERE completed_at IS NULL;
CREATE INDEX dota_prediction_outbox_match_order_idx
    ON dota_prediction_outbox (channel_id, match_id, sequence)
    WHERE completed_at IS NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS dota_prediction_outbox;
DROP TABLE IF EXISTS dota_channel_match_states;
-- +goose StatementEnd
```

- [ ] **Step 2: Write failing model-validation tests**

Add tests before implementation:

```go
func TestOutboxActionInputValidation(t *testing.T) {
    err := ValidateOutboxActionInput(OutboxActionInput{
        ChannelID: uuid.New(),
        MatchID:   0,
        Action:    OutboxActionCreate,
        Sequence:  1,
    })

    require.ErrorContains(t, err, "match ID")
}

func TestOutboxActionInputRejectsUnknownAction(t *testing.T) {
    err := ValidateOutboxActionInput(OutboxActionInput{
        ChannelID: uuid.New(),
        MatchID:   42,
        Action:    OutboxAction("unknown"),
        Sequence:  1,
    })

    require.ErrorContains(t, err, "action")
}
```

- [ ] **Step 3: Run the validation tests and verify RED**

Run:

```bash
cd libs/repositories && go test ./dota/... -run 'TestOutboxActionInput'
```

Expected: FAIL because outbox action types and validation do not exist.

- [ ] **Step 4: Add lifecycle models and contracts**

Create `model/match_lifecycle.go`:

```go
type MatchState struct {
    ChannelID         uuid.UUID
    Revision          int64
    ProviderTimestamp int64
    Snapshot          json.RawMessage
    UpdatedAt         time.Time
}

type OutboxAction string

const (
    OutboxActionCreate  OutboxAction = "create"
    OutboxActionResolve OutboxAction = "resolve"
    OutboxActionCancel  OutboxAction = "cancel"
)

type OutboxActionInput struct {
    ChannelID uuid.UUID
    MatchID   int64
    Action    OutboxAction
    Sequence  int64
    Payload   json.RawMessage
}

type ClaimedOutboxAction struct {
    ID        uuid.UUID
    LockToken uuid.UUID
    OutboxActionInput
    Attempts int
}

type ApplyMatchStateTransitionInput struct {
    ChannelID         uuid.UUID
    ExpectedRevision  int64
    ProviderTimestamp int64
    Snapshot          json.RawMessage
    Actions           []OutboxActionInput
}

type ClaimPredictionActionsInput struct {
    Limit int
    Lease time.Duration
}
```

Add these repository members:

```go
GetMatchState(context.Context, uuid.UUID) (model.MatchState, error)
ApplyMatchStateTransition(context.Context, ApplyMatchStateTransitionInput) (bool, error)
ClaimPredictionActions(context.Context, ClaimPredictionActionsInput) ([]model.ClaimedOutboxAction, error)
CompletePredictionAction(context.Context, uuid.UUID, uuid.UUID) error
RetryPredictionAction(context.Context, uuid.UUID, uuid.UUID, time.Time) error
```

Validate positive match IDs, action kinds, unique actions, and sequence values
before entering PostgreSQL.

- [ ] **Step 5: Run validation tests and repository build**

Run:

```bash
cd libs/repositories && go test ./dota/... && go build ./...
```

Expected: PASS.

- [ ] **Step 6: Commit schema and contracts**

```bash
git add libs/migrations/postgres libs/repositories/dota
git commit -m "feat: add dota match state outbox schema"
```

### Task 2: Implement Transactional State Transition and Outbox Claims

**Files:**
- Modify: `libs/repositories/dota/pgx/pgx.go`
- Create: `libs/repositories/dota/pgx/match_lifecycle_test.go`
- Modify: `libs/repositories/dota/pgx/pgx_test.go`

- [ ] **Step 1: Write failing transaction-flow tests**

Use a narrow fake transaction executor that records SQL decisions. Add these behavioral tests:

```go
func TestApplyTransitionUpdatesStateAndOutboxInOneTransaction(t *testing.T) {
    executor := newFakeLifecycleExecutor(matchState{Revision: 4})
    input := transitionInput(channelID, 4, 100, snapshotJSON(5), []model.OutboxActionInput{
        createActionInput(channelID, 42, 5),
    })

    committed, err := executor.ApplyMatchStateTransition(ctx, input)

    require.NoError(t, err)
    require.True(t, committed)
    require.Equal(t, int64(5), executor.savedState.Revision)
    require.Equal(t, []model.OutboxActionInput{createActionInput(channelID, 42, 5)}, executor.insertedActions)
}

func TestApplyTransitionLosesRevisionWithoutOutboxInsert(t *testing.T) {
    executor := newFakeLifecycleExecutor(matchState{Revision: 5})
    committed, err := executor.ApplyMatchStateTransition(ctx, transitionInput(channelID, 4, 100, snapshotJSON(5), nil))

    require.NoError(t, err)
    require.False(t, committed)
    require.Empty(t, executor.insertedActions)
}

func TestClaimReturnsOnlyEarliestUnfinishedActionPerMatch(t *testing.T) {
    executor := newFakeLifecycleExecutor(matchState{Revision: 2})
    executor.outbox = []outboxRow{
        pendingOutboxRow(channelID, 42, model.OutboxActionCreate, 10),
        pendingOutboxRow(channelID, 42, model.OutboxActionResolve, 11),
    }

    claimed, err := executor.ClaimPredictionActions(ctx, ClaimPredictionActionsInput{Limit: 10, Lease: time.Minute})

    require.NoError(t, err)
    require.Len(t, claimed, 1)
    require.Equal(t, model.OutboxActionCreate, claimed[0].Action)
}
```

- [ ] **Step 2: Run transaction tests and verify RED**

Run:

```bash
cd libs/repositories && go test ./dota/pgx/... -run 'TestApplyTransition|TestClaim'
```

Expected: FAIL because transition and claim methods do not exist.

- [ ] **Step 3: Implement the locked state transition**

Inside `ApplyMatchStateTransition`, use `p.trManager.Do` and this ordering:

```go
_, err := conn.Exec(ctx, `
    INSERT INTO dota_channel_match_states (channel_id, snapshot)
    VALUES ($1, $2)
    ON CONFLICT (channel_id) DO NOTHING`,
    input.ChannelID, idleSnapshotJSON(input.ChannelID),
)
if err != nil { return fmt.Errorf("ensure dota match state: %w", err) }

var errTransitionConflict = errors.New("dota match state revision conflict")

func idleSnapshotJSON(channelID uuid.UUID) []byte {
    return []byte(fmt.Sprintf(`{"channelId":%q,"state":"idle"}`, channelID.String()))
}

var revision int64
err = conn.QueryRow(ctx, `
    SELECT revision
    FROM dota_channel_match_states
    WHERE channel_id = $1
    FOR UPDATE`, input.ChannelID,
).Scan(&revision)
if err != nil { return fmt.Errorf("lock dota match state: %w", err) }
if revision != input.ExpectedRevision { return errTransitionConflict }

// UPDATE the locked row to revision + 1, then INSERT all validated actions.
```

Map `errTransitionConflict` to `(false, nil)` only after the transaction ends.
For a committed transition, update the row and insert every outbox action in
the same transaction. A failing action insert rolls back the state update.

- [ ] **Step 4: Implement ordered, token-owned claims**

Claim rows using one transaction with a stale-lock condition and an earlier
unfinished-action guard:

```sql
WITH candidates AS (
    SELECT o.id
    FROM dota_prediction_outbox o
    WHERE o.completed_at IS NULL
      AND o.available_at <= now()
      AND (o.locked_at IS NULL OR o.locked_at < now() - $1::interval)
      AND NOT EXISTS (
          SELECT 1
          FROM dota_prediction_outbox earlier
          WHERE earlier.channel_id = o.channel_id
            AND earlier.match_id = o.match_id
            AND earlier.sequence < o.sequence
            AND earlier.completed_at IS NULL
      )
    ORDER BY o.created_at
    FOR UPDATE SKIP LOCKED
    LIMIT $2
)
UPDATE dota_prediction_outbox o
SET locked_at = now(), lock_token = $3::uuid, attempts = attempts + 1
FROM candidates
WHERE o.id = candidates.id
RETURNING o.id, o.channel_id, o.match_id, o.action, o.sequence,
          o.payload, o.attempts, o.lock_token;
```

`CompletePredictionAction` and `RetryPredictionAction` must include both
action ID and lock token in their `WHERE` clause. Retry clears the token and
sets the supplied next availability time.

- [ ] **Step 5: Run repository tests and build**

Run:

```bash
cd libs/repositories && go test ./dota/... && go vet ./... && go build ./...
```

Expected: PASS. Tests prove a revision conflict creates no outbox row and a
later action cannot overtake its match's earlier action.

- [ ] **Step 6: Commit the repository implementation**

```bash
git add libs/repositories/dota
git commit -m "feat: persist ordered dota lifecycle actions"
```

### Task 3: Replace Local Redis State Machine With Repository State

**Files:**
- Modify: `apps/dota/internal/match/actions.go`
- Modify: `apps/dota/internal/match/state_machine.go`
- Modify: `apps/dota/internal/match/state_machine_test.go`
- Modify: `apps/dota/internal/processor/processor.go`
- Modify: `apps/dota/internal/processor/processor_test.go`
- Delete: `apps/dota/internal/match/redis_store.go`
- Delete: `apps/dota/internal/match/redis_store_test.go`

- [ ] **Step 1: Write failing transition tests against a fake repository state store**

Add tests before replacing transition code:

```go
func TestDelayedPostGameCannotSettleNewerMatch(t *testing.T) {
    repo := newFakeLifecycleRepository()
    machine := New(repo, emitter, logger)
    require.NoError(t, machine.Process(ctx, channelID, inGamePayloadAt(2006, 100)))
    require.NoError(t, machine.Process(ctx, channelID, inGamePayloadAt(2007, 200)))
    require.NoError(t, machine.Process(ctx, channelID, postGamePayloadAt(2006, 150, gsi.WinTeamRadiant)))

    require.Equal(t, int64(2007), repo.snapshot(channelID).MatchID)
    require.Empty(t, repo.actionsFor(channelID, 2007, ActionResolve))
}

func TestConcurrentTerminalInputsPersistOneAction(t *testing.T) {
    repo := newFakeLifecycleRepository()
    first := New(repo, emitter, logger)
    second := New(repo, emitter, logger)
    require.NoError(t, first.Process(ctx, channelID, inGamePayloadAt(2007, 100)))

    runConcurrently(
        func() { _ = first.Process(ctx, channelID, postGamePayloadAt(2007, 200, gsi.WinTeamRadiant)) },
        func() { _ = second.Process(ctx, channelID, gsi.Payload{Provider: gsi.Provider{Timestamp: 200}}) },
    )

    require.Len(t, repo.terminalActions(channelID, 2007), 1)
}

func TestStatsUpdatePreservesNewerMatchIdentity(t *testing.T) {
    repo := newFakeLifecycleRepository()
    machine := New(repo, emitter, logger)
    require.NoError(t, machine.Process(ctx, channelID, inGamePayloadAt(2007, 200)))
    require.NoError(t, machine.UpdateStats(ctx, channelID, 4200, 8, 3))

    snapshot := repo.snapshot(channelID)
    require.Equal(t, int64(2007), snapshot.MatchID)
    require.Equal(t, 4200, snapshot.Mmr)
}
```

- [ ] **Step 2: Run state tests and verify RED**

Run:

```bash
cd apps/dota && go test ./internal/match/... -run 'TestDelayedPostGame|TestConcurrentTerminal|TestStatsUpdate'
```

Expected: FAIL because the current machine uses process-local and Redis CAS state.

- [ ] **Step 3: Implement side-effect-free repository transitions**

Replace `channelState`, `kv.KV`, Redis CAS, and local mutex maps. The machine:

```go
for attempt := 0; attempt < maxTransitionRetries; attempt++ {
    current, err := m.repository.GetMatchState(ctx, channelID)
    if err != nil { return fmt.Errorf("get dota match state: %w", err) }

    next, actions, changed := m.transition(current, payload)
    if !changed { return nil }

    committed, err := m.repository.ApplyMatchStateTransition(ctx, toTransitionInput(current, next, actions))
    if err != nil { return fmt.Errorf("apply dota match transition: %w", err) }
    if committed {
        m.publishNonCriticalStartAndAbandonEvents(ctx, current, next, actions)
        return nil
    }
}
return errors.New("dota match transition retry limit reached")
```

`transition` accepts only newer `(provider timestamp, same-match game time)`
input. A post-game payload must match the current match ID. An inactive/map-less
payload can cancel only with strictly newer provider timestamp. Start returns a
create action; valid post-game returns resolve with `Win`; replacement returns
cancel(old) followed by create(new). The state machine never updates MMR/W/L
directly.

- [ ] **Step 4: Refactor async win probability and stats update**

`GetSnapshot` reads repository state. `UpdateWinProbability` and the later
settlement stats update use optimistic transitions that preserve the current
match ID, revision, and source-order fields. A stale async result for match A
must be a no-op when match B is current.

- [ ] **Step 5: Delete superseded Redis state store**

Delete only the unused authoritative Redis state files from Task 2. Keep
`apps/dota/internal/predictions/redis_store.go`, which owns Twitch prediction
dedupe and terminal claims.

- [ ] **Step 6: Run Dota state and processor tests**

Run:

```bash
cd apps/dota && go test ./internal/match/... && go test ./internal/processor/... && go build ./...
```

Expected: PASS.

- [ ] **Step 7: Commit repository-backed state machine**

```bash
git add apps/dota/internal/match apps/dota/internal/processor
git commit -m "fix: persist dota match state transactionally"
```

### Task 4: Add Token-Owned Prediction Terminal Claims and Direct APIs

**Files:**
- Modify: `apps/dota/internal/predictions/predictions.go`
- Modify: `apps/dota/internal/predictions/redis_store.go`
- Modify: `apps/dota/internal/predictions/predictions_test.go`
- Modify: `apps/dota/internal/predictions/redis_store_test.go`

- [ ] **Step 1: Write failing direct-action and terminal-race tests**

```go
func TestTerminalClaimAllowsOnlyOneTwitchEnd(t *testing.T) {
    fixture := newFixture(t)
    fixture.store.Store(predictionKey(fixture.channelID, 42), managedPrediction())
    fixture.client.getResponse = predictionResponse("prediction-42", "ACTIVE")

    runConcurrently(
        func() { _ = fixture.predictions.Resolve(ctx, resolveAction(fixture.channelID, 42, true)) },
        func() { _ = fixture.predictions.Cancel(ctx, cancelAction(fixture.channelID, 42)) },
    )

    require.Len(t, fixture.client.EndCalls(), 1)
}

func TestTerminalFailureReleasesClaimForOutboxRetry(t *testing.T) {
    fixture := newFixture(t)
    fixture.client.endErr = errors.New("temporary Twitch error")
    require.Error(t, fixture.predictions.Cancel(ctx, cancelAction(fixture.channelID, 42)))
    fixture.client.endErr = nil
    require.NoError(t, fixture.predictions.Cancel(ctx, cancelAction(fixture.channelID, 42)))
}
```

- [ ] **Step 2: Run tests and verify RED**

Run:

```bash
cd apps/dota && go test ./internal/predictions/... -run 'TestTerminalClaim|TestTerminalFailure'
```

Expected: FAIL because the current service is NATS-subscriber shaped.

- [ ] **Step 3: Implement direct idempotent methods**

Expose methods consumed by outbox workers:

```go
func (p *Predictions) Create(ctx context.Context, action match.LifecycleAction) error
func (p *Predictions) Resolve(ctx context.Context, action match.LifecycleAction) (dotamodel.ChannelDotaSettings, error)
func (p *Predictions) Cancel(ctx context.Context, action match.LifecycleAction) error
```

`Resolve` calls `ApplyMatchResultOnce` before resolving Twitch. Add
`ClaimTerminal`, `CompleteTerminal`, and `ReleaseTerminal` to the Redis Store;
each Lua cleanup checks the token. On retry, a non-active Twitch prediction is
a safe terminal cleanup rather than an error. Remove prediction lifecycle Core
NATS subscriptions but preserve the existing managed-title correlation logic.

- [ ] **Step 4: Run prediction tests and race detector**

Run:

```bash
cd apps/dota && go test ./internal/predictions/... && go test -race ./internal/predictions/...
```

Expected: PASS.

- [ ] **Step 5: Commit direct prediction actions**

```bash
git add apps/dota/internal/predictions
git commit -m "fix: claim durable dota terminal actions"
```

### Task 5: Implement PostgreSQL Outbox Worker

**Files:**
- Create: `apps/dota/internal/predictions/lifecycle_worker.go`
- Create: `apps/dota/internal/predictions/lifecycle_worker_test.go`
- Modify: `apps/dota/internal/predictions/predictions.go`
- Modify: `apps/dota/internal/match/state_machine.go`

- [ ] **Step 1: Write failing worker tests**

```go
func TestWorkerCompletesClaimedCreateAction(t *testing.T) {
    worker, repo, service := newWorkerFixture(t)
    action := claimedCreateAction(channelID, 42, lockToken)
    repo.claimed = []model.ClaimedOutboxAction{action}

    require.NoError(t, worker.ProcessOnce(ctx))
    require.Equal(t, []match.LifecycleAction{toLifecycleAction(action)}, service.createCalls)
    require.Equal(t, []claim{{action.ID, lockToken}}, repo.completed)
}

func TestWorkerRetriesFailureWithSameLockToken(t *testing.T) {
    worker, repo, service := newWorkerFixture(t)
    service.cancelErr = errors.New("temporary Twitch error")
    require.NoError(t, worker.ProcessOnce(ctx))
    require.Len(t, repo.retried, 1)
    require.True(t, repo.retried[0].availableAt.After(time.Now()))
}

func TestWorkerSkipsCreateForTerminalMatch(t *testing.T) {
    worker, repo, service := newWorkerFixture(t)
    action := claimedCreateAction(channelID, 42, lockToken)
    repo.claimed = []model.ClaimedOutboxAction{action}
    repo.states[channelID] = stateForDifferentMatch(channelID, 43)

    require.NoError(t, worker.ProcessOnce(ctx))
    require.Empty(t, service.createCalls)
    require.Equal(t, []claim{{action.ID, lockToken}}, repo.completed)
}
```

- [ ] **Step 2: Run worker tests and verify RED**

Run:

```bash
cd apps/dota && go test ./internal/predictions/... -run TestWorker
```

Expected: FAIL because the polling worker does not exist.

- [ ] **Step 3: Implement lifecycle and claims**

Use Fx lifecycle hooks with a cancellable poll loop. Each iteration claims a
bounded batch using `ClaimPredictionActions(ctx, ClaimPredictionActionsInput{
Limit: 10, Lease: 2 * time.Minute})`. For each action:

```go
func retryDelay(attempt int) time.Duration {
    if attempt < 1 {
        return time.Second
    }
    delay := time.Second << min(attempt-1, 6)
    if delay > time.Minute {
        return time.Minute
    }
    return delay
}

func toLifecycleAction(action model.ClaimedOutboxAction) (match.LifecycleAction, error) {
    var payload match.LifecycleAction
    if err := json.Unmarshal(action.Payload, &payload); err != nil {
        return match.LifecycleAction{}, fmt.Errorf("decode dota outbox payload: %w", err)
    }
    return payload, nil
}

if err := worker.dispatch(ctx, action); err != nil {
    next := time.Now().Add(retryDelay(action.Attempts))
    return repository.RetryPredictionAction(ctx, action.ID, action.LockToken, next)
}
return repository.CompletePredictionAction(ctx, action.ID, action.LockToken)
```

`dispatch` validates create actions against current repository state, calls
direct prediction APIs, and for resolve emits `MatchEnded` plus API state data
only after settlement succeeds. A bus publish failure is retryable, so the
outbox row remains unfinished.

- [ ] **Step 4: Run worker and Dota tests**

Run:

```bash
cd apps/dota && go test ./internal/predictions/... && go test ./...
```

Expected: PASS. Tests prove no terminal action overtakes creation, failures get
rescheduled, and stale leases are reclaimable by another worker.

- [ ] **Step 5: Commit durable worker**

```bash
git add apps/dota/internal/predictions apps/dota/internal/match
git commit -m "feat: process durable dota prediction outbox"
```

### Task 6: Wire Fx and Run Distributed Regression Coverage

**Files:**
- Modify: `apps/dota/app/app.go`
- Modify: `apps/dota/app/app_test.go`
- Modify: `apps/dota/internal/buslistener/bus_listener_test.go`
- Modify: `apps/dota/internal/chatalerts/chat_alerts_test.go`
- Modify: `apps/dota/internal/match/state_machine_test.go`
- Modify: `apps/dota/internal/predictions/lifecycle_worker_test.go`

- [ ] **Step 1: Write failing Fx wiring and sequence tests**

```go
func TestAppProvidesPostgresLifecycleWorker(t *testing.T) {
    app := fx.New(App, fx.Replace(newTestBus(), newTestPgxPool()), fx.Populate(&worker))
    require.NoError(t, app.Err())
    require.NotNil(t, worker)
}

func TestReplacementThenDelayedPostGameDoesNotSettleNewMatch(t *testing.T) {
    require.NoError(t, machine.Process(ctx, channelID, inGamePayloadAt(100, 100)))
    require.NoError(t, machine.Process(ctx, channelID, inGamePayloadAt(101, 200)))
    require.NoError(t, machine.Process(ctx, channelID, postGamePayloadAt(100, 150, gsi.WinTeamDire)))

    require.Equal(t, []model.OutboxAction{
        model.OutboxActionCancel,
        model.OutboxActionCreate,
    }, repository.actionKinds(channelID))
}

func TestResolveReplayUpdatesMmrOnceAndEndsPredictionOnce(t *testing.T) {
    action := claimedResolveAction(channelID, 42, true, lockToken)
    repository.claimed = []model.ClaimedOutboxAction{action}
    require.NoError(t, worker.ProcessOnce(ctx))
    repository.claimed = []model.ClaimedOutboxAction{action}
    require.NoError(t, worker.ProcessOnce(ctx))

    require.Equal(t, 1, repository.settlementCount(channelID, 42))
    require.Len(t, twitch.EndCalls(), 1)
}
```

- [ ] **Step 2: Run tests and verify RED**

Run:

```bash
cd apps/dota && go test ./app/... ./internal/match/... ./internal/predictions/...
```

Expected: new tests fail before final wiring.

- [ ] **Step 3: Wire repository state and worker**

In `app.go`, retain `dotarepositorypgx.NewFx` and provide the worker. Remove
the unused `match.NewRedisStateStore` provider and old prediction NATS
lifecycle invocation. Keep `BusEmitter`, ChatAlerts, GSI, and GetData listener
registered. Ensure `GetData` reads the repository-backed state snapshot.

- [ ] **Step 4: Run full verification**

Run:

```bash
cd apps/dota && go test -race ./internal/match/... ./internal/predictions/...
cd apps/dota && go test ./... && go vet ./... && go build ./...
cd libs/repositories && go test ./... && go vet ./... && go build ./...
cd libs/bus-core && go test ./... && go vet ./...
bun cli build app dota
rtk git status --short
```

Expected: all commands pass. Leave unrelated `cli/internal/cmds/proxy/proxy.go` untouched if it remains modified.

- [ ] **Step 5: Commit final wiring and regression coverage**

```bash
git add apps/dota libs/repositories
git commit -m "feat: wire transactional dota lifecycle"
```

## Final Verification

- [ ] `cd apps/dota && go test ./... && go vet ./... && go build ./...`
- [ ] `cd libs/repositories && go test ./... && go vet ./... && go build ./...`
- [ ] `cd libs/bus-core && go test ./... && go vet ./...`
- [ ] `bun cli build app dota`
- [ ] `git diff --check`
- [ ] `rtk git status --short`
- [ ] Update `docs/superpowers/specs/2026-07-21-dota-multi-replica-lifecycle-design.md` status to implemented and record any deliberate deviation.
