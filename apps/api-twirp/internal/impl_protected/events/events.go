package events

import (
	"context"
	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/api-twirp/internal/impl_deps"
	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/grpc/generated/api/events"
	"google.golang.org/protobuf/types/known/emptypb"
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
			Delay:          uint64(operation.Delay),
			Repeat:         uint64(operation.Repeat),
			UseAnnounce:    operation.UseAnnounce,
			TimeoutTime:    uint64(operation.TimeoutTime),
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

func (c *Events) EventsGetAll(ctx context.Context, _ *emptypb.Empty) (*events.GetAllResponse, error) {
	dashboardId := ctx.Value("dashboardId").(string)
	var evnts []*model.Event
	if err := c.Db.
		WithContext(ctx).
		Preload("Operations").
		Preload("Operations.Filters").
		Where(`"channelId" = ?`, dashboardId).Find(&evnts).Error; err != nil {
		return nil, err
	}

	return &events.GetAllResponse{
		Events: lo.Map(evnts, func(entity *model.Event, _ int) *events.Event {
			return c.convertEntity(entity)
		}),
	}, nil
}

func (c *Events) EventsGetById(ctx context.Context, request *events.GetByIdRequest) (*events.Event, error) {
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

func (c *Events) EventsCreate(ctx context.Context, request *events.CreateRequest) (*events.Event, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Events) EventsDelete(ctx context.Context, request *events.DeleteRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Events) EventsUpdate(ctx context.Context, request *events.PutRequest) (*events.Event, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Events) EventsEnableOrDisable(ctx context.Context, request *events.PatchRequest) (*events.Event, error) {
	//TODO implement me
	panic("implement me")
}
