# Superseded: Dota Multi-Replica Lifecycle Implementation Plan

> Superseded on 2026-07-21 by `2026-07-21-dota-postgres-state-outbox.md` after review showed that a latest Redis snapshot cannot durably fence older queued actions across later stats updates.

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Make Dota GSI state transitions, match settlement, and Twitch prediction lifecycle processing safe across six service replicas.

**Architecture:** Redis is the authoritative live-state store. A Lua compare-and-swap writes each accepted snapshot revision and durable Redis Stream lifecycle action atomically. A Dota consumer group retries Stream actions until acknowledged. PostgreSQL records one settlement per `(channel_id, match_id)` and applies MMR/W/L in the same transaction, so action replay is safe.

**Tech Stack:** Go 1.26, Redis Lua and Streams (`XADD`, `XREADGROUP`, `XACK`, `XAUTOCLAIM`), go-redis/v9, pgx v5, PostgreSQL 18, Fx, Twitch Helix.

---

## File Structure

- `libs/migrations/postgres/<timestamp>_dota_match_settlements.sql`: durable idempotency ledger.
- `libs/repositories/dota/repository.go`: settlement repository contract.
- `libs/repositories/dota/pgx/pgx.go`: transactional `ApplyMatchResultOnce` implementation.
- `apps/dota/internal/match/actions.go`: Stream action contract shared by state producer and worker.
- `apps/dota/internal/match/redis_store.go`: authoritative snapshot load, revision-fenced CAS, and atomic `XADD` Lua implementation.
- `apps/dota/internal/match/state_machine.go`: pure transition loop over authoritative Redis state, with timestamp and match identity guards.
- `apps/dota/internal/predictions/lifecycle_worker.go`: durable Stream consumer-group worker.
- `apps/dota/internal/predictions/redis_store.go`: terminal action ownership claim in addition to existing create ownership.
- `apps/dota/internal/predictions/predictions.go`: direct create/resolve/cancel APIs used by the worker, no Core NATS lifecycle subscription.
- `apps/dota/app/app.go`: Fx wiring for the authoritative store and Stream worker.

### Task 1: Add Match Settlement Ledger and Repository Contract

**Files:**
- Create: `libs/migrations/postgres/<timestamp>_dota_match_settlements.sql`
- Modify: `libs/repositories/dota/repository.go`
- Modify: `libs/repositories/dota/pgx/pgx.go`
- Modify: `libs/repositories/dota/repository_test.go`
- Create: `libs/repositories/dota/pgx/pgx_test.go`

- [ ] **Step 1: Create the migration through the project CLI**

Verify the migration directory exists, then run:

```bash
ls libs/migrations/postgres
bun cli m create --name dota_match_settlements --db postgres --type sql
```

Fill the generated migration with an immutable settlement ledger. The primary key must make replay impossible and the channel-leading unique key supports channel lookup:

```sql
-- +goose Up
-- +goose StatementBegin
CREATE TABLE dota_match_settlements (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    channel_id UUID NOT NULL REFERENCES channels(id) ON DELETE CASCADE,
    match_id BIGINT NOT NULL,
    won BOOLEAN NOT NULL,
    mmr_delta INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (channel_id, match_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS dota_match_settlements;
-- +goose StatementEnd
```

- [ ] **Step 2: Write a failing repository-contract test**

Add an explicit input and expected behavior test in `repository_test.go` before implementation:

```go
func TestApplyMatchResultInputRejectsMissingMatchID(t *testing.T) {
    err := ValidateMatchResultInput(ApplyMatchResultInput{
        ChannelID: uuid.New(),
        MatchID:   0,
        Won:       true,
        MmrDelta:  25,
    })

    require.ErrorContains(t, err, "match ID")
}
```

Also add a pgx SQL test with a fake transactional executor that asserts this exact sequence: first delivery returns `RowsAffected() == 1`, performs one settings `UPDATE`, and returns the updated settings; replay returns `RowsAffected() == 0`, performs no settings `UPDATE`, and returns the existing settings.

