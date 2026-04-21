package channels_overlays

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/mappers"
	"github.com/twirapp/twir/libs/audit"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/parser"
	customoverlayentity "github.com/twirapp/twir/libs/entities/custom_overlay"
	apperrors "github.com/twirapp/twir/libs/errors"
	"github.com/twirapp/twir/libs/repositories/channels"
	channels_overlays "github.com/twirapp/twir/libs/repositories/channels_overlays"
	"github.com/twirapp/twir/libs/repositories/channels_overlays/model"
	"github.com/twirapp/twir/libs/repositories/plans"
	"github.com/twirapp/twir/libs/wsrouter"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	OverlaysRepository channels_overlays.Repository
	ChannelsRepository channels.Repository
	PlansRepository    plans.Repository
	AuditRecorder      audit.Recorder
	Bus                *buscore.Bus
	WsRouter           wsrouter.WsRouter
}

func New(opts Opts) *Service {
	return &Service{
		overlaysRepository: opts.OverlaysRepository,
		channelsRepository: opts.ChannelsRepository,
		plansRepository:    opts.PlansRepository,
		auditRecorder:      opts.AuditRecorder,
		bus:                opts.Bus,
		wsRouter:           opts.WsRouter,
	}
}

type Service struct {
	overlaysRepository channels_overlays.Repository
	channelsRepository channels.Repository
	plansRepository    plans.Repository
	auditRecorder      audit.Recorder
	bus                *buscore.Bus
	wsRouter           wsrouter.WsRouter
}

func (s *Service) modelToEntity(m model.Overlay) customoverlayentity.ChannelOverlay {
	layers := make([]customoverlayentity.ChannelOverlayLayer, len(m.Layers))
	for i, l := range m.Layers {
		layers[i] = customoverlayentity.ChannelOverlayLayer{
			ID:   l.ID,
			Type: customoverlayentity.ChannelOverlayType(l.Type),
			Settings: customoverlayentity.ChannelOverlayLayerSettings{
				HtmlOverlayHTML:                    l.Settings.HtmlOverlayHTML,
				HtmlOverlayCSS:                     l.Settings.HtmlOverlayCSS,
				HtmlOverlayJS:                      l.Settings.HtmlOverlayJS,
				HtmlOverlayDataPollSecondsInterval: l.Settings.HtmlOverlayDataPollSecondsInterval,
				ImageUrl:                           l.Settings.ImageUrl,
			},
			OverlayID:               l.OverlayID,
			PosX:                    l.PosX,
			PosY:                    l.PosY,
			Width:                   l.Width,
			Height:                  l.Height,
			Rotation:                l.Rotation,
			CreatedAt:               l.CreatedAt,
			UpdatedAt:               l.UpdatedAt,
			PeriodicallyRefetchData: l.PeriodicallyRefetchData,
			Locked:                  l.Locked,
			Visible:                 l.Visible,
			Opacity:                 l.Opacity,
		}
	}

	return customoverlayentity.ChannelOverlay{
		ID:        m.ID,
		ChannelID: m.ChannelID,
		Name:      m.Name,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		Width:     m.Width,
		Height:    m.Height,
		InstaSave: m.InstaSave,
		Layers:    layers,
	}
}

