package channelsgamesvotebanprogressstate

import (
	"context"
	"time"

	"github.com/twirapp/twir/libs/repositories/channels_games_voteban_progress_state/model"
)

type Repository interface {
	// Get returns the current voting state for a channel.
	// Returns ErrNotFound if no voting is in progress.
	Get(ctx context.Context, channelID string) (model.VoteState, error)

	// Exists checks if a voting session is in progress for a channel.
	Exists(ctx context.Context, channelID string) (bool, error)

	// Create creates a new voting session for a channel with the specified TTL.
	Create(ctx context.Context, channelID string, state model.VoteState, ttl time.Duration) error

	// Update updates the voting state for a channel.
	Update(ctx context.Context, channelID string, state model.VoteState) error

	// Delete removes the voting session for a channel.
	Delete(ctx context.Context, channelID string) error

	// UserHasVoted checks if a user has already voted in the current session.
	UserHasVoted(ctx context.Context, channelID, userID string) (bool, error)

	// MarkUserVoted marks a user as having voted in the current session.
	MarkUserVoted(ctx context.Context, channelID, userID string, ttl time.Duration) error

	// ClearUserVotes removes all user vote markers for a channel.
	ClearUserVotes(ctx context.Context, channelID string) error
}
