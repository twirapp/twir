package gamesvoteban

import (
	"context"

	"github.com/lib/pq"
	"github.com/samber/lo"
	"github.com/twirapp/twir/libs/audit"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	votebanentity "github.com/twirapp/twir/libs/entities/voteban"
	apperrors "github.com/twirapp/twir/libs/errors"
	channelsgamesvoteban "github.com/twirapp/twir/libs/repositories/channels_games_voteban"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	Repository    channelsgamesvoteban.Repository
	AuditRecorder audit.Recorder
	Cacher        *generic_cacher.GenericCacher[votebanentity.Voteban]
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
	cacher        *generic_cacher.GenericCacher[votebanentity.Voteban]
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
	NeededVotes:              2,
	VoteDuration:             10,
	VotingMode:               votebanentity.VotingModeChat,
	ChatVotesWordsPositive:   pq.StringArray{"Yay"},
	ChatVotesWordsNegative:   pq.StringArray{"Nay"},
}

func (s *Service) GetByChannelID(ctx context.Context, channelID string) (
	votebanentity.Voteban,
	error,
) {
	result, err := s.repository.GetOrCreateByChannelID(ctx, channelID, defaultSettings)
	if err != nil {
		return votebanentity.Nil, apperrors.NewInternalError("Failed to get or create voteban settings", err)
	}

	if err = s.cacher.Invalidate(ctx, result.ChannelID); err != nil {
		return votebanentity.Nil, apperrors.NewInternalError("Failed to invalidate cache", err)
	}

	return result, nil
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
	VotingMode               *votebanentity.VotingMode
	ChatVotesWordsPositive   []string
	ChatVotesWordsNegative   []string
}

func (s *Service) Update(ctx context.Context, input UpdateInput) (votebanentity.Voteban, error) {
	// Get current entity first (or create with defaults)
	currentEntity, err := s.repository.GetOrCreateByChannelID(ctx, input.ChannelID, defaultSettings)
	if err != nil {
		return votebanentity.Nil, apperrors.NewInternalError("Failed to get or create voteban settings", err)
	}

	// Check if vote options words are unique (usually there are no more than 2-3 of them)
	for _, positive := range input.ChatVotesWordsPositive {
		for _, negative := range input.ChatVotesWordsNegative {
			if negative == positive {
				return votebanentity.Voteban{}, apperrors.NewValidationError(
					"Vote option words must be unique",
					map[string]any{
						"positive": input.ChatVotesWordsPositive,
						"negative": input.ChatVotesWordsNegative,
					},
				)
			}
		}
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
		VotingMode:               input.VotingMode,
	}

	updatedEntity, err := s.repository.Update(ctx, currentEntity.ID, updateInput)
	if err != nil {
		return votebanentity.Nil, apperrors.NewInternalError("Failed to update voteban settings", err)
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
		return votebanentity.Nil, apperrors.NewInternalError("Failed to invalidate cache", err)
	}

	return updatedEntity, nil
}
