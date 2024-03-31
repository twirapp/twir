package spam

import (
	"context"
	"strconv"
	"strings"

	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/satont/twir/apps/parser/internal/types"
	model "github.com/satont/twir/libs/gomodels"
)

var Command = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "spam",
		Description: null.StringFrom("Spam into chat. Example usage: <b>!spam 5 Follow me on twitter"),
		RolesIDS:    pq.StringArray{model.ChannelRoleTypeModerator.String()},
		Module:      "MODERATION",
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		result := &types.CommandsHandlerResult{}

		count := 1
		params := strings.Split(*parseCtx.Text, " ")

		paramsLen := len(params)
		if paramsLen < 2 {
			result.Result = []string{"you have type count and message"}
			return result, nil
		}

		newCount, err := strconv.Atoi(params[0])
		if err == nil {
			count = newCount
		}

		if count > 10 || count <= 0 {
			result.Result = []string{"count cannot be more than 20 and fewer then 1"}
			return result, nil
		}

		message := strings.Join(params[1:], " ")

		for i := 0; i < count; i++ {
			result.Result = append(result.Result, message)
		}

		return result, nil
	},
}
