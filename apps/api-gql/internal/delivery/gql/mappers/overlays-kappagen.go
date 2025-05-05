package mappers

import (
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
)

func MapKappagenEntityToGQL(entity entity.KappagenOverlay) gqlmodel.KappagenOverlay {
	animations := make([]gqlmodel.KappagenOverlayAnimationsSettings, 0, len(entity.Animations))
	for _, a := range entity.Animations {
		animations = append(
			animations, gqlmodel.KappagenOverlayAnimationsSettings{
				Style: a.Style,
				Prefs: &gqlmodel.KappagenOverlayAnimationsPrefsSettings{
					Size:    a.Prefs.Size,
					Center:  a.Prefs.Center,
					Speed:   a.Prefs.Speed,
					Faces:   a.Prefs.Faces,
					Message: a.Prefs.Message,
					Time:    a.Prefs.Time,
				},
				Count:   a.Count,
				Enabled: a.Enabled,
			},
		)
	}

	return gqlmodel.KappagenOverlay{
		ID:             entity.ID,
		EnableSpawn:    entity.EnableSpawn,
		ExcludedEmotes: entity.ExcludedEmotes,
		EnableRave:     entity.EnableRave,
		Animation: &gqlmodel.KappagenOverlayAnimationSettings{
			FadeIn:  entity.Animation.FadeIn,
			FadeOut: entity.Animation.FadeOut,
			ZoomIn:  entity.Animation.ZoomIn,
			ZoomOut: entity.Animation.ZoomOut,
		},
		Animations: animations,
	}
}
