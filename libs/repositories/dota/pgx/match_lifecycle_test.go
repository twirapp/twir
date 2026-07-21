package pgx

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	jackcpgx "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/twirapp/twir/libs/repositories/dota"
	"github.com/twirapp/twir/libs/repositories/dota/model"
)

func TestApplyMatchStateTransition(t *testing.T) {
	channelID := uuid.New()

	t.Run("inserts state and actions then increments revision", func(t *testing.T) {
		executor := &matchLifecycleExecutorFake{}
		input := lifecycleTransitionInput(channelID, 0)

		applied, err := (&Pgx{}).applyMatchStateTransition(context.Background(), executor, input)
		if err != nil {
			t.Fatalf("applyMatchStateTransition() error = %v", err)
		}
		if !applied {
			t.Fatal("applyMatchStateTransition() applied = false, want true")
		}
		if !executor.hasState {
			t.Fatal("state was not inserted")
		}
		if got, want := executor.state.Revision, int64(1); got != want {
			t.Errorf("state revision = %d, want %d", got, want)
		}
		if got, want := executor.state.ProviderTimestamp, input.ProviderTimestamp; got != want {
			t.Errorf("provider timestamp = %d, want %d", got, want)
		}
		if got, want := string(executor.state.Snapshot), string(input.Snapshot); got != want {
			t.Errorf("state snapshot = %s, want %s", got, want)
		}
		if got, want := executor.actions, input.Actions; !equalOutboxActions(got, want) {
			t.Errorf("inserted actions = %#v, want %#v", got, want)
		}
	})

	t.Run("does not insert actions when the revision conflicts", func(t *testing.T) {
		executor := &matchLifecycleExecutorFake{
			hasState: true,
			state: model.MatchState{
				ChannelID: channelID,
				Revision:  1,
				Snapshot:  json.RawMessage(`{"state":"idle"}`),
			},
		}

		applied, err := (&Pgx{}).applyMatchStateTransition(
			context.Background(),
			executor,
			lifecycleTransitionInput(channelID, 0),
		)
		if err != nil {
			t.Fatalf("applyMatchStateTransition() error = %v", err)
		}
		if applied {
			t.Fatal("applyMatchStateTransition() applied = true, want false")
		}
		if got, want := executor.state.Revision, int64(1); got != want {
			t.Errorf("state revision = %d, want %d", got, want)
		}
		if len(executor.actions) != 0 {
			t.Errorf("inserted actions = %#v, want none", executor.actions)
		}
	})

	t.Run("rolls back the state when an action insert fails", func(t *testing.T) {
		store := matchLifecycleStore{executor: matchLifecycleExecutorFake{
			outboxInsertErr: errors.New("outbox insert failed"),
		}}

		err := store.transaction(func(executor *matchLifecycleExecutorFake) error {
			_, err := (&Pgx{}).applyMatchStateTransition(
				context.Background(),
				executor,
				lifecycleTransitionInput(channelID, 0),
			)
			return err
		})
		if err == nil {
			t.Fatal("transaction error = nil, want action insert error")
		}
		if store.executor.hasState {
			t.Errorf("state persisted after rollback = %#v, want none", store.executor.state)
		}
		if len(store.executor.actions) != 0 {
			t.Errorf("actions persisted after rollback = %#v, want none", store.executor.actions)
		}
	})
}

func TestScanMatchStateReturnsIdleStateWhenNoRowExists(t *testing.T) {
	channelID := uuid.New()

	state, err := (&Pgx{}).scanMatchState(matchLifecycleRow{err: jackcpgx.ErrNoRows}, channelID)
	if err != nil {
		t.Fatalf("scanMatchState() error = %v", err)
	}
	if got, want := state.ChannelID, channelID; got != want {
		t.Errorf("idle channel ID = %s, want %s", got, want)
	}
	if state.Revision != 0 || state.ProviderTimestamp != 0 {
		t.Errorf("idle state = %#v, want zero revision and provider timestamp", state)
	}
	if got, want := string(state.Snapshot), `{}`; got != want {
		t.Errorf("idle snapshot = %s, want %s", got, want)
	}
}

