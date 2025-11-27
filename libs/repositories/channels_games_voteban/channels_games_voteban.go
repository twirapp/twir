package channels_games_voteban

import (
	"context"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/twirapp/twir/libs/repositories/channels_games_voteban/model"
)

type Repository interface {
	GetByChannelID(ctx context.Context, channelID string) (model.VoteBan, error)
	GetOrCreateByChannelID(ctx context.Context, channelID string, input CreateInput) (model.VoteBan, error)
	Update(ctx context.Context, id uuid.UUID, input UpdateInput) (model.VoteBan, error)
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
	VotingMode               model.VotingMode
	ChatVotesWordsPositive   pq.StringArray
	ChatVotesWordsNegative   pq.StringArray
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
	VotingMode               *model.VotingMode
	ChatVotesWordsPositive   []string
	ChatVotesWordsNegative   []string
}
