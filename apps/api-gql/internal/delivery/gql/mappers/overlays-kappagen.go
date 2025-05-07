package mappers

import (
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
)

func MapKappagenEntityToGQL(entity entity.KappagenOverlay) gqlmodel.KappagenOverlay {
	animations := make(
		[]gqlmodel.KappagenOverlayAnimationsSettings,
		0,
		len(entity.Settings.Animations),
	)
	for _, a := range entity.Settings.Animations {
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
	events := make([]gqlmodel.KappagenOverlayEvent, 0, len(entity.Settings.Events))
	for _, e := range entity.Settings.Events {
		events = append(
			events, gqlmodel.KappagenOverlayEvent{
				Event:              gqlmodel.EventType(e.Event),
				DisabledAnimations: e.DisabledAnimations,
				Enabled:            e.Enabled,
			},
		)
	}

	return gqlmodel.KappagenOverlay{
		ID:             entity.ID,
		EnableSpawn:    entity.Settings.EnableSpawn,
		ExcludedEmotes: entity.Settings.ExcludedEmotes,
		EnableRave:     entity.Settings.EnableRave,
		Animation: &gqlmodel.KappagenOverlayAnimationSettings{
			FadeIn:  entity.Settings.Animation.FadeIn,
			FadeOut: entity.Settings.Animation.FadeOut,
			ZoomIn:  entity.Settings.Animation.ZoomIn,
			ZoomOut: entity.Settings.Animation.ZoomOut,
		},
		Animations: animations,
		Emotes: &gqlmodel.KappagenEmoteSettings{
			Time:           entity.Settings.Emotes.Time,
			Max:            entity.Settings.Emotes.Max,
			Queue:          entity.Settings.Emotes.Queue,
			FfzEnabled:     entity.Settings.Emotes.FfzEnabled,
			BttvEnabled:    entity.Settings.Emotes.BttvEnabled,
			SevenTvEnabled: entity.Settings.Emotes.SevenTvEnabled,
			EmojiStyle:     mapKappagenEmojiStyleToGql[entity.Settings.Emotes.EmojiStyle],
		},
		Size: &gqlmodel.KappagenSizeSettings{
			RationNormal: entity.Settings.Size.RatioNormal,
			RationSmall:  entity.Settings.Size.RatioSmall,
			Min:          entity.Settings.Size.Min,
			Max:          entity.Settings.Size.Max,
		},
		Events:    events,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}
}

var mapKappagenEmojiStyleToGql = map[entity.KappagenEmojiStyle]gqlmodel.KappagenEmojiStyle{
	entity.KappagenEmojiStyleNone:     gqlmodel.KappagenEmojiStyleNone,
	entity.KappagenEmojiStyleTwemoji:  gqlmodel.KappagenEmojiStyleTwemoji,
	entity.KappagenEmojiStyleOpenmoji: gqlmodel.KappagenEmojiStyleOpenmoji,
	entity.KappagenEmojiStyleNoto:     gqlmodel.KappagenEmojiStyleNoto,
	entity.KappagenEmojiStyleBlobmoji: gqlmodel.KappagenEmojiStyleBlobmoji,
}

func MapGqlKappagenEmoteStyleToEntity(style gqlmodel.KappagenEmojiStyle) entity.KappagenEmojiStyle {
	for k, v := range mapKappagenEmojiStyleToGql {
		if v == style {
			return k
		}
	}

	return entity.KappagenEmojiStyleNone
}
