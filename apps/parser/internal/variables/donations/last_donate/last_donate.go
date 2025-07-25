package last_donate

import (
	"context"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
	model "github.com/twirapp/twir/libs/gomodels"
	"gorm.io/gorm"
)

func getLatestDonateData(
	ctx context.Context, db *gorm.DB,
	channelId string,
) *model.ChannelsEventsListItemData {
	entity := model.ChannelsEventsListItem{}
	if err := db.
		WithContext(ctx).
		Where(
			"channel_id = ? AND type = ?",
			channelId,
			model.ChannelEventListItemTypeDonation,
		).
		Order("created_at DESC").
		First(&entity).Error; err != nil {
		return nil
	}

	return entity.Data
}

var UserName = &types.Variable{
	Name:                "donations.latest.username",
	Description:         lo.ToPtr("Latest donate username"),
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
	Name:                "donations.latest.amount",
	Description:         lo.ToPtr("Latest donate amount"),
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
	Name:                "donations.latest.currency",
	Description:         lo.ToPtr("Latest donate currency"),
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
