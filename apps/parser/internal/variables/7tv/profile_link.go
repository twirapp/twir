package seventv

import (
	"context"
	"fmt"

	"github.com/samber/lo"
	"github.com/satont/twir/apps/parser/internal/types"
)

var ProfileLink = &types.Variable{
	Name:         "7tv.profile.link",
	Description:  lo.ToPtr("Link to 7tv profile"),
	CommandsOnly: true,
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		result := types.VariableHandlerResult{}

		profile, err := parseCtx.Services.SevenTvCache.Get(ctx, parseCtx.Channel.ID)
		if err != nil {
			result.Result = fmt.Sprintf("[Twir err] Failed to get 7tv profile: %s", err)
		} else {
			result.Result = fmt.Sprintf("https://7tv.app/users/%s", profile.User.Id)
		}

		return &result, nil
	},
}
