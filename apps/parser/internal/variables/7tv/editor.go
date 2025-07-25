package seventv

import (
	"context"
	"fmt"
	"strconv"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
)

var EditorForCount = &types.Variable{
	Name:         "7tv.editorfor.count",
	Description:  lo.ToPtr("Count of channels where user is editor"),
	CommandsOnly: true,
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		result := types.VariableHandlerResult{}

		profile, err := parseCtx.Cacher.GetSeventvProfileGetTwitchId(ctx, parseCtx.Channel.ID)
		if err != nil {
			result.Result = fmt.Sprintf("[Twir err] Failed to get 7tv profile: %s", err)
			return &result, nil
		}

		result.Result = strconv.Itoa(len(profile.EditorFor))

		return &result, nil
	},
}