func TestClaimPredictionActions(t *testing.T) {
	channelID := uuid.New()
	now := time.Date(2026, time.July, 21, 2, 0, 0, 0, time.UTC)
	input := dota.ClaimPredictionActionsInput{Limit: 2, Lease: time.Minute}

	t.Run("returns create then resolve after create completes", func(t *testing.T) {
		create := testOutboxAction(channelID, 42, model.OutboxActionCreate, 10, now)
		resolve := testOutboxAction(channelID, 42, model.OutboxActionResolve, 11, now.Add(time.Second))
		executor := newPredictionActionStoreExecutor(now, create, resolve)

		first, err := (&Pgx{}).claimPredictionActions(context.Background(), executor, input)
		if err != nil {
			t.Fatalf("first claimPredictionActions() error = %v", err)
		}
		if len(first) != 1 {
			t.Fatalf("first claimed actions = %#v, want only create", first)
		}
		if got := first[0]; got.ID != create.ID || got.Action != model.OutboxActionCreate || got.Sequence != 10 {
			t.Errorf("first claimed action = %#v, want create sequence 10", got)
		}
		if !strings.Contains(executor.query, "DISTINCT ON (channel_id, match_id)") {
			t.Errorf("claim query does not select the first action per match: %s", executor.query)
		}
		if !strings.Contains(executor.query, "FOR UPDATE OF outbox SKIP LOCKED") {
			t.Errorf("claim query does not lock candidates with SKIP LOCKED: %s", executor.query)
		}

		if err := (&Pgx{}).completePredictionAction(context.Background(), executor, first[0].ID, first[0].LockToken); err != nil {
			t.Fatalf("completePredictionAction() error = %v", err)
		}

		second, err := (&Pgx{}).claimPredictionActions(context.Background(), executor, input)
		if err != nil {
			t.Fatalf("second claimPredictionActions() error = %v", err)
		}
		if len(second) != 1 {
			t.Fatalf("second claimed actions = %#v, want only resolve", second)
		}
		if got := second[0]; got.ID != resolve.ID || got.Action != model.OutboxActionResolve || got.Sequence != 11 {
			t.Errorf("second claimed action = %#v, want resolve sequence 11", got)
		}
	})

	t.Run("does not reclaim a non-expired lease", func(t *testing.T) {
		action := testOutboxAction(channelID, 43, model.OutboxActionCreate, 10, now)
		lockToken := uuid.New()
		lockedAt := now.Add(-30 * time.Second)
		action.Attempts = 1
		action.LockToken = lockToken
		action.LockedAt = &lockedAt
		executor := newPredictionActionStoreExecutor(now, action)

		claimed, err := (&Pgx{}).claimPredictionActions(context.Background(), executor, input)
		if err != nil {
			t.Fatalf("claimPredictionActions() error = %v", err)
		}
		if len(claimed) != 0 {
			t.Fatalf("claimed actions = %#v, want none while lease is active", claimed)
		}
		stored := executor.action(action.ID)
		if stored.LockToken != lockToken || stored.Attempts != 1 || stored.CompletedAt != nil {
			t.Errorf("stored action changed during active lease = %#v", stored)
		}
	})

	t.Run("reclaims an expired lease with a new token and incremented attempts", func(t *testing.T) {
		action := testOutboxAction(channelID, 44, model.OutboxActionCreate, 10, now)
		previousToken := uuid.New()
		lockedAt := now.Add(-61 * time.Second)
		action.Attempts = 1
		action.LockToken = previousToken
		action.LockedAt = &lockedAt
		executor := newPredictionActionStoreExecutor(now, action)

		claimed, err := (&Pgx{}).claimPredictionActions(context.Background(), executor, input)
		if err != nil {
			t.Fatalf("claimPredictionActions() error = %v", err)
		}
		if len(claimed) != 1 {
			t.Fatalf("claimed actions = %#v, want one reclaimed action", claimed)
		}
		if got := claimed[0]; got.LockToken == previousToken || got.Attempts != 2 {
			t.Errorf("reclaimed action = %#v, want new token and attempts 2", got)
		}
		stored := executor.action(action.ID)
		if stored.LockToken != claimed[0].LockToken || stored.Attempts != claimed[0].Attempts {
			t.Errorf("stored reclaimed action = %#v, want %#v", stored, claimed[0])
		}
	})
}

