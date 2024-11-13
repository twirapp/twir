package moderation

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/imroc/req/v3"
	"github.com/samber/lo"
	"github.com/satont/twir/apps/api/internal/impl_deps"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/api/messages/moderation"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Moderation struct {
	*impl_deps.Deps
}

func (c *Moderation) convertDbToGrpc(i model.ChannelModerationSettings) moderation.ItemWithId {
	return moderation.ItemWithId{
		Id: i.ID,
		Data: &moderation.Item{
			Type:                i.Type.String(),
			ChannelId:           i.ChannelID,
			Enabled:             i.Enabled,
			BanTime:             i.BanTime,
			BanMessage:          i.BanMessage,
			WarningMessage:      i.WarningMessage,
			CheckClips:          i.CheckClips,
			TriggerLength:       int32(i.TriggerLength),
			MaxPercentage:       int32(i.MaxPercentage),
			DenyList:            i.DenyList,
			DeniedChatLanguages: i.DeniedChatLanguages,
			ExcludedRoles:       i.ExcludedRoles,
			MaxWarnings:         int32(i.MaxWarnings),
			CreatedAt:           i.CreatedAt.String(),
			UpdatedAt:           i.UpdatedAt.String(),
		},
	}
}

func (c *Moderation) convertGrpcToDb(i *moderation.Item) model.ChannelModerationSettings {
	return model.ChannelModerationSettings{
		Type:                model.ModerationSettingsType(i.Type),
		ChannelID:           i.ChannelId,
		Enabled:             i.Enabled,
		BanTime:             i.BanTime,
		BanMessage:          i.BanMessage,
		WarningMessage:      i.WarningMessage,
		CheckClips:          i.CheckClips,
		TriggerLength:       int(i.TriggerLength),
		MaxPercentage:       int(i.MaxPercentage),
		DenyList:            i.DenyList,
		DeniedChatLanguages: i.DeniedChatLanguages,
		ExcludedRoles:       i.ExcludedRoles,
		MaxWarnings:         int(i.MaxWarnings),
	}
}

func (c *Moderation) clearCache(ctx context.Context, channelId string) {
	cacheKey := fmt.Sprintf("channels:%s:moderation_settings", channelId)

	c.Redis.Del(ctx, cacheKey)
}

func (c *Moderation) ModerationGetAll(
	ctx context.Context,
	_ *emptypb.Empty,
) (*moderation.GetAllResponse, error) {
	dashboardId := ctx.Value("dashboardId").(string)

	var entities []model.ChannelModerationSettings
	if err := c.Db.WithContext(ctx).Where(
		"channel_id = ?",
		dashboardId,
	).Order("created_at desc").Find(&entities).Error; err != nil {
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

func (c *Moderation) convertCreateGrpcTypeToDb(item *moderation.ItemCreateMessage) model.
	ChannelModerationSettings {
	return model.ChannelModerationSettings{
		Type:                model.ModerationSettingsType(item.Type),
		Enabled:             item.Enabled,
		BanTime:             item.BanTime,
		BanMessage:          item.BanMessage,
		WarningMessage:      item.WarningMessage,
		CheckClips:          item.CheckClips,
		TriggerLength:       int(item.TriggerLength),
		MaxPercentage:       int(item.MaxPercentage),
		DenyList:            item.DenyList,
		DeniedChatLanguages: item.DeniedChatLanguages,
		ExcludedRoles:       item.ExcludedRoles,
		MaxWarnings:         int(item.MaxWarnings),
	}
}

func (c *Moderation) ModerationCreate(
	ctx context.Context,
	req *moderation.CreateRequest,
) (*moderation.ItemWithId, error) {
	dashboardId := ctx.Value("dashboardId").(string)

	entity := c.convertCreateGrpcTypeToDb(req.Data)
	entity.ChannelID = dashboardId
	entity.ID = uuid.NewString()
	if err := c.Db.WithContext(ctx).Create(&entity).Error; err != nil {
		return nil, err
	}

	c.clearCache(ctx, dashboardId)

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

	c.clearCache(ctx, dashboardId)

	return &emptypb.Empty{}, nil
}

func (c *Moderation) ModerationUpdate(
	ctx context.Context,
	req *moderation.UpdateRequest,
) (*moderation.ItemWithId, error) {
	dashboardId := ctx.Value("dashboardId").(string)

	entity := c.convertCreateGrpcTypeToDb(req.Data)
	entity.ChannelID = dashboardId
	entity.ID = req.Id
	if err := c.Db.WithContext(ctx).Where(
		"channel_id = ? AND id = ?",
		dashboardId,
		req.Id,
	).Save(&entity).Error; err != nil {
		return nil, err
	}

	c.clearCache(ctx, dashboardId)

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
	entity.UpdatedAt = time.Now()

	if err := c.Db.WithContext(ctx).Where(
		"channel_id = ? AND id = ?",
		dashboardId,
		req.Id,
	).Save(&entity).Error; err != nil {
		return nil, err
	}

	c.clearCache(ctx, dashboardId)

	return &moderation.ItemWithId{
		Id:   entity.ID,
		Data: c.convertDbToGrpc(entity).Data,
	}, nil
}

type availableLanguage struct {
	Code    int    `json:"code"`
	Iso6933 int    `json:"iso_693_3"`
	Name    string `json:"name"`
}

func (c *Moderation) ModerationAvailableLanguages(
	ctx context.Context,
	_ *emptypb.Empty,
) (*moderation.AvailableLanguagesResponse, error) {
	var reqUrl string
	if c.Config.AppEnv == "production" {
		reqUrl = fmt.Sprint("http://language-detector:3012")
	} else {
		reqUrl = "http://localhost:3012"
	}

	var resp []availableLanguage
	res, err := req.R().SetSuccessResult(&resp).Get(reqUrl + "/languages")
	if err != nil {
		return nil, err
	}
	if !res.IsSuccessState() {
		return nil, errors.New("cannot get response")
	}

	langs := make([]*moderation.AvailableLanguagesResponse_Lang, len(resp))
	for i, lang := range resp {
		langs[i] = &moderation.AvailableLanguagesResponse_Lang{
			Code:      int32(lang.Code),
			Iso_693_3: int32(lang.Iso6933),
			Name:      lang.Name,
		}
	}

	return &moderation.AvailableLanguagesResponse{
		Langs: langs,
	}, nil
}
