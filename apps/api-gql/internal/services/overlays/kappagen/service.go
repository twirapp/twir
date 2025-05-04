package kappagen

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/satont/twir/apps/api-gql/internal/entity"
	"github.com/twirapp/twir/libs/repositories/overlays_kappagen"
	"github.com/twirapp/twir/libs/repositories/overlays_kappagen/model"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	Repository overlays_kappagen.Repository
}

func New(opts Opts) *Service {
	return &Service{
		repository: opts.Repository,
	}
}

type Service struct {
	repository overlays_kappagen.Repository
}

// GetOrCreate gets the kappagen overlay for the given channel ID or creates a new one with default settings if it doesn't exist
func (s *Service) GetOrCreate(ctx context.Context, channelID string) (entity.KappagenOverlay, error) {
	overlay, err := s.repository.Get(ctx, channelID)
	if err != nil {
		if errors.Is(err, overlays_kappagen.ErrNotFound) {
			// Create a new overlay with default settings
			created, err := s.repository.Create(ctx, createDefaultOverlayInput(channelID))
			if err != nil {
				return entity.KappagenOverlay{}, err
			}
			return mapModelToEntity(created), nil
		}
		return entity.KappagenOverlay{}, err
	}

	return mapModelToEntity(overlay), nil
}

// Update updates the kappagen overlay for the given channel ID
func (s *Service) Update(ctx context.Context, channelID string, input entity.KappagenOverlay) (entity.KappagenOverlay, error) {
	// Try to get the overlay first
	_, err := s.repository.Get(ctx, channelID)
	if err != nil {
		if errors.Is(err, overlays_kappagen.ErrNotFound) {
			// Create a new overlay with the provided settings
			createInput := createDefaultOverlayInput(channelID)
			// Override default values with provided input
			createInput.EnableSpawn = input.EnableSpawn
			createInput.ExcludedEmotes = input.ExcludedEmotes
			createInput.EnableRave = input.EnableRave
			createInput.Animation = mapEntityAnimationToModel(input.Animation)
			createInput.Animations = mapEntityAnimationsToModel(input.Animations)

			created, err := s.repository.Create(ctx, createInput)
			if err != nil {
				return entity.KappagenOverlay{}, err
			}
			return mapModelToEntity(created), nil
		}
		return entity.KappagenOverlay{}, err
	}

	// Update the existing overlay
	updateInput := overlays_kappagen.UpdateInput{
		EnableSpawn:    input.EnableSpawn,
		ExcludedEmotes: input.ExcludedEmotes,
		EnableRave:     input.EnableRave,
		Animation:      mapEntityAnimationToModel(input.Animation),
		Animations:     mapEntityAnimationsToModel(input.Animations),
	}

	updated, err := s.repository.Update(ctx, channelID, updateInput)
	if err != nil {
		return entity.KappagenOverlay{}, err
	}

	return mapModelToEntity(updated), nil
}

// Mappers between repository model and entity
func mapModelToEntity(model model.KappagenOverlay) entity.KappagenOverlay {
	animations := make([]entity.KappagenOverlayAnimationsSettings, 0, len(model.Animations))
	for _, a := range model.Animations {
		var prefs entity.KappagenOverlayAnimationsPrefsSettings
		if a.Prefs != nil {
			var size float64
			if a.Prefs.Size != nil {
				size = *a.Prefs.Size
			}

			var center float64
			if a.Prefs.Center != nil && *a.Prefs.Center {
				center = 1
			}

			var speed int
			if a.Prefs.Speed != nil {
				speed = int(*a.Prefs.Speed)
			}

			var faces bool
			if a.Prefs.Faces != nil {
				faces = *a.Prefs.Faces
			}

			var time int
			if a.Prefs.Time != nil {
				time = int(*a.Prefs.Time)
			}

			prefs = entity.KappagenOverlayAnimationsPrefsSettings{
				Size:    size,
				Center:  center,
				Speed:   speed,
				Faces:   faces,
				Message: a.Prefs.Message,
				Time:    time,
			}
		}

		var count int
		if a.Count != nil {
			count = int(*a.Count)
		}

		animations = append(animations, entity.KappagenOverlayAnimationsSettings{
			Style:   a.Style,
			Prefs:   prefs,
			Count:   count,
			Enabled: a.Enabled,
		})
	}

	return entity.KappagenOverlay{
		ID:             model.ID,
		EnableSpawn:    model.EnableSpawn,
		ExcludedEmotes: model.ExcludedEmotes,
		EnableRave:     model.EnableRave,
		Animation: entity.KappagenOverlayAnimationSettings{
			FadeIn:  model.Animation.FadeIn,
			FadeOut: model.Animation.FadeOut,
			ZoomIn:  model.Animation.ZoomIn,
			ZoomOut: model.Animation.ZoomOut,
		},
		Animations: animations,
	}
}

