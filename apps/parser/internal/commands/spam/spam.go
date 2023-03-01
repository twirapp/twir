package spam

import (
	"github.com/guregu/null"
	"github.com/lib/pq"
	model "github.com/satont/tsuwari/libs/gomodels"
	"strconv"
	"strings"

	"github.com/satont/tsuwari/apps/parser/internal/types"
	variables_cache "github.com/satont/tsuwari/apps/parser/internal/variablescache"
)

var Command = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "spam",
		Description: null.StringFrom("Spam into chat. Example usage: <b>!spam 5 https://tsuwari.tk</b>"),
		RolesIDS:    pq.StringArray{model.ChannelRoleTypeModerator.String()},
		Module:      "MODERATION",
	},
	Handler: func(ctx *variables_cache.ExecutionContext) *types.CommandsHandlerResult {
		result := &types.CommandsHandlerResult{}

		count := 1
		params := strings.Split(*ctx.Text, " ")

		paramsLen := len(params)
		if paramsLen < 2 {
			result.Result = []string{"you have type count and message"}
			return result
		}

		newCount, err := strconv.Atoi(params[0])
		if err == nil {
			count = newCount
		}

		if count > 20 || count <= 0 {
			result.Result = []string{"count cannot be more then 20 and fewer then 1"}
			return result
		}

		message := strings.Join(params[1:], " ")

		for i := 0; i < count; i++ {
			result.Result = append(result.Result, message)
		}

		return result
	},
}
