package keywords

import (
	"context"

	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/samber/lo"
	"github.com/satont/twir/apps/api/internal/impl_deps"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/grpc/generated/api/keywords"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Keywords struct {
	*impl_deps.Deps
}

func (c *Keywords) convertEntity(entity *model.ChannelsKeywords) *keywords.Keyword {
	return &keywords.Keyword{
		Id:        entity.ID,
		ChannelId: entity.ChannelID,
		Text:      entity.Text,
		Response:  entity.Response,
		Enabled:   entity.Enabled,
		Cooldown:  int32(entity.Cooldown),
		IsReply:   entity.IsReply,
		IsRegular: entity.IsRegular,
		Usages:    int32(entity.Usages),
	}
}

func (c *Keywords) KeywordsGetAll(ctx context.Context, _ *emptypb.Empty) (*keywords.GetAllResponse, error) {
	dashboardId := ctx.Value("dashboardId").(string)
	var entities []*model.ChannelsKeywords
	if err := c.Db.WithContext(ctx).Where(`"channelId" = ?`, dashboardId).Find(&entities).Error; err != nil {
		return nil, err
	}

	return &keywords.GetAllResponse{
		Keywords: lo.Map(
			entities, func(item *model.ChannelsKeywords, _ int) *keywords.Keyword {
				return c.convertEntity(item)
			},
		),
	}, nil
}

func (c *Keywords) KeywordsGetById(ctx context.Context, request *keywords.GetByIdRequest) (*keywords.Keyword, error) {
	keyword := &model.ChannelsKeywords{}
	if err := c.Db.WithContext(ctx).Where("id = ?", request.Id).First(keyword).Error; err != nil {
		return nil, err
	}

	return c.convertEntity(keyword), nil
}

func (c *Keywords) KeywordsCreate(ctx context.Context, request *keywords.CreateRequest) (*keywords.Keyword, error) {
	dashboardId := ctx.Value("dashboardId").(string)
	keyword := &model.ChannelsKeywords{
		ID:               uuid.New().String(),
		ChannelID:        dashboardId,
		Text:             request.Text,
		Response:         request.Response,
		Enabled:          true,
		Cooldown:         int(request.Cooldown),
		CooldownExpireAt: null.Time{},
		IsReply:          request.IsReply,
		IsRegular:        request.IsRegular,
		Usages:           0,
	}

	if err := c.Db.WithContext(ctx).Create(keyword).Error; err != nil {
		return nil, err
	}

	return c.convertEntity(keyword), nil
}

func (c *Keywords) KeywordsDelete(ctx context.Context, request *keywords.DeleteRequest) (*emptypb.Empty, error) {
	if err := c.Db.WithContext(ctx).Where("id = ?", request.Id).Delete(&model.ChannelsKeywords{}).Error; err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (c *Keywords) KeywordsUpdate(ctx context.Context, request *keywords.PutRequest) (*keywords.Keyword, error) {
	keyword := &model.ChannelsKeywords{}
	if err := c.Db.WithContext(ctx).Where("id = ?", request.Id).First(keyword).Error; err != nil {
		return nil, err
	}

	keyword.Text = request.Keyword.Text
	keyword.Response = request.Keyword.Response
	keyword.Enabled = request.Keyword.Enabled
	keyword.Cooldown = int(request.Keyword.Cooldown)
	keyword.CooldownExpireAt = null.Time{}
	keyword.IsReply = request.Keyword.IsReply
	keyword.IsRegular = request.Keyword.IsRegular

	if err := c.Db.WithContext(ctx).Save(keyword).Error; err != nil {
		return nil, err
	}

	return c.convertEntity(keyword), nil
}

func (c *Keywords) KeywordsEnableOrDisable(
	ctx context.Context,
	request *keywords.PatchRequest,
) (*keywords.Keyword, error) {
	keyword := &model.ChannelsKeywords{}
	if err := c.Db.WithContext(ctx).Where("id = ?", request.Id).First(keyword).Error; err != nil {
		return nil, err
	}

	keyword.Enabled = request.Enabled

	if err := c.Db.WithContext(ctx).Save(keyword).Error; err != nil {
		return nil, err
	}

	return c.convertEntity(keyword), nil
}
