package repeat

import (
	"context"
	"math"
	"strconv"
	"strings"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
	deprecatedgormmodel "github.com/twirapp/twir/libs/gomodels"
)

var Variable = &types.Variable{
	Name:        "repeat",
	Description: lo.ToPtr("If you use $(repeat) in response, then it will be repeated as many times as user typed in chat. Max repeat is 20"),
	Example:     lo.ToPtr("repeat"),
	Priority:    math.MaxInt,
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		repeatCount := 1
		splittedRawText := strings.Fields(parseCtx.RawText)
		if len(splittedRawText) >= 2 {
			parsedCount, err := strconv.Atoi(splittedRawText[1])
			if err == nil {
				repeatCount = parsedCount
			}
		}

		if repeatCount > 20 {
			repeatCount = 20
		}

		var hasAccess bool
		for _, r := range parseCtx.Sender.Roles {
			if r.Type == deprecatedgormmodel.ChannelRoleTypeBroadcaster || r.Type == deprecatedgormmodel.ChannelRoleTypeModerator || parseCtx.Sender.DbUser.IsBotAdmin {
				hasAccess = true
				break
			}
		}

		if !hasAccess {
			repeatCount = 1
		}

		result := &types.VariableHandlerResult{
			Result:               "",
			RepeatVariableResult: &repeatCount,
		}

		return result, nil
	},
}