- [ ] **Step 3: Run the new tests and verify RED**

Run:

```bash
cd libs/repositories && go test ./dota/... -run 'TestApplyMatchResult'
```

Expected: FAIL because `ApplyMatchResultInput`, validation, and the repository method do not exist.

- [ ] **Step 4: Define the contract and validate input**

Add these types and method to `repository.go`:

```go
type ApplyMatchResultInput struct {
    ChannelID uuid.UUID
    MatchID   int64
    Won       bool
    MmrDelta  int
}

func ValidateMatchResultInput(input ApplyMatchResultInput) error {
    if input.ChannelID == uuid.Nil {
        return errors.New("channel ID is required")
    }
    if input.MatchID <= 0 {
        return errors.New("match ID must be positive")
    }
    return nil
}

// Add this member to the existing Repository interface.
ApplyMatchResultOnce(context.Context, ApplyMatchResultInput) (model.ChannelDotaSettings, error)
```

- [ ] **Step 5: Implement transactional first-apply semantics**

Change `pgx.NewFx` to receive the transaction manager supplied by `baseapp`, then use one transaction. Insert the ledger first and update settings only when the insert returns a row:

```go
func (p *Pgx) ApplyMatchResultOnce(
    ctx context.Context,
    input dota.ApplyMatchResultInput,
) (model.ChannelDotaSettings, error) {
    if err := dota.ValidateMatchResultInput(input); err != nil {
        return model.Nil, err
    }

    var result model.ChannelDotaSettings
    err := p.trManager.Do(ctx, func(ctx context.Context) error {
        conn := p.getter.DefaultTrOrDB(ctx, p.pool)
        inserted, err := conn.Exec(ctx, `
            INSERT INTO dota_match_settlements (channel_id, match_id, won, mmr_delta)
            VALUES ($1, $2, $3, $4)
            ON CONFLICT (channel_id, match_id) DO NOTHING`,
            input.ChannelID, input.MatchID, input.Won, input.MmrDelta,
        )
        if err != nil {
            return fmt.Errorf("insert dota match settlement: %w", err)
        }
        if inserted.RowsAffected() == 1 {
            wins, losses := 0, 1
            if input.Won {
                wins, losses = 1, 0
            }
            if _, err := conn.Exec(ctx, `
                UPDATE channels_dota_settings
                SET mmr = mmr + $2,
                    session_wins = session_wins + $3,
                    session_losses = session_losses + $4,
                    updated_at = now()
                WHERE channel_id = $1`,
                input.ChannelID, input.MmrDelta, wins, losses,
            ); err != nil {
                return fmt.Errorf("update settled dota result: %w", err)
            }
        }
        settings, err := p.getByChannelIDWithConn(ctx, conn, input.ChannelID)
        if err != nil {
            return err
        }
        result = settings
        return nil
    })
    if err != nil {
        return model.Nil, fmt.Errorf("apply dota match result once: %w", err)
    }
    return result, nil
}
```

The update must set `updated_at = now()` and increment exactly one of
`session_wins` or `session_losses`. Do not call the old unconditional
`UpdateMatchResult` from new code.

- [ ] **Step 6: Run repository tests and build**

Run:

```bash
cd libs/repositories && go test ./dota/... && go build ./...
```

Expected: PASS. The SQL test proves first application updates once and replay fetches without a second update.

- [ ] **Step 7: Commit the ledger work**

```bash
git add libs/migrations/postgres libs/repositories/dota
git commit -m "feat: add idempotent dota match settlements"
```

### Task 2: Define Durable Lifecycle Actions and Redis CAS Store

**Files:**
- Create: `apps/dota/internal/match/actions.go`
- Create: `apps/dota/internal/match/redis_store.go`
- Create: `apps/dota/internal/match/redis_store_test.go`
- Modify: `apps/dota/internal/match/state_machine.go`
- Modify: `apps/dota/internal/match/state_machine_test.go`

- [ ] **Step 1: Write failing Redis-store tests**

Add tests using a narrow fake Redis script client. Cover all atomic boundaries:

