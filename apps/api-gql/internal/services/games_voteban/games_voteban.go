package gamesvoteban

import (
	"context"
	"fmt"

	"github.com/lib/pq"
	"github.com/samber/lo"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	"go.uber.org/fx"

	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	"github.com/twirapp/twir/libs/audit"
	channelsgamesvoteban "github.com/twirapp/twir/libs/repositories/channels_games_voteban"
	"github.com/twirapp/twir/libs/repositories/channels_games_voteban/model"
)

type Opts struct {
	fx.In

	Repository    channelsgamesvoteban.Repository
	AuditRecorder audit.Recorder
	Cacher        *generic_cacher.GenericCacher[model.VoteBan]
}

func New(opts Opts) *Service {
	return &Service{
		repository:    opts.Repository,
		auditRecorder: opts.AuditRecorder,
		cacher:        opts.Cacher,
	}
}

type Service struct {
	repository    channelsgamesvoteban.Repository
	auditRecorder audit.Recorder
	cacher        *generic_cacher.GenericCacher[model.VoteBan]
}

func (s *Service) mapToEntity(m model.VoteBan) entity.GamesVoteBan {
	return entity.GamesVoteBan{
		ID:                       m.ID,
		ChannelID:                m.ChannelID,
		Enabled:                  m.Enabled,
		TimeoutSeconds:           m.TimeoutSeconds,
		TimeoutModerators:        m.TimeoutModerators,
		InitMessage:              m.InitMessage,
		BanMessage:               m.BanMessage,
		BanMessageModerators:     m.BanMessageModerators,
		SurviveMessage:           m.SurviveMessage,
		SurviveMessageModerators: m.SurviveMessageModerators,
		NeededVotes:              m.NeededVotes,
		VoteDuration:             m.VoteDuration,
		VotingMode:               entity.VotingMode(m.VotingMode),
		ChatVotesWordsPositive:   m.ChatVotesWordsPositive,
		ChatVotesWordsNegative:   m.ChatVotesWordsNegative,
	}
}

func (s *Service) mapVotingModeToModel(mode entity.VotingMode) model.VotingMode {
	switch mode {
	case entity.VotingModeChat:
		return model.VotingModeChat
	case entity.VotingModeTwitchPolls:
		return model.VotingModeTwitchPolls
	default:
		return model.VotingModeChat
	}
}

var defaultSettings = channelsgamesvoteban.CreateInput{
	Enabled:        false,
	TimeoutSeconds: 60,
	InitMessage: "The Twitch Police have decided that {targetUser} is not worthy of" +
		" being in chat for not knowing memes. Write \"{positiveTexts}\" to support, " +
		"or \"{negativeTexts}\" if you disagree.",
	BanMessage:               "User {targetUser} is not worthy of being in chat.",
	BanMessageModerators:     "User {targetUser} is not worthy of being in chat.",
	SurviveMessage:           "Looks like something is mixed up, {targetUser} is the kindest and most knowledgeable chat user.",
	SurviveMessageModerators: "Looks like something is mixed up, {targetUser} is the kindest and most knowledgeable chat user.",
	NeededVotes:              1,
	VoteDuration:             1,
	VotingMode:               model.VotingModeChat,
	ChatVotesWordsPositive:   pq.StringArray{"Yay"},
	ChatVotesWordsNegative:   pq.StringArray{"Nay"},
}

func (s *Service) GetByChannelID(ctx context.Context, channelID string) (
	entity.GamesVoteBan,
	error,
) {
	result, err := s.repository.GetOrCreateByChannelID(ctx, channelID, defaultSettings)
	if err != nil {
		return entity.GamesVoteBanNil, fmt.Errorf("failed to get or create games voteban: %w", err)
	}

	if err = s.cacher.Invalidate(ctx, result.ChannelID); err != nil {
		return entity.GamesVoteBanNil, fmt.Errorf("failed to invalidate cache: %w", err)
	}

	return s.mapToEntity(result), nil
}

type UpdateInput struct {
	ActorID   string
	ChannelID string

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
	VotingMode               *entity.VotingMode
	ChatVotesWordsPositive   []string
	ChatVotesWordsNegative   []string
}

func (s *Service) Update(ctx context.Context, input UpdateInput) (entity.GamesVoteBan, error) {
	// Get current entity first (or create with defaults)
	currentEntity, err := s.repository.GetOrCreateByChannelID(ctx, input.ChannelID, defaultSettings)
	if err != nil {
		return entity.GamesVoteBanNil, fmt.Errorf("failed to get or create games voteban: %w", err)
	}

	// Build update input
	updateInput := channelsgamesvoteban.UpdateInput{
		Enabled:                  input.Enabled,
		TimeoutSeconds:           input.TimeoutSeconds,
		TimeoutModerators:        input.TimeoutModerators,
		InitMessage:              input.InitMessage,
		BanMessage:               input.BanMessage,
		BanMessageModerators:     input.BanMessageModerators,
		SurviveMessage:           input.SurviveMessage,
		SurviveMessageModerators: input.SurviveMessageModerators,
		NeededVotes:              input.NeededVotes,
		VoteDuration:             input.VoteDuration,
		ChatVotesWordsPositive:   input.ChatVotesWordsPositive,
		ChatVotesWordsNegative:   input.ChatVotesWordsNegative,
	}

	if input.VotingMode != nil {
		votingMode := s.mapVotingModeToModel(*input.VotingMode)
		updateInput.VotingMode = &votingMode
	}

	updatedEntity, err := s.repository.Update(ctx, currentEntity.ID, updateInput)
	if err != nil {
		return entity.GamesVoteBanNil, fmt.Errorf("failed to update games voteban: %w", err)
	}

	_ = s.auditRecorder.RecordUpdateOperation(
		ctx,
		audit.UpdateOperation{
			Metadata: audit.OperationMetadata{
				System:    "channels_games_voteban",
				ActorID:   lo.ToPtr(input.ActorID),
				ChannelID: lo.ToPtr(input.ChannelID),
				ObjectID:  lo.ToPtr(updatedEntity.ID.String()),
			},
			NewValue: updatedEntity,
			OldValue: currentEntity,
		},
	)

	if err := s.cacher.Invalidate(ctx, input.ChannelID); err != nil {
		return entity.GamesVoteBanNil, fmt.Errorf("failed to invalidate cache: %w", err)
	}

	return s.mapToEntity(updatedEntity), nil
}
