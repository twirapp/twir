package processor

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/twirapp/twir/apps/dota/internal/gsi"
	"github.com/twirapp/twir/apps/dota/internal/match"
	"github.com/twirapp/twir/libs/logger"
)

type WinProbabilityProvider interface {
	WinProbability(ctx context.Context, matchID int64) (float64, error)
}

type Processor struct {
	match  *match.StateMachine
	stats  WinProbabilityProvider
	logger *slog.Logger
}

var _ gsi.MatchProcessor = (*Processor)(nil)

func New(
	matchState *match.StateMachine,
	stats WinProbabilityProvider,
	logger *slog.Logger,
) *Processor {
	return &Processor{
		match:  matchState,
		stats:  stats,
		logger: logger,
	}
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

	probability, err := p.stats.WinProbability(ctx, snapshot.MatchID)
	if err != nil {
		p.logger.WarnContext(
			ctx,
			"dota processor: failed to fetch win probability",
			logger.Error(err),
			slog.Int64("match_id", snapshot.MatchID),
		)
		return nil
	}

	if err := p.match.UpdateWinProbability(ctx, channelID, probability); err != nil {
		p.logger.WarnContext(
			ctx,
			"dota processor: failed to update win probability",
			logger.Error(err),
			slog.Int64("match_id", snapshot.MatchID),
		)
	}

	return nil
}