func TestPredictionActionOwnershipLoss(t *testing.T) {
	now := time.Date(2026, time.July, 21, 2, 0, 0, 0, time.UTC)
	action := testOutboxAction(uuid.New(), 42, model.OutboxActionCreate, 10, now)
	ownerToken := uuid.New()
	lockedAt := now.Add(-time.Second)
	action.Attempts = 1
	action.LockToken = ownerToken
	action.LockedAt = &lockedAt
	executor := newPredictionActionStoreExecutor(now, action)
	wrongToken := uuid.New()

	t.Run("completion with a wrong token reports ownership loss and leaves the row unfinished", func(t *testing.T) {
		err := (&Pgx{}).completePredictionAction(context.Background(), executor, action.ID, wrongToken)
		if !errors.Is(err, dota.ErrPredictionActionOwnershipLost) {
			t.Fatalf("completePredictionAction() error = %v, want ownership loss", err)
		}
		stored := executor.action(action.ID)
		if stored.CompletedAt != nil || stored.LockToken != ownerToken {
			t.Errorf("stored action changed after failed completion = %#v", stored)
		}
	})

	t.Run("retry with a wrong token reports ownership loss and leaves the row unfinished", func(t *testing.T) {
		availableAt := now.Add(time.Minute)
		err := (&Pgx{}).retryPredictionAction(
			context.Background(),
			executor,
			action.ID,
			wrongToken,
			availableAt,
		)
		if !errors.Is(err, dota.ErrPredictionActionOwnershipLost) {
			t.Fatalf("retryPredictionAction() error = %v, want ownership loss", err)
		}
		stored := executor.action(action.ID)
		if stored.CompletedAt != nil || stored.LockToken != ownerToken || !stored.AvailableAt.Equal(action.AvailableAt) {
			t.Errorf("stored action changed after failed retry = %#v", stored)
		}
	})
}

func lifecycleTransitionInput(channelID uuid.UUID, expectedRevision int64) dota.ApplyMatchStateTransitionInput {
	return dota.ApplyMatchStateTransitionInput{
		ChannelID:         channelID,
		ExpectedRevision:  expectedRevision,
		ProviderTimestamp: 1_700_000_000,
		Snapshot:          json.RawMessage(`{"state":"active"}`),
		Actions: []model.OutboxActionInput{
			{
				ChannelID: channelID,
				MatchID:   42,
				Action:    model.OutboxActionCreate,
				Sequence:  10,
				Payload:   json.RawMessage(`{"kind":"create"}`),
			},
			{
				ChannelID: channelID,
				MatchID:   42,
				Action:    model.OutboxActionResolve,
				Sequence:  11,
				Payload:   json.RawMessage(`{"kind":"resolve"}`),
			},
		},
	}
}

func equalOutboxActions(got, want []model.OutboxActionInput) bool {
	if len(got) != len(want) {
		return false
	}

	for index := range got {
		if got[index].ChannelID != want[index].ChannelID ||
			got[index].MatchID != want[index].MatchID ||
			got[index].Action != want[index].Action ||
			got[index].Sequence != want[index].Sequence ||
			string(got[index].Payload) != string(want[index].Payload) {
			return false
		}
	}

	return true
}

type matchLifecycleStore struct {
	executor matchLifecycleExecutorFake
}

func (s *matchLifecycleStore) transaction(fn func(*matchLifecycleExecutorFake) error) error {
	tx := s.executor.clone()
	if err := fn(&tx); err != nil {
		return err
	}

	s.executor = tx
	return nil
}

type matchLifecycleExecutorFake struct {
	hasState        bool
	state           model.MatchState
	actions         []model.OutboxActionInput
	outboxInsertErr error
}

func (f *matchLifecycleExecutorFake) clone() matchLifecycleExecutorFake {
	clone := *f
	clone.state.Snapshot = append(json.RawMessage(nil), f.state.Snapshot...)
	clone.actions = append([]model.OutboxActionInput(nil), f.actions...)
	return clone
}

func (f *matchLifecycleExecutorFake) Exec(
	_ context.Context,
	query string,
	arguments ...any,
) (pgconn.CommandTag, error) {
	switch {
	case strings.Contains(query, "INSERT INTO dota_channel_match_states"):
		if !f.hasState {
			f.hasState = true
			f.state = model.MatchState{
				ChannelID: arguments[0].(uuid.UUID),
				Snapshot:  json.RawMessage(`{}`),
			}
		}
		return pgconn.NewCommandTag("INSERT 0 1"), nil
	case strings.Contains(query, "UPDATE dota_channel_match_states"):
		if !f.hasState {
			return pgconn.CommandTag{}, errors.New("state does not exist")
		}
		f.state.Revision = arguments[1].(int64)
		f.state.ProviderTimestamp = arguments[2].(int64)
		f.state.Snapshot = json.RawMessage(arguments[3].(string))
		return pgconn.NewCommandTag("UPDATE 1"), nil
	case strings.Contains(query, "INSERT INTO dota_prediction_outbox"):
		if f.outboxInsertErr != nil {
			return pgconn.CommandTag{}, f.outboxInsertErr
		}
		f.actions = append(f.actions, model.OutboxActionInput{
			ChannelID: arguments[0].(uuid.UUID),
			MatchID:   arguments[1].(int64),
			Action:    model.OutboxAction(arguments[2].(string)),
			Sequence:  arguments[3].(int64),
			Payload:   json.RawMessage(arguments[4].(string)),
		})
		return pgconn.NewCommandTag("INSERT 0 1"), nil
	default:
		return pgconn.CommandTag{}, fmt.Errorf("unexpected exec query: %s", query)
	}
}

