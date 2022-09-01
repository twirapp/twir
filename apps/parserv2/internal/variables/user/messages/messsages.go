package usermessages

import (
	"strconv"
	types "tsuwari/parser/internal/types"
	variablescache "tsuwari/parser/internal/variablescache"
)

const Name = "user.messages"

func Handler(ctx *variablescache.VariablesCacheService, data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
	result := types.VariableHandlerResult{}

	if ctx.Cache.DbUserStats != nil {
		result.Result = strconv.Itoa(int(ctx.Cache.DbUserStats.Messages))
	} else {
		result.Result = "0"
	}

	return &result, nil
}
