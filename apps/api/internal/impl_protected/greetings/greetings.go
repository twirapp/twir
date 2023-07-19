package greetings

import (
	"context"

	"github.com/samber/lo"
	"github.com/satont/twir/apps/api/internal/impl_deps"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/grpc/generated/api/greetings"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Greetings struct {
	*impl_deps.Deps
}

func (c *Greetings) convertEntity(entity *model.ChannelsGreetings) *greetings.Greeting {
	return &greetings.Greeting{
		Id:        entity.ID,
		ChannelId: entity.ChannelID,
		UserId:    entity.UserID,
		Enabled:   entity.Enabled,
		Text:      entity.Text,
		IsReply:   entity.IsReply,
		Processed: entity.Processed,
	}
}

func (c *Greetings) GreetingsGetAll(ctx context.Context, _ *emptypb.Empty) (*greetings.GetAllResponse, error) {
	dashboardId := ctx.Value("dashboardId").(string)
	var dbGreetings []*model.ChannelsGreetings
	err := c.Db.WithContext(ctx).Where(`"channelId" = ?`, dashboardId).Find(&dbGreetings).Error
	if err != nil {
		return nil, err
	}

	return &greetings.GetAllResponse{
		Greetings: lo.Map(
			dbGreetings, func(entity *model.ChannelsGreetings, _ int) *greetings.Greeting {
				return c.convertEntity(entity)
			},
		),
	}, nil
}

func (c *Greetings) GreetingsGetById(
	ctx context.Context,
	request *greetings.GetByIdRequest,
) (*greetings.Greeting, error) {
	dashboardId := ctx.Value("dashboardId").(string)
	var dbGreetings model.ChannelsGreetings
	err := c.Db.WithContext(ctx).Where(`"channelId" = ? AND "id" = ?`, dashboardId, request.Id).First(&dbGreetings).Error
	if err != nil {
		return nil, err
	}

	return c.convertEntity(&dbGreetings), nil
}

func (c *Greetings) GreetingsCreate(
	ctx context.Context,
	request *greetings.CreateRequest,
) (*greetings.Greeting, error) {
	dashboardId := ctx.Value("dashboardId").(string)
	entity := &model.ChannelsGreetings{
		ChannelID: dashboardId,
		UserID:    request.UserId,
		Enabled:   request.Enabled,
		Text:      request.Text,
		IsReply:   request.IsReply,
		Processed: false,
	}
	if err := c.Db.WithContext(ctx).Create(entity).Error; err != nil {
		return nil, err
	}

	return c.convertEntity(entity), nil
}

func (c *Greetings) GreetingsDelete(ctx context.Context, request *greetings.DeleteRequest) (*emptypb.Empty, error) {
	dashboardId := ctx.Value("dashboardId").(string)
	err := c.Db.
		WithContext(ctx).
		Where(`"channelId" = ? AND "id" = ?`, dashboardId, request.Id).
		Delete(&model.ChannelsGreetings{}).Error
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (c *Greetings) GreetingsUpdate(ctx context.Context, request *greetings.PutRequest) (*greetings.Greeting, error) {
	dashboardId := ctx.Value("dashboardId").(string)
	entity := &model.ChannelsGreetings{
		ChannelID: dashboardId,
		UserID:    request.Greeting.UserId,
		Enabled:   request.Greeting.Enabled,
		Text:      request.Greeting.Text,
		IsReply:   request.Greeting.IsReply,
		Processed: false,
	}
	err := c.Db.
		WithContext(ctx).
		Where(`"channelId" = ? AND "id" = ?`, dashboardId, request.Id).
		Updates(entity).Error
	if err != nil {
		return nil, err
	}

	return c.convertEntity(entity), nil
}

func (c *Greetings) GreetingsEnableOrDisable(
	ctx context.Context,
	request *greetings.PatchRequest,
) (*greetings.Greeting, error) {
	dashboardId := ctx.Value("dashboardId").(string)
	entity := &model.ChannelsGreetings{}
	err := c.Db.
		WithContext(ctx).
		Find(`"channelId" = ? AND "id" = ?`, dashboardId, request.Id).Error
	if err != nil {
		return nil, err
	}

	entity.Enabled = request.Enabled
	if err := c.Db.WithContext(ctx).Save(entity).Error; err != nil {
		return nil, err
	}

	return c.convertEntity(entity), nil
}
