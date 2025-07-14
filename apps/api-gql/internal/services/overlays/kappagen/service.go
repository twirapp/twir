package kappagen

import (
	"context"
	"errors"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	"github.com/twirapp/twir/apps/api-gql/internal/wsrouter"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/api"
	eventmodel "github.com/twirapp/twir/libs/repositories/events/model"
	"github.com/twirapp/twir/libs/repositories/overlays_kappagen"
	"github.com/twirapp/twir/libs/repositories/overlays_kappagen/model"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In
	LC fx.Lifecycle

	Repository overlays_kappagen.Repository
	WsRouter   wsrouter.WsRouter
	TwirBus    *buscore.Bus
}

func New(opts Opts) *Service {
	s := &Service{
		repository: opts.Repository,
		wsRouter:   opts.WsRouter,
		twirBus:    opts.TwirBus,
	}

	opts.LC.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				opts.TwirBus.Api.TriggerKappagen.SubscribeGroup(
					"api",
					s.handleTriggerKappagenEvent,
				)

				return nil
			},
			OnStop: func(ctx context.Context) error {
				opts.TwirBus.Api.TriggerKappagen.Unsubscribe()
				return nil
			},
		},
	)

	return s
}

type Service struct {
	repository overlays_kappagen.Repository
	wsRouter   wsrouter.WsRouter
	twirBus    *buscore.Bus
}

// GetOrCreate gets the kappagen overlay for the given channel ID or creates a new one with default settings if it doesn't exist
func (s *Service) GetOrCreate(ctx context.Context, channelID string) (
	entity.KappagenOverlay,
	error,
) {
	overlay, err := s.repository.GetByChannelID(ctx, channelID)
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
	Settings  entity.KappagenOverlaySettings
}

// Update updates the kappagen overlay for the given channel ID
func (s *Service) Update(
	ctx context.Context,
	input UpdateInput,
) (entity.KappagenOverlay, error) {
	defaultInput := createDefaultOverlayInput(input.ChannelID)
	if len(input.Settings.Animations) != len(defaultInput.Settings.Animations) {
		return entity.KappagenOverlay{}, errors.New("invalid number of animations")
	}

	if len(input.Settings.Events) != len(defaultInput.Settings.Events) {
		return entity.KappagenOverlay{}, errors.New("invalid number of events")
	}

	updated, err := s.repository.Update(
		ctx, input.ChannelID, overlays_kappagen.UpdateInput{
			Settings: mapSettingsEntityToModel(input.Settings),
		},
	)
	if err != nil {
		return entity.KappagenOverlay{}, err
	}

	if err := s.wsRouter.Publish(
		CreateSettingsSubscriptionKey(input.ChannelID),
		mapModelToEntity(updated),
	); err != nil {
		return entity.KappagenOverlay{}, err
	}

	return mapModelToEntity(updated), nil
}

// Mappers between repository model and entity
func mapModelToEntity(m model.KappagenOverlay) entity.KappagenOverlay {
	animations := make([]entity.KappagenOverlayAnimationsSettings, 0, len(m.Settings.Animations))
	for _, a := range m.Settings.Animations {
		var prefs *entity.KappagenOverlayAnimationsPrefsSettings = nil
		if a.Prefs != nil {
			prefs = &entity.KappagenOverlayAnimationsPrefsSettings{
				Size:    a.Prefs.Size,
				Center:  a.Prefs.Center,
				Speed:   a.Prefs.Speed,
				Faces:   a.Prefs.Faces,
				Message: a.Prefs.Message,
				Time:    a.Prefs.Time,
			}
		}

		animations = append(
			animations, entity.KappagenOverlayAnimationsSettings{
				Style:   entity.KappagenOverlayAnimationStyle(a.Style),
				Prefs:   prefs,
				Count:   a.Count,
				Enabled: a.Enabled,
			},
		)
	}

	events := make([]entity.KappagenOverlayEvent, 0, len(m.Settings.Events))
	for _, e := range m.Settings.Events {
		events = append(
			events, entity.KappagenOverlayEvent{
				Event:              entity.EventType(e.Event),
				DisabledAnimations: e.DisabledAnimations,
				Enabled:            e.Enabled,
			},
		)
	}

	return entity.KappagenOverlay{
		ID:        m.ID,
		ChannelID: m.ChannelID,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		Settings: entity.KappagenOverlaySettings{
			EnableSpawn:    m.Settings.EnableSpawn,
			ExcludedEmotes: m.Settings.ExcludedEmotes,
			EnableRave:     m.Settings.EnableRave,
			Animation: entity.KappagenOverlayAnimationSettings{
				FadeIn:  m.Settings.Animation.FadeIn,
				FadeOut: m.Settings.Animation.FadeOut,
				ZoomIn:  m.Settings.Animation.ZoomIn,
				ZoomOut: m.Settings.Animation.ZoomOut,
			},
			Animations: animations,
			Emotes: entity.KappagenOverlayEmotesSettings{
				Time:           m.Settings.Emotes.Time,
				Max:            m.Settings.Emotes.Max,
				Queue:          m.Settings.Emotes.Queue,
				FfzEnabled:     m.Settings.Emotes.FfzEnabled,
				BttvEnabled:    m.Settings.Emotes.BttvEnabled,
				SevenTvEnabled: m.Settings.Emotes.SevenTvEnabled,
				EmojiStyle:     entity.KappagenEmojiStyle(m.Settings.Emotes.EmojiStyle),
			},
			Size: entity.KappagenOverlaySizeSettings{
				RatioNormal: m.Settings.Size.RatioNormal,
				RatioSmall:  m.Settings.Size.RatioSmall,
				Min:         m.Settings.Size.Min,
				Max:         m.Settings.Size.Max,
			},
			Events: events,
		},
	}
}

