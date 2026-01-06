package channels_moderation_settings

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	"github.com/twirapp/twir/libs/repositories/channels_moderation_settings"
	"github.com/twirapp/twir/libs/repositories/channels_moderation_settings/model"
	"github.com/twirapp/twir/libs/repositories/plans"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	Repo            channels_moderation_settings.Repository
	PlansRepository plans.Repository
	Cacher          *generic_cacher.GenericCacher[[]model.ChannelModerationSettings]
}

func New(opts Opts) *Service {
	return &Service{
		repo:            opts.Repo,
		plansRepository: opts.PlansRepository,
		cacher:          opts.Cacher,
	}
}

type Service struct {
	repo            channels_moderation_settings.Repository
	plansRepository plans.Repository
	cacher          *generic_cacher.GenericCacher[[]model.ChannelModerationSettings]
}

func (c *Service) modelToEntity(m model.ChannelModerationSettings) entity.ChannelModerationSettings {
	return entity.ChannelModerationSettings{
		ID:                              m.ID,
		Type:                            entity.ModerationSettingsType(m.Type.String()),
		ChannelID:                       m.ChannelID,
		Enabled:                         m.Enabled,
		Name:                            m.Name,
		BanTime:                         m.BanTime,
		BanMessage:                      m.BanMessage,
		WarningMessage:                  m.WarningMessage,
		CheckClips:                      m.CheckClips,
		TriggerLength:                   m.TriggerLength,
		MaxPercentage:                   m.MaxPercentage,
		DeniedChatLanguages:             m.DeniedChatLanguages,
		ExcludedRoles:                   m.ExcludedRoles,
		MaxWarnings:                     m.MaxWarnings,
		DenyList:                        m.DenyList,
		DenyListRegexpEnabled:           m.DenyListRegexpEnabled,
		DenyListWordBoundaryEnabled:     m.DenyListWordBoundaryEnabled,
		DenyListSensitivityEnabled:      m.DenyListSensitivityEnabled,
		OneManSpamMinimumStoredMessages: m.OneManSpamMinimumStoredMessages,
		OneManSpamMessageMemorySeconds:  m.OneManSpamMessageMemorySeconds,
		CreatedAt:                       m.CreatedAt,
		UpdatedAt:                       m.UpdatedAt,
		LanguageExcludedWords:           m.LanguageExcludedWords,
	}
}

func (c *Service) GetByID(ctx context.Context, id uuid.UUID) (
	entity.ChannelModerationSettings,
	error,
) {
	item, err := c.repo.GetByID(ctx, id)
	if err != nil {
		return entity.ChannelModerationSettings{}, err
	}

	return c.modelToEntity(item), nil
}

func (c *Service) GetByChannelID(ctx context.Context, channelID string) (
	[]entity.ChannelModerationSettings,
	error,
) {
	items, err := c.repo.GetByChannelID(ctx, channelID)
	if err != nil {
		return nil, err
	}

	converted := make([]entity.ChannelModerationSettings, len(items))
	for idx, item := range items {
		converted[idx] = c.modelToEntity(item)
	}

	return converted, nil
}

type CreateOrUpdateInput struct {
	Name                            *string
	ChannelID                       string
	Type                            entity.ModerationSettingsType
	Enabled                         bool
	BanTime                         int32
	BanMessage                      string
	WarningMessage                  string
	CheckClips                      bool
	TriggerLength                   int
	MaxPercentage                   int
	DenyList                        []string
	DenyListRegexpEnabled           bool
	DenyListWordBoundaryEnabled     bool
	DenyListSensitivityEnabled      bool
	DeniedChatLanguages             []string
	ExcludedRoles                   []string
	MaxWarnings                     int
	OneManSpamMinimumStoredMessages int
	OneManSpamMessageMemorySeconds  int
	LanguageExcludedWords           []string
}

