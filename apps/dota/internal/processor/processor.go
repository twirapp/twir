package processor

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/twirapp/twir/apps/dota/internal/gsi"
	"github.com/twirapp/twir/apps/dota/internal/match"
	"github.com/twirapp/twir/libs/logger"
	"go.uber.org/fx"
)

const winProbabilityTimeout = 5 * time.Second

type winProbabilityKey struct {
	channelID uuid.UUID
	matchID   int64
}

type WinProbabilityProvider interface {
	WinProbability(ctx context.Context, matchID int64) (float64, error)
}

type Processor struct {
	match      *match.StateMachine
	stats      WinProbabilityProvider
	logger     *slog.Logger
	serviceCtx context.Context

	jobsMu   sync.Mutex
	jobs     sync.WaitGroup
	stopping bool
	inFlight map[winProbabilityKey]struct{}
}

var _ gsi.MatchProcessor = (*Processor)(nil)

func New(
	matchState *match.StateMachine,
	stats WinProbabilityProvider,
	logger *slog.Logger,
	lifecycle fx.Lifecycle,
) *Processor {
	serviceCtx, cancel := context.WithCancel(context.Background())
	p := &Processor{
		match:      matchState,
		stats:      stats,
		logger:     logger,
		serviceCtx: serviceCtx,
		inFlight:   make(map[winProbabilityKey]struct{}),
	}
	lifecycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			p.jobsMu.Lock()
			p.stopping = true
			cancel()
			p.jobsMu.Unlock()

			jobsDone := make(chan struct{})
			go func() {
				p.jobs.Wait()
				close(jobsDone)
			}()

			select {
			case <-jobsDone:
				return nil
			case <-ctx.Done():
				return ctx.Err()
			}
		},
	})

	return p
}

func (p *Processor) Process(ctx context.Context, channelID uuid.UUID, payload gsi.Payload) error {
	if err := p.match.Process(ctx, channelID, payload); err != nil {
		return err
	}

	snapshot, err := p.match.GetSnapshot(ctx, channelID)
	if err != nil {
		return fmt.Errorf("get Dota match snapshot: %w", err)
	}
	if !snapshot.InGame || snapshot.MatchID == 0 {
		return nil
	}

	key := winProbabilityKey{channelID: channelID, matchID: snapshot.MatchID}
	p.jobsMu.Lock()
	if p.stopping {
		p.jobsMu.Unlock()
		return nil
	}
	if _, exists := p.inFlight[key]; exists {
		p.jobsMu.Unlock()
		return nil
	}
	p.inFlight[key] = struct{}{}
	p.jobs.Add(1)
	p.jobsMu.Unlock()

	go p.updateWinProbability(ctx, key)

	return nil
}

func (p *Processor) updateWinProbability(ctx context.Context, key winProbabilityKey) {
	defer func() {
		p.jobsMu.Lock()
		delete(p.inFlight, key)
		p.jobsMu.Unlock()
		p.jobs.Done()
	}()

	requestCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), winProbabilityTimeout)
	stopCancel := context.AfterFunc(p.serviceCtx, cancel)
	defer func() {
		stopCancel()
		cancel()
	}()

	probability, err := p.stats.WinProbability(requestCtx, key.matchID)
	if err != nil {
		p.logger.WarnContext(
			requestCtx,
			"dota processor: failed to fetch win probability",
			logger.Error(err),
			slog.Int64("match_id", key.matchID),
		)
		return
	}

	p.jobsMu.Lock()
	if p.stopping {
		p.jobsMu.Unlock()
		return
	}
	err = p.match.UpdateWinProbability(requestCtx, key.channelID, key.matchID, probability)
	p.jobsMu.Unlock()
	if err != nil {
		p.logger.WarnContext(
			requestCtx,
			"dota processor: failed to update win probability",
			logger.Error(err),
			slog.Int64("match_id", key.matchID),
		)
	}
}
