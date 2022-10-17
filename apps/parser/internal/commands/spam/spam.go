package spam

import (
	"strconv"
	"strings"
	"tsuwari/parser/internal/types"
	variables_cache "tsuwari/parser/internal/variablescache"

	"github.com/samber/lo"
)

var Command = types.DefaultCommand{
	Command: types.Command{
		Name:        "spam",
		Description: lo.ToPtr("Spam into chat. Example usage: <b>!spam 5 https://tsuwari.tk</b>"),
		Permission:  "MODERATOR",
		Visible:     false,
		Module:      lo.ToPtr("CHANNEL"),
		IsReply:     false,
	},
	IsReply: lo.ToPtr(false),
	Handler: func(ctx variables_cache.ExecutionContext) *types.CommandsHandlerResult {
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