func (c *Service) Create(
	ctx context.Context,
	input CreateOrUpdateInput,
) (entity.ChannelModerationSettings, error) {
	plan, err := c.plansRepository.GetByChannelID(ctx, input.ChannelID)
	if err != nil {
		return entity.ChannelModerationSettings{}, fmt.Errorf("failed to get plan: %w", err)
	}
	if plan.IsNil() {
		return entity.ChannelModerationSettings{}, fmt.Errorf("plan not found for channel")
	}

	existingSettings, err := c.repo.GetByChannelID(ctx, input.ChannelID)
	if err != nil {
		return entity.ChannelModerationSettings{}, fmt.Errorf("failed to get moderation settings: %w", err)
	}

	if len(existingSettings) >= plan.MaxModerationRules {
		return entity.ChannelModerationSettings{}, fmt.Errorf("you can have only %v moderation rules", plan.MaxModerationRules)
	}

	item, err := c.repo.Create(
		ctx,
		channels_moderation_settings.CreateOrUpdateInput{
			Name:                            input.Name,
			ChannelID:                       input.ChannelID,
			Type:                            model.ModerationSettingsType(input.Type.String()),
			Enabled:                         input.Enabled,
			BanTime:                         input.BanTime,
			BanMessage:                      input.BanMessage,
			WarningMessage:                  input.WarningMessage,
			CheckClips:                      input.CheckClips,
			TriggerLength:                   input.TriggerLength,
			MaxPercentage:                   input.MaxPercentage,
			DenyList:                        input.DenyList,
			DenyListRegexpEnabled:           input.DenyListRegexpEnabled,
			DenyListWordBoundaryEnabled:     input.DenyListWordBoundaryEnabled,
			DenyListSensitivityEnabled:      input.DenyListSensitivityEnabled,
			DeniedChatLanguages:             input.DeniedChatLanguages,
			ExcludedRoles:                   input.ExcludedRoles,
			MaxWarnings:                     input.MaxWarnings,
			OneManSpamMinimumStoredMessages: input.OneManSpamMinimumStoredMessages,
			OneManSpamMessageMemorySeconds:  input.OneManSpamMessageMemorySeconds,
			LanguageExcludedWords:           input.LanguageExcludedWords,
		},
	)
	if err != nil {
		return entity.ChannelModerationSettings{}, err
	}

	if err := c.cacher.Invalidate(ctx, input.ChannelID); err != nil {
		return entity.ChannelModerationSettings{}, err
	}

	return c.modelToEntity(item), nil
}

func (c *Service) Update(ctx context.Context, id uuid.UUID, input CreateOrUpdateInput) (
	entity.ChannelModerationSettings,
	error,
) {
	item, err := c.GetByID(ctx, id)
	if err != nil {
		return entity.ChannelModerationSettings{}, err
	}

	newItem, err := c.repo.Update(
		ctx,
		id,
		channels_moderation_settings.CreateOrUpdateInput{
			Name:                            input.Name,
			ChannelID:                       item.ChannelID,
			Type:                            model.ModerationSettingsType(input.Type.String()),
			Enabled:                         input.Enabled,
			BanTime:                         input.BanTime,
			BanMessage:                      input.BanMessage,
			WarningMessage:                  input.WarningMessage,
			CheckClips:                      input.CheckClips,
			TriggerLength:                   input.TriggerLength,
			MaxPercentage:                   input.MaxPercentage,
			DenyList:                        input.DenyList,
			DenyListRegexpEnabled:           input.DenyListRegexpEnabled,
			DenyListWordBoundaryEnabled:     input.DenyListWordBoundaryEnabled,
			DenyListSensitivityEnabled:      input.DenyListSensitivityEnabled,
			DeniedChatLanguages:             input.DeniedChatLanguages,
			ExcludedRoles:                   input.ExcludedRoles,
			MaxWarnings:                     input.MaxWarnings,
			OneManSpamMinimumStoredMessages: input.OneManSpamMinimumStoredMessages,
			OneManSpamMessageMemorySeconds:  input.OneManSpamMessageMemorySeconds,
			LanguageExcludedWords:           input.LanguageExcludedWords,
		},
	)
	if err != nil {
		return entity.ChannelModerationSettings{}, err
	}

	if err := c.cacher.Invalidate(ctx, input.ChannelID); err != nil {
		return entity.ChannelModerationSettings{}, err
	}

	return c.modelToEntity(newItem), nil
}

func (c *Service) Delete(ctx context.Context, id uuid.UUID, channelID string) error {
	item, err := c.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if item.ChannelID != channelID {
		return fmt.Errorf("item not found")
	}

	if err := c.cacher.Invalidate(ctx, channelID); err != nil {
		return err
	}

	return c.repo.Delete(ctx, id)
}
