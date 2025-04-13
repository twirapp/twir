package seventv

import (
	"context"
	"fmt"

	"github.com/samber/lo"
	"github.com/satont/twir/apps/parser/internal/types"
	"github.com/satont/twir/apps/parser/pkg/helpers"
)

var ProfileCreatedAt = &types.Variable{
	Name:         "7tv.profile.createdAt",
	Description:  lo.ToPtr("Date when profile created on 7tv"),
	CommandsOnly: true,
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		result := types.VariableHandlerResult{}

		profile, err := parseCtx.Cacher.GetSeventvProfileGetTwitchId(ctx, parseCtx.Channel.ID)
		if err != nil {
			result.Result = fmt.Sprintf("[Twir err] Failed to get 7tv profile: %s", err)
		} else {
			result.Result = helpers.Duration(profile.MainConnection.LinkedAt, &helpers.DurationOpts{})
		}

		return &result, nil
	},
}