func (f *matchLifecycleExecutorFake) QueryRow(
	_ context.Context,
	query string,
	_ ...any,
) jackcpgx.Row {
	if !strings.Contains(query, "SELECT revision") {
		return matchLifecycleRow{err: fmt.Errorf("unexpected query row query: %s", query)}
	}
	if !f.hasState {
		return matchLifecycleRow{err: jackcpgx.ErrNoRows}
	}

	return matchLifecycleRow{revision: f.state.Revision}
}

type matchLifecycleRow struct {
	revision int64
	err      error
}

func (r matchLifecycleRow) Scan(dest ...any) error {
	if r.err != nil {
		return r.err
	}
	if len(dest) != 1 {
		return fmt.Errorf("scan destinations = %d, want 1", len(dest))
	}

	*dest[0].(*int64) = r.revision
	return nil
}

type predictionActionStoreExecutorFake struct {
	now     time.Time
	actions []storedPredictionAction
	query   string
}

type storedPredictionAction struct {
	ID uuid.UUID
	model.OutboxActionInput
	Attempts    int
	AvailableAt time.Time
	LockedAt    *time.Time
	LockToken   uuid.UUID
	CompletedAt *time.Time
	CreatedAt   time.Time
}

func newPredictionActionStoreExecutor(
	now time.Time,
	actions ...storedPredictionAction,
) *predictionActionStoreExecutorFake {
	return &predictionActionStoreExecutorFake{
		now:     now,
		actions: append([]storedPredictionAction(nil), actions...),
	}
}

func testOutboxAction(
	channelID uuid.UUID,
	matchID int64,
	action model.OutboxAction,
	sequence int64,
	createdAt time.Time,
) storedPredictionAction {
	return storedPredictionAction{
		ID: uuid.New(),
		OutboxActionInput: model.OutboxActionInput{
			ChannelID: channelID,
			MatchID:   matchID,
			Action:    action,
			Sequence:  sequence,
			Payload:   json.RawMessage(`{"kind":"prediction"}`),
		},
		AvailableAt: createdAt.Add(-time.Second),
		CreatedAt:   createdAt,
	}
}

func (f *predictionActionStoreExecutorFake) Query(
	_ context.Context,
	query string,
	arguments ...any,
) (jackcpgx.Rows, error) {
	f.query = query
	if len(arguments) != 3 {
		return nil, fmt.Errorf("claim arguments = %d, want 3", len(arguments))
	}

	leaseMicros, ok := arguments[0].(int64)
	if !ok {
		return nil, fmt.Errorf("claim lease = %T, want int64", arguments[0])
	}
	limit, ok := arguments[1].(int)
	if !ok {
		return nil, fmt.Errorf("claim limit = %T, want int", arguments[1])
	}
	lockToken, ok := arguments[2].(uuid.UUID)
	if !ok {
		return nil, fmt.Errorf("claim lock token = %T, want uuid.UUID", arguments[2])
	}

	lease := time.Duration(leaseMicros) * time.Microsecond
	earliestByMatch := make(map[struct {
		channelID uuid.UUID
		matchID   int64
	}]*storedPredictionAction)
	for index := range f.actions {
		action := &f.actions[index]
		if action.CompletedAt != nil {
			continue
		}

		key := struct {
			channelID uuid.UUID
			matchID   int64
		}{channelID: action.ChannelID, matchID: action.MatchID}
		current, exists := earliestByMatch[key]
		if !exists || action.Sequence < current.Sequence ||
			(action.Sequence == current.Sequence && action.CreatedAt.Before(current.CreatedAt)) {
			earliestByMatch[key] = action
		}
	}

	candidates := make([]*storedPredictionAction, 0, len(earliestByMatch))
	for _, action := range earliestByMatch {
		if action.AvailableAt.After(f.now) {
			continue
		}
		if action.LockedAt != nil && !action.LockedAt.Before(f.now.Add(-lease)) {
			continue
		}
		candidates = append(candidates, action)
	}
	sort.Slice(candidates, func(left, right int) bool {
		if !candidates[left].AvailableAt.Equal(candidates[right].AvailableAt) {
			return candidates[left].AvailableAt.Before(candidates[right].AvailableAt)
		}
		if candidates[left].Sequence != candidates[right].Sequence {
			return candidates[left].Sequence < candidates[right].Sequence
		}
		return candidates[left].CreatedAt.Before(candidates[right].CreatedAt)
	})
	if len(candidates) > limit {
		candidates = candidates[:limit]
	}

	claimed := make([]model.ClaimedOutboxAction, 0, len(candidates))
	for _, action := range candidates {
		lockedAt := f.now
		action.LockedAt = &lockedAt
		action.LockToken = lockToken
		action.Attempts++
		claimed = append(claimed, model.ClaimedOutboxAction{
			ID:                action.ID,
			LockToken:         action.LockToken,
			OutboxActionInput: action.OutboxActionInput,
			Attempts:          action.Attempts,
		})
	}

	return &claimedOutboxRows{actions: claimed}, nil
}

