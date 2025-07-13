package mappers

import (
	"fmt"

	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
)

func MapKappagenEntityToGQL(e entity.KappagenOverlay) gqlmodel.KappagenOverlay {
	animations := make(
		[]gqlmodel.KappagenOverlayAnimationsSettings,
		0,
		len(e.Settings.Animations),
	)
	for _, a := range e.Settings.Animations {
		var prefs *gqlmodel.KappagenOverlayAnimationsPrefsSettings = nil
		if a.Prefs != nil {
			switch a.Style {
			case entity.KappagenOverlayAnimationStyleTheCube:
				prefs = &gqlmodel.KappagenOverlayAnimationsPrefsSettings{
					Size:   a.Prefs.Size,
					Center: a.Prefs.Center,
					Speed:  a.Prefs.Speed,
					Faces:  a.Prefs.Faces,
				}
			case entity.KappagenOverlayAnimationStyleText:
				prefs = &gqlmodel.KappagenOverlayAnimationsPrefsSettings{
					Message: a.Prefs.Message,
				}
			}
		}

		animationStyle := mapKappagenAnimationStyleToGql[a.Style]

		animations = append(
			animations, gqlmodel.KappagenOverlayAnimationsSettings{
				Style:   animationStyle,
				Prefs:   prefs,
				Count:   a.Count,
				Enabled: a.Enabled,
			},
		)
	}
	events := make([]gqlmodel.KappagenOverlayEvent, 0, len(e.Settings.Events))
	for _, e := range e.Settings.Events {
		events = append(
			events, gqlmodel.KappagenOverlayEvent{
				Event:              gqlmodel.EventType(e.Event),
				DisabledAnimations: e.DisabledAnimations,
				Enabled:            e.Enabled,
			},
		)
	}

	return gqlmodel.KappagenOverlay{
		ID:             e.ID,
		EnableSpawn:    e.Settings.EnableSpawn,
		ExcludedEmotes: e.Settings.ExcludedEmotes,
		EnableRave:     e.Settings.EnableRave,
		Animation: &gqlmodel.KappagenOverlayAnimationSettings{
			FadeIn:  e.Settings.Animation.FadeIn,
			FadeOut: e.Settings.Animation.FadeOut,
			ZoomIn:  e.Settings.Animation.ZoomIn,
			ZoomOut: e.Settings.Animation.ZoomOut,
		},
		Animations: animations,
		Emotes: &gqlmodel.KappagenEmoteSettings{
			Time:           e.Settings.Emotes.Time,
			Max:            e.Settings.Emotes.Max,
			Queue:          e.Settings.Emotes.Queue,
			FfzEnabled:     e.Settings.Emotes.FfzEnabled,
			BttvEnabled:    e.Settings.Emotes.BttvEnabled,
			SevenTvEnabled: e.Settings.Emotes.SevenTvEnabled,
			EmojiStyle:     mapKappagenEmojiStyleToGql[e.Settings.Emotes.EmojiStyle],
		},
		Size: &gqlmodel.KappagenSizeSettings{
			RationNormal: e.Settings.Size.RatioNormal,
			RationSmall:  e.Settings.Size.RatioSmall,
			Min:          e.Settings.Size.Min,
			Max:          e.Settings.Size.Max,
		},
		Events:    events,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
		ChannelID: e.ChannelID,
	}
}

var mapKappagenEmojiStyleToGql = map[entity.KappagenEmojiStyle]gqlmodel.KappagenEmojiStyle{
	entity.KappagenEmojiStyleNone:     gqlmodel.KappagenEmojiStyleNone,
	entity.KappagenEmojiStyleTwemoji:  gqlmodel.KappagenEmojiStyleTwemoji,
	entity.KappagenEmojiStyleOpenmoji: gqlmodel.KappagenEmojiStyleOpenmoji,
	entity.KappagenEmojiStyleNoto:     gqlmodel.KappagenEmojiStyleNoto,
	entity.KappagenEmojiStyleBlobmoji: gqlmodel.KappagenEmojiStyleBlobmoji,
}

var mapKappagenAnimationStyleToGql = map[entity.KappagenOverlayAnimationStyle]gqlmodel.KappagenOverlayAnimationStyle{
	entity.KappagenOverlayAnimationStyleConfetti:     gqlmodel.KappagenOverlayAnimationStyleConfetti,
	entity.KappagenOverlayAnimationStyleSpiral:       gqlmodel.KappagenOverlayAnimationStyleSpiral,
	entity.KappagenOverlayAnimationStyleStampede:     gqlmodel.KappagenOverlayAnimationStyleStampede,
	entity.KappagenOverlayAnimationStyleFireworks:    gqlmodel.KappagenOverlayAnimationStyleFireworks,
	entity.KappagenOverlayAnimationStyleFountain:     gqlmodel.KappagenOverlayAnimationStyleFountain,
	entity.KappagenOverlayAnimationStyleBurst:        gqlmodel.KappagenOverlayAnimationStyleBurst,
	entity.KappagenOverlayAnimationStyleTheCube:      gqlmodel.KappagenOverlayAnimationStyleTheCube,
	entity.KappagenOverlayAnimationStyleText:         gqlmodel.KappagenOverlayAnimationStyleText,
	entity.KappagenOverlayAnimationStyleConga:        gqlmodel.KappagenOverlayAnimationStyleConga,
	entity.KappagenOverlayAnimationStyleSmallPyramid: gqlmodel.KappagenOverlayAnimationStyleSmallPyramid,
	entity.KappagenOverlayAnimationStylePyramid:      gqlmodel.KappagenOverlayAnimationStylePyramid,
}

func MapGqlKappagenEmoteStyleToEntity(style gqlmodel.KappagenEmojiStyle) entity.KappagenEmojiStyle {
	for k, v := range mapKappagenEmojiStyleToGql {
		if v == style {
			return k
		}
	}

	return entity.KappagenEmojiStyleNone
}

func MapGqlKappagenAnimationStyleToEntity(style gqlmodel.KappagenOverlayAnimationStyle) (
	entity.KappagenOverlayAnimationStyle,
	error,
) {
	for k, v := range mapKappagenAnimationStyleToGql {
		if v == style {
			return k, nil
		}
	}

	return entity.KappagenOverlayAnimationStyleConfetti, fmt.Errorf(
		"unknown kappagen animation style: %s",
		style,
	)
}
