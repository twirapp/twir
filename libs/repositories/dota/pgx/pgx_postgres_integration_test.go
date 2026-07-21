package pgx

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"testing"
	"time"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	trmmanager "github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
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
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		t.Fatalf("start PostgreSQL container: %v", err)
	}
	t.Cleanup(func() {
		if err := container.Terminate(context.Background()); err != nil {
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

		if _, err := pool.Exec(ctx, `UPDATE dota_prediction_outbox SET locked_at = now() - INTERVAL '2 minutes' WHERE id = $1`, first[0].ID); err != nil {
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
	statements := []string{
		`CREATE TABLE channels (id UUID PRIMARY KEY)`,
		`CREATE TABLE dota_channel_match_states (
			channel_id UUID PRIMARY KEY REFERENCES channels(id) ON DELETE CASCADE,
			revision BIGINT NOT NULL DEFAULT 0 CHECK (revision >= 0),
			provider_timestamp BIGINT NOT NULL DEFAULT 0,
			snapshot JSONB NOT NULL,
			updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
		)`,
		`CREATE TABLE dota_prediction_outbox (
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
		)`,
	}
	for _, statement := range statements {
		if _, err := pool.Exec(ctx, statement); err != nil {
			return fmt.Errorf("execute schema statement: %w", err)
		}
	}

	return nil
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
