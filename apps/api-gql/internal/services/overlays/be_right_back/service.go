package be_right_back

import (
	"context"
	"errors"
	"log/slog"

	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	"github.com/twirapp/twir/apps/api-gql/internal/wsrouter"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/api"
	"github.com/twirapp/twir/libs/repositories/overlays_be_right_back"
	"github.com/twirapp/twir/libs/repositories/overlays_be_right_back/model"
	"github.com/twirapp/twir/libs/repositories/users"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In
	LC fx.Lifecycle

	Repository      overlays_be_right_back.Repository
	WsRouter        wsrouter.WsRouter
	TwirBus         *buscore.Bus
	Logger          *slog.Logger
	UsersRepository users.Repository
}

func New(opts Opts) *Service {
	s := &Service{
		repository:      opts.Repository,
		wsRouter:        opts.WsRouter,
		twirBus:         opts.TwirBus,
		usersRepository: opts.UsersRepository,
	}

	opts.LC.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				s.twirBus.Api.TriggerBrbStart.SubscribeGroup(
					"api",
					func(ctx context.Context, data api.TriggerBrbStart) (struct{}, error) {
						return struct{}{}, s.wsRouter.Publish(createStartSubscriptionKey(data.ChannelId), data)
					},
				)

				opts.Logger.Info("Subscribed to TriggerBrbStart events")

				s.twirBus.Api.TriggerBrbStop.SubscribeGroup(
					"api",
					func(ctx context.Context, data api.TriggerBrbStop) (struct{}, error) {
						return struct{}{}, s.wsRouter.Publish(createStopSubscriptionKey(data.ChannelId), data)
					},
				)

				opts.Logger.Info("Subscribed to TriggerBrbStop events")

				return nil
			},
			OnStop: func(ctx context.Context) error {
				s.twirBus.Api.TriggerBrbStart.Unsubscribe()
				s.twirBus.Api.TriggerBrbStop.Unsubscribe()

				opts.Logger.Info("Unsubscribed from TriggerBrbStart and TriggerBrbStop events")

				return nil
			},
		},
	)

	return s
}

type Service struct {
	repository overlays_be_right_back.Repository
	wsRouter   wsrouter.WsRouter
	twirBus    *buscore.Bus

	usersRepository users.Repository
}

// GetOrCreate gets the be right back overlay for the given channel ID or creates a new one with default settings if it doesn't exist
func (s *Service) GetOrCreate(ctx context.Context, channelID string) (
	entity.BeRightBackOverlay,
	error,
) {
	overlay, err := s.repository.GetByChannelID(ctx, channelID)
	if err != nil {
		if errors.Is(err, overlays_be_right_back.ErrNotFound) {
			// Create a new overlay with default settings
			created, err := s.repository.Create(ctx, createDefaultOverlayInput(channelID))
			if err != nil {
				return entity.BeRightBackOverlay{}, err
			}
			return mapModelToEntity(created), nil
		}
		return entity.BeRightBackOverlay{}, err
	}

	return mapModelToEntity(overlay), nil
}

type UpdateInput struct {
	ChannelID string
	Settings  entity.BeRightBackOverlaySettings
}

// Update updates the be right back overlay for the given channel ID
func (s *Service) Update(
	ctx context.Context,
	input UpdateInput,
) (entity.BeRightBackOverlay, error) {
	updated, err := s.repository.Update(
		ctx, input.ChannelID, overlays_be_right_back.UpdateInput{
			Settings: mapSettingsEntityToModel(input.Settings),
		},
	)
	if err != nil {
		return entity.BeRightBackOverlay{}, err
	}

	if err := s.wsRouter.Publish(
		createSettingsSubscriptionKey(input.ChannelID),
		mapModelToEntity(updated),
	); err != nil {
		return entity.BeRightBackOverlay{}, err
	}

	return mapModelToEntity(updated), nil
}

// Mappers between repository model and entity
func mapModelToEntity(m model.BeRightBackOverlay) entity.BeRightBackOverlay {
	return entity.BeRightBackOverlay{
		ID:        m.ID,
		ChannelID: m.ChannelID,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		Settings: entity.BeRightBackOverlaySettings{
			Text: m.Settings.Text,
			Late: entity.BeRightBackOverlayLateSettings{
				Enabled:        m.Settings.Late.Enabled,
				Text:           m.Settings.Late.Text,
				DisplayBrbTime: m.Settings.Late.DisplayBrbTime,
			},
			BackgroundColor: m.Settings.BackgroundColor,
			FontSize:        m.Settings.FontSize,
			FontColor:       m.Settings.FontColor,
			FontFamily:      m.Settings.FontFamily,
		},
	}
}

func mapSettingsEntityToModel(e entity.BeRightBackOverlaySettings) model.BeRightBackOverlaySettings {
	return model.BeRightBackOverlaySettings{
		Text: e.Text,
		Late: model.BeRightBackOverlayLateSettings{
			Enabled:        e.Late.Enabled,
			Text:           e.Late.Text,
			DisplayBrbTime: e.Late.DisplayBrbTime,
		},
		BackgroundColor: e.BackgroundColor,
		FontSize:        e.FontSize,
		FontColor:       e.FontColor,
		FontFamily:      e.FontFamily,
	}
}

func createDefaultOverlayInput(channelID string) overlays_be_right_back.CreateInput {
	return overlays_be_right_back.CreateInput{
		ChannelID: channelID,
		Settings: model.BeRightBackOverlaySettings{
			Text: "BRB",
			Late: model.BeRightBackOverlayLateSettings{
				Enabled:        false,
				Text:           "Streamer is late",
				DisplayBrbTime: false,
			},
			BackgroundColor: "#000000",
			FontSize:        48,
			FontColor:       "#FFFFFF",
			FontFamily:      "Roboto",
		},
	}
}
