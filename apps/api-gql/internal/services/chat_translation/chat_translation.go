package chat_translation

import (
	"context"
	"errors"
	"fmt"

	"github.com/oklog/ulid/v2"
	"github.com/samber/lo"
	"github.com/satont/twir/libs/logger"
	"github.com/satont/twir/libs/logger/audit"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	repo "github.com/twirapp/twir/libs/repositories/chat_translation"
	"github.com/twirapp/twir/libs/repositories/chat_translation/model"
	"go.uber.org/fx"

	"github.com/twirapp/twir/apps/api-gql/internal/entity"
)

type Opts struct {
	fx.In

	ChatTranslationRepository repo.Repository
	Logger                    logger.Logger
	TranslationsSettingsCache *generic_cacher.GenericCacher[model.ChatTranslation]
}

func New(opts Opts) *Service {
	return &Service{
		chatTranslationRepository: opts.ChatTranslationRepository,
		logger:                    opts.Logger,
		translationsSettingsCache: opts.TranslationsSettingsCache,
	}
}

type Service struct {
	chatTranslationRepository repo.Repository
	logger                    logger.Logger
	translationsSettingsCache *generic_cacher.GenericCacher[model.ChatTranslation]
}

func chatTranslationModelToEntity(m model.ChatTranslation) entity.ChatTranslation {
	return entity.ChatTranslation{
		ID:                m.ID,
		ChannelID:         m.ChannelID,
		Enabled:           m.Enabled,
		TargetLanguage:    m.TargetLanguage,
		ExcludedLanguages: m.ExcludedLanguages,
		UseItalic:         m.UseItalic,
		ExcludedUsersIDs:  m.ExcludedUsersIDs,
		CreatedAt:         m.CreatedAt,
		UpdatedAt:         m.UpdatedAt,
	}
}

func (s *Service) GetByChannelID(ctx context.Context, channelID string) (
	entity.ChatTranslation,
	error,
) {
	translation, err := s.chatTranslationRepository.GetByChannelID(ctx, channelID)
	if err != nil {
		return entity.ChatTranslation{}, fmt.Errorf("failed to get chat translation: %w", err)
	}

	return chatTranslationModelToEntity(translation), nil
}

type CreateInput struct {
	ChannelID string
	ActorID   string

	Enabled           bool
	TargetLanguage    string
	ExcludedLanguages []string
	UseItalic         bool
	ExcludedUsersIDs  []string
}

func (s *Service) Create(ctx context.Context, input CreateInput) (
	entity.ChatTranslation,
	error,
) {
	translation, err := s.chatTranslationRepository.Create(
		ctx, repo.CreateInput{
			ChannelID:         input.ChannelID,
			Enabled:           input.Enabled,
			TargetLanguage:    input.TargetLanguage,
			ExcludedLanguages: input.ExcludedLanguages,
			UseItalic:         input.UseItalic,
			ExcludedUsersIDs:  input.ExcludedUsersIDs,
		},
	)
	if err != nil {
		return entity.ChatTranslation{}, fmt.Errorf("failed to create chat translation: %w", err)
	}

	s.logger.Audit(
		"Chat translation created", audit.Fields{
			System:        "channels_chat_translation",
			OperationType: audit.OperationCreate,
			OldValue:      nil,
			NewValue:      translation,
			ActorID:       &input.ActorID,
			ChannelID:     &input.ChannelID,
			ObjectID:      lo.ToPtr(translation.ID.String()),
		},
	)

	if err := s.translationsSettingsCache.Invalidate(ctx, input.ChannelID); err != nil {
		return entity.ChatTranslation{}, fmt.Errorf("failed to invalidate cache: %w", err)
	}

	return chatTranslationModelToEntity(translation), nil
}

type UpdateInput struct {
	ChannelID string
	ActorID   string

	Enabled           *bool
	TargetLanguage    *string
	ExcludedLanguages *[]string
	UseItalic         *bool
	ExcludedUsersIDs  *[]string
}

func (s *Service) Update(
	ctx context.Context,
	id ulid.ULID,
	input UpdateInput,
) (entity.ChatTranslation, error) {
	// First get the existing translation to get the ID
	existingTranslation, err := s.GetByChannelID(ctx, input.ChannelID)
	if err != nil {
		if errors.Is(err, repo.ErrSettingsNotFound) {
			return entity.ChatTranslation{}, errors.New("chat translation settings not found")
		}
		return entity.ChatTranslation{}, err
	}

	// Check if the settings belong to the current channel
	if existingTranslation.ChannelID != input.ChannelID {
		return entity.ChatTranslation{}, errors.New("chat translation settings do not belong to this channel")
	}

	translation, err := s.chatTranslationRepository.Update(
		ctx,
		id,
		repo.UpdateInput{
			Enabled:           input.Enabled,
			TargetLanguage:    input.TargetLanguage,
			ExcludedLanguages: input.ExcludedLanguages,
			UseItalic:         input.UseItalic,
			ExcludedUsersIDs:  input.ExcludedUsersIDs,
		},
	)
	if err != nil {
		return entity.ChatTranslation{}, fmt.Errorf("failed to update chat translation: %w", err)
	}

	s.logger.Audit(
		"Chat translation updated", audit.Fields{
			OldValue:      translation,
			ActorID:       &input.ActorID,
			ChannelID:     &input.ChannelID,
			System:        "channels_chat_translation",
			OperationType: audit.OperationUpdate,
			ObjectID:      lo.ToPtr(translation.ID.String()),
		},
	)

	if err := s.translationsSettingsCache.Invalidate(ctx, input.ChannelID); err != nil {
		return entity.ChatTranslation{}, fmt.Errorf("failed to invalidate cache: %w", err)
	}

	return chatTranslationModelToEntity(translation), nil
}

type DeleteInput struct {
	ID        ulid.ULID
	ChannelID string
	ActorID   string
}

func (s *Service) Delete(ctx context.Context, input DeleteInput) error {
	oldTranslation, err := s.chatTranslationRepository.GetByChannelID(ctx, input.ChannelID)
	if err != nil {
		return fmt.Errorf("failed to get chat translation: %w", err)
	}
	// Check if the settings belong to the current channel
	if oldTranslation.ChannelID != input.ChannelID {
		return fmt.Errorf("chat translation settings do not belong to this channel")
	}

	if err := s.chatTranslationRepository.Delete(ctx, input.ID); err != nil {
		return fmt.Errorf("failed to delete chat translation: %w", err)
	}

	s.logger.Audit(
		"Chat translation deleted", audit.Fields{
			System:        "channels_chat_translation",
			OperationType: audit.OperationDelete,
			OldValue:      &oldTranslation,
			NewValue:      nil,
			ActorID:       &input.ActorID,
			ChannelID:     &input.ChannelID,
			ObjectID:      lo.ToPtr(oldTranslation.ID.String()),
		},
	)

	if err := s.translationsSettingsCache.Invalidate(ctx, input.ChannelID); err != nil {
		return fmt.Errorf("failed to invalidate cache: %w", err)
	}

	return nil
}
