package predictions

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/twirapp/twir/apps/dota/internal/match"
	busdota "github.com/twirapp/twir/libs/bus-core/dota"
	"github.com/twirapp/twir/libs/logger"
	dotarepository "github.com/twirapp/twir/libs/repositories/dota"
	"github.com/twirapp/twir/libs/repositories/dota/model"
	"go.uber.org/fx"
)

const (
	predictionActionBatchSize    = 10
	predictionActionLease        = 2 * time.Minute
	predictionActionRenewEvery   = predictionActionLease / 3
	predictionActionRenewTimeout = 5 * time.Second
	predictionActionPoll         = time.Second
	predictionActionRetryCap     = time.Minute
)

type predictionActionRepository interface {
	GetMatchState(context.Context, uuid.UUID) (model.MatchState, error)
	ClaimPredictionActions(
		context.Context,
		dotarepository.ClaimPredictionActionsInput,
	) ([]model.ClaimedOutboxAction, error)
	RenewPredictionAction(context.Context, uuid.UUID, uuid.UUID, time.Duration) error
	CompletePredictionAction(context.Context, uuid.UUID, uuid.UUID) error
	RetryPredictionAction(context.Context, uuid.UUID, uuid.UUID, time.Time) error
}

type predictionActionDispatcher interface {
	Create(context.Context, match.LifecycleAction) error
	Resolve(context.Context, match.LifecycleAction) (model.ChannelDotaSettings, error)
	Cancel(context.Context, match.LifecycleAction) error
}

type predictionStatsUpdater interface {
	UpdateStats(context.Context, uuid.UUID, model.ChannelDotaSettings) error
}

type matchEndedEmitter interface {
	MatchEnded(context.Context, busdota.MatchEndedMessage) error
}

type matchEndedDeliveryStore interface {
	ClaimMatchEndedDelivery(context.Context, string, string, time.Duration) (matchEndedDeliveryState, error)
	CompleteMatchEndedDelivery(context.Context, string, string, time.Duration) (bool, error)
	RenewMatchEndedDelivery(context.Context, string, string, time.Duration) (bool, error)
	ReleaseMatchEndedDelivery(context.Context, string, string) error
}

var (
	errMatchEndedDeliveryPending       = errors.New("match ended delivery is pending")
	errMatchEndedDeliveryOwnershipLost = errors.New("match ended delivery ownership lost")
	errMatchEndedDeliveryHeartbeatStop = errors.New("match ended delivery heartbeat stopped")
)

type LifecycleWorkerOpts struct {
	fx.In

	Lifecycle   fx.Lifecycle
	Repository  dotarepository.Repository
	Predictions *Predictions
	State       *match.StateMachine
	Store       Store
	Emitter     match.EventEmitter
	Logger      *slog.Logger
}

// LifecycleWorker drains durable prediction actions without coupling their delivery to GSI processing.
type LifecycleWorker struct {
	repository       predictionActionRepository
	predictions      predictionActionDispatcher
	stats            predictionStatsUpdater
	emitter          matchEndedEmitter
	deliveryStore    matchEndedDeliveryStore
	logger           *slog.Logger
	pollEvery        time.Duration
	actionLease      time.Duration
	actionRenewEvery time.Duration

	mu     sync.Mutex
	cancel context.CancelFunc
	done   chan struct{}
}

func NewLifecycleWorker(opts LifecycleWorkerOpts) *LifecycleWorker {
	return newLifecycleWorker(
		opts.Repository,
		opts.Predictions,
		opts.State,
		opts.Store,
		opts.Emitter,
		opts.Logger,
		opts.Lifecycle,
		predictionActionPoll,
	)
}

func newLifecycleWorker(
	repository predictionActionRepository,
	predictions predictionActionDispatcher,
	stats predictionStatsUpdater,
	deliveryStore matchEndedDeliveryStore,
	emitter matchEndedEmitter,
	logger *slog.Logger,
	lifecycle fx.Lifecycle,
	pollEvery time.Duration,
) *LifecycleWorker {
	if logger == nil {
		logger = slog.Default()
	}
	if pollEvery <= 0 {
		pollEvery = predictionActionPoll
	}

	worker := &LifecycleWorker{
		repository:       repository,
		predictions:      predictions,
		stats:            stats,
		deliveryStore:    deliveryStore,
		emitter:          emitter,
		logger:           logger,
		pollEvery:        pollEvery,
		actionLease:      predictionActionLease,
		actionRenewEvery: predictionActionRenewEvery,
	}
	if lifecycle != nil {
		lifecycle.Append(fx.Hook{
			OnStart: worker.start,
			OnStop:  worker.stop,
		})
	}

	return worker
}