```go
func TestCompareAndSwapWritesSnapshotAndCreateAction(t *testing.T) {
    store := newFakeAuthoritativeStore()
    current := Snapshot{ChannelID: channelID, State: StateIdle, Revision: 1}
    next := Snapshot{ChannelID: channelID, State: StateInGame, MatchID: 42, Revision: 2}

    committed, err := store.CompareAndSwap(ctx, current, next, []LifecycleAction{{
        Kind: ActionCreate, ChannelID: channelID, MatchID: 42, Revision: 2,
    }})

    require.NoError(t, err)
    require.True(t, committed)
    require.Equal(t, next, store.snapshot(channelID))
    require.Equal(t, []LifecycleAction{{Kind: ActionCreate, ChannelID: channelID, MatchID: 42, Revision: 2}}, store.actions())
}

func TestCompareAndSwapRejectsStaleRevisionWithoutAction(t *testing.T) {
    store := newFakeAuthoritativeStore()
    current := Snapshot{ChannelID: channelID, State: StateInGame, MatchID: 42, Revision: 2}
    store.seed(current)
    stale := Snapshot{ChannelID: channelID, State: StateIdle, Revision: 1}
    next := Snapshot{ChannelID: channelID, State: StateInGame, MatchID: 43, Revision: 2}

    committed, err := store.CompareAndSwap(ctx, stale, next, []LifecycleAction{{
        Kind: ActionCreate, ChannelID: channelID, MatchID: 43, Revision: 2,
    }})

    require.NoError(t, err)
    require.False(t, committed)
    require.Equal(t, current, store.snapshot(channelID))
    require.Empty(t, store.actions())
}
```

Add a test for replacement transition with ordered `cancel(old)` then `create(new)` actions in one CAS.

- [ ] **Step 2: Run the store tests and verify RED**

Run:

```bash
cd apps/dota && go test ./internal/match/... -run 'TestCompareAndSwap'
```

Expected: FAIL because `LifecycleAction`, `Revision`, and `CompareAndSwap` do not exist.

- [ ] **Step 3: Define the action and snapshot contracts**

In `actions.go`, define JSON-safe actions shared by producer and worker:

```go
type ActionKind string

const (
    ActionCreate  ActionKind = "create"
    ActionResolve ActionKind = "resolve"
    ActionCancel  ActionKind = "cancel"
)

type LifecycleAction struct {
    Kind      ActionKind `json:"kind"`
    ChannelID uuid.UUID  `json:"channelId"`
    MatchID   int64      `json:"matchId"`
    Revision  uint64     `json:"revision"`
    Win       bool       `json:"win,omitempty"`
    HeroName  string     `json:"heroName,omitempty"`
}
```

Extend `Snapshot` with `Revision uint64`, `LastProviderTimestamp int64`, and
`LastGameTime int`. Preserve JSON field names. Add this authoritative-store
interface in `redis_store.go`:

```go
type StateStore interface {
    Load(context.Context, uuid.UUID) (Snapshot, error)
    CompareAndSwap(context.Context, Snapshot, Snapshot, []LifecycleAction) (bool, error)
    UpdateStats(context.Context, uuid.UUID, int, int, int) error
}
```

`UpdateStats` must load/CAS-retry only the MMR/session fields and preserve a
newer match's identity and state.

- [ ] **Step 4: Implement one authoritative Lua CAS boundary**

`RedisStateStore` must use a single script that compares the raw expected JSON
state, writes the next JSON state with the six-hour TTL, and `XADD`s every
encoded action before returning success:

```lua
local current = redis.call("GET", KEYS[1])
if current ~= ARGV[1] then return 0 end
redis.call("SET", KEYS[1], ARGV[2], "PX", ARGV[3])
for i = 5, #ARGV do
  redis.call("XADD", KEYS[2], "MAXLEN", "~", ARGV[4], "*", "action", ARGV[i])
end
return 1
```

