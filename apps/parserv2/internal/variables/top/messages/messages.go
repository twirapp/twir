package messages

import (
	"fmt"
	"strconv"
	"strings"
	"tsuwari/parser/internal/types"
	"tsuwari/parser/internal/variables/top"
	"tsuwari/parser/internal/variablescache"

	"github.com/samber/lo"
)

const Name = "top.messages"

func Handler(ctx *variablescache.VariablesCacheService, data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
	result := &types.VariableHandlerResult{}
	var page int = 1

	if ctx.Context.Text != nil {
		p, err := strconv.Atoi(*ctx.Context.Text)
		if err != nil {
			page = p
		}

		if page <= 0 {
			page = 1
		}
	}

	fmt.Println(page)
	topUsers := top.GetTop(ctx, ctx.Context.ChannelId, "messages", &page)

	if topUsers == nil || len(*topUsers) == 0 {
		return result, nil
	}

	mappedTop := lo.Map(*topUsers, func(user *top.UserStats, idx int) string {
		return fmt.Sprintf(
			"%v. %s — %v",
			(idx+1)+(page-1)*10,
			user.UserName,
			user.Value,
		)
	})

	result.Result = strings.Join(mappedTop, ", ")
	return result, nil
}