func (f *predictionActionStoreExecutorFake) Exec(
	_ context.Context,
	query string,
	arguments ...any,
) (pgconn.CommandTag, error) {
	if len(arguments) < 2 {
		return pgconn.CommandTag{}, fmt.Errorf("action mutation arguments = %d, want at least 2", len(arguments))
	}

	actionID, ok := arguments[0].(uuid.UUID)
	if !ok {
		return pgconn.CommandTag{}, fmt.Errorf("action ID = %T, want uuid.UUID", arguments[0])
	}
	lockToken, ok := arguments[1].(uuid.UUID)
	if !ok {
		return pgconn.CommandTag{}, fmt.Errorf("lock token = %T, want uuid.UUID", arguments[1])
	}
	action := f.action(actionID)
	if action == nil || action.CompletedAt != nil || action.LockToken != lockToken {
		return pgconn.NewCommandTag("UPDATE 0"), nil
	}

	switch {
	case strings.Contains(query, "SET completed_at = now()"):
		completedAt := f.now
		action.CompletedAt = &completedAt
		action.LockedAt = nil
		action.LockToken = uuid.Nil
		return pgconn.NewCommandTag("UPDATE 1"), nil
	case strings.Contains(query, "SET available_at = $3"):
		availableAt, ok := arguments[2].(time.Time)
		if !ok {
			return pgconn.CommandTag{}, fmt.Errorf("available at = %T, want time.Time", arguments[2])
		}
		action.AvailableAt = availableAt
		action.LockedAt = nil
		action.LockToken = uuid.Nil
		return pgconn.NewCommandTag("UPDATE 1"), nil
	default:
		return pgconn.CommandTag{}, fmt.Errorf("unexpected action mutation query: %s", query)
	}
}

func (f *predictionActionStoreExecutorFake) action(id uuid.UUID) *storedPredictionAction {
	for index := range f.actions {
		if f.actions[index].ID == id {
			return &f.actions[index]
		}
	}

	return nil
}

type claimedOutboxRows struct {
	actions []model.ClaimedOutboxAction
	index   int
	err     error
}

func (r *claimedOutboxRows) Close() {}

func (r *claimedOutboxRows) Err() error {
	return r.err
}

func (r *claimedOutboxRows) CommandTag() pgconn.CommandTag {
	return pgconn.NewCommandTag("SELECT")
}

func (r *claimedOutboxRows) FieldDescriptions() []pgconn.FieldDescription {
	return nil
}

func (r *claimedOutboxRows) Next() bool {
	if r.index >= len(r.actions) {
		return false
	}

	r.index++
	return true
}

func (r *claimedOutboxRows) Scan(dest ...any) error {
	if r.index == 0 || r.index > len(r.actions) {
		return errors.New("Scan called without a current row")
	}
	if len(dest) != 8 {
		return fmt.Errorf("scan destinations = %d, want 8", len(dest))
	}

	action := r.actions[r.index-1]
	*dest[0].(*uuid.UUID) = action.ID
	*dest[1].(*uuid.UUID) = action.ChannelID
	*dest[2].(*int64) = action.MatchID
	*dest[3].(*string) = string(action.Action)
	*dest[4].(*int64) = action.Sequence
	*dest[5].(*[]byte) = append([]byte(nil), action.Payload...)
	*dest[6].(*int) = action.Attempts
	*dest[7].(*uuid.UUID) = action.LockToken
	return nil
}

func (r *claimedOutboxRows) Values() ([]any, error) {
	return nil, errors.New("Values is not implemented")
}

func (r *claimedOutboxRows) RawValues() [][]byte {
	return nil
}

func (r *claimedOutboxRows) Conn() *jackcpgx.Conn {
	return nil
}