Pass the action payloads as JSON. Keep the Stream key and max length as named
constants. Pass `expectedState`, `nextState`, `ttlMilliseconds`, `maxLen`, then
zero or more action JSON values as `ARGV[1:]` in that order. `Load` returns an
idle snapshot when the key is absent. Do not use
`kv.KV` for authoritative mutation because it has no CAS primitive.

- [ ] **Step 5: Run the CAS tests and refactor only after GREEN**

Run:

```bash
cd apps/dota && go test ./internal/match/... -run 'TestCompareAndSwap'
```

Expected: PASS. Then run the complete package to identify old local-cache tests that need conversion:

```bash
cd apps/dota && go test ./internal/match/...
```

- [ ] **Step 6: Commit the store boundary**

```bash
git add apps/dota/internal/match
git commit -m "feat: add authoritative dota state store"
```

### Task 3: Refactor GSI State Transitions to Revision-Fenced Semantics

**Files:**
- Modify: `apps/dota/internal/match/state_machine.go`
- Modify: `apps/dota/internal/match/state_machine_test.go`
- Modify: `apps/dota/internal/processor/processor.go`
- Modify: `apps/dota/internal/processor/processor_test.go`

- [ ] **Step 1: Add failing multi-replica and delayed-payload tests**

Construct two state machines sharing one fake authoritative store. Add these
tests before changing transition code:

```go
func TestDelayedPostGameCannotSettleNewerMatch(t *testing.T) {
    store := newFakeAuthoritativeStore()
    first := New(repo, emitter, store, logger)
    second := New(repo, emitter, store, logger)
    require.NoError(t, first.Process(ctx, channelID, inGamePayloadAt(2006, 100)))
    require.NoError(t, second.Process(ctx, channelID, inGamePayloadAt(2007, 200)))

    require.NoError(t, first.Process(ctx, channelID, postGamePayloadAt(2006, 150, gsi.WinTeamRadiant)))

    snapshot, err := store.Load(ctx, channelID)
    require.NoError(t, err)
    require.Equal(t, int64(2007), snapshot.MatchID)
    require.NotContains(t, store.actions(), resolveAction(channelID, 2007, true))
}

func TestStaleReplicaCannotCancelSettledMatch(t *testing.T) {
    store := newFakeAuthoritativeStore()
    first := New(repo, emitter, store, logger)
    second := New(repo, emitter, store, logger)
    require.NoError(t, first.Process(ctx, channelID, inGamePayloadAt(2006, 100)))
    require.NoError(t, first.Process(ctx, channelID, postGamePayloadAt(2006, 200, gsi.WinTeamRadiant)))
    require.NoError(t, second.Process(ctx, channelID, inGamePayloadAt(2007, 300)))

    require.NotContains(t, store.actions(), cancelAction(channelID, 2006))
}

func TestMaplessPayloadCannotCancelWithOlderProviderTimestamp(t *testing.T) {
    require.NoError(t, machine.Process(ctx, channelID, inGamePayloadAt(2007, 300)))
    require.NoError(t, machine.Process(ctx, channelID, gsi.Payload{Provider: gsi.Provider{Timestamp: 299}}))

    snapshot, err := store.Load(ctx, channelID)
    require.NoError(t, err)
    require.Equal(t, int64(2007), snapshot.MatchID)
}

func TestConcurrentTerminalTransitionsEmitOneAction(t *testing.T) {
    require.NoError(t, machine.Process(ctx, channelID, inGamePayloadAt(2007, 100)))
    runConcurrently(
        func() { _ = first.Process(ctx, channelID, postGamePayloadAt(2007, 200, gsi.WinTeamRadiant)) },
        func() { _ = second.Process(ctx, channelID, gsi.Payload{Provider: gsi.Provider{Timestamp: 200}}) },
    )

    require.Len(t, terminalActions(store.actions(), channelID, 2007), 1)
}
```

- [ ] **Step 2: Run the targeted tests and verify RED**

Run:

```bash
cd apps/dota && go test ./internal/match/... -run 'TestDelayedPostGame|TestStaleReplica|TestMaplessPayload|TestConcurrentTerminal'
```

Expected: FAIL against the process-local cache and blind persistence.

