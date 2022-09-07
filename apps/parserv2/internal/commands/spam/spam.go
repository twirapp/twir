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
		Visible:     true,
		Module:      lo.ToPtr("CHANNEL"),
	},
	Handler: func(ctx variables_cache.ExecutionContext) []string {
		count := 1
		params := strings.Split(*ctx.Text, " ")

		paramsLen := len(params)
		if paramsLen < 2 {
			return []string{"you have type count and message"}
		}

		newCount, err := strconv.Atoi(params[0])
		if err == nil {
			count = newCount
		}

		if count > 20 || count <= 0 {
			return []string{"count cannot be more then 20 and fewer then 1"}
		}

		message := strings.Join(params[1:], " ")
		result := make([]string, count)

		for i := 0; i < count; i++ {
			result[i] = message
		}

		return result
	},
}