func mapSettingsEntityToModel(e entity.KappagenOverlaySettings) model.KappagenOverlaySettings {
	animations := make([]model.KappagenOverlayAnimationsSettings, 0, len(e.Animations))
	for _, a := range e.Animations {
		var prefs *model.KappagenOverlayAnimationsPrefsSettings
		if a.Prefs != nil {
			prefs = &model.KappagenOverlayAnimationsPrefsSettings{
				Size:    a.Prefs.Size,
				Center:  a.Prefs.Center,
				Speed:   a.Prefs.Speed,
				Faces:   a.Prefs.Faces,
				Message: a.Prefs.Message,
				Time:    a.Prefs.Time,
			}
		}

		animations = append(
			animations,
			model.KappagenOverlayAnimationsSettings{
				Style:   string(a.Style),
				Prefs:   prefs,
				Count:   a.Count,
				Enabled: a.Enabled,
			},
		)
	}

	events := make([]model.KappagenOverlayEvent, 0, len(e.Events))
	for _, event := range e.Events {
		events = append(
			events, model.KappagenOverlayEvent{
				Event:              eventmodel.EventType(event.Event),
				DisabledAnimations: event.DisabledAnimations,
				Enabled:            event.Enabled,
			},
		)
	}

	return model.KappagenOverlaySettings{
		EnableSpawn:    e.EnableSpawn,
		ExcludedEmotes: e.ExcludedEmotes,
		EnableRave:     e.EnableRave,
		Animation: model.KappagenOverlayAnimationSettings{
			FadeIn:  e.Animation.FadeIn,
			FadeOut: e.Animation.FadeOut,
			ZoomIn:  e.Animation.ZoomIn,
			ZoomOut: e.Animation.ZoomOut,
		},
		Animations: animations,
		Emotes: model.KappagenOverlayEmotesSettings{
			Time:           e.Emotes.Time,
			Max:            e.Emotes.Max,
			Queue:          e.Emotes.Queue,
			FfzEnabled:     e.Emotes.FfzEnabled,
			BttvEnabled:    e.Emotes.BttvEnabled,
			SevenTvEnabled: e.Emotes.SevenTvEnabled,
			EmojiStyle:     model.KappagenEmojiStyle(e.Emotes.EmojiStyle),
		},
		Size: model.KappagenOverlaySizeSettings{
			RatioNormal: e.Size.RatioNormal,
			RatioSmall:  e.Size.RatioSmall,
			Min:         e.Size.Min,
			Max:         e.Size.Max,
		},
		Events: events,
	}
}

var (
	oneHundredFifty = 150
	fifty           = 50
)

