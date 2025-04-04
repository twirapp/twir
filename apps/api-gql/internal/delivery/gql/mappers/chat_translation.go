package mappers

import (
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
)

func ChatTranslationEntityTo(e entity.ChatTranslation) gqlmodel.ChatTranslation {
	return gqlmodel.ChatTranslation{
		ID:                e.ID.String(),
		ChannelID:         e.ChannelID,
		Enabled:           e.Enabled,
		TargetLanguage:    e.TargetLanguage,
		ExcludedLanguages: e.ExcludedLanguages,
		UseItalic:         e.UseItalic,
		ExcludedUsersIDs:  e.ExcludedUsersIDs,
		CreatedAt:         e.CreatedAt,
		UpdatedAt:         e.UpdatedAt,
	}
}