- [ ] **Step 3: Replace process-local channels with a CAS retry loop**

Remove `channelState`, `channels`, and one-time `restore`. The state machine
loads a fresh snapshot on every mutation and retries only CAS conflicts:

```go
for attempt := 0; attempt < maxTransitionRetries; attempt++ {
    current, err := m.store.Load(ctx, channelID)
    if err != nil {
        return fmt.Errorf("load dota match state: %w", err)
    }
    next, actions, changed := m.transition(ctx, current, payload)
    if !changed {
        return nil
    }
    committed, err := m.store.CompareAndSwap(ctx, current, next, actions)
    if err != nil {
        return fmt.Errorf("commit dota match state: %w", err)
    }
    if committed {
        m.publishNonCriticalEvents(ctx, current, next, payload, actions)
        return nil
    }
}
return errors.New("dota match state changed too frequently")
```

`transition` must be side-effect free. It accepts newer source data only, uses
`payload.Provider.Timestamp` plus same-match game time for ordering, and emits
actions only through the returned slice.

- [ ] **Step 4: Apply strict match identity rules**

Implement these rules in `transition`:

```go
if payload.Map != nil && payload.Map.GameState == gsi.GameStatePostGame &&
    payload.Map.MatchID != current.MatchID {
    return current, nil, false
}

if payload.Map == nil || payload.Player == nil ||
    payload.Player.Activity != gsi.PlayerActivityPlaying {
    if payload.Provider.Timestamp <= current.LastProviderTimestamp {
        return current, nil, false
    }
    // Emit ActionCancel only for a current nonzero match.
}
```

For a newer in-game replacement, return `ActionCancel` for the old match
followed by `ActionCreate` for the new one. For a valid post-game result,
return one `ActionResolve`; do not update MMR in the state machine.

- [ ] **Step 5: Keep win-probability writes fenced**

Update `UpdateWinProbability` to load/CAS the snapshot and require the same
`expectedMatchID`. It must never overwrite a snapshot for a newer match.
Update processor tests so an async result from match A cannot change match B.

- [ ] **Step 6: Run state and processor tests**

Run:

```bash
cd apps/dota && go test ./internal/match/... && go test ./internal/processor/...
```

Expected: PASS, including existing GSI event and snapshot behavior.

- [ ] **Step 7: Commit state-machine changes**

```bash
git add apps/dota/internal/match apps/dota/internal/processor
git commit -m "fix: fence dota match transitions across replicas"
```

### Task 4: Add Atomic Prediction Terminal Claims and Direct Action APIs

**Files:**
- Modify: `apps/dota/internal/predictions/predictions.go`
- Modify: `apps/dota/internal/predictions/redis_store.go`
- Modify: `apps/dota/internal/predictions/predictions_test.go`
- Modify: `apps/dota/internal/predictions/redis_store_test.go`

- [ ] **Step 1: Write failing terminal-race tests**

Add tests that invoke resolve and cancel concurrently against one fake Store:

```go
func TestTerminalClaimLetsOnlyOneActionEndPrediction(t *testing.T) {
    fixture := newFixture(t)
    fixture.store.Store(predictionKey(fixture.channelID, 42), managedPrediction())
    fixture.client.getResponse = predictionResponse("prediction-42", "ACTIVE")

    runConcurrently(
        func() { _ = fixture.predictions.Resolve(ctx, resolveAction(fixture.channelID, 42, true)) },
        func() { _ = fixture.predictions.Cancel(ctx, cancelAction(fixture.channelID, 42)) },
    )

    require.Len(t, fixture.client.EndCalls(), 1)
}

func TestTerminalFailureReleasesClaimForRetry(t *testing.T) {
    fixture := newFixture(t)
    fixture.store.Store(predictionKey(fixture.channelID, 42), managedPrediction())
    fixture.client.getResponse = predictionResponse("prediction-42", "ACTIVE")
    fixture.client.endErr = errors.New("temporary Twitch failure")
    require.Error(t, fixture.predictions.Cancel(ctx, cancelAction(fixture.channelID, 42)))

    fixture.client.endErr = nil
    require.NoError(t, fixture.predictions.Cancel(ctx, cancelAction(fixture.channelID, 42)))
    require.Len(t, fixture.client.EndCalls(), 2)
}
```

