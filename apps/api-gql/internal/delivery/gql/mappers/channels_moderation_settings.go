package mappers

import (
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
)

var mapModerationTypeToGql = map[entity.ModerationSettingsType]gqlmodel.ModerationSettingsType{
	entity.ModerationSettingsTypeLinks:       gqlmodel.ModerationSettingsTypeLinks,
	entity.ModerationSettingsTypeDenylist:    gqlmodel.ModerationSettingsTypeDenyList,
	entity.ModerationSettingsTypeSymbols:     gqlmodel.ModerationSettingsTypeSymbols,
	entity.ModerationSettingsTypeLongMessage: gqlmodel.ModerationSettingsTypeLongMessage,
	entity.ModerationSettingsTypeCaps:        gqlmodel.ModerationSettingsTypeCaps,
	entity.ModerationSettingsTypeEmotes:      gqlmodel.ModerationSettingsTypeEmotes,
	entity.ModerationSettingsTypeLanguage:    gqlmodel.ModerationSettingsTypeLanguage,
	entity.ModerationSettingsTypeOneManSpam:  gqlmodel.ModerationSettingsTypeOneManSpam,
}

func ModerationSettingsEntityToGql(m entity.ChannelModerationSettings) gqlmodel.ModerationSettingsItem {
	return gqlmodel.ModerationSettingsItem{
		ID:                              m.ID,
		Type:                            ModerationSettingsEntityTypeToGql(m.Type),
		Name:                            m.Name,
		Enabled:                         m.Enabled,
		BanTime:                         int(m.BanTime),
		BanMessage:                      m.BanMessage,
		WarningMessage:                  m.WarningMessage,
		CheckClips:                      m.CheckClips,
		TriggerLength:                   m.TriggerLength,
		MaxPercentage:                   m.MaxPercentage,
		DenyList:                        m.DenyList,
		DeniedChatLanguages:             m.DeniedChatLanguages,
		ExcludedRoles:                   m.ExcludedRoles,
		MaxWarnings:                     m.MaxWarnings,
		CreatedAt:                       m.CreatedAt,
		UpdatedAt:                       m.UpdatedAt,
		DenyListSensitivityEnabled:      m.DenyListSensitivityEnabled,
		DenyListRegexpEnabled:           m.DenyListRegexpEnabled,
		DenyListWordBoundaryEnabled:     m.DenyListWordBoundaryEnabled,
		OneManSpamMinimumStoredMessages: m.OneManSpamMinimumStoredMessages,
		OneManSpamMessageMemorySeconds:  m.OneManSpamMessageMemorySeconds,
	}
}

func ModerationSettingsEntityTypeToGql(t entity.ModerationSettingsType) gqlmodel.ModerationSettingsType {
	return mapModerationTypeToGql[t]
}

func ModerationSettingsTypeToEntity(t gqlmodel.ModerationSettingsType) entity.ModerationSettingsType {
	for k, v := range mapModerationTypeToGql {
		if v == t {
			return k
		}
	}

	return ""
}
