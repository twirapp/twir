package top_donate_stream

import (
	"context"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
	model "github.com/twirapp/twir/libs/gomodels"
	channelservice "github.com/twirapp/twir/libs/services/channels"
	"gorm.io/gorm"
)

func getLatestDonateData(
	ctx context.Context,
	db *gorm.DB,
	channelService *channelservice.ChannelService,
	channelId string,
) *model.ChannelsEventsListItemData {
	channelUUID, err := uuid.Parse(channelId)
	if err != nil {
		return nil
	}

	streams, err := channelService.GetChannelStreams(ctx, channelUUID)
	if err != nil || len(streams) == 0 {
		return nil
	}

	entity := model.ChannelsEventsListItem{}
	if err := db.
		WithContext(ctx).
		Where(
			"channel_id = ? AND type = ? AND created_at >= ?",
			channelId,
			model.ChannelEventListItemTypeDonation,
			streams[0].StartedAt,
		).
		Order(`(data->>'donationAmount')::numeric DESC`).
		First(&entity).Error; err != nil {
		return nil
	}

	return entity.Data
}

var UserName = &types.Variable{
	Name:                "donations.top.stream.userName",
	Description:         lo.ToPtr("Stream top donate username"),
	CanBeUsedInRegistry: true,
	Handler: func(
		ctx context.Context,
		parseCtx *types.VariableParseContext,
		variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		result := &types.VariableHandlerResult{}

		data := getLatestDonateData(ctx, parseCtx.Services.Gorm, parseCtx.Services.ChannelService, parseCtx.Channel.ID)
		if data != nil {
			result.Result = data.DonationUsername
		}

		return result, nil
	},
}

var Amount = &types.Variable{
	Name:                "donations.top.stream.amount",
	Description:         lo.ToPtr("Stream top donate amount"),
	CanBeUsedInRegistry: true,
	Handler: func(
		ctx context.Context,
		parseCtx *types.VariableParseContext,
		variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		result := &types.VariableHandlerResult{}

		data := getLatestDonateData(ctx, parseCtx.Services.Gorm, parseCtx.Services.ChannelService, parseCtx.Channel.ID)
		if data != nil {
			result.Result = data.DonationAmount
		}

		return result, nil
	},
}

var Currency = &types.Variable{
	Name:                "donations.top.stream.currency",
	Description:         lo.ToPtr("Stream top donate currency"),
	CanBeUsedInRegistry: true,
	Handler: func(
		ctx context.Context,
		parseCtx *types.VariableParseContext,
		variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		result := &types.VariableHandlerResult{}

		data := getLatestDonateData(ctx, parseCtx.Services.Gorm, parseCtx.Services.ChannelService, parseCtx.Channel.ID)
		if data != nil {
			result.Result = data.DonationCurrency
		}

		return result, nil
	},
}
