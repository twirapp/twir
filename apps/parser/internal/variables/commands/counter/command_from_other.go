package command_counter

import (
	"fmt"
	"github.com/samber/do"
	"github.com/satont/tsuwari/apps/parser/internal/di"
	"gorm.io/gorm"

	"github.com/satont/tsuwari/apps/parser/internal/types"
	variables_cache "github.com/satont/tsuwari/apps/parser/internal/variablescache"
	model "github.com/satont/tsuwari/libs/gomodels"

	"github.com/samber/lo"
)

var CommandVariableFromOther = types.Variable{
	Name:         "command.counter.fromother",
	Description:  lo.ToPtr("Counter saying how many times OTHER command was used"),
	CommandsOnly: lo.ToPtr(true),
	Example:      lo.ToPtr("command.counter.fromother|commandName"),
	Handler: func(ctx *variables_cache.VariablesCacheService, data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
		result := &types.VariableHandlerResult{}
		db := do.MustInvoke[gorm.DB](di.Provider)

		if data.Params == nil {
			result.Result = "Have not passed params to variable. "
			return result, nil
		}

		cmd := model.ChannelsCommands{}
		err := db.
			Where(`"channelId" = ? AND "name" = ?`, ctx.ChannelId, *data.Params).
			First(&cmd).Error

		if err != nil || cmd.ID == "" {
			result.Result = fmt.Sprintf(`Command with name "%s" not found`, *data.Params)
			return result, nil
		}

		count, err := getCount(cmd.ID, nil)
		if err != nil {
			result.Result = "cannot get count"
			return result, nil
		}
		result.Result = count

		return result, nil
	},
}
