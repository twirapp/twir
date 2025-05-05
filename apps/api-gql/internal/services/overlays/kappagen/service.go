package kappagen

import (
	"context"
	"errors"

	"github.com/twirapp/twir/apps/api-gql/internal/entity"
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
func (s *Service) GetOrCreate(ctx context.Context, channelID string) (
	entity.KappagenOverlay,
	error,
) {
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

type UpdateInput struct {
	ChannelID string

	EnableSpawn    bool
	ExcludedEmotes []string
	EnableRave     bool
	Animation      UpdateInputAnimationSettings
	Animations     []entity.KappagenOverlayAnimationsSettings
}
type UpdateInputAnimationSettings struct {
	FadeIn  bool
	FadeOut bool
	ZoomIn  bool
	ZoomOut bool
}

type UpdateInputAnimationsSettings struct {
	Style   string
	Prefs   UpdateInputAnimationsPrefsSettings
	Count   int
	Enabled bool
}

type UpdateInputAnimationsPrefsSettings struct {
	Size    float64
	Center  bool
	Speed   int
	Faces   bool
	Message []string
	Time    int
}

// Update updates the kappagen overlay for the given channel ID
func (s *Service) Update(
	ctx context.Context,
	input UpdateInput,
) (entity.KappagenOverlay, error) {
	animations := make([]model.KappagenOverlayAnimationsSettings, 0, len(input.Animations))
	for _, a := range input.Animations {
		animations = append(
			animations, model.KappagenOverlayAnimationsSettings{
				Style: a.Style,
				Prefs: model.KappagenOverlayAnimationsPrefsSettings{
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

	// Update the existing overlay
	updateInput := overlays_kappagen.UpdateInput{
		EnableSpawn:    input.EnableSpawn,
		ExcludedEmotes: input.ExcludedEmotes,
		EnableRave:     input.EnableRave,
		Animation: model.KappagenOverlayAnimationSettings{
			FadeIn:  input.Animation.FadeIn,
			FadeOut: input.Animation.FadeOut,
			ZoomIn:  input.Animation.ZoomIn,
			ZoomOut: input.Animation.ZoomOut,
		},
		Animations: animations,
	}

	updated, err := s.repository.Update(ctx, input.ChannelID, updateInput)
	if err != nil {
		return entity.KappagenOverlay{}, err
	}

	return mapModelToEntity(updated), nil
}

// Mappers between repository model and entity
func mapModelToEntity(model model.KappagenOverlay) entity.KappagenOverlay {
	animations := make([]entity.KappagenOverlayAnimationsSettings, 0, len(model.Animations))
	for _, a := range model.Animations {
		prefs := entity.KappagenOverlayAnimationsPrefsSettings{
			Size:    a.Prefs.Size,
			Center:  a.Prefs.Center,
			Speed:   a.Prefs.Speed,
			Faces:   a.Prefs.Faces,
			Message: a.Prefs.Message,
			Time:    a.Prefs.Time,
		}

		animations = append(
			animations, entity.KappagenOverlayAnimationsSettings{
				Style:   a.Style,
				Prefs:   prefs,
				Count:   a.Count,
				Enabled: a.Enabled,
			},
		)
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

// createDefaultOverlayInput creates input for default kappagen overlay
func createDefaultOverlayInput(channelID string) overlays_kappagen.CreateInput {
	return overlays_kappagen.CreateInput{
		ChannelID:      channelID,
		EnableSpawn:    true,
		ExcludedEmotes: []string{},
		EnableRave:     false,
		Animation: model.KappagenOverlayAnimationSettings{
			FadeIn:  true,
			FadeOut: true,
			ZoomIn:  true,
			ZoomOut: true,
		},
		Animations: []model.KappagenOverlayAnimationsSettings{
			{
				Style: "TheCube",
				Prefs: model.KappagenOverlayAnimationsPrefsSettings{
					Size:    0.2,
					Center:  false,
					Faces:   false,
					Speed:   6,
					Message: []string{},
				},
				Enabled: true,
			},
			{
				Style: "Text",
				Prefs: model.KappagenOverlayAnimationsPrefsSettings{
					Message: []string{"Twir"},
					Time:    3,
				},
				Enabled: true,
			},
			{
				Style:   "Confetti",
				Count:   150,
				Enabled: true,
			},
			{
				Style:   "Spiral",
				Count:   150,
				Enabled: true,
			},
			{
				Style:   "Stampede",
				Count:   150,
				Enabled: true,
			},
			{
				Style:   "Burst",
				Count:   50,
				Enabled: true,
			},
			{
				Style:   "Fountain",
				Count:   50,
				Enabled: true,
			},
			{
				Style:   "SmallPyramid",
				Enabled: true,
			},
			{
				Style:   "Pyramid",
				Enabled: true,
			},
			{
				Style:   "Fireworks",
				Count:   150,
				Enabled: true,
			},
			{
				Style: "Conga",
				Prefs: model.KappagenOverlayAnimationsPrefsSettings{
					Message: []string{},
				},
				Enabled: true,
			},
		},
		Emotes: model.KappagenOverlayEmotesSettings{
			Time:           5,
			Max:            0,
			Queue:          0,
			FfzEnabled:     true,
			BttvEnabled:    true,
			SevenTvEnabled: true,
			EmojiStyle:     model.KappagenEmojiStyleTwemoji,
		},
		Size: model.KappagenOverlaySizeSettings{
			RatioNormal: 0.05,
			RatioSmall:  0.02,
			Min:         1,
			Max:         256,
		},
		Cube: model.KappagenOverlayCubeSettings{
			Speed: 6,
		},
	}
}
