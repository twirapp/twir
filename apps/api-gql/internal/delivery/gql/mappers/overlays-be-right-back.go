package mappers

import (
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
)

func MapBeRightBackEntityToGQL(e entity.BeRightBackOverlay) gqlmodel.BeRightBackOverlay {
	late := MapBeRightBackLateEntityToGQL(e.Settings.Late)
	return gqlmodel.BeRightBackOverlay{
		ID:              e.ID,
		Text:            e.Settings.Text,
		Late:            &late,
		BackgroundColor: e.Settings.BackgroundColor,
		FontSize:        int(e.Settings.FontSize),
		FontColor:       e.Settings.FontColor,
		FontFamily:      e.Settings.FontFamily,
		CreatedAt:       e.CreatedAt,
		UpdatedAt:       e.UpdatedAt,
		ChannelID:       e.ChannelID,
	}
}

func MapBeRightBackLateEntityToGQL(e entity.BeRightBackOverlayLateSettings) gqlmodel.BeRightBackLate {
	return gqlmodel.BeRightBackLate{
		Enabled:        e.Enabled,
		Text:           e.Text,
		DisplayBrbTime: e.DisplayBrbTime,
	}
}
