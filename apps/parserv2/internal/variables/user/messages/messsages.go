package usermessages

import (
	"strconv"
	types "tsuwari/parser/internal/types"
	variablescache "tsuwari/parser/internal/variablescache"
)

const Name = "user.messages"
const Description = "User messages"

func Handler(ctx *variablescache.VariablesCacheService, data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
	result := types.VariableHandlerResult{}

	dbUser := ctx.GetGbUser()
	if dbUser != nil {
		result.Result = strconv.Itoa(int(dbUser.Messages))
	} else {
		result.Result = "0"
	}

	return &result, nil
}
