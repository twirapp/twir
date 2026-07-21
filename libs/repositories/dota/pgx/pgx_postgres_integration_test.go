package pgx

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	trmmanager "github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/twirapp/twir/libs/repositories/dota"
	"github.com/twirapp/twir/libs/repositories/dota/model"
)

func TestPgxLifecyclePostgres(t *testing.T) {
	t.Log("requires a healthy Docker provider; Testcontainers skips this test when Docker is unavailable")
	testcontainers.SkipIfProviderIsNotHealthy(t)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	t.Cleanup(cancel)

	container, err := postgres.Run(
		ctx,
		"postgres:18",
		postgres.WithDatabase("dota"),
		postgres.WithUsername("dota"),
		postgres.WithPassword("dota"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").WithOccurrence(2),
		),
	)
	if err != nil {
		t.Fatalf("start PostgreSQL container: %v", err)
	}
	t.Cleanup(func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cleanupCancel()
		if err := container.Terminate(cleanupCtx); err != nil {
			t.Errorf("terminate PostgreSQL container: %v", err)
		}
	})

	connectionString, err := container.ConnectionString(ctx, "sslmode=disable", "connect_timeout=5")
	if err != nil {
		t.Fatalf("get PostgreSQL connection string: %v", err)
	}
	pool, err := pgxpool.New(ctx, connectionString)
	if err != nil {
		t.Fatalf("create PostgreSQL pool: %v", err)
	}
	t.Cleanup(pool.Close)

	if err := initializePgxLifecycleSchema(ctx, pool); err != nil {
		logs, logsErr := container.Logs(context.Background())
		if logsErr == nil {
			defer logs.Close()
			if output, readErr := io.ReadAll(logs); readErr == nil {
				t.Logf("PostgreSQL logs after schema failure:\n%s", output)
			}
		}
		t.Fatalf("initialize lifecycle schema: %v", err)
	}
	trManager, err := trmmanager.New(trmpgx.NewDefaultFactory(pool))
	if err != nil {
		t.Fatalf("create transaction manager: %v", err)
	}
	repository := New(Opts{PgxPool: pool, TrManager: trManager})

	t.Run("concurrent state transitions commit exactly once", func(t *testing.T) {
		resetPgxLifecycleData(t, ctx, pool)
		channelID := insertPgxLifecycleChannel(t, ctx, pool)
		input := pgxLifecycleTransition(
			channelID,
			0,
			101,
			[]model.OutboxActionInput{pgxLifecycleAction(channelID, 101, model.OutboxActionCreate, 10)},
		)

		start := make(chan struct{})
		results := make(chan struct {
			committed bool
			err       error
		}, 2)
		for range 2 {
			go func() {
				<-start
				committed, err := repository.ApplyMatchStateTransition(ctx, input)
				results <- struct {
					committed bool
					err       error
				}{committed: committed, err: err}
			}()
		}
		close(start)

		committed, conflicted := 0, 0
		for range 2 {
			result := <-results
			if result.err != nil {
				t.Fatalf("ApplyMatchStateTransition() error = %v", result.err)
			}
			if result.committed {
				committed++
			} else {
				conflicted++
			}
		}
		if committed != 1 || conflicted != 1 {
			t.Fatalf("transition outcomes = committed %d, conflicted %d; want 1 each", committed, conflicted)
		}

		var revision int64
		if err := pool.QueryRow(ctx, `SELECT revision FROM dota_channel_match_states WHERE channel_id = $1`, channelID).Scan(&revision); err != nil {
			t.Fatalf("get state revision: %v", err)
		}
		if revision != 1 {
			t.Errorf("state revision = %d, want 1", revision)
		}
		var actionCount int
		if err := pool.QueryRow(ctx, `SELECT count(*) FROM dota_prediction_outbox WHERE channel_id = $1`, channelID).Scan(&actionCount); err != nil {
			t.Fatalf("count outbox actions: %v", err)
		}
		if actionCount != 1 {
			t.Errorf("outbox actions = %d, want 1", actionCount)
		}
	})

	t.Run("rejected nonzero transition does not create state", func(t *testing.T) {
		resetPgxLifecycleData(t, ctx, pool)
		channelID := insertPgxLifecycleChannel(t, ctx, pool)

		committed, err := repository.ApplyMatchStateTransition(
			ctx,
			pgxLifecycleTransition(channelID, 1, 151, nil),
		)
		if err != nil {
			t.Fatalf("ApplyMatchStateTransition() error = %v", err)
		}
		if committed {
			t.Fatal("ApplyMatchStateTransition() committed = true, want false")
		}

		var stateCount int
		if err := pool.QueryRow(
			ctx,
			`SELECT count(*) FROM dota_channel_match_states WHERE channel_id = $1`,
			channelID,
		).Scan(&stateCount); err != nil {
			t.Fatalf("count match states: %v", err)
		}
		if stateCount != 0 {
			t.Errorf("match state rows = %d, want 0", stateCount)
		}
	})

	t.Run("migration rejects duplicate sequences across matches in a channel", func(t *testing.T) {
		resetPgxLifecycleData(t, ctx, pool)
		channelID := insertPgxLifecycleChannel(t, ctx, pool)

		if _, err := pool.Exec(
			ctx,
			`INSERT INTO dota_prediction_outbox (channel_id, match_id, action, sequence, payload) VALUES ($1, $2, 'create', $3, '{}'::jsonb)`,
			channelID,
			161,
			10,
		); err != nil {
			t.Fatalf("insert first outbox action: %v", err)
		}
		if _, err := pool.Exec(
			ctx,
			`INSERT INTO dota_prediction_outbox (channel_id, match_id, action, sequence, payload) VALUES ($1, $2, 'resolve', $3, '{}'::jsonb)`,
			channelID,
			162,
			10,
		); err == nil {
			t.Fatal("insert duplicate sequence error = nil, want unique constraint violation")
		}
	})

	t.Run("migration indexes the complete per-channel claim order", func(t *testing.T) {
		var definition string
		if err := pool.QueryRow(
			ctx,
			`SELECT indexdef FROM pg_indexes WHERE schemaname = current_schema() AND indexname = 'dota_prediction_outbox_channel_order_idx'`,
		).Scan(&definition); err != nil {
			t.Fatalf("get channel order index: %v", err)
		}
		if !strings.Contains(definition, "channel_id, sequence") {
			t.Errorf("channel order index = %q, want channel_id before sequence", definition)
		}
	})

	t.Run("claims resolve only after create completes", func(t *testing.T) {
		resetPgxLifecycleData(t, ctx, pool)
		channelID := insertPgxLifecycleChannel(t, ctx, pool)
		committed, err := repository.ApplyMatchStateTransition(
			ctx,
			pgxLifecycleTransition(
				channelID,
				0,
				201,
				[]model.OutboxActionInput{
					pgxLifecycleAction(channelID, 201, model.OutboxActionCreate, 10),
					pgxLifecycleAction(channelID, 201, model.OutboxActionResolve, 11),
				},
			),
		)
		if err != nil || !committed {
			t.Fatalf("ApplyMatchStateTransition() = (%t, %v), want (true, nil)", committed, err)
		}

		claimInput := dota.ClaimPredictionActionsInput{Limit: 2, Lease: time.Minute}
		first, err := repository.ClaimPredictionActions(ctx, claimInput)
		if err != nil {
			t.Fatalf("first ClaimPredictionActions() error = %v", err)
		}
		if len(first) != 1 || first[0].Action != model.OutboxActionCreate || first[0].Sequence != 10 {
			t.Fatalf("first claimed actions = %#v, want only create sequence 10", first)
		}
		if err := repository.CompletePredictionAction(ctx, first[0].ID, first[0].LockToken); err != nil {
			t.Fatalf("CompletePredictionAction() error = %v", err)
		}

		second, err := repository.ClaimPredictionActions(ctx, claimInput)
		if err != nil {
			t.Fatalf("second ClaimPredictionActions() error = %v", err)
		}
		if len(second) != 1 || second[0].Action != model.OutboxActionResolve || second[0].Sequence != 11 {
			t.Fatalf("second claimed actions = %#v, want only resolve sequence 11", second)
		}
	})

	t.Run("claims replacement cancel before replacement create", func(t *testing.T) {
		resetPgxLifecycleData(t, ctx, pool)
		channelID := insertPgxLifecycleChannel(t, ctx, pool)
		committed, err := repository.ApplyMatchStateTransition(
			ctx,
			pgxLifecycleTransition(
				channelID,
				0,
				211,
				[]model.OutboxActionInput{pgxLifecycleAction(channelID, 211, model.OutboxActionCreate, 1)},
			),
		)
		if err != nil || !committed {
			t.Fatalf("initial ApplyMatchStateTransition() = (%t, %v), want (true, nil)", committed, err)
		}

		claimInput := dota.ClaimPredictionActionsInput{Limit: 2, Lease: time.Minute}
		initial, err := repository.ClaimPredictionActions(ctx, claimInput)
		if err != nil || len(initial) != 1 {
			t.Fatalf("initial ClaimPredictionActions() = (%#v, %v), want one create", initial, err)
		}
		if err := repository.CompletePredictionAction(ctx, initial[0].ID, initial[0].LockToken); err != nil {
			t.Fatalf("complete initial create: %v", err)
		}

		committed, err = repository.ApplyMatchStateTransition(
			ctx,
			pgxLifecycleTransition(
				channelID,
				1,
				212,
				[]model.OutboxActionInput{
					pgxLifecycleAction(channelID, 211, model.OutboxActionCancel, 3),
					pgxLifecycleAction(channelID, 212, model.OutboxActionCreate, 4),
				},
			),
		)
		if err != nil || !committed {
			t.Fatalf("replacement ApplyMatchStateTransition() = (%t, %v), want (true, nil)", committed, err)
		}

		first, err := repository.ClaimPredictionActions(ctx, claimInput)
		if err != nil {
			t.Fatalf("first replacement ClaimPredictionActions() error = %v", err)
		}
		if len(first) != 1 || first[0].Action != model.OutboxActionCancel || first[0].MatchID != 211 {
			t.Fatalf("first replacement claimed actions = %#v, want only old-match cancel", first)
		}
		if err := repository.CompletePredictionAction(ctx, first[0].ID, first[0].LockToken); err != nil {
			t.Fatalf("complete replacement cancel: %v", err)
		}

		second, err := repository.ClaimPredictionActions(ctx, claimInput)
		if err != nil {
			t.Fatalf("second replacement ClaimPredictionActions() error = %v", err)
		}
		if len(second) != 1 || second[0].Action != model.OutboxActionCreate || second[0].MatchID != 212 {
			t.Fatalf("second replacement claimed actions = %#v, want only new-match create", second)
		}
	})

	t.Run("does not reclaim active leases but reclaims expired leases", func(t *testing.T) {
		resetPgxLifecycleData(t, ctx, pool)
		channelID := insertPgxLifecycleChannel(t, ctx, pool)
		committed, err := repository.ApplyMatchStateTransition(
			ctx,
			pgxLifecycleTransition(
				channelID,
				0,
				301,
				[]model.OutboxActionInput{pgxLifecycleAction(channelID, 301, model.OutboxActionCreate, 10)},
			),
		)
		if err != nil || !committed {
			t.Fatalf("ApplyMatchStateTransition() = (%t, %v), want (true, nil)", committed, err)
		}

		claimInput := dota.ClaimPredictionActionsInput{Limit: 1, Lease: time.Minute}
		first, err := repository.ClaimPredictionActions(ctx, claimInput)
		if err != nil {
			t.Fatalf("first ClaimPredictionActions() error = %v", err)
		}
		if len(first) != 1 {
			t.Fatalf("first claimed actions = %#v, want one action", first)
		}
		second, err := repository.ClaimPredictionActions(ctx, claimInput)
		if err != nil {
			t.Fatalf("second ClaimPredictionActions() error = %v", err)
		}
		if len(second) != 0 {
			t.Fatalf("second claimed actions = %#v, want none while lease is active", second)
		}

		if _, err := pool.Exec(ctx, `UPDATE dota_prediction_outbox SET locked_at = now() - INTERVAL '2 minutes', lease_expires_at = now() - INTERVAL '1 minute' WHERE id = $1`, first[0].ID); err != nil {
			t.Fatalf("expire outbox lease: %v", err)
		}
		reclaimed, err := repository.ClaimPredictionActions(ctx, claimInput)
		if err != nil {
			t.Fatalf("reclaim ClaimPredictionActions() error = %v", err)
		}
		if len(reclaimed) != 1 {
			t.Fatalf("reclaimed actions = %#v, want one action", reclaimed)
		}
		if reclaimed[0].ID != first[0].ID || reclaimed[0].LockToken == first[0].LockToken || reclaimed[0].Attempts <= first[0].Attempts {
			t.Errorf("reclaimed action = %#v, want same ID, a new token, and higher attempts than %#v", reclaimed[0], first[0])
		}
	})

	t.Run("owner renews an active lease and stale token is rejected", func(t *testing.T) {
		resetPgxLifecycleData(t, ctx, pool)
		channelID := insertPgxLifecycleChannel(t, ctx, pool)
		committed, err := repository.ApplyMatchStateTransition(
			ctx,
			pgxLifecycleTransition(
				channelID,
				0,
				326,
				[]model.OutboxActionInput{pgxLifecycleAction(channelID, 326, model.OutboxActionCreate, 10)},
			),
		)
		if err != nil || !committed {
			t.Fatalf("ApplyMatchStateTransition() = (%t, %v), want (true, nil)", committed, err)
		}

		claimed, err := repository.ClaimPredictionActions(
			ctx,
			dota.ClaimPredictionActionsInput{Limit: 1, Lease: time.Minute},
		)
		if err != nil || len(claimed) != 1 {
			t.Fatalf("ClaimPredictionActions() = (%#v, %v), want one action", claimed, err)
		}

		var before time.Time
		if err := pool.QueryRow(ctx, `SELECT lease_expires_at FROM dota_prediction_outbox WHERE id = $1`, claimed[0].ID).Scan(&before); err != nil {
			t.Fatalf("get lease before renewal: %v", err)
		}
		if err := repository.RenewPredictionAction(ctx, claimed[0].ID, claimed[0].LockToken, time.Nanosecond); err == nil {
			t.Fatal("RenewPredictionAction() error = nil, want sub-microsecond lease validation error")
		}
		if err := repository.RenewPredictionAction(ctx, claimed[0].ID, claimed[0].LockToken, 2*time.Minute); err != nil {
			t.Fatalf("RenewPredictionAction() error = %v", err)
		}

		var after time.Time
		if err := pool.QueryRow(ctx, `SELECT lease_expires_at FROM dota_prediction_outbox WHERE id = $1`, claimed[0].ID).Scan(&after); err != nil {
			t.Fatalf("get lease after renewal: %v", err)
		}
		if !after.After(before) {
			t.Errorf("renewed lease expiry = %v, want after %v", after, before)
		}

		if err := repository.RenewPredictionAction(ctx, claimed[0].ID, uuid.New(), time.Minute); !errors.Is(err, dota.ErrPredictionActionOwnershipLost) {
			t.Fatalf("stale RenewPredictionAction() error = %v, want ownership loss", err)
		}
		var storedToken uuid.UUID
		var storedExpiry time.Time
		if err := pool.QueryRow(ctx, `SELECT lock_token, lease_expires_at FROM dota_prediction_outbox WHERE id = $1`, claimed[0].ID).Scan(&storedToken, &storedExpiry); err != nil {
			t.Fatalf("get lease after stale renewal: %v", err)
		}
		if storedToken != claimed[0].LockToken || !storedExpiry.Equal(after) {
			t.Errorf("stored lease after stale renewal = token %s expiry %v, want token %s expiry %v", storedToken, storedExpiry, claimed[0].LockToken, after)
		}
	})

	t.Run("expired same-token owner cannot mutate an action", func(t *testing.T) {
		tests := []struct {
			name   string
			mutate func(model.ClaimedOutboxAction) error
		}{
			{
				name: "renew",
				mutate: func(action model.ClaimedOutboxAction) error {
					return repository.RenewPredictionAction(ctx, action.ID, action.LockToken, time.Minute)
				},
			},
			{
				name: "complete",
				mutate: func(action model.ClaimedOutboxAction) error {
					return repository.CompletePredictionAction(ctx, action.ID, action.LockToken)
				},
			},
			{
				name: "retry",
				mutate: func(action model.ClaimedOutboxAction) error {
					return repository.RetryPredictionAction(ctx, action.ID, action.LockToken, time.Now().Add(time.Minute))
				},
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				resetPgxLifecycleData(t, ctx, pool)
				channelID := insertPgxLifecycleChannel(t, ctx, pool)
				committed, err := repository.ApplyMatchStateTransition(
					ctx,
					pgxLifecycleTransition(
						channelID,
						0,
						341,
						[]model.OutboxActionInput{pgxLifecycleAction(channelID, 341, model.OutboxActionCreate, 10)},
					),
				)
				if err != nil || !committed {
					t.Fatalf("ApplyMatchStateTransition() = (%t, %v), want (true, nil)", committed, err)
				}

				claimed, err := repository.ClaimPredictionActions(
					ctx,
					dota.ClaimPredictionActionsInput{Limit: 1, Lease: time.Minute},
				)
				if err != nil || len(claimed) != 1 {
					t.Fatalf("ClaimPredictionActions() = (%#v, %v), want one action", claimed, err)
				}

				var beforeAvailable time.Time
				if err := pool.QueryRow(ctx, `SELECT available_at FROM dota_prediction_outbox WHERE id = $1`, claimed[0].ID).Scan(&beforeAvailable); err != nil {
					t.Fatalf("get available time before %s: %v", test.name, err)
				}
				if _, err := pool.Exec(ctx, `UPDATE dota_prediction_outbox SET lease_expires_at = now() - INTERVAL '1 second' WHERE id = $1`, claimed[0].ID); err != nil {
					t.Fatalf("expire lease before %s: %v", test.name, err)
				}

				if err := test.mutate(claimed[0]); !errors.Is(err, dota.ErrPredictionActionOwnershipLost) {
					t.Fatalf("expired %s error = %v, want ownership loss", test.name, err)
				}

				var afterAvailable, afterLease time.Time
				var afterToken uuid.UUID
				var unfinished bool
				if err := pool.QueryRow(
					ctx,
					`SELECT available_at, lease_expires_at, lock_token, completed_at IS NULL FROM dota_prediction_outbox WHERE id = $1`,
					claimed[0].ID,
				).Scan(&afterAvailable, &afterLease, &afterToken, &unfinished); err != nil {
					t.Fatalf("get action after expired %s: %v", test.name, err)
				}
				if !unfinished || afterToken != claimed[0].LockToken || !afterAvailable.Equal(beforeAvailable) || !afterLease.Before(time.Now()) {
					t.Errorf(
						"stored action after expired %s = unfinished %t token %s available %v lease %v",
						test.name,
						unfinished,
						afterToken,
						afterAvailable,
						afterLease,
					)
				}
			})
		}
	})

	t.Run("does not steal an active lease when a later claimant uses a shorter lease", func(t *testing.T) {
		resetPgxLifecycleData(t, ctx, pool)
		channelID := insertPgxLifecycleChannel(t, ctx, pool)
		committed, err := repository.ApplyMatchStateTransition(
			ctx,
			pgxLifecycleTransition(
				channelID,
				0,
				351,
				[]model.OutboxActionInput{pgxLifecycleAction(channelID, 351, model.OutboxActionCreate, 10)},
			),
		)
		if err != nil || !committed {
			t.Fatalf("ApplyMatchStateTransition() = (%t, %v), want (true, nil)", committed, err)
		}

		first, err := repository.ClaimPredictionActions(
			ctx,
			dota.ClaimPredictionActionsInput{Limit: 1, Lease: 5 * time.Minute},
		)
		if err != nil || len(first) != 1 {
			t.Fatalf("first ClaimPredictionActions() = (%#v, %v), want one action", first, err)
		}
		if _, err := pool.Exec(
			ctx,
			`UPDATE dota_prediction_outbox SET locked_at = now() - INTERVAL '2 minutes' WHERE id = $1`,
			first[0].ID,
		); err != nil {
			t.Fatalf("age lock timestamp: %v", err)
		}

		second, err := repository.ClaimPredictionActions(
			ctx,
			dota.ClaimPredictionActionsInput{Limit: 1, Lease: time.Minute},
		)
		if err != nil {
			t.Fatalf("second ClaimPredictionActions() error = %v", err)
		}
		if len(second) != 0 {
			t.Fatalf("second claimed actions = %#v, want none while original lease remains active", second)
		}
	})

	t.Run("concurrent claims receive distinct channel heads", func(t *testing.T) {
		resetPgxLifecycleData(t, ctx, pool)
		firstChannelID := insertPgxLifecycleChannel(t, ctx, pool)
		secondChannelID := insertPgxLifecycleChannel(t, ctx, pool)
		committed, err := repository.ApplyMatchStateTransition(
			ctx,
			pgxLifecycleTransition(
				firstChannelID,
				0,
				361,
				[]model.OutboxActionInput{pgxLifecycleAction(firstChannelID, 361, model.OutboxActionCreate, 1)},
			),
		)
		if err != nil || !committed {
			t.Fatalf("first ApplyMatchStateTransition() = (%t, %v), want (true, nil)", committed, err)
		}
		committed, err = repository.ApplyMatchStateTransition(
			ctx,
			pgxLifecycleTransition(
				secondChannelID,
				0,
				362,
				[]model.OutboxActionInput{pgxLifecycleAction(secondChannelID, 362, model.OutboxActionCreate, 1)},
			),
		)
		if err != nil || !committed {
			t.Fatalf("second ApplyMatchStateTransition() = (%t, %v), want (true, nil)", committed, err)
		}

		start := make(chan struct{})
		results := make(chan struct {
			actions []model.ClaimedOutboxAction
			err     error
		}, 2)
		for range 2 {
			go func() {
				<-start
				actions, err := repository.ClaimPredictionActions(
					ctx,
					dota.ClaimPredictionActionsInput{Limit: 1, Lease: time.Minute},
				)
				results <- struct {
					actions []model.ClaimedOutboxAction
					err     error
				}{actions: actions, err: err}
			}()
		}
		close(start)

		var claimed [2]model.ClaimedOutboxAction
		for index := range 2 {
			result := <-results
			if result.err != nil {
				t.Fatalf("ClaimPredictionActions() error = %v", result.err)
			}
			if len(result.actions) != 1 {
				t.Fatalf("claimed actions = %#v, want one action", result.actions)
			}
			claimed[index] = result.actions[0]
		}
		if claimed[0].ID == claimed[1].ID {
			t.Errorf("both workers claimed action %s", claimed[0].ID)
		}
		if claimed[0].ChannelID == claimed[1].ChannelID {
			t.Errorf("both workers claimed channel %s, want distinct channels", claimed[0].ChannelID)
		}
	})

	t.Run("outbox insert failure rolls back its state update", func(t *testing.T) {
		resetPgxLifecycleData(t, ctx, pool)
		channelID := insertPgxLifecycleChannel(t, ctx, pool)
		if _, err := pool.Exec(
			ctx,
			`INSERT INTO dota_prediction_outbox (channel_id, match_id, action, sequence, payload) VALUES ($1, $2, 'create', $3, '{}'::jsonb)`,
			channelID,
			371,
			10,
		); err != nil {
			t.Fatalf("seed outbox action: %v", err)
		}

		committed, err := repository.ApplyMatchStateTransition(
			ctx,
			pgxLifecycleTransition(
				channelID,
				0,
				371,
				[]model.OutboxActionInput{pgxLifecycleAction(channelID, 371, model.OutboxActionCreate, 10)},
			),
		)
		if err == nil {
			t.Fatal("ApplyMatchStateTransition() error = nil, want outbox uniqueness error")
		}
		if committed {
			t.Fatal("ApplyMatchStateTransition() committed = true, want false")
		}

		var stateCount int
		if err := pool.QueryRow(
			ctx,
			`SELECT count(*) FROM dota_channel_match_states WHERE channel_id = $1`,
			channelID,
		).Scan(&stateCount); err != nil {
			t.Fatalf("count match states after rollback: %v", err)
		}
		if stateCount != 0 {
			t.Errorf("match state rows after rollback = %d, want 0", stateCount)
		}
	})

	t.Run("wrong ownership tokens leave actions unfinished and locked", func(t *testing.T) {
		resetPgxLifecycleData(t, ctx, pool)
		channelID := insertPgxLifecycleChannel(t, ctx, pool)
		committed, err := repository.ApplyMatchStateTransition(
			ctx,
			pgxLifecycleTransition(
				channelID,
				0,
				401,
				[]model.OutboxActionInput{pgxLifecycleAction(channelID, 401, model.OutboxActionCreate, 10)},
			),
		)
		if err != nil || !committed {
			t.Fatalf("ApplyMatchStateTransition() = (%t, %v), want (true, nil)", committed, err)
		}

		claimed, err := repository.ClaimPredictionActions(ctx, dota.ClaimPredictionActionsInput{Limit: 1, Lease: time.Minute})
		if err != nil {
			t.Fatalf("ClaimPredictionActions() error = %v", err)
		}
		if len(claimed) != 1 {
			t.Fatalf("claimed actions = %#v, want one action", claimed)
		}
		wrongToken := uuid.New()
		if err := repository.CompletePredictionAction(ctx, claimed[0].ID, wrongToken); !errors.Is(err, dota.ErrPredictionActionOwnershipLost) {
			t.Fatalf("CompletePredictionAction() error = %v, want ownership loss", err)
		}
		assertPgxLifecycleActionUnfinishedAndLocked(t, ctx, pool, claimed[0].ID, claimed[0].LockToken)

		if err := repository.RetryPredictionAction(ctx, claimed[0].ID, wrongToken, time.Now().Add(time.Minute)); !errors.Is(err, dota.ErrPredictionActionOwnershipLost) {
			t.Fatalf("RetryPredictionAction() error = %v, want ownership loss", err)
		}
		assertPgxLifecycleActionUnfinishedAndLocked(t, ctx, pool, claimed[0].ID, claimed[0].LockToken)
	})
}