func (s *Service) GetManyByChannelID(ctx context.Context, channelID string) (
	[]customoverlayentity.ChannelOverlay,
	error,
) {
	overlays, err := s.overlaysRepository.GetManyByChannelID(ctx, channelID)
	if err != nil {
		return nil, err
	}

	entities := make([]customoverlayentity.ChannelOverlay, len(overlays))
	for i, o := range overlays {
		entities[i] = s.modelToEntity(o)
	}

	return entities, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (customoverlayentity.ChannelOverlay, error) {
	overlay, err := s.overlaysRepository.GetByID(ctx, id)
	if err != nil {
		return customoverlayentity.ChannelOverlayNil, err
	}

	return s.modelToEntity(overlay), nil
}

type CreateLayerInput struct {
	Type                    customoverlayentity.ChannelOverlayType
	Settings                customoverlayentity.ChannelOverlayLayerSettings
	PosX                    int
	PosY                    int
	Width                   int
	Height                  int
	Rotation                int
	PeriodicallyRefetchData bool
	Locked                  bool
	Visible                 bool
	Opacity                 float64
}

type UpdateLayerInput struct {
	ID                      *uuid.UUID
	Type                    customoverlayentity.ChannelOverlayType
	Settings                customoverlayentity.ChannelOverlayLayerSettings
	PosX                    int
	PosY                    int
	Width                   int
	Height                  int
	Rotation                int
	PeriodicallyRefetchData bool
	Locked                  bool
	Visible                 bool
	Opacity                 float64
}

type CreateInput struct {
	ChannelID string
	ActorID   string
	Name      string
	Width     int
	Height    int
	InstaSave bool
	Layers    []CreateLayerInput
}

func (s *Service) Create(ctx context.Context, input CreateInput) (customoverlayentity.ChannelOverlay, error) {
	plan, err := s.plansRepository.GetByChannelID(ctx, input.ChannelID)
	if err != nil {
		return customoverlayentity.ChannelOverlayNil, apperrors.NewInternalError("Failed to retrieve plan information", err)
	}
	if plan.IsNil() {
		return customoverlayentity.ChannelOverlayNil, apperrors.NewNotFoundError("Plan not found for this channel")
	}

	existingOverlays, err := s.overlaysRepository.GetManyByChannelID(ctx, input.ChannelID)
	if err != nil {
		return customoverlayentity.ChannelOverlayNil, apperrors.NewInternalError("Failed to retrieve existing overlays", err)
	}

	if len(existingOverlays) >= plan.MaxCustomOverlays {
		return customoverlayentity.ChannelOverlayNil, apperrors.NewBadRequestError(
			fmt.Sprintf("You have reached the maximum limit of %d custom overlays for your plan", plan.MaxCustomOverlays),
		)
	}

	repoLayers := make([]channels_overlays.CreateLayerInput, len(input.Layers))
	for i, l := range input.Layers {
		repoLayers[i] = channels_overlays.CreateLayerInput{
			Type: model.OverlayType(l.Type),
			Settings: model.OverlayLayerSettings{
				HtmlOverlayHTML:                    l.Settings.HtmlOverlayHTML,
				HtmlOverlayCSS:                     l.Settings.HtmlOverlayCSS,
				HtmlOverlayJS:                      l.Settings.HtmlOverlayJS,
				HtmlOverlayDataPollSecondsInterval: l.Settings.HtmlOverlayDataPollSecondsInterval,
				ImageUrl:                           l.Settings.ImageUrl,
			},
			PosX:                    l.PosX,
			PosY:                    l.PosY,
			Width:                   l.Width,
			Height:                  l.Height,
			Rotation:                l.Rotation,
			PeriodicallyRefetchData: l.PeriodicallyRefetchData,
			Locked:                  l.Locked,
			Visible:                 l.Visible,
			Opacity:                 l.Opacity,
		}
	}

	overlay, err := s.overlaysRepository.Create(
		ctx,
		channels_overlays.CreateInput{
			ChannelID: input.ChannelID,
			Name:      input.Name,
			Width:     input.Width,
			Height:    input.Height,
			InstaSave: input.InstaSave,
			Layers:    repoLayers,
		},
	)
	if err != nil {
		return customoverlayentity.ChannelOverlayNil, err
	}

	_ = s.auditRecorder.RecordCreateOperation(
		ctx,
		audit.CreateOperation{
			Metadata: audit.OperationMetadata{
				System:    mappers.AuditSystemToTableName(gqlmodel.AuditLogSystemChannelsOverlays),
				ActorID:   &input.ActorID,
				ChannelID: &input.ChannelID,
				ObjectID:  lo.ToPtr(overlay.ID.String()),
			},
			NewValue: overlay,
		},
	)

	return s.modelToEntity(overlay), nil
}

type UpdateInput struct {
	ChannelID string
	ActorID   string
	Name      string
	Width     int
	Height    int
	InstaSave bool
	Layers    []UpdateLayerInput
}

func (s *Service) Update(ctx context.Context, id uuid.UUID, input UpdateInput) (
	customoverlayentity.ChannelOverlay,
	error,
) {
	dbOverlay, err := s.overlaysRepository.GetByID(ctx, id)
	if err != nil {
		return customoverlayentity.ChannelOverlayNil, err
	}

	if dbOverlay.ChannelID != input.ChannelID {
		return customoverlayentity.ChannelOverlayNil, apperrors.NewNotFoundError("Overlay with this ID was not found for your channel")
	}

	repoLayers := make([]channels_overlays.UpdateLayerInputWithID, len(input.Layers))
	for i, l := range input.Layers {
		repoLayers[i] = channels_overlays.UpdateLayerInputWithID{
			ID:   l.ID,
			Type: model.OverlayType(l.Type),
			Settings: model.OverlayLayerSettings{
				HtmlOverlayHTML:                    l.Settings.HtmlOverlayHTML,
				HtmlOverlayCSS:                     l.Settings.HtmlOverlayCSS,
				HtmlOverlayJS:                      l.Settings.HtmlOverlayJS,
				HtmlOverlayDataPollSecondsInterval: l.Settings.HtmlOverlayDataPollSecondsInterval,
				ImageUrl:                           l.Settings.ImageUrl,
			},
			PosX:                    l.PosX,
			PosY:                    l.PosY,
			Width:                   l.Width,
			Height:                  l.Height,
			Rotation:                l.Rotation,
			PeriodicallyRefetchData: l.PeriodicallyRefetchData,
			Locked:                  l.Locked,
			Visible:                 l.Visible,
			Opacity:                 l.Opacity,
		}
	}

	newOverlay, err := s.overlaysRepository.Update(
		ctx,
		id,
		channels_overlays.UpdateInput{
			Name:      input.Name,
			Width:     input.Width,
			Height:    input.Height,
			InstaSave: input.InstaSave,
			Layers:    repoLayers,
		},
	)
	if err != nil {
		return customoverlayentity.ChannelOverlayNil, err
	}

	_ = s.auditRecorder.RecordUpdateOperation(
		ctx,
		audit.UpdateOperation{
			Metadata: audit.OperationMetadata{
				System:    mappers.AuditSystemToTableName(gqlmodel.AuditLogSystemChannelsOverlays),
				ActorID:   &input.ActorID,
				ChannelID: &input.ChannelID,
				ObjectID:  lo.ToPtr(newOverlay.ID.String()),
			},
			NewValue: newOverlay,
			OldValue: dbOverlay,
		},
	)

	entityOverlay := s.modelToEntity(newOverlay)

	// Publish update event via WebSocket
	if err := s.wsRouter.Publish(
		CreateCustomOverlayWsRouterKey(input.ChannelID, id),
		entityOverlay,
	); err != nil {
		// Log error but don't fail the update
		// The overlay was updated successfully in the database
	}

	return entityOverlay, nil
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID, channelID, actorID string) error {
	dbOverlay, err := s.overlaysRepository.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if dbOverlay.ChannelID != channelID {
		return apperrors.NewNotFoundError("Overlay with this ID was not found for your channel")
	}

	if err := s.overlaysRepository.Delete(ctx, id); err != nil {
		return err
	}

	_ = s.auditRecorder.RecordDeleteOperation(
		ctx,
		audit.DeleteOperation{
			Metadata: audit.OperationMetadata{
				System:    mappers.AuditSystemToTableName(gqlmodel.AuditLogSystemChannelsOverlays),
				ActorID:   &actorID,
				ChannelID: &channelID,
				ObjectID:  lo.ToPtr(dbOverlay.ID.String()),
			},
			OldValue: dbOverlay,
		},
	)

	return nil
}

type ParseHtmlInput struct {
	ChannelID string
	Html      string
}

func (s *Service) ParseHtml(ctx context.Context, input ParseHtmlInput) (string, error) {
	parsedChannelID, err := uuid.Parse(input.ChannelID)
	if err != nil {
		return "", fmt.Errorf("cannot parse channel id: %w", err)
	}

	channel, err := s.channelsRepository.GetByID(ctx, parsedChannelID)
	if err != nil {
		return "", fmt.Errorf("cannot get channel: %w", err)
	}

	var channelTwitchUserID string
	if channel.TwitchUserID != nil {
		channelTwitchUserID = *channel.TwitchPlatformID
	}

	var channelPlatformID string
	if channel.TwitchPlatformID != nil {
		channelPlatformID = *channel.TwitchPlatformID
	}

	res, err := s.bus.Parser.ParseVariablesInText.Request(
		ctx, parser.ParseVariablesInTextRequest{
			ChannelID:           channelPlatformID,
			ChannelTwitchUserID: channelTwitchUserID,
			Text:                input.Html,
		},
	)
	if err != nil {
		return "", err
	}

	return res.Data.Text, nil
}
