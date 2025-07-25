package top_donate_stream

import (
	"context"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
	model "github.com/twirapp/twir/libs/gomodels"
	"gorm.io/gorm"
)

func getLatestDonateData(
	ctx context.Context,
	db *gorm.DB,
	channelId string,
) *model.ChannelsEventsListItemData {
	stream := model.ChannelsStreams{}
	if err := db.
		WithContext(ctx).
		Where(
			`"userId" = ?`,
			channelId,
		).First(&stream).Error; err != nil {
		return nil
	}

	entity := model.ChannelsEventsListItem{}
	if err := db.
		WithContext(ctx).
		Where(
			"channel_id = ? AND type = ? AND created_at >= ?",
			channelId,
			model.ChannelEventListItemTypeDonation,
			stream.StartedAt,
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

		data := getLatestDonateData(ctx, parseCtx.Services.Gorm, parseCtx.Channel.ID)
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

		data := getLatestDonateData(ctx, parseCtx.Services.Gorm, parseCtx.Channel.ID)
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

		data := getLatestDonateData(ctx, parseCtx.Services.Gorm, parseCtx.Channel.ID)
		if data != nil {
			result.Result = data.DonationCurrency
		}

		return result, nil
	},
}
