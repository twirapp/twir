package channels_games_voteban

import (
	"context"

	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/entities/voteban"
)

type Repository interface {
	GetByChannelID(ctx context.Context, channelID string) (voteban.Voteban, error)
	GetOrCreateByChannelID(ctx context.Context, channelID string, input CreateInput) (
		voteban.Voteban,
		error,
	)
	Update(ctx context.Context, id uuid.UUID, input UpdateInput) (voteban.Voteban, error)
}

type CreateInput struct {
	Enabled                  bool
	TimeoutSeconds           int
	TimeoutModerators        bool
	InitMessage              string
	BanMessage               string
	BanMessageModerators     string
	SurviveMessage           string
	SurviveMessageModerators string
	NeededVotes              int
	VoteDuration             int
	VotingMode               voteban.VotingMode
	ChatVotesWordsPositive   []string
	ChatVotesWordsNegative   []string
}

type UpdateInput struct {
	Enabled                  *bool
	TimeoutSeconds           *int
	TimeoutModerators        *bool
	InitMessage              *string
	BanMessage               *string
	BanMessageModerators     *string
	SurviveMessage           *string
	SurviveMessageModerators *string
	NeededVotes              *int
	VoteDuration             *int
	VotingMode               *voteban.VotingMode
	ChatVotesWordsPositive   []string
	ChatVotesWordsNegative   []string
}
