package channels_overlays

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/mappers"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	"github.com/twirapp/twir/libs/audit"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/parser"
	channels_overlays "github.com/twirapp/twir/libs/repositories/channels_overlays"
	"github.com/twirapp/twir/libs/repositories/channels_overlays/model"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	OverlaysRepository channels_overlays.Repository
	AuditRecorder      audit.Recorder
	Bus                *buscore.Bus
}

func New(opts Opts) *Service {
	return &Service{
		overlaysRepository: opts.OverlaysRepository,
		auditRecorder:      opts.AuditRecorder,
		bus:                opts.Bus,
	}
}

type Service struct {
	overlaysRepository channels_overlays.Repository
	auditRecorder      audit.Recorder
	bus                *buscore.Bus
}

func (s *Service) modelToEntity(m model.Overlay) entity.ChannelOverlay {
	layers := make([]entity.ChannelOverlayLayer, len(m.Layers))
	for i, l := range m.Layers {
		layers[i] = entity.ChannelOverlayLayer{
			ID:   l.ID,
			Type: entity.ChannelOverlayType(l.Type),
			Settings: entity.ChannelOverlayLayerSettings{
				HtmlOverlayHTML:                    l.Settings.HtmlOverlayHTML,
				HtmlOverlayCSS:                     l.Settings.HtmlOverlayCSS,
				HtmlOverlayJS:                      l.Settings.HtmlOverlayJS,
				HtmlOverlayDataPollSecondsInterval: l.Settings.HtmlOverlayDataPollSecondsInterval,
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
		}
	}

	return entity.ChannelOverlay{
		ID:        m.ID,
		ChannelID: m.ChannelID,
		Name:      m.Name,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		Width:     m.Width,
		Height:    m.Height,
		Layers:    layers,
	}
}

func (s *Service) GetManyByChannelID(ctx context.Context, channelID string) (
	[]entity.ChannelOverlay,
	error,
) {
	overlays, err := s.overlaysRepository.GetManyByChannelID(ctx, channelID)
	if err != nil {
		return nil, err
	}

	entities := make([]entity.ChannelOverlay, len(overlays))
	for i, o := range overlays {
		entities[i] = s.modelToEntity(o)
	}

	return entities, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (entity.ChannelOverlay, error) {
	overlay, err := s.overlaysRepository.GetByID(ctx, id)
	if err != nil {
		return entity.ChannelOverlayNil, err
	}

	return s.modelToEntity(overlay), nil
}

type CreateLayerInput struct {
	Type                    entity.ChannelOverlayType
	Settings                entity.ChannelOverlayLayerSettings
	PosX                    int
	PosY                    int
	Width                   int
	Height                  int
	Rotation                int
	PeriodicallyRefetchData bool
}

type CreateInput struct {
	ChannelID string
	ActorID   string
	Name      string
	Width     int
	Height    int
	Layers    []CreateLayerInput
}

func (s *Service) Create(ctx context.Context, input CreateInput) (entity.ChannelOverlay, error) {
	repoLayers := make([]channels_overlays.CreateLayerInput, len(input.Layers))
	for i, l := range input.Layers {
		repoLayers[i] = channels_overlays.CreateLayerInput{
			Type: model.OverlayType(l.Type),
			Settings: model.OverlayLayerSettings{
				HtmlOverlayHTML:                    l.Settings.HtmlOverlayHTML,
				HtmlOverlayCSS:                     l.Settings.HtmlOverlayCSS,
				HtmlOverlayJS:                      l.Settings.HtmlOverlayJS,
				HtmlOverlayDataPollSecondsInterval: l.Settings.HtmlOverlayDataPollSecondsInterval,
			},
			PosX:                    l.PosX,
			PosY:                    l.PosY,
			Width:                   l.Width,
			Height:                  l.Height,
			Rotation:                l.Rotation,
			PeriodicallyRefetchData: l.PeriodicallyRefetchData,
		}
	}

	overlay, err := s.overlaysRepository.Create(
		ctx,
		channels_overlays.CreateInput{
			ChannelID: input.ChannelID,
			Name:      input.Name,
			Width:     input.Width,
			Height:    input.Height,
			Layers:    repoLayers,
		},
	)
	if err != nil {
		return entity.ChannelOverlayNil, err
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
	Layers    []CreateLayerInput
}

func (s *Service) Update(ctx context.Context, id uuid.UUID, input UpdateInput) (
	entity.ChannelOverlay,
	error,
) {
	dbOverlay, err := s.overlaysRepository.GetByID(ctx, id)
	if err != nil {
		return entity.ChannelOverlayNil, err
	}

	if dbOverlay.ChannelID != input.ChannelID {
		return entity.ChannelOverlayNil, fmt.Errorf("overlay not found")
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
			},
			PosX:                    l.PosX,
			PosY:                    l.PosY,
			Width:                   l.Width,
			Height:                  l.Height,
			Rotation:                l.Rotation,
			PeriodicallyRefetchData: l.PeriodicallyRefetchData,
		}
	}

	newOverlay, err := s.overlaysRepository.Update(
		ctx,
		id,
		channels_overlays.UpdateInput{
			Name:   input.Name,
			Width:  input.Width,
			Height: input.Height,
			Layers: repoLayers,
		},
	)
	if err != nil {
		return entity.ChannelOverlayNil, err
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

	return s.modelToEntity(newOverlay), nil
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID, channelID, actorID string) error {
	dbOverlay, err := s.overlaysRepository.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if dbOverlay.ChannelID != channelID {
		return fmt.Errorf("overlay not found")
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
	res, err := s.bus.Parser.ParseVariablesInText.Request(
		ctx, parser.ParseVariablesInTextRequest{
			ChannelID: input.ChannelID,
			Text:      input.Html,
		},
	)
	if err != nil {
		return "", err
	}

	return res.Data.Text, nil
}