Also add direct API tests for `Create`, `Resolve`, and `Cancel` without NATS
subscription wrappers.

- [ ] **Step 2: Run terminal tests and verify RED**

Run:

```bash
cd apps/dota && go test ./internal/predictions/... -run 'TestTerminalClaim|TestTerminalFailure|TestCreate|TestResolve|TestCancel'
```

Expected: FAIL because terminal claims and direct APIs do not exist.

- [ ] **Step 3: Extend the Store with ownership-scoped terminal methods**

Use a separate terminal key. The Lua claim must use `SET NX PX`; completion and
release must compare the token before deleting:

```go
type TerminalClaimResult int

const (
    TerminalClaimed TerminalClaimResult = iota
    TerminalClaimBusy
    TerminalClaimComplete
)

type Store interface {
    ClaimTerminal(context.Context, string, string, time.Duration) (TerminalClaimResult, error)
    CompleteTerminal(context.Context, string, string) error
    ReleaseTerminal(context.Context, string, string) error
}
```

Do not let a failed worker delete another worker's claim. A timeout after
Twitch accepts an end is safe because retry fetches the prediction, sees it no
longer active, and completes the record cleanup.

- [ ] **Step 4: Remove prediction lifecycle Core NATS subscriptions**

Replace `handleMatchStarted`, `handleMatchEnded`, and `handleMatchAbandoned`
with direct methods consumed by the Stream worker:

```go
func (p *Predictions) Create(ctx context.Context, action match.LifecycleAction) error
func (p *Predictions) Resolve(ctx context.Context, action match.LifecycleAction) (dotamodel.ChannelDotaSettings, error)
func (p *Predictions) Cancel(ctx context.Context, action match.LifecycleAction) error
```

`Resolve` calls `ApplyMatchResultOnce` before resolving Twitch. It returns the
settled settings on both first processing and replay. `Create` must not create
after a terminal state; that validation is supplied by the worker.

- [ ] **Step 5: Run predictions tests and race detector**

Run:

```bash
cd apps/dota && go test ./internal/predictions/... && go test -race ./internal/predictions/...
```

Expected: PASS. Exactly one terminal Twitch call is observable in the race test.

- [ ] **Step 6: Commit terminal ownership**

```bash
git add apps/dota/internal/predictions
git commit -m "fix: claim dota prediction terminal actions"
```

### Task 5: Implement Redis Stream Lifecycle Worker

**Files:**
- Create: `apps/dota/internal/predictions/lifecycle_worker.go`
- Create: `apps/dota/internal/predictions/lifecycle_worker_test.go`
- Modify: `apps/dota/internal/predictions/predictions.go`
- Modify: `apps/dota/internal/match/redis_store.go`

- [ ] **Step 1: Write failing worker behavior tests**

Use a narrow Stream client interface and fakes; do not require a live Redis
server. Cover acknowledgement, retry, ordering validation, and takeover:

```go
func TestWorkerAcknowledgesOnlySuccessfulAction(t *testing.T) {
    worker, stream, service := newWorkerFixture(t)
    action := createAction(channelID, 42)
    service.createErr = errors.New("temporary Twitch failure")
    require.Error(t, worker.handleEntry(ctx, stream.entry(action)))
    require.Empty(t, stream.acked())

    service.createErr = nil
    require.NoError(t, worker.handleEntry(ctx, stream.entry(action)))
    require.Equal(t, []string{stream.entryID(action)}, stream.acked())
}

func TestWorkerSkipsLateCreateAfterTerminalState(t *testing.T) {
    worker, stream, service := newWorkerFixture(t)
    stream.state.seed(Snapshot{ChannelID: channelID, State: StateIdle, Revision: 3})
    action := createAction(channelID, 42)

    require.NoError(t, worker.handleEntry(ctx, stream.entry(action)))
    require.Empty(t, service.createCalls)
    require.Equal(t, []string{stream.entryID(action)}, stream.acked())
}

func TestWorkerClaimsAbandonedPendingEntry(t *testing.T) {
    worker, stream, service := newWorkerFixture(t)
    action := cancelAction(channelID, 42)
    stream.pendingForDeadConsumer(action)

    require.NoError(t, worker.claimPending(ctx))
    require.Equal(t, []match.LifecycleAction{action}, service.cancelCalls)
    require.Equal(t, []string{stream.entryID(action)}, stream.acked())
}
```