func initializePgxLifecycleSchema(ctx context.Context, pool *pgxpool.Pool) error {
	if _, err := pool.Exec(ctx, `CREATE TABLE channels (id UUID PRIMARY KEY)`); err != nil {
		return fmt.Errorf("create channels table: %w", err)
	}

	migration, err := os.ReadFile(dotaLifecycleMigrationPath())
	if err != nil {
		return fmt.Errorf("read lifecycle migration: %w", err)
	}
	upMigration, err := gooseUpMigration(migration)
	if err != nil {
		return fmt.Errorf("extract lifecycle migration: %w", err)
	}
	if _, err := pool.Exec(ctx, upMigration); err != nil {
		return fmt.Errorf("apply lifecycle migration: %w", err)
	}

	return nil
}

func dotaLifecycleMigrationPath() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("find lifecycle integration test file")
	}

	return filepath.Join(filepath.Dir(file), "../../../migrations/postgres/20260721012653_dota_match_state_outbox.sql")
}

func gooseUpMigration(migration []byte) (string, error) {
	contents := string(migration)
	upStart := strings.Index(contents, "-- +goose Up")
	downStart := strings.Index(contents, "-- +goose Down")
	if upStart == -1 || downStart == -1 || downStart <= upStart {
		return "", errors.New("migration must contain ordered goose Up and Down sections")
	}

	return contents[upStart:downStart], nil
}

