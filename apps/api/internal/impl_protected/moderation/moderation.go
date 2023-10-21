package moderation

import (
	"context"

	"github.com/samber/lo"
	"github.com/satont/twir/apps/api/internal/impl_deps"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/grpc/generated/api/moderation"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Moderation struct {
	*impl_deps.Deps
}

func (c *Moderation) convertDbToGrpc(i model.ChannelModerationSettings) moderation.ItemWithId {
	return moderation.ItemWithId{
		Id: i.ID,
		Data: &moderation.Item{
			Type:                  i.Type.String(),
			ChannelId:             i.ChannelID,
			Enabled:               i.Enabled,
			BanTime:               i.BanTime,
			BanMessage:            i.BanMessage,
			WarningMessage:        i.WarningMessage,
			CheckClips:            i.CheckClips,
			TriggerLength:         int32(i.TriggerLength),
			MaxPercentage:         int32(i.MaxPercentage),
			DenyList:              i.DenyList,
			AcceptedChatLanguages: i.AcceptedChatLanguages,
			ExcludedRoles:         i.ExcludedRoles,
			MaxWarnings:           int32(i.MaxWarnings),
			CreatedAt:             i.CreatedAt.String(),
			UpdatedAt:             i.UpdatedAt.String(),
		},
	}
}

func (c *Moderation) convertGrpcToDb(i *moderation.Item) model.ChannelModerationSettings {
	return model.ChannelModerationSettings{
		Type:                  model.ModerationSettingsType(i.Type),
		ChannelID:             i.ChannelId,
		Enabled:               i.Enabled,
		BanTime:               i.BanTime,
		BanMessage:            i.BanMessage,
		WarningMessage:        i.WarningMessage,
		CheckClips:            i.CheckClips,
		TriggerLength:         int(i.TriggerLength),
		MaxPercentage:         int(i.MaxPercentage),
		DenyList:              i.DenyList,
		AcceptedChatLanguages: i.AcceptedChatLanguages,
		ExcludedRoles:         i.ExcludedRoles,
		MaxWarnings:           int(i.MaxWarnings),
	}
}

func (c *Moderation) ModerationGetAll(ctx context.Context) (*moderation.GetAllResponse, error) {
	dashboardId := ctx.Value("dashboardId").(string)

	var entities []model.ChannelModerationSettings
	if err := c.Db.WithContext(ctx).Where(
		"channel_id = ?",
		dashboardId,
	).Find(&entities).Error; err != nil {
		return nil, err
	}

	return &moderation.GetAllResponse{
		Body: lo.Map(
			entities, func(item model.ChannelModerationSettings, _ int) *moderation.ItemWithId {
				converted := c.convertDbToGrpc(item)
				return &converted
			},
		),
	}, nil
}

func (c *Moderation) ModerationCreate(
	ctx context.Context,
	req *moderation.CreateRequest,
) (*moderation.ItemWithId, error) {
	dashboardId := ctx.Value("dashboardId").(string)

	entity := c.convertGrpcToDb(req.Data)
	entity.ChannelID = dashboardId
	if err := c.Db.WithContext(ctx).Create(&entity).Error; err != nil {
		return nil, err
	}

	return &moderation.ItemWithId{
		Id:   entity.ID,
		Data: c.convertDbToGrpc(entity).Data,
	}, nil
}

func (c *Moderation) ModerationDelete(
	ctx context.Context,
	req *moderation.DeleteRequest,
) (*emptypb.Empty, error) {
	dashboardId := ctx.Value("dashboardId").(string)

	if err := c.Db.WithContext(ctx).Where(
		"channel_id = ? AND id = ?",
		dashboardId,
		req.Id,
	).Delete(&model.ChannelModerationSettings{}).Error; err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (c *Moderation) ModerationUpdate(
	ctx context.Context,
	req *moderation.UpdateRequest,
) (*moderation.ItemWithId, error) {
	dashboardId := ctx.Value("dashboardId").(string)

	entity := c.convertGrpcToDb(req.Data)
	entity.ChannelID = dashboardId
	if err := c.Db.WithContext(ctx).Where(
		"channel_id = ? AND id = ?",
		dashboardId,
		req.Id,
	).Updates(&entity).Error; err != nil {
		return nil, err
	}

	return &moderation.ItemWithId{
		Id:   entity.ID,
		Data: c.convertDbToGrpc(entity).Data,
	}, nil
}

func (c *Moderation) ModerationEnableOrDisable(
	ctx context.Context,
	req *moderation.PatchRequest,
) (*moderation.ItemWithId, error) {
	dashboardId := ctx.Value("dashboardId").(string)

	entity := model.ChannelModerationSettings{}
	if err := c.Db.WithContext(ctx).Where(
		"channel_id = ? AND id = ?",
		dashboardId,
		req.Id,
	).First(&entity).Error; err != nil {
		return nil, err
	}

	entity.Enabled = req.Enabled

	if err := c.Db.WithContext(ctx).Where(
		"channel_id = ? AND id = ?",
		dashboardId,
		req.Id,
	).Updates(&entity).Error; err != nil {
		return nil, err
	}

	return &moderation.ItemWithId{
		Id:   entity.ID,
		Data: c.convertDbToGrpc(entity).Data,
	}, nil
}