- [ ] **Step 2: Run worker tests and verify RED**

Run:

```bash
cd apps/dota && go test ./internal/predictions/... -run 'TestWorker'
```

Expected: FAIL because the Stream worker does not exist.

- [ ] **Step 3: Implement consumer-group lifecycle**

Create the group with `XGroupCreateMkStream`, ignoring `BUSYGROUP`. Use a
unique consumer name per process, a stop context, and a wait group. The read
loop uses new entries and a bounded block:

```go
streams, err := client.XReadGroup(ctx, &redis.XReadGroupArgs{
    Group:    lifecycleConsumerGroup,
    Consumer: worker.consumer,
    Streams:  []string{match.LifecycleStreamKey, ">"},
    Count:    lifecycleReadCount,
    Block:    lifecycleReadBlock,
}).Result()
```

For every action, call the direct `Predictions` method. Call `XACK` only after
success or a safe no-op. Return errors to leave the entry in the pending list.

- [ ] **Step 4: Recover crashed consumers**

On a ticker, call `XAutoClaim` with a conservative idle threshold. Process
claimed messages through the same handler and acknowledge only after success:

```go
messages, next, err := client.XAutoClaim(ctx, &redis.XAutoClaimArgs{
    Stream:   match.LifecycleStreamKey,
    Group:    lifecycleConsumerGroup,
    Consumer: worker.consumer,
    MinIdle:  lifecycleClaimIdle,
    Start:    cursor,
    Count:    lifecycleClaimCount,
}).Result()
```

Do not use Core NATS for prediction correctness. The worker should still let
the bus handle existing chat-alert consumers independently.

- [ ] **Step 5: Publish settled state and MatchEnded after successful resolve**

After `Resolve` returns settings, publish `bus.Dota.MatchEnded` and
`bus.Api.DotaStateUpdate` with the settled MMR/session values. Treat a publish
failure as retryable so the Stream entry is not acknowledged. The action must
be idempotent: a retry settles via the ledger, finds Twitch terminal/non-active,
and safely retries the publish.

- [ ] **Step 6: Run worker and full Dota tests**

Run:

```bash
cd apps/dota && go test ./internal/predictions/... && go test ./...
```

Expected: PASS. Worker tests prove failure remains pending and stale consumer
entries are reclaimed.

- [ ] **Step 7: Commit durable worker**

```bash
git add apps/dota/internal/predictions apps/dota/internal/match
git commit -m "feat: add durable dota prediction worker"
```

### Task 6: Wire Fx and Preserve Public Behavior

**Files:**
- Modify: `apps/dota/app/app.go`
- Modify: `apps/dota/app/app_test.go`
- Modify: `apps/dota/internal/buslistener/bus_listener_test.go`
- Modify: `apps/dota/internal/chatalerts/chat_alerts_test.go`

- [ ] **Step 1: Add failing app-wiring test**

Add a test that constructs the Dota Fx module with fakes and asserts the
authoritative state store and lifecycle worker are registered:

```go
func TestAppProvidesLifecycleWorker(t *testing.T) {
    app := fx.New(
        App,
        fx.Replace(newTestRedis(), newTestBus(), newTestPgxPool()),
        fx.Populate(&worker),
    )
    require.NoError(t, app.Err())
    require.NotNil(t, worker)
}
```

- [ ] **Step 2: Run wiring test and verify RED**

Run:

```bash
cd apps/dota && go test ./app/... -run TestAppProvidesLifecycleWorker
```

Expected: FAIL before worker and store wiring exists.

