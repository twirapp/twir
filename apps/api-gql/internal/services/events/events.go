package events

import (
	"context"
	"errors"
	"fmt"

	deprecatedmodel "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/logger"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	"github.com/twirapp/twir/libs/repositories/events"
	"github.com/twirapp/twir/libs/repositories/events/model"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	EventsRepository events.Repository
	Logger           logger.Logger
	Cacher           *generic_cacher.GenericCacher[[]deprecatedmodel.Event]
}

func New(opts Opts) *Service {
	return &Service{
		eventsRepository: opts.EventsRepository,
		logger:           opts.Logger,
		cacher:           opts.Cacher,
	}
}

type Service struct {
	eventsRepository events.Repository
	logger           logger.Logger
	cacher           *generic_cacher.GenericCacher[[]deprecatedmodel.Event]
}

func (s *Service) mapToEntity(m model.Event) entity.Event {
	operations := make([]entity.EventOperation, 0, len(m.Operations))
	for _, op := range m.Operations {
		filters := make([]entity.EventOperationFilter, 0, len(op.Filters))
		for _, f := range op.Filters {
			filters = append(
				filters, entity.EventOperationFilter{
					ID:    f.ID,
					Type:  f.Type,
					Left:  f.Left,
					Right: f.Right,
				},
			)
		}

		operations = append(
			operations, entity.EventOperation{
				ID:             op.ID,
				Type:           entity.EventOperationType(op.Type),
				Input:          op.Input,
				Delay:          op.Delay,
				Repeat:         op.Repeat,
				UseAnnounce:    op.UseAnnounce,
				TimeoutTime:    op.TimeoutTime,
				TimeoutMessage: op.TimeoutMessage,
				Target:         op.Target,
				Enabled:        op.Enabled,
				Filters:        filters,
			},
		)
	}

	return entity.Event{
		ID:          m.ID,
		ChannelID:   m.ChannelID,
		Type:        entity.EventType(m.Type),
		RewardID:    m.RewardID,
		CommandID:   m.CommandID,
		KeywordID:   m.KeywordID,
		Description: m.Description,
		Enabled:     m.Enabled,
		OnlineOnly:  m.OnlineOnly,
		Operations:  operations,
	}
}

func (s *Service) GetAll(ctx context.Context, channelID string) ([]entity.Event, error) {
	channelEvents, err := s.eventsRepository.GetManyByChannelID(ctx, channelID)
	if err != nil {
		return nil, fmt.Errorf("failed to get channelEvents: %w", err)
	}

	result := make([]entity.Event, 0, len(channelEvents))
	for _, e := range channelEvents {
		result = append(result, s.mapToEntity(e))
	}

	return result, nil
}

func (s *Service) GetByID(ctx context.Context, id string) (entity.Event, error) {
	event, err := s.eventsRepository.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, events.ErrNotFound) {
			return entity.EventNil, fmt.Errorf("event not found")
		}
		return entity.EventNil, fmt.Errorf("failed to get event: %w", err)
	}

	return s.mapToEntity(event), nil
}

func (s *Service) Create(ctx context.Context, input CreateInput) (entity.Event, error) {
	channelEvents, err := s.eventsRepository.GetManyByChannelID(ctx, input.ChannelID)
	if err != nil {
		return entity.EventNil, fmt.Errorf("failed to get channelEvents: %w", err)
	}

	if len(channelEvents) >= 10 {
		return entity.EventNil, fmt.Errorf("maximum number of channelEvents reached")
	}

	repoInput := events.CreateInput{
		ChannelID:   input.ChannelID,
		Type:        model.EventType(input.Type),
		RewardID:    input.RewardID,
		CommandID:   input.CommandID,
		KeywordID:   input.KeywordID,
		Description: input.Description,
		Enabled:     input.Enabled,
		OnlineOnly:  input.OnlineOnly,
		Operations:  make([]events.OperationInput, 0, len(input.Operations)),
	}

	for _, op := range input.Operations {
		filters := make([]events.OperationFilterInput, 0, len(op.Filters))
		for _, f := range op.Filters {
			filters = append(
				filters, events.OperationFilterInput{
					Type:  f.Type,
					Left:  f.Left,
					Right: f.Right,
				},
			)
		}

		repoInput.Operations = append(
			repoInput.Operations, events.OperationInput{
				Type:           model.EventOperationType(op.Type),
				Input:          op.Input,
				Delay:          op.Delay,
				Repeat:         op.Repeat,
				UseAnnounce:    op.UseAnnounce,
				TimeoutTime:    op.TimeoutTime,
				TimeoutMessage: op.TimeoutMessage,
				Target:         op.Target,
				Enabled:        op.Enabled,
				Filters:        filters,
			},
		)
	}

	event, err := s.eventsRepository.Create(ctx, repoInput)
	if err != nil {
		return entity.EventNil, fmt.Errorf("failed to create event: %w", err)
	}

	if err := s.cacher.Invalidate(ctx, input.ChannelID); err != nil {
		return entity.EventNil, fmt.Errorf("failed to invalidate cache: %w", err)
	}

	return s.mapToEntity(event), nil
}

