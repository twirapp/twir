package pgx

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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
	create := model.ClaimedOutboxAction{
		ID: uuid.New(),
		OutboxActionInput: model.OutboxActionInput{
			ChannelID: channelID,
			MatchID:   42,
			Action:    model.OutboxActionCreate,
			Sequence:  10,
			Payload:   json.RawMessage(`{"kind":"create"}`),
		},
	}
	resolve := model.ClaimedOutboxAction{
		ID: uuid.New(),
		OutboxActionInput: model.OutboxActionInput{
			ChannelID: channelID,
			MatchID:   42,
			Action:    model.OutboxActionResolve,
			Sequence:  11,
			Payload:   json.RawMessage(`{"kind":"resolve"}`),
		},
	}

	t.Run("returns only the first unfinished action for a match", func(t *testing.T) {
		executor := &claimPredictionActionsExecutorFake{rows: []model.ClaimedOutboxAction{create}}

		claimed, err := (&Pgx{}).claimPredictionActions(context.Background(), executor, dota.ClaimPredictionActionsInput{
			Limit: 2,
			Lease: time.Minute,
		})
		if err != nil {
			t.Fatalf("claimPredictionActions() error = %v", err)
		}
		if len(claimed) != 1 {
			t.Fatalf("claimed actions = %#v, want only create action", claimed)
		}
		if got := claimed[0]; got.Action != model.OutboxActionCreate || got.Sequence != 10 {
			t.Errorf("claimed action = %#v, want create sequence 10", got)
		}
		for _, action := range claimed {
			if action.ID == resolve.ID {
				t.Fatalf("claimed resolve sequence 11 before create sequence 10")
			}
		}
		if !strings.Contains(executor.query, "DISTINCT ON (channel_id, match_id)") {
			t.Errorf("claim query does not select the first action per match: %s", executor.query)
		}
		if !strings.Contains(executor.query, "FOR UPDATE OF outbox SKIP LOCKED") {
			t.Errorf("claim query does not lock candidates with SKIP LOCKED: %s", executor.query)
		}
	})

	t.Run("reclaims an expired lease with a new token", func(t *testing.T) {
		executor := &claimPredictionActionsExecutorFake{rows: []model.ClaimedOutboxAction{create}}
		input := dota.ClaimPredictionActionsInput{Limit: 1, Lease: time.Minute}

		first, err := (&Pgx{}).claimPredictionActions(context.Background(), executor, input)
		if err != nil {
			t.Fatalf("first claimPredictionActions() error = %v", err)
		}
		second, err := (&Pgx{}).claimPredictionActions(context.Background(), executor, input)
		if err != nil {
			t.Fatalf("second claimPredictionActions() error = %v", err)
		}
		if len(first) != 1 || len(second) != 1 {
			t.Fatalf("claimed rows = %#v, %#v; want one reclaimed row each", first, second)
		}
		if first[0].LockToken == second[0].LockToken {
			t.Errorf("reclaimed lock token = %s, want a new token", second[0].LockToken)
		}
		if second[0].Attempts <= first[0].Attempts {
			t.Errorf("reclaimed attempts = %d, want greater than %d", second[0].Attempts, first[0].Attempts)
		}
		if !strings.Contains(executor.query, "locked_at < now()") {
			t.Errorf("claim query does not permit expired leases: %s", executor.query)
		}
	})
}

func TestPredictionActionOwnershipLoss(t *testing.T) {
	actionID := uuid.New()
	wrongToken := uuid.New()
	executor := &predictionActionMutationExecutorFake{}

	t.Run("completion reports ownership loss for a wrong token", func(t *testing.T) {
		err := (&Pgx{}).completePredictionAction(context.Background(), executor, actionID, wrongToken)
		if !errors.Is(err, dota.ErrPredictionActionOwnershipLost) {
			t.Fatalf("completePredictionAction() error = %v, want ownership loss", err)
		}
	})

	t.Run("retry reports ownership loss for a wrong token", func(t *testing.T) {
		err := (&Pgx{}).retryPredictionAction(
			context.Background(),
			executor,
			actionID,
			wrongToken,
			time.Now().Add(time.Minute),
		)
		if !errors.Is(err, dota.ErrPredictionActionOwnershipLost) {
			t.Fatalf("retryPredictionAction() error = %v, want ownership loss", err)
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

type claimPredictionActionsExecutorFake struct {
	rows       []model.ClaimedOutboxAction
	query      string
	attempts   int
	queryError error
}

func (f *claimPredictionActionsExecutorFake) Query(
	_ context.Context,
	query string,
	arguments ...any,
) (jackcpgx.Rows, error) {
	f.query = query
	if f.queryError != nil {
		return nil, f.queryError
	}
	if len(arguments) != 3 {
		return nil, fmt.Errorf("claim arguments = %d, want 3", len(arguments))
	}

	lockToken, ok := arguments[2].(uuid.UUID)
	if !ok {
		return nil, fmt.Errorf("claim lock token = %T, want uuid.UUID", arguments[2])
	}
	f.attempts++

	rows := append([]model.ClaimedOutboxAction(nil), f.rows...)
	for index := range rows {
		rows[index].LockToken = lockToken
		rows[index].Attempts = f.attempts
	}

	return &claimedOutboxRows{actions: rows}, nil
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

type predictionActionMutationExecutorFake struct {
	query string
}

func (f *predictionActionMutationExecutorFake) Exec(
	_ context.Context,
	query string,
	_ ...any,
) (pgconn.CommandTag, error) {
	f.query = query
	return pgconn.NewCommandTag("UPDATE 0"), nil
}