- [ ] **Step 3: Wire concrete dependencies**

In `app.go`, provide:

```go
fx.Annotate(match.NewRedisStateStore, fx.As(new(match.StateStore))),
fx.Annotate(predictions.NewRedisPredictionStore, fx.As(new(predictions.Store))),
predictions.NewLifecycleWorker,
```

Invoke the worker so its Fx lifecycle starts. Remove the old prediction NATS
subscription lifecycle invocation. Keep `BusEmitter`, `ChatAlerts`, GSI server,
and GetData listener wired as before.

- [ ] **Step 4: Update compatibility tests**

Ensure existing chat alerts still subscribe to `MatchStarted`, `MatchEnded`,
and `MatchAbandoned`. Verify GetData continues to return current snapshot data
after a CAS transition. Keep bus Gob compatibility tests unchanged unless a
new action type is intentionally transported over Gob.

- [ ] **Step 5: Run service build and focused verification**

Run:

```bash
cd apps/dota && go test ./... && go vet ./... && go build ./...
cd libs/repositories && go test ./... && go vet ./...
```

Expected: PASS.

- [ ] **Step 6: Commit wiring**

```bash
git add apps/dota libs/repositories
git commit -m "feat: wire durable dota lifecycle processing"
```

### Task 7: End-to-End Concurrency Regression Verification

**Files:**
- Modify: `apps/dota/internal/match/state_machine_test.go`
- Modify: `apps/dota/internal/predictions/lifecycle_worker_test.go`
- Modify: `apps/dota/internal/predictions/predictions_test.go`

- [ ] **Step 1: Add full failure-sequence tests**

Write tests for complete sequences, not source-text assertions:

```go
func TestReplacementThenDelayedPostGameLeavesNewMatchUntouched(t *testing.T) {
    require.NoError(t, machine.Process(ctx, channelID, inGamePayloadAt(100, 100)))
    require.NoError(t, machine.Process(ctx, channelID, inGamePayloadAt(101, 200)))
    require.NoError(t, machine.Process(ctx, channelID, postGamePayloadAt(100, 150, gsi.WinTeamDire)))

    require.Equal(t, []match.LifecycleAction{
        cancelAction(channelID, 100),
        createAction(channelID, 101),
    }, store.actions())
}

func TestResolveActionReplaySettlesAndEndsExactlyOnce(t *testing.T) {
    action := resolveAction(channelID, 42, true)
    require.NoError(t, worker.handleEntry(ctx, stream.entry(action)))
    require.NoError(t, worker.handleEntry(ctx, stream.entry(action)))

    require.Equal(t, 1, repository.appliedMatches(channelID, 42))
    require.Len(t, twitch.EndCalls(), 1)
}
```

- [ ] **Step 2: Run tests and verify any new assertions are RED first**

Run:

```bash
cd apps/dota && go test ./internal/match/... ./internal/predictions/...
```

Expected: new tests fail before their final implementation fixes, then pass.

- [ ] **Step 3: Run race and broad verification**

Run:

```bash
cd apps/dota && go test -race ./internal/match/... ./internal/predictions/...
cd apps/dota && go test ./...
cd libs/bus-core && go test ./... && go vet ./...
cd libs/repositories && go test ./... && go vet ./...
bun cli build app dota
git diff --check
rtk git status --short
```

Expected: all commands pass and the worktree is clean.

- [ ] **Step 4: Commit regression coverage**

```bash
git add apps/dota
git commit -m "test: cover distributed dota lifecycle recovery"
```

## Final Verification

- [ ] `cd apps/dota && go test ./... && go vet ./... && go build ./...`
- [ ] `cd libs/repositories && go test ./... && go vet ./...`
- [ ] `cd libs/bus-core && go test ./... && go vet ./...`
- [ ] `bun cli build app dota`
- [ ] `git diff --check`
- [ ] `rtk git status --short`
- [ ] Update `docs/superpowers/specs/2026-07-21-dota-multi-replica-lifecycle-design.md` from approved to implemented, recording any deliberate deviations.