func (s *Service) Update(ctx context.Context, id string, input UpdateInput) (entity.Event, error) {
	var convertedType *model.EventType
	if input.Type != nil {
		convertedType = (*model.EventType)(input.Type)
	}

	repoInput := events.UpdateInput{
		Type:        convertedType,
		RewardID:    input.RewardID,
		CommandID:   input.CommandID,
		KeywordID:   input.KeywordID,
		Description: input.Description,
		Enabled:     input.Enabled,
		OnlineOnly:  input.OnlineOnly,
	}

	if input.Operations != nil {
		operations := make([]events.OperationInput, 0, len(*input.Operations))
		for _, op := range *input.Operations {
			filters := make([]events.OperationFilterInput, 0, len(op.Filters))
			for _, f := range op.Filters {
				filters = append(
					filters, events.OperationFilterInput{
						Type:  f.Type,
						Left:  f.Left,
						Right: f.Right,
					},
				)
			}

			operations = append(
				operations, events.OperationInput{
					Type:           model.EventOperationType(op.Type),
					Input:          op.Input,
					Delay:          op.Delay,
					Repeat:         op.Repeat,
					UseAnnounce:    op.UseAnnounce,
					TimeoutTime:    op.TimeoutTime,
					TimeoutMessage: op.TimeoutMessage,
					Target:         op.Target,
					Enabled:        op.Enabled,
					Filters:        filters,
				},
			)
		}
		repoInput.Operations = &operations
	}

	event, err := s.eventsRepository.Update(ctx, id, repoInput)
	if err != nil {
		if errors.Is(err, events.ErrNotFound) {
			return entity.EventNil, fmt.Errorf("event not found")
		}
		return entity.EventNil, fmt.Errorf("failed to update event: %w", err)
	}

	if err := s.cacher.Invalidate(ctx, input.ChannelID); err != nil {
		return entity.EventNil, fmt.Errorf("failed to invalidate cache: %w", err)
	}

	return s.mapToEntity(event), nil
}

func (s *Service) Delete(ctx context.Context, id, channelID string) error {
	err := s.eventsRepository.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete event: %w", err)
	}

	if err := s.cacher.Invalidate(ctx, channelID); err != nil {
		return fmt.Errorf("failed to invalidate cache: %w", err)
	}

	return nil
}

type CreateInput struct {
	ChannelID   string
	Type        entity.EventType
	RewardID    *string
	CommandID   *string
	KeywordID   *string
	Description string
	Enabled     bool
	OnlineOnly  bool
	Operations  []OperationInput
}

type UpdateInput struct {
	ChannelID   string
	Type        *entity.EventType
	RewardID    *string
	CommandID   *string
	KeywordID   *string
	Description *string
	Enabled     *bool
	OnlineOnly  *bool
	Operations  *[]OperationInput
}

type OperationInput struct {
	Type           entity.EventOperationType
	Input          *string
	Delay          int
	Repeat         int
	UseAnnounce    bool
	TimeoutTime    int
	TimeoutMessage *string
	Target         *string
	Enabled        bool
	Filters        []OperationFilterInput
}

type OperationFilterInput struct {
	Type  string
	Left  string
	Right string
}
