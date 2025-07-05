package events

import (
	"context"

	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/samber/lo"
	"github.com/satont/twir/apps/api/internal/impl_deps"
	model "github.com/satont/twir/libs/gomodels"
	eventstypes "github.com/satont/twir/libs/types/types/events"
	"github.com/twirapp/twir/libs/api/messages/events"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

type Events struct {
	*impl_deps.Deps
}

func (c *Events) convertEntity(entity *model.Event) *events.Event {
	event := &events.Event{
		Id:          entity.ID,
		ChannelId:   entity.ChannelID,
		Type:        entity.Type.String(),
		RewardId:    entity.RewardID.Ptr(),
		CommandId:   entity.CommandID.Ptr(),
		KeywordId:   entity.KeywordID.Ptr(),
		Description: entity.Description.String,
		Enabled:     entity.Enabled,
		OnlineOnly:  entity.OnlineOnly,
		Operations:  make([]*events.Event_Operation, len(entity.Operations)),
	}

	for i, operation := range entity.Operations {
		event.Operations[i] = &events.Event_Operation{
			Type:           operation.Type.String(),
			Input:          operation.Input.Ptr(),
			Delay:          uint32(operation.Delay),
			Repeat:         uint32(operation.Repeat),
			UseAnnounce:    operation.UseAnnounce,
			TimeoutTime:    uint32(operation.TimeoutTime),
			TimeoutMessage: operation.TimeoutMessage.Ptr(),
			Target:         operation.Target.Ptr(),
			Filters:        make([]*events.Event_OperationFilter, len(operation.Filters)),
			Enabled:        operation.Enabled,
		}

		for j, filter := range operation.Filters {
			event.Operations[i].Filters[j] = &events.Event_OperationFilter{
				Type:  filter.Type.String(),
				Left:  filter.Left,
				Right: filter.Right,
			}
		}
	}

	return event
}

func (c *Events) EventsGetAll(ctx context.Context, _ *emptypb.Empty) (
	*events.GetAllResponse,
	error,
) {
	dashboardId := ctx.Value("dashboardId").(string)
	var evnts []*model.Event
	if err := c.Db.
		WithContext(ctx).
		Preload("Operations").
		Preload("Operations.Filters").
		Group(`"id"`).
		Where(`"channelId" = ?`, dashboardId).Find(&evnts).Error; err != nil {
		return nil, err
	}

	return &events.GetAllResponse{
		Events: lo.Map(
			evnts, func(entity *model.Event, _ int) *events.Event {
				return c.convertEntity(entity)
			},
		),
	}, nil
}

func (c *Events) EventsGetById(ctx context.Context, request *events.GetByIdRequest) (
	*events.Event,
	error,
) {
	dashboardId := ctx.Value("dashboardId").(string)
	entity := &model.Event{}
	if err := c.Db.
		WithContext(ctx).
		Preload("Operations").
		Preload("Operations.Filters").
		Where(`"id" = ? AND "channelId" = ?`, request.Id, dashboardId).First(entity).Error; err != nil {
		return nil, err
	}

	return c.convertEntity(entity), nil
}

func (c *Events) EventsCreate(ctx context.Context, request *events.CreateRequest) (
	*events.Event,
	error,
) {
	dashboardId := ctx.Value("dashboardId").(string)

	entity := &model.Event{
		ID:          uuid.New().String(),
		ChannelID:   dashboardId,
		Type:        model.EventType(request.Event.Type),
		RewardID:    null.StringFromPtr(request.Event.RewardId),
		CommandID:   null.StringFromPtr(request.Event.CommandId),
		KeywordID:   null.StringFromPtr(request.Event.KeywordId),
		Description: null.StringFrom(request.Event.Description),
		Enabled:     true,
		OnlineOnly:  request.Event.OnlineOnly,
		Operations:  make([]model.EventOperation, len(request.Event.Operations)),
	}

	for i, operation := range request.Event.Operations {
		entity.Operations[i] = model.EventOperation{
			ID:             uuid.New().String(),
			Type:           model.EventOperationType(operation.Type),
			Delay:          int(operation.Delay),
			Input:          null.StringFromPtr(operation.Input),
			Repeat:         int(operation.Repeat),
			Order:          i,
			UseAnnounce:    operation.UseAnnounce,
			TimeoutTime:    int(operation.TimeoutTime),
			TimeoutMessage: null.StringFromPtr(operation.TimeoutMessage),
			Target:         null.StringFromPtr(operation.Target),
			Enabled:        operation.Enabled,
			Filters:        make([]*model.EventOperationFilter, len(operation.Filters)),
		}

		for j, filter := range operation.Filters {
			entity.Operations[i].Filters[j] = &model.EventOperationFilter{
				Type:  eventstypes.EventOperationFilterType(filter.Type),
				Left:  filter.Left,
				Right: filter.Right,
			}
		}
	}

	if err := c.Db.WithContext(ctx).Create(entity).Error; err != nil {
		return nil, err
	}

	if err := c.ChannelsEventsWithOperationsCache.Invalidate(ctx, dashboardId); err != nil {
		return nil, err
	}

	return c.convertEntity(entity), nil
}

func (c *Events) EventsDelete(ctx context.Context, request *events.DeleteRequest) (
	*emptypb.Empty,
	error,
) {
	dashboardId := ctx.Value("dashboardId").(string)
	if err := c.Db.
		WithContext(ctx).
		Where(
			`"id" = ? AND "channelId" = ?`,
			request.Id,
			dashboardId,
		).Delete(&model.Event{}).Error; err != nil {
		return nil, err
	}

	if err := c.ChannelsEventsWithOperationsCache.Invalidate(ctx, dashboardId); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (c *Events) EventsUpdate(ctx context.Context, request *events.PutRequest) (
	*events.Event,
	error,
) {
	dashboardId := ctx.Value("dashboardId").(string)
	entity := &model.Event{}
	if err := c.Db.
		WithContext(ctx).
		Preload("Operations").
		Preload("Operations.Filters").
		Where(`"id" = ? AND "channelId" = ?`, request.Id, dashboardId).First(entity).Error; err != nil {
		return nil, err
	}

	entity.RewardID = null.StringFromPtr(request.Event.RewardId)
	entity.CommandID = null.StringFromPtr(request.Event.CommandId)
	entity.KeywordID = null.StringFromPtr(request.Event.KeywordId)
	entity.Description = null.StringFrom(request.Event.Description)
	entity.OnlineOnly = request.Event.OnlineOnly
	entity.Type = model.EventType(request.Event.Type)
	entity.Operations = make([]model.EventOperation, len(request.Event.Operations))

	for i, operation := range request.Event.Operations {
		entity.Operations[i] = model.EventOperation{
			ID:             uuid.New().String(),
			EventID:        entity.ID,
			Type:           model.EventOperationType(operation.Type),
			Delay:          int(operation.Delay),
			Input:          null.StringFromPtr(operation.Input),
			Repeat:         int(operation.Repeat),
			Order:          i,
			UseAnnounce:    operation.UseAnnounce,
			TimeoutTime:    int(operation.TimeoutTime),
			TimeoutMessage: null.StringFromPtr(operation.TimeoutMessage),
			Target:         null.StringFromPtr(operation.Target),
			Enabled:        operation.Enabled,
		}

		entity.Operations[i].Filters = make([]*model.EventOperationFilter, len(operation.Filters))

		for j, filter := range operation.Filters {
			entity.Operations[i].Filters[j] = &model.EventOperationFilter{
				ID:          uuid.New().String(),
				OperationID: entity.Operations[i].ID,
				Type:        eventstypes.EventOperationFilterType(filter.Type),
				Left:        filter.Left,
				Right:       filter.Right,
			}
		}
	}

	err := c.Db.Transaction(
		func(tx *gorm.DB) error {
			for _, operation := range entity.Operations {
				if err := tx.
					WithContext(ctx).
					Where(`"operationId" = ?`, operation.ID).
					Delete(&model.EventOperationFilter{}).Error; err != nil {
					return err
				}
			}

			if err := tx.WithContext(ctx).Where(
				`"eventId" = ?`,
				entity.ID,
			).Delete(&model.EventOperation{}).Error; err != nil {
				return err
			}

			return tx.WithContext(ctx).Save(entity).Error
		},
	)

	if err != nil {
		return nil, err
	}

	if err := c.ChannelsEventsWithOperationsCache.Invalidate(ctx, dashboardId); err != nil {
		return nil, err
	}

	return c.convertEntity(entity), nil
}

func (c *Events) EventsEnableOrDisable(
	ctx context.Context,
	request *events.PatchRequest,
) (*events.Event, error) {
	dashboardId := ctx.Value("dashboardId").(string)
	entity := &model.Event{}
	if err := c.Db.
		WithContext(ctx).
		Where(`"id" = ? AND "channelId" = ?`, request.Id, dashboardId).First(entity).Error; err != nil {
		return nil, err
	}

	entity.Enabled = request.Enabled
	if err := c.Db.WithContext(ctx).Save(entity).Error; err != nil {
		return nil, err
	}

	if err := c.ChannelsEventsWithOperationsCache.Invalidate(ctx, dashboardId); err != nil {
		return nil, err
	}

	return c.convertEntity(entity), nil
}