func mapEntityAnimationToModel(animation entity.KappagenOverlayAnimationSettings) model.KappagenOverlayAnimationSettings {
	return model.KappagenOverlayAnimationSettings{
		FadeIn:  animation.FadeIn,
		FadeOut: animation.FadeOut,
		ZoomIn:  animation.ZoomIn,
		ZoomOut: animation.ZoomOut,
	}
}

func mapEntityAnimationsToModel(animations []entity.KappagenOverlayAnimationsSettings) []model.KappagenOverlayAnimationsSettings {
	result := make([]model.KappagenOverlayAnimationsSettings, 0, len(animations))
	for _, a := range animations {
		size := float64(a.Prefs.Size)
		center := a.Prefs.Center > 0
		speed := int32(a.Prefs.Speed)
		faces := a.Prefs.Faces
		time := int32(a.Prefs.Time)
		count := int32(a.Count)

		result = append(result, model.KappagenOverlayAnimationsSettings{
			Style: a.Style,
			Prefs: &model.KappagenOverlayAnimationsPrefsSettings{
				Size:    &size,
				Center:  &center,
				Speed:   &speed,
				Faces:   &faces,
				Message: a.Prefs.Message,
				Time:    &time,
			},
			Count:   &count,
			Enabled: a.Enabled,
		})
	}
	return result
}

// createDefaultOverlayInput creates input for default kappagen overlay
func createDefaultOverlayInput(channelID string) overlays_kappagen.CreateInput {
	size := float64(1)
	speedInt32 := int32(5)
	timeInt32 := int32(5000)
	countInt32 := int32(10)
	falseValue := false
	trueValue := true
	centerFalse := false
	centerTrue := true

	return overlays_kappagen.CreateInput{
		ChannelID:      channelID,
		EnableSpawn:    true,
		ExcludedEmotes: []string{},
		EnableRave:     false,
		Animation: model.KappagenOverlayAnimationSettings{
			FadeIn:  true,
			FadeOut: true,
			ZoomIn:  false,
			ZoomOut: false,
		},
		Animations: []model.KappagenOverlayAnimationsSettings{
			{
				Style: "rain",
				Prefs: &model.KappagenOverlayAnimationsPrefsSettings{
					Size:    &size,
					Center:  &centerFalse,
					Speed:   &speedInt32,
					Faces:   &falseValue,
					Message: []string{},
					Time:    &timeInt32,
				},
				Count:   &countInt32,
				Enabled: true,
			},
			{
				Style: "explode",
				Prefs: &model.KappagenOverlayAnimationsPrefsSettings{
					Size:    &size,
					Center:  &centerTrue,
					Speed:   &speedInt32,
					Faces:   &falseValue,
					Message: []string{},
					Time:    &timeInt32,
				},
				Count:   &countInt32,
				Enabled: true,
			},
		},
		Emotes: model.KappagenOverlayEmotesSettings{
			Time:           5000,
			Max:            50,
			Queue:          10,
			FfzEnabled:     true,
			BttvEnabled:    true,
			SevenTvEnabled: true,
			EmojiStyle:     model.KappagenEmojiStyleTwemoji,
		},
		Size: model.KappagenOverlaySizeSettings{
			RatioNormal: 7,
			RatioSmall:  14,
			Min:         30,
			Max:         100,
		},
		Cube: model.KappagenOverlayCubeSettings{
			Speed: 5,
		},
	}
}
