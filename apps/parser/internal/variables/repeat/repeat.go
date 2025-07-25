package repeat

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
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

		result := &types.VariableHandlerResult{
			Result: fmt.Sprintf("__REPEAT_MARKER_%d__", repeatCount),
		}

		return result, nil
	},
}