func (w *LifecycleWorker) start(context.Context) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.cancel != nil {
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	w.cancel = cancel
	w.done = make(chan struct{})
	go func() {
		defer close(w.done)
		w.run(ctx)
	}()

	return nil
}

func (w *LifecycleWorker) stop(ctx context.Context) error {
	w.mu.Lock()
	cancel := w.cancel
	done := w.done
	w.mu.Unlock()
	if cancel == nil {
		return nil
	}

	cancel()
	select {
	case <-done:
		w.mu.Lock()
		if w.done == done {
			w.done = nil
			w.cancel = nil
		}
		w.mu.Unlock()
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (w *LifecycleWorker) run(ctx context.Context) {
	ticker := time.NewTicker(w.pollEvery)
	defer ticker.Stop()

	for {
		if err := w.ProcessOnce(ctx); err != nil && !errors.Is(err, context.Canceled) {
			w.logger.ErrorContext(ctx, "dota prediction outbox worker failed", logger.Error(err))
		}

		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

// ProcessOnce claims and processes one bounded batch of durable prediction actions.
func (w *LifecycleWorker) ProcessOnce(ctx context.Context) error {
	actions, err := w.repository.ClaimPredictionActions(ctx, dotarepository.ClaimPredictionActionsInput{
		Limit: predictionActionBatchSize,
		Lease: w.actionLease,
	})
	if err != nil {
		return fmt.Errorf("claim prediction actions: %w", err)
	}

	heartbeats := make([]*actionHeartbeat, len(actions))
	for index, action := range actions {
		heartbeats[index] = w.startActionHeartbeat(ctx, action)
	}
	defer func() {
		for _, heartbeat := range heartbeats {
			heartbeat.stop()
		}
	}()

	var firstErr error
	for index, action := range actions {
		if err := w.processAction(ctx, action, heartbeats[index]); err != nil && firstErr == nil {
			firstErr = err
		}
	}

	return firstErr
}

func (w *LifecycleWorker) processAction(
	ctx context.Context,
	claimed model.ClaimedOutboxAction,
	heartbeat *actionHeartbeat,
) error {
	action, err := lifecycleActionFromClaim(claimed)
	if err == nil {
		err = w.dispatch(heartbeat, action)
	}
	if heartbeatErr := heartbeat.stop(); heartbeatErr != nil {
		err = heartbeatErr
	}
	if err != nil {
		return w.retry(ctx, claimed, err)
	}

	if err := w.repository.CompletePredictionAction(ctx, claimed.ID, claimed.LockToken); err != nil {
		if errors.Is(err, dotarepository.ErrPredictionActionOwnershipLost) {
			w.logger.DebugContext(ctx, "dota prediction action completion ownership lost")
			return nil
		}
		return fmt.Errorf("complete prediction action: %w", err)
	}

	return nil
}

func (w *LifecycleWorker) retry(ctx context.Context, claimed model.ClaimedOutboxAction, cause error) error {
	availableAt := time.Now().Add(retryDelay(claimed.Attempts))
	if err := w.repository.RetryPredictionAction(ctx, claimed.ID, claimed.LockToken, availableAt); err != nil {
		if errors.Is(err, dotarepository.ErrPredictionActionOwnershipLost) {
			w.logger.DebugContext(ctx, "dota prediction action retry ownership lost")
			return nil
		}
		return fmt.Errorf("retry prediction action: %w", err)
	}

	w.logger.WarnContext(ctx, "dota prediction action failed; retry scheduled", logger.Error(cause))
	return nil
}

func (w *LifecycleWorker) dispatch(heartbeat *actionHeartbeat, action match.LifecycleAction) error {
	if err := heartbeat.operationError(); err != nil {
		return err
	}
	ctx := heartbeat.ctx

	switch action.Kind {
	case match.ActionCreate:
		state, err := w.repository.GetMatchState(ctx, action.ChannelID)
		if err != nil {
			return fmt.Errorf("get current match state: %w", err)
		}

		var snapshot match.Snapshot
		if err := json.Unmarshal(state.Snapshot, &snapshot); err != nil {
			return fmt.Errorf("decode current match snapshot: %w", err)
		}
		if snapshot.ChannelID != uuid.Nil && snapshot.ChannelID != action.ChannelID {
			return fmt.Errorf(
				"current match snapshot channel ID %s does not match action channel ID %s",
				snapshot.ChannelID,
				action.ChannelID,
			)
		}
		if snapshot.State != match.StateInGame || !snapshot.InGame || snapshot.MatchID != action.MatchID {
			return nil
		}
		if err := heartbeat.operationError(); err != nil {
			return err
		}

		return w.predictions.Create(ctx, action)
	case match.ActionResolve:
		settings, err := w.predictions.Resolve(ctx, action)
		if err != nil {
			return err
		}
		if err := heartbeat.operationError(); err != nil {
			return err
		}
		if err := w.stats.UpdateStats(ctx, action.ChannelID, settings); err != nil {
			return fmt.Errorf("update match stats: %w", err)
		}
		if err := heartbeat.operationError(); err != nil {
			return err
		}
		return w.publishMatchEnded(heartbeat, action, settings)
	case match.ActionCancel:
		return w.predictions.Cancel(ctx, action)
	default:
		return fmt.Errorf("unsupported lifecycle action kind %q", action.Kind)
	}
}

func (w *LifecycleWorker) publishMatchEnded(
	heartbeat *actionHeartbeat,
	action match.LifecycleAction,
	settings model.ChannelDotaSettings,
) error {
	if w.deliveryStore == nil {
		return errors.New("match ended delivery store is required")
	}
	if err := heartbeat.operationError(); err != nil {
		return err
	}

	ctx := heartbeat.ctx
	key := matchEndedDeliveryKey(action.ChannelID, action.MatchID)
	token := uuid.NewString()
	state, err := w.deliveryStore.ClaimMatchEndedDelivery(ctx, key, token, w.matchEndedDeliveryLease())
	if err != nil {
		return fmt.Errorf("claim match ended delivery: %w", err)
	}
	switch state {
	case matchEndedDeliveryPending:
		return errMatchEndedDeliveryPending
	case matchEndedDeliveryDelivered:
		return nil
	case matchEndedDeliveryAcquired:
	default:
		return fmt.Errorf("unexpected match ended delivery state: %d", state)
	}
	deliveryHeartbeat := w.startMatchEndedDeliveryHeartbeat(heartbeat, key, token)

	if err := heartbeat.operationError(); err != nil {
		return w.releaseMatchEndedDeliveryFailure(ctx, key, token, deliveryHeartbeat, err)
	}
	if err := w.emitter.MatchEnded(ctx, busdota.MatchEndedMessage{
		ChannelID:      action.ChannelID.String(),
		SteamAccountID: action.SteamAccountID,
		Win:            action.Win,
		HeroName:       action.HeroName,
		Mmr:            settings.Mmr,
		SessionWins:    settings.SessionWins,
		SessionLosses:  settings.SessionLosses,
		MatchID:        action.MatchID,
	}); err != nil {
		return w.releaseMatchEndedDeliveryFailure(ctx, key, token, deliveryHeartbeat, err)
	}
	if err := deliveryHeartbeat.stop(); err != nil {
		return w.releaseMatchEndedDeliveryFailure(ctx, key, token, nil, err)
	}

	return w.completeMatchEndedDelivery(ctx, key, token)
}

func (w *LifecycleWorker) matchEndedDeliveryLease() time.Duration {
	if w.actionLease > 0 {
		return w.actionLease
	}

	return predictionActionLease
}

func (w *LifecycleWorker) completeMatchEndedDelivery(ctx context.Context, key string, token string) error {
	cleanupCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 2*time.Second)
	defer cancel()
	completed, err := w.deliveryStore.CompleteMatchEndedDelivery(
		cleanupCtx,
		key,
		token,
		matchEndedDeliveryTTL,
	)
	if err != nil {
		return fmt.Errorf("complete match ended delivery: %w", err)
	}
	if !completed {
		return errMatchEndedDeliveryOwnershipLost
	}

	return nil
}

func (w *LifecycleWorker) releaseMatchEndedDeliveryFailure(
	ctx context.Context,
	key string,
	token string,
	heartbeat *matchEndedDeliveryHeartbeat,
	cause error,
) error {
	if heartbeat != nil {
		if err := heartbeat.stop(); err != nil {
			cause = fmt.Errorf("%w; stop match ended delivery heartbeat: %v", cause, err)
		}
	}
	cleanupCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 2*time.Second)
	defer cancel()
	if err := w.deliveryStore.ReleaseMatchEndedDelivery(cleanupCtx, key, token); err != nil {
		return fmt.Errorf("%w; release match ended delivery: %v", cause, err)
	}

	return cause
}

type actionHeartbeat struct {
	ctx    context.Context
	cancel context.CancelFunc
	stopCh chan struct{}
	done   chan struct{}

	stopOnce sync.Once
	errMu    sync.Mutex
	err      error
}

type matchEndedDeliveryHeartbeat struct {
	parent *actionHeartbeat
	key    string
	token  string

	stopCh chan struct{}
	done   chan struct{}

	stopOnce sync.Once
	errMu    sync.Mutex
	err      error
}

var errActionHeartbeatStopped = errors.New("action heartbeat stopped")

func (w *LifecycleWorker) startActionHeartbeat(ctx context.Context, claimed model.ClaimedOutboxAction) *actionHeartbeat {
	operationCtx, cancel := context.WithCancel(ctx)
	heartbeat := &actionHeartbeat{
		ctx:    operationCtx,
		cancel: cancel,
		stopCh: make(chan struct{}),
		done:   make(chan struct{}),
	}

	go func() {
		defer close(heartbeat.done)
		defer cancel()

		ticker := time.NewTicker(w.actionRenewEvery)
		defer ticker.Stop()
		for {
			select {
			case <-heartbeat.stopCh:
				return
			case <-operationCtx.Done():
				return
			case <-ticker.C:
				if err := w.renewPredictionAction(operationCtx, heartbeat.stopCh, claimed); err != nil {
					if errors.Is(err, errActionHeartbeatStopped) ||
						(errors.Is(err, context.Canceled) && operationCtx.Err() != nil) {
						return
					}
					heartbeat.setError(fmt.Errorf("renew prediction action lease: %w", err))
					cancel()
					return
				}
			}
		}
	}()

	return heartbeat
}

func (w *LifecycleWorker) startMatchEndedDeliveryHeartbeat(
	parent *actionHeartbeat,
	key string,
	token string,
) *matchEndedDeliveryHeartbeat {
	heartbeat := &matchEndedDeliveryHeartbeat{
		parent: parent,
		key:    key,
		token:  token,
		stopCh: make(chan struct{}),
		done:   make(chan struct{}),
	}

	go func() {
		defer close(heartbeat.done)

		ticker := time.NewTicker(w.actionRenewEvery)
		defer ticker.Stop()
		for {
			select {
			case <-heartbeat.stopCh:
				return
			case <-parent.ctx.Done():
				return
			case <-ticker.C:
				if err := w.renewMatchEndedDelivery(parent.ctx, heartbeat.stopCh, key, token); err != nil {
					if errors.Is(err, errMatchEndedDeliveryHeartbeatStop) ||
						(errors.Is(err, context.Canceled) && parent.ctx.Err() != nil) {
						return
					}
					heartbeat.fail(fmt.Errorf("renew match ended delivery lease: %w", err))
					return
				}
			}
		}
	}()

	return heartbeat
}

func (w *LifecycleWorker) renewPredictionAction(
	ctx context.Context,
	stop <-chan struct{},
	claimed model.ClaimedOutboxAction,
) error {
	renewCtx, cancel := context.WithTimeout(ctx, w.renewTimeout())
	defer cancel()

	result := make(chan error, 1)
	// Keep the heartbeat responsive if a repository implementation ignores cancellation.
	go func() {
		result <- w.repository.RenewPredictionAction(renewCtx, claimed.ID, claimed.LockToken, w.actionLease)
	}()

	select {
	case err := <-result:
		if err != nil {
			return err
		}
		return renewCtx.Err()
	case <-renewCtx.Done():
		return renewCtx.Err()
	case <-stop:
		return errActionHeartbeatStopped
	}
}

func (w *LifecycleWorker) renewMatchEndedDelivery(
	ctx context.Context,
	stop <-chan struct{},
	key string,
	token string,
) error {
	renewCtx, cancel := context.WithTimeout(ctx, w.renewTimeout())
	defer cancel()

	type result struct {
		renewed bool
		err     error
	}
	response := make(chan result, 1)
	// Keep the heartbeat responsive if a Redis implementation ignores cancellation.
	go func() {
		renewed, err := w.deliveryStore.RenewMatchEndedDelivery(
			renewCtx,
			key,
			token,
			w.matchEndedDeliveryLease(),
		)
		response <- result{renewed: renewed, err: err}
	}()

	select {
	case result := <-response:
		if result.err != nil {
			return result.err
		}
		if !result.renewed {
			return errMatchEndedDeliveryOwnershipLost
		}
		return renewCtx.Err()
	case <-renewCtx.Done():
		return renewCtx.Err()
	case <-stop:
		return errMatchEndedDeliveryHeartbeatStop
	}
}

func (w *LifecycleWorker) renewTimeout() time.Duration {
	timeout := predictionActionRenewTimeout
	if w.actionLease > 0 && w.actionLease/3 < timeout {
		timeout = w.actionLease / 3
	}

	return timeout
}

func (h *actionHeartbeat) stop() error {
	h.stopOnce.Do(func() { close(h.stopCh) })
	<-h.done

	h.errMu.Lock()
	defer h.errMu.Unlock()
	return h.err
}

func (h *actionHeartbeat) setError(err error) {
	h.errMu.Lock()
	defer h.errMu.Unlock()
	if h.err == nil {
		h.err = err
	}
}

func (h *actionHeartbeat) fail(err error) {
	h.setError(err)
	h.cancel()
}

func (h *matchEndedDeliveryHeartbeat) stop() error {
	h.stopOnce.Do(func() { close(h.stopCh) })
	<-h.done

	h.errMu.Lock()
	defer h.errMu.Unlock()
	return h.err
}

func (h *matchEndedDeliveryHeartbeat) setError(err error) {
	h.errMu.Lock()
	defer h.errMu.Unlock()
	if h.err == nil {
		h.err = err
	}
}

func (h *matchEndedDeliveryHeartbeat) fail(err error) {
	h.setError(err)
	h.parent.fail(err)
}

func (h *actionHeartbeat) operationError() error {
	h.errMu.Lock()
	defer h.errMu.Unlock()
	if h.err != nil {
		return h.err
	}

	return h.ctx.Err()
}

func lifecycleActionFromClaim(claimed model.ClaimedOutboxAction) (match.LifecycleAction, error) {
	if claimed.ChannelID == uuid.Nil {
		return match.LifecycleAction{}, errors.New("claimed action channel ID is required")
	}
	if claimed.MatchID <= 0 {
		return match.LifecycleAction{}, errors.New("claimed action match ID must be positive")
	}

	var action match.LifecycleAction
	if err := json.Unmarshal(claimed.Payload, &action); err != nil {
		return match.LifecycleAction{}, fmt.Errorf("decode lifecycle action payload: %w", err)
	}
	if action.ChannelID != claimed.ChannelID {
		return match.LifecycleAction{}, errors.New("lifecycle action channel ID does not match claim")
	}
	if action.MatchID != claimed.MatchID {
		return match.LifecycleAction{}, errors.New("lifecycle action match ID does not match claim")
	}

	var expectedKind match.ActionKind
	switch claimed.Action {
	case model.OutboxActionCreate:
		expectedKind = match.ActionCreate
	case model.OutboxActionResolve:
		expectedKind = match.ActionResolve
	case model.OutboxActionCancel:
		expectedKind = match.ActionCancel
	default:
		return match.LifecycleAction{}, fmt.Errorf("unsupported claimed action %q", claimed.Action)
	}
	if action.Kind != expectedKind {
		return match.LifecycleAction{}, fmt.Errorf(
			"lifecycle action kind %q does not match claim %q",
			action.Kind,
			claimed.Action,
		)
	}

	return action, nil
}

func retryDelay(attempts int) time.Duration {
	delay := time.Second
	for attempt := 1; attempt < attempts && delay < predictionActionRetryCap; attempt++ {
		delay *= 2
	}
	if delay > predictionActionRetryCap {
		return predictionActionRetryCap
	}

	return delay
}