func insertPgxLifecycleChannel(t *testing.T, ctx context.Context, pool *pgxpool.Pool) uuid.UUID {
	t.Helper()

	channelID := uuid.New()
	if _, err := pool.Exec(ctx, `INSERT INTO channels (id) VALUES ($1)`, channelID); err != nil {
		t.Fatalf("insert channel: %v", err)
	}

	return channelID
}

func resetPgxLifecycleData(t *testing.T, ctx context.Context, pool *pgxpool.Pool) {
	t.Helper()

	if _, err := pool.Exec(ctx, `TRUNCATE TABLE channels CASCADE`); err != nil {
		t.Fatalf("reset lifecycle test data: %v", err)
	}
}

func pgxLifecycleTransition(
	channelID uuid.UUID,
	expectedRevision int64,
	matchID int64,
	actions []model.OutboxActionInput,
) dota.ApplyMatchStateTransitionInput {
	return dota.ApplyMatchStateTransitionInput{
		ChannelID:         channelID,
		ExpectedRevision:  expectedRevision,
		ProviderTimestamp: 1_700_000_000 + matchID,
		Snapshot:          json.RawMessage(`{"state":"active"}`),
		Actions:           actions,
	}
}

func pgxLifecycleAction(
	channelID uuid.UUID,
	matchID int64,
	action model.OutboxAction,
	sequence int64,
) model.OutboxActionInput {
	return model.OutboxActionInput{
		ChannelID: channelID,
		MatchID:   matchID,
		Action:    action,
		Sequence:  sequence,
		Payload:   json.RawMessage(`{"kind":"prediction"}`),
	}
}

func assertPgxLifecycleActionUnfinishedAndLocked(
	t *testing.T,
	ctx context.Context,
	pool *pgxpool.Pool,
	actionID uuid.UUID,
	lockToken uuid.UUID,
) {
	t.Helper()

	var unfinished, locked bool
	var storedToken uuid.UUID
	if err := pool.QueryRow(
		ctx,
		`SELECT completed_at IS NULL, locked_at IS NOT NULL, lock_token FROM dota_prediction_outbox WHERE id = $1`,
		actionID,
	).Scan(&unfinished, &locked, &storedToken); err != nil {
		t.Fatalf("get outbox action state: %v", err)
	}
	if !unfinished || !locked || storedToken != lockToken {
		t.Errorf("outbox action state = unfinished %t, locked %t, token %s; want unfinished, locked, token %s", unfinished, locked, storedToken, lockToken)
	}
}