var defaultAnimations = []model.KappagenOverlayAnimationsSettings{
	{
		Style: "TheCube",
		Prefs: &model.KappagenOverlayAnimationsPrefsSettings{
			Size:    lo.ToPtr(0.2),
			Center:  lo.ToPtr(false),
			Faces:   lo.ToPtr(false),
			Speed:   lo.ToPtr(6),
			Message: []string{},
		},
		Enabled: true,
	},
	{
		Style: "Text",
		Prefs: &model.KappagenOverlayAnimationsPrefsSettings{
			Message: []string{"Twir"},
			Time:    lo.ToPtr(3),
		},
		Enabled: true,
	},
	{
		Style:   "Confetti",
		Count:   &oneHundredFifty,
		Enabled: true,
	},
	{
		Style:   "Spiral",
		Count:   &oneHundredFifty,
		Enabled: true,
	},
	{
		Style:   "Stampede",
		Count:   &oneHundredFifty,
		Enabled: true,
	},
	{
		Style:   "Burst",
		Count:   &fifty,
		Enabled: true,
	},
	{
		Style:   "Fountain",
		Count:   &fifty,
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
		Count:   &oneHundredFifty,
		Enabled: true,
	},
	{
		Style: "Conga",
		Prefs: &model.KappagenOverlayAnimationsPrefsSettings{
			Message: []string{},
		},
		Enabled: true,
	},
}
var defaultEvents = []model.KappagenOverlayEvent{
	{
		Event:   eventmodel.EventTypeFollow,
		Enabled: true,
	},
	{
		Event:   eventmodel.EventTypeSubscribe,
		Enabled: true,
	},
	{
		Event:   eventmodel.EventTypeResubscribe,
		Enabled: true,
	},
	{
		Event:   eventmodel.EventTypeSubGift,
		Enabled: true,
	},
	{
		Event:   eventmodel.EventTypeRedemptionCreated,
		Enabled: true,
	},
	{
		Event:   eventmodel.EventTypeFirstUserMessage,
		Enabled: true,
	},
	{
		Event:   eventmodel.EventTypeRaided,
		Enabled: true,
	},
	{
		Event:   eventmodel.EventTypeTitleOrCategoryChanged,
		Enabled: true,
	},
	{
		Event:   eventmodel.EventTypeStreamOnline,
		Enabled: true,
	},
	{
		Event:   eventmodel.EventTypeStreamOffline,
		Enabled: true,
	},
	{
		Event:   eventmodel.EventTypeOnChatClear,
		Enabled: true,
	},
	{
		Event:   eventmodel.EventTypeDonate,
		Enabled: true,
	},
	{
		Event:   eventmodel.EventTypeKeywordMatched,
		Enabled: true,
	},
	{
		Event:   eventmodel.EventTypeGreetingSended,
		Enabled: true,
	},
	{
		Event:   eventmodel.EventTypePollBegin,
		Enabled: true,
	},
	{
		Event:   eventmodel.EventTypePollProgress,
		Enabled: true,
	},
	{
		Event:   eventmodel.EventTypePollEnd,
		Enabled: true,
	},
	{
		Event:   eventmodel.EventTypePredictionBegin,
		Enabled: true,
	},
	{
		Event:   eventmodel.EventTypePredictionProgress,
		Enabled: true,
	},
	{
		Event:   eventmodel.EventTypePredictionEnd,
		Enabled: true,
	},
	{
		Event:   eventmodel.EventTypePredictionLock,
		Enabled: true,
	},
	{
		Event:   eventmodel.EventTypeChannelBan,
		Enabled: true,
	},
	{
		Event:   eventmodel.EventTypeChannelUnbanRequestCreate,
		Enabled: true,
	},
	{
		Event:   eventmodel.EventTypeChannelUnbanRequestResolve,
		Enabled: true,
	},
	{
		Event:   eventmodel.EventTypeChannelMessageDelete,
		Enabled: true,
	},
	{
		Event:   eventmodel.EventTypeCommandUsed,
		Enabled: true,
	},
}

// createDefaultOverlayInput creates input for default kappagen overlay
func createDefaultOverlayInput(channelID string) overlays_kappagen.CreateInput {
	return overlays_kappagen.CreateInput{
		ChannelID: channelID,
		Settings: model.KappagenOverlaySettings{
			EnableSpawn:    true,
			ExcludedEmotes: []string{},
			EnableRave:     false,
			Animation: model.KappagenOverlayAnimationSettings{
				FadeIn:  true,
				FadeOut: true,
				ZoomIn:  true,
				ZoomOut: true,
			},
			Animations: defaultAnimations,
			Emotes: model.KappagenOverlayEmotesSettings{
				Time:           5,
				Max:            100,
				Queue:          100,
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
			Events: defaultEvents,
		},
	}
}

func CreateSettingsSubscriptionKey(channelId string) string {
	return "api.kappagensettings." + channelId
}

func CreateTriggerSubscriptionKey(channelId string) string {
	return "api.kappagen.trigger." + channelId
}

func (s *Service) handleTriggerKappagenEvent(
	ctx context.Context,
	msg api.TriggerKappagenMessage,
) (struct{}, error) {
	return struct{}{}, s.wsRouter.Publish(CreateTriggerSubscriptionKey(msg.ChannelId), msg)
}

func (s *Service) GetAvailableAnimations() []string {
	animations := make([]string, len(defaultAnimations))
	for i, a := range defaultAnimations {
		animations[i] = a.Style
	}

	return animations
}
